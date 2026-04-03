package logger

import (
	"sync"
	"testing"
)

func reset() {
	global = nil
	once = sync.Once{}
}

func TestLReturnsNoOpBeforeInit(t *testing.T) {
	reset()
	l := L()
	if l == nil {
		t.Fatal("expected non-nil logger before Init")
	}
	l.Info("this should not panic")
}

func TestInitDebugMode(t *testing.T) {
	reset()
	Init(true)
	if global == nil {
		t.Fatal("expected global logger to be set")
	}
	L().Debug("debug mode test")
}

func TestInitProductionMode(t *testing.T) {
	reset()
	Init(false)
	if global == nil {
		t.Fatal("expected global logger to be set")
	}
	L().Info("production mode test")
}

func TestInitOnlyOnce(t *testing.T) {
	reset()
	Init(true)
	first := global
	Init(false)
	if global != first {
		t.Fatal("expected Init to be called only once")
	}
}

func TestSync(t *testing.T) {
	reset()
	Init(false)
	Sync()
}

func TestSyncBeforeInit(t *testing.T) {
	reset()
	Sync()
}
