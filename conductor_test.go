package conductor

import (
	"context"
	"testing"
)

type testFlags struct {
	started bool
	stopped bool
}

type testService struct {
	flags *testFlags
}

func (t *testService) Run(started, stopped chan bool, stop chan context.Context) error {
	go func() {
		started <- true
		t.flags.started = true
		select {
		case <-stop:
			t.flags.stopped = true
			stopped <- true
		}
	}()
	return nil
}

func TestConductor(t *testing.T) {

	f := testFlags{}
	s := testService{&f}
	c := New()
	c.Service("test service", &s)
	stopped := c.Start()
	c.Stop()
	<-stopped

	if !f.started {
		t.Fatal("Failed to start service")
	}

	if !f.stopped {
		t.Fatal("Failed to stop service")
	}
}
