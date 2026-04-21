package ua

import "testing"

func TestCSeqManagerNextStartsFromSeedAndIncrements(t *testing.T) {
	manager := NewCSeqManager(100)

	if got := manager.Next(); got != 100 {
		t.Fatalf("first cseq = %d, want 100", got)
	}
	if got := manager.Next(); got != 101 {
		t.Fatalf("second cseq = %d, want 101", got)
	}
}


func TestCSeqManagerDefaultsSeedToOne(t *testing.T) {
	manager := NewCSeqManager(0)

	if got := manager.Next(); got != 1 {
		t.Fatalf("first cseq with default seed = %d, want 1", got)
	}
}
