package main

import (
    "sync"
    "log"
    "time"
    "encoding/json"
    "bytes"
    "github.com/warmans/fluentd-api-client/monitoring"
)

//-----------------------------------------------
// Fuentd Status Monitor
//-----------------------------------------------

type State struct {
    ID                       string
    Host                     string
    HostUp                   bool
    HostError                string
    Timestamp                int64
    PluginID                 string
    PluginCategory           string
    PluginType               string
    PluginConfig             map[string]interface{}
    OutputPlugin             bool
    BufferQueueLength        []int
    BufferTotalQueuedSize    []int
    RetryCount               []int
}

type FMonitor struct {
    Hub             *Hub
    Hosts           []*monitoring.Host
    States          []*State
    StatesLock      sync.RWMutex
}

func (mon *FMonitor) Run() {

    go func(){
        updateTicker := time.NewTicker(time.Second * 5)
        for range updateTicker.C {

            //update all hosts in parallel
            for _, host := range mon.Hosts {
                go func() {
                    host.Update()
                    mon.AddState(host);
                }();
            }
        }
    }()

    broadcastTicker := time.NewTicker(time.Second * 1)
    for range broadcastTicker.C {

        //Create a JSON payload
        buffer := &bytes.Buffer{}
        encoder := json.NewEncoder(buffer)

        if err := encoder.Encode(mon.GetStateUpdate()); err != nil {
            log.Print("Failed to encode host stats payload");
            continue
        }

        //Broadcast to all clients
        mon.Hub.Broadcast(buffer.Bytes())
    }
}

func (mon *FMonitor) GetStateUpdate() []*State {
    mon.StatesLock.RLock()
    defer mon.StatesLock.RUnlock()

    return mon.States
}

func (mon *FMonitor) AddState(host *monitoring.Host) {
    mon.StatesLock.Lock()
    defer mon.StatesLock.Unlock()

    //flatten trees into rows
    for _, plugin := range host.Plugins.Plugins {

        newState := &State{
            ID: host.Address+"::"+plugin.PluginId,
            Host: host.Address,
            HostUp: host.Online,
            HostError: host.LastError,
            Timestamp: time.Now().Unix(),
            PluginID: plugin.PluginId,
            PluginCategory: plugin.PluginCategory,
            PluginType: plugin.Type,
            PluginConfig: plugin.Config,
            OutputPlugin: plugin.OutputPlugin,

            BufferQueueLength: []int{plugin.BufferQueueLength},
            BufferTotalQueuedSize: []int{plugin.BufferTotalQueuedSize},
            RetryCount: []int{plugin.RetryCount}}

        //attempt existing row update
        update := false
        for key, state := range mon.States {
            if newState.ID == state.ID {

                update = true
                oldState := state

                //replace old state
                mon.States[key] = newState

                //some values should hold a running history so...
                if len(oldState.BufferQueueLength) >= 59 {
                    oldState.BufferQueueLength = oldState.BufferQueueLength[1:]
                }
                mon.States[key].BufferQueueLength = append(oldState.BufferQueueLength, newState.BufferQueueLength[0])

                if len(oldState.BufferTotalQueuedSize) >= 59 {
                    oldState.BufferTotalQueuedSize = oldState.BufferTotalQueuedSize[1:]
                }
                mon.States[key].BufferTotalQueuedSize = append(oldState.BufferTotalQueuedSize, newState.BufferTotalQueuedSize[0])

                if len(oldState.RetryCount) >= 59 {
                    oldState.RetryCount = oldState.RetryCount[1:]
                }
                mon.States[key].RetryCount = append(oldState.RetryCount, newState.RetryCount[0])
            }
        }
        //else append new
        if update == false {
            mon.States = append(mon.States, newState)
        }
    }
}

func NewMonitor(hub *Hub, hosts []*monitoring.Host) *FMonitor {
    return &FMonitor{
        Hub: hub,
        Hosts: hosts,
        States: make([]*State, 0)}
}
