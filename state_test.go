package main

import "testing"

type TestingStateVisitor struct {
	StatesVisited []*PluginState
}

func (v *TestingStateVisitor) VisitState(state *PluginState) {
	v.StatesVisited = append(v.StatesVisited, state)
}

func TestPluginStateAcceptsVisitor(t *testing.T) {

	testingVisitor := &TestingStateVisitor{}

	state := &PluginState{}
	state.Accept(testingVisitor)

	if len(testingVisitor.StatesVisited) != 1 || testingVisitor.StatesVisited[0] != state {
		t.Error("Testing visitor did not visit the target state")
	}
}

func TestBufferWarningVisitorDetectsEverIncreasingBuffer(t *testing.T) {

	check := &BufferWarningVisitor{}

	state := &PluginState{BufferTotalQueuedSize: []int{1,2,3,4,5}}
	state.Accept(check)

	if len(state.Warnings) != 1 || state.Warnings[0] != "buffer is not clearing fast enough" {
		t.Error("BufferWarningVisitor did not detect increasing buffer")
	}
}

// a toothed increase is buffer that is cleared at a slower rate than it is filled leading to a jagged increase
func TestBufferWarningVisitorDetectsToothedIncrease(t *testing.T) {

	check := &BufferWarningVisitor{}

	state := &PluginState{BufferTotalQueuedSize: []int{1,2,3,2,4,5,4,5}}
	state.Accept(check)

	if len(state.Warnings) != 1 || state.Warnings[0] != "buffer is not clearing fast enough" {
		t.Error("BufferWarningVisitor did not detect increasing buffer")
	}
}

func TestBufferWarningVisitorDoesNotAlertOnDecreasingBuffer(t *testing.T) {

	check := &BufferWarningVisitor{}

	state := &PluginState{BufferTotalQueuedSize: []int{5,4,3,2,1}}
	state.Accept(check)

	if len(state.Warnings) != 0  {
		t.Error("BufferWarningVisitor mis-identified buffer clearing issue")
	}
}

func TestBufferWarningVisitorDoesNotAlertOnToothedDecrease(t *testing.T) {

	check := &BufferWarningVisitor{}

	state := &PluginState{BufferTotalQueuedSize: []int{1,2,3,1,2,3,1,2}}
	state.Accept(check)

	if len(state.Warnings) != 0 {
		t.Error("BufferWarningVisitor mis-identified buffer clearing issue with toothed buffer profile")
	}
}

func TestBufferWarningVisitorDoesNotAlertOnFixedZeroLengthBuffer(t *testing.T) {

	check := &BufferWarningVisitor{}

	state := &PluginState{BufferTotalQueuedSize: []int{0,0,0,0,0,0}}
	state.Accept(check)

	if len(state.Warnings) != 0 {
		t.Error("BufferWarningVisitor mis-identified buffer clearing issue with flat zero buffer profile")
	}
}
