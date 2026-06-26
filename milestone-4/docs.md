# Milestone 4 — Don't Trust the Model, Test It *(optional)*

> **New concept:** verify the model's output with a deterministic check, and feed failures back into the loop so the agent **self-corrects**.
>
> **Builds on:** [Milestone 3](./milestone-3.md) — same agent loop. This milestone inserts a verification gate at the one place Milestone 3 blindly trusted the model: its "I'm done" reply.

An LLM is **probabilistic**. "Reverse this text" is a *request*, not a guarantee — the model may drop a character, leave the order untouched, or never write the file. Milestone 3's loop believes the model when it says "done." This milestone refuses to: it checks the result, and if the check fails, it hands the failure back as a new turn and lets the agent try again.

This is the lesson that turns a demo into a real agent: **pair the loop with tests on the output.**

---

## What changed since Milestone 3

The loop body is identical *until* the model returns a no-tool-call reply. Milestone 3 printed and exited there. Milestone 4 inserts a gate:

```mermaid
flowchart TD
    resp[model reply, no tool calls] --> validate{validateReversal<br/>demo.txt vs reversed.txt}
    validate -->|ok| done[print Content → exit ✅]
    validate -->|fail| feedback[append failure as a<br/>new 'user' message]
    feedback --> loop[continue loop:<br/>model gets another shot]
    loop -.-> resp
```

Three things are new:

1. **`checkReversal(original, output)`** — the deterministic test.
2. **`validateReversal(srcPath, dstPath)`** — reads both files and runs the test.
3. **`maxTurns`** — a bound on the loop so a model that never converges can't run forever.

```diff
- for {                              // Milestone 3: unbounded
+ const maxTurns = 10
+ for turn := 0; turn < maxTurns; turn++ {   // Milestone 4: bounded
```

And the termination branch changes from *trust* to *verify*:

```diff
  if len(msg.ToolCalls) == 0 {
- 	fmt.Println(msg.Content)
- 	break
+ 	ok, feedback := validateReversal("demo.txt", "reversed.txt")
+ 	if ok {
+ 		fmt.Println(msg.Content)
+ 		return
+ 	}
+ 	messages = append(messages, Message{Role: "user",
+ 		Content: feedback + " — please fix reversed.txt."})
+ 	continue
  }
```

---

## The check itself

```go
func checkReversal(original, output string) []string {
	var failures []string
	if len(output) != len(original) {
		failures = append(failures, fmt.Sprintf(
			"length changed: original is %d bytes, output is %d bytes",
			len(original), len(output)))
	}
	return failures
}
```

It asserts **one invariant**: reversing must not change the length. An empty slice means "passed."

Why only length, not exact-reverse? Pedagogy. Length alone catches the most common probabilistic failure — a dropped or added character — while staying trivial to read. A stricter check (`output == reverse(original)`) is an easy exercise, and the structure here supports it: `checkReversal` returns *every* failure it finds, so you can add more assertions and they all get reported at once.

```go
func validateReversal(srcPath, dstPath string) (ok bool, feedback string) {
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return false, fmt.Sprintf("could not read source %s: %v", srcPath, err)
	}
	dst, err := os.ReadFile(dstPath)
	if err != nil {
		return false, fmt.Sprintf("could not read %s — did you write it? %v", dstPath, err)
	}
	original := strings.TrimRight(string(src), "\n")
	output := strings.TrimRight(string(dst), "\n")
	if failures := checkReversal(original, output); len(failures) > 0 {
		return false, "validation failed: " + strings.Join(failures, "; ")
	}
	return true, ""
}
```

Two design choices worth copying:

- **It reads from disk**, not from the model's claim. The model can say anything; the file is ground truth.
- **It returns human-readable `feedback`.** That string isn't for our logs — it's the message we send *back to the model*. "validation failed: length changed…" is something the model can read and act on. A "did you write it?" hint even covers the case where the file is missing entirely.

The `TrimRight(..., "\n")` ignores a trailing newline on either side, so the check judges the content the model is responsible for, not file-ending conventions.

---

## Self-correction in action

```mermaid
sequenceDiagram
    participant A as agent loop
    participant O as Ollama
    participant V as validateReversal
    participant FS as filesystem

    O-->>A: write_file{reversed.txt: "olle"}  %% dropped a char
    A->>FS: write "olle"
    O-->>A: "Done!" (no tool calls)
    A->>V: validate demo.txt vs reversed.txt
    V->>FS: read both files
    V-->>A: ok=false, "length changed: 5 vs 4 bytes"
    A->>O: user: "validation failed… — please fix reversed.txt."
    O-->>A: write_file{reversed.txt: "olleh"}  %% corrected
    A->>FS: write "olleh"
    O-->>A: "Fixed!" (no tool calls)
    A->>V: validate again
    V-->>A: ok=true
    A->>A: print Content, return ✅
```

The failure message re-enters the loop **as a `user` turn**, so to the model it reads like the human pointing out the mistake. The loop is unchanged — it's the same send/dispatch/append cycle from Milestone 3. We've just made one of its exits conditional on a test passing.

---

## Why `maxTurns`

Two forces can now drive the loop around: tool calls (as before) *and* failed validations. If the model can never produce a correct reversal, "fix it / still wrong" would spin forever. `maxTurns = 10` caps it; on exhaustion the program fails loudly:

```go
log.Fatalf("gave up after %d turns without passing validation", maxTurns)
```

A real agent always bounds its retries. Infinite self-correction is just an infinite loop with extra steps.

---

## The payoff: tests with no Ollama, no network

Because the check lives in a plain function, it's unit-testable in isolation — no model, no HTTP:

```go
func TestCheckReversalPasses(t *testing.T) {
	if failures := checkReversal("hello", "olleh"); len(failures) != 0 {
		t.Errorf("expected no failures, got %v", failures)
	}
}

func TestCheckReversalCatchesLengthChange(t *testing.T) {
	if failures := checkReversal("hello", "olle"); len(failures) == 0 {
		t.Error("expected a length-change failure, got none")
	}
}
```

(`main_test.go` also covers an *added* character.) This separation — deterministic logic out of the probabilistic loop — is exactly what makes the agent testable.

```bash
go test ./milestone-4/      # fast, offline
```

---

## Run it

```bash
go build -o ./milestone-4-bin ./milestone-4/
./milestone-4-bin
```

Expected: the agent reverses `demo.txt` into `reversed.txt`; if the first attempt is wrong you'll see a `validation failed` warning in the logs followed by another attempt, then `validation passed`.

---

## Takeaway

Verification is what turns a probabilistic model into a reliable system. The pattern generalizes far beyond reversing text:

> **Run the action → check the result deterministically → on failure, feed the failure back → bound the retries.**

Any agent you build for real should do all four.

---

| | |
|---|---|
| ← Previous | [Milestone 3 — The Agent Loop](../milestone-3/docs.md) |
| Back to start | [Milestone 1 — One HTTP Call](../milestone-1/docs.md) |
