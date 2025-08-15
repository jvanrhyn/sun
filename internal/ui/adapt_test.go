package ui

import "testing"

func TestAdaptiveColumns(t *testing.T) {
	cols := AdaptiveColumns(80)
	if len(cols) != 6 {
		t.Fatalf("want 6 cols, got %d", len(cols))
	}
	if cols[2].Width < 20 {
		t.Fatalf("conditions col too small: %d", cols[2].Width)
	}
	cols2 := AdaptiveColumns(40)
	if cols2[2].Width != 20 {
		t.Fatalf("expected min 20, got %d", cols2[2].Width)
	}
}
