package main

type PluginState struct {
	ID                    string
	Host                  string
	HostUp                bool
	HostError             string
	Timestamp             int64
	PluginID              string
	PluginCategory        string
	PluginType            string
	PluginConfig          map[string]interface{}
	OutputPlugin          bool
	BufferQueueLength     []int
	BufferTotalQueuedSize []int
	RetryCount            []int
	Warnings              []string
}

func (ps *PluginState) Accept(visitor PluginStateVisitor) {
	visitor.VisitState(ps)
}

//PluginStateVisitor implementations can collect additional information about a plugin's state based on its properties
type PluginStateVisitor interface {
	VisitState(state *PluginState)
}

//BufferWarningVisitor detects buffer size issues. Specifically when a buffer is filling faster than it can be cleared
type BufferWarningVisitor struct {
}

func (v *BufferWarningVisitor) VisitState(state *PluginState) {
	if len(state.BufferTotalQueuedSize) == 0 {
		return
	}

	lowPoint := -1
	for _, point := range state.BufferTotalQueuedSize[1:] {
		if lowPoint == -1 || point < lowPoint {
			lowPoint = point;
		}
	}

	if lowPoint > state.BufferTotalQueuedSize[0] {
		state.Warnings = append(state.Warnings, "buffer is not clearing fast enough")
	}
}
