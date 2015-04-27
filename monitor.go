package main

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/warmans/fluentd-api-client/monitoring"
)

type FMonitor struct {
	Hub                *Hub
	Hosts              []*monitoring.Host
	States             []PluginState
	StateChecks        []PluginStateVisitor
	StatesLock         sync.RWMutex
	PushTickSeconds    int
	HistoryTickSeconds int
	HistorySize        int
}

func (mon *FMonitor) AddStateCheck(vi PluginStateVisitor) {
	mon.StateChecks = append(mon.StateChecks, vi)
}

func (mon *FMonitor) Run() {

	go func() {
		updateTicker := time.NewTicker(time.Second * time.Duration(mon.HistoryTickSeconds))
		for range updateTicker.C {

			//update all hosts in parallel
			for _, host := range mon.Hosts {
				go func(host *monitoring.Host) {
					host.Update()
					mon.AddState(host)
				}(host)
			}
		}
	}()

	broadcastTicker := time.NewTicker(time.Second * time.Duration(mon.PushTickSeconds))
	for range broadcastTicker.C {

		//Create a JSON payload
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)

		if err := encoder.Encode(mon.GetStateUpdate()); err != nil {
			log.Print("Failed to encode host stats payload")
			continue
		}

		//Broadcast to all clients
		mon.Hub.Broadcast(buffer.Bytes())
	}
}

func (mon *FMonitor) GetStateUpdate() []PluginState {
	mon.StatesLock.RLock()
	defer mon.StatesLock.RUnlock()

	return mon.States
}

func (mon *FMonitor) AddState(host *monitoring.Host) {
	mon.StatesLock.Lock()
	defer mon.StatesLock.Unlock()

	activePlugins := make(map[string]bool, len(host.Plugins.Plugins))

	//flatten trees into rows and merge metrics
	for _, plugin := range host.Plugins.Plugins {

		newState := PluginState{
			ID:             host.Address + "::" + plugin.PluginId,
			Host:           host.Address,
			HostUp:         host.Online,
			HostError:      host.LastError,
			Timestamp:      time.Now().Unix(),
			PluginID:       plugin.PluginId,
			PluginCategory: plugin.PluginCategory,
			PluginType:     plugin.Type,
			PluginConfig:   plugin.Config,
			OutputPlugin:   plugin.OutputPlugin,

			BufferQueueLength:     []int{plugin.BufferQueueLength},
			BufferTotalQueuedSize: []int{plugin.BufferTotalQueuedSize},
			RetryCount:            []int{plugin.RetryCount},

			Warnings: 			   []string{}}

		//keep track of currently active states
		activePlugins[newState.ID] = true

		//update existing states
		stateHandled := false
		for oldKey, oldState := range mon.States {
			if newState.ID == oldState.ID {

				if len(oldState.BufferQueueLength) >= mon.HistorySize-1 {
					oldState.BufferQueueLength = oldState.BufferQueueLength[1:]
				}
				newState.BufferQueueLength = append(oldState.BufferQueueLength, newState.BufferQueueLength[0])

				if len(oldState.BufferTotalQueuedSize) >= mon.HistorySize-1 {
					oldState.BufferTotalQueuedSize = oldState.BufferTotalQueuedSize[1:]
				}
				newState.BufferTotalQueuedSize = append(oldState.BufferTotalQueuedSize, newState.BufferTotalQueuedSize[0])

				if len(oldState.RetryCount) >= mon.HistorySize-1 {
					oldState.RetryCount = oldState.RetryCount[1:]
				}
				newState.RetryCount = append(oldState.RetryCount, newState.RetryCount[0])

				//replace in state list
				mon.States[oldKey] = newState
				stateHandled = true
			}
		}

		//or append new
		if stateHandled == false {
			mon.States = append(mon.States, newState)
		}

		//apply checks
		for _, check := range mon.StateChecks {
			newState.Accept(check)
		}
	}

	//clean up redundant states (discard pluginIDs that are no longer reported by the host. PluginIDs are regenerated
	//on restart by fluentd so you cannot easily identify plugins between restarts)
	for i := 0; i < len(mon.States); i++ {
		if host.Address == mon.States[i].Host {
			if _, ok := activePlugins[mon.States[i].ID]; ok == false {
				//splice out redundant state
				mon.States = append(mon.States[:i], mon.States[i+1:]...)
				//slice is now 1 element smaller (re-indexed by append) so we don't need to advance
				i--
			}
		}
	}
}

func NewMonitor(hub *Hub, hosts []*monitoring.Host) *FMonitor {
	return &FMonitor{
		Hub:                hub,
		Hosts:              hosts,
		States:             make([]PluginState, 0),
		PushTickSeconds:    1,
		HistoryTickSeconds: 10,
		HistorySize:        60}
}
