package conductor

import (
	"context"
	"testing"
	"time"
)

type testFlags struct {
	started   bool
	stopped   bool
	startTime time.Time
}

func (f *testFlags) check(t *testing.T) {
	if !f.started {
		t.Fatal("Failed to start service")
	}

	if !f.stopped {
		t.Fatal("Failed to stop service")
	}
}

type testService struct {
	flags *testFlags
	test  *testing.T
}

func (t *testService) Run(started, stopped chan bool, stop chan context.Context) error {
	go func() {
		started <- true
		t.test.Log("start")
		t.flags.started = true
		t.flags.startTime = time.Now()
		select {
		case <-stop:
			t.test.Log("stopping because signalled ")
			t.flags.stopped = true
			stopped <- true
		}
	}()
	return nil
}

func TestConductor(t *testing.T) {

	f := testFlags{}
	s := testService{&f, t}
	c := New(Noisy())
	c.Service("test service", &s)
	stopped := c.Start()
	c.Stop()
	<-stopped
	f.check(t)
}

func TestOrderedStartup(t *testing.T) {

	f1 := testFlags{}
	s1 := testService{&f1, t}

	f2 := testFlags{}
	s2 := testService{&f2, t}

	c := New(Noisy())
	c.Service("test service", &s1)
	c.Service("test service 2", &s2)
	stopped := c.Start()
	c.Stop()
	<-stopped

	f1.check(t)
	f2.check(t)
	if !f2.startTime.After(f1.startTime) {
		t.Fatal("Services not started in order")
	}

}
