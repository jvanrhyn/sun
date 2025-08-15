package main

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jvanrhyn/sun/cmd/sun"
)

func TestStateTransitions(t *testing.T) {
	m := initialModel("X", 1, &sun.Client{})
	_, cmd := m.Update(fetchMsg{})
	if cmd == nil {
		t.Fatal("expected command to start fetch")
	}
	_ = tea.Batch() // just ensure tea imported
	_ = time.Now()
}
