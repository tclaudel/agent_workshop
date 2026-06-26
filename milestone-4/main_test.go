package main

import "testing"

// These tests need no Ollama and no network — they exercise the deterministic
// check the agent runs on the model's output. Run them with:
//
//	go test ./milestone-4/

// Same length but different content still passes — the check only asserts that
// reversing preserved the length, not the exact order.
func TestCheckReversalPasses(t *testing.T) {
	if failures := checkReversal("hello", "olleh"); len(failures) != 0 {
		t.Errorf("expected no failures, got %v", failures)
	}
}

// A dropped character changes the length — the kind of bug a probabilistic
// model produces and a deterministic check catches.
func TestCheckReversalCatchesLengthChange(t *testing.T) {
	if failures := checkReversal("hello", "olle"); len(failures) == 0 {
		t.Error("expected a length-change failure, got none")
	}
}

// An added character also changes the length.
func TestCheckReversalCatchesAddedChar(t *testing.T) {
	if failures := checkReversal("hello", "ollehh"); len(failures) == 0 {
		t.Error("expected a length-change failure, got none")
	}
}
