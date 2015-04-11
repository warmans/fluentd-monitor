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
    Host                     string
    HostUp                   bool
    HostError                string
    Timestamp                int64
    PluginID                 string
    PluginCategory           string
    PluginType               string
    OutputPlugin             bool
    CurBufferQueueLength     int
    CurBufferTotalQueuedSize int
    CurRetryCount            int
}

type FMonitor struct {
    Hub             *Hub
    Hosts           []*monitoring.Host
    States          []*State
    StatesLock      sync.RWMutex
}

func (mon *FMonitor) Run() {

    go func(){
        updateTicker := time.NewTicker(time.Second * 10)
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
            Host: host.Address,
            HostUp: host.Online,
            HostError: host.LastError,
            Timestamp: time.Now().Unix(),
            PluginID: plugin.PluginId,
            PluginCategory: plugin.PluginCategory,
            PluginType: plugin.Type,
            OutputPlugin: plugin.OutputPlugin,
            CurBufferQueueLength: plugin.BufferQueueLength,
            CurBufferTotalQueuedSize: plugin.BufferTotalQueuedSize,
            CurRetryCount: plugin.RetryCount}

        //attempt existing row update
        for key, state := range mon.States {
            if newState.Host == state.Host && newState.PluginID == state.PluginID {
                mon.States[key] = newState
                return
            }
        }
        //else append new
        mon.States = append(mon.States, newState)
    }
}

func NewMonitor(hub *Hub, hosts []*monitoring.Host) *FMonitor {
    return &FMonitor{
        Hub: hub,
        Hosts: hosts,
        States: make([]*State, 0)}
}
