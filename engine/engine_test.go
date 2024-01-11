package engine

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testServer struct {
	step    int
	running chan struct{}
}

func (s *testServer) Name() string {
	return "testServer"
}

func (s *testServer) BeforeRunHandler() BehaveHandler {
	return func(engine *Engine, srv Server) error {
		s.step++
		return nil
	}
}

func (s *testServer) BeforeShutdownHandler() BehaveHandler {
	return func(engine *Engine, srv Server) error {
		if s.step != 1 {
			return fmt.Errorf("s.step should be 1 but %d", s.step)
		}
		s.step++
		return nil
	}
}

func (s *testServer) AfterStopHandler() BehaveHandler {
	return func(engine *Engine, srv Server) error {
		if s.step != 2 {
			return fmt.Errorf("s.step should be 2 but %d", s.step)
		}
		return nil
	}
}

func (s *testServer) Run() error {
	<-s.running
	return nil
}

func (s *testServer) Shutdown() error {
	close(s.running)
	return nil
}

func TestEngine(t *testing.T) {
	opts := Option{
		Name:            "testEngine",
		ShutdownSignal:  []os.Signal{syscall.SIGINT},
		ShutdownHandler: nil,
		Console:         true,
		BeforeHandler:   nil,
		AfterHandler:    nil,
	}
	srv := &testServer{
		running: make(chan struct{}),
	}
	engine := New(opts, srv)
	assert.Equal(t, "testEngine", engine.Name())

	go func() {
		time.Sleep(1 * time.Second)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(t, err)
	}()

	require.NoError(t, engine.Run())
}
