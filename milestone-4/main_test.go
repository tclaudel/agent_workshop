package main

import "testing"

// These tests need no Ollama and no network — they exercise the deterministic
// check the agent runs on the model's output. Run them with:
//
//	go test ./milestone-4/

// Correct uppercasing passes — output equals the original with every letter
// uppercased.
func TestCheckUppercasePasses(t *testing.T) {
	if failures := checkUppercase("hello agent workshop", "HELLO AGENT WORKSHOP"); len(failures) != 0 {
		t.Errorf("expected no failures, got %v", failures)
	}
}

// A lowercase letter left untouched is the kind of bug a probabilistic model
// produces and a deterministic check catches.
func TestCheckUppercaseCatchesMissedLetter(t *testing.T) {
	if failures := checkUppercase("hello agent workshop", "HELLO AGENT WORKSHOp"); len(failures) == 0 {
		t.Error("expected a failure for the un-uppercased letter, got none")
	}
}

// Dropped or altered content also fails — the output no longer matches the
// uppercased original.
func TestCheckUppercaseCatchesChangedContent(t *testing.T) {
	if failures := checkUppercase("hello agent workshop", "HELLO WORKSHOP"); len(failures) == 0 {
		t.Error("expected a failure for changed content, got none")
	}
}
