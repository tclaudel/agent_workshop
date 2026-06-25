---
baseline_commit: 1-3-milestone-2-tool-struct-and-tools-slice
---

# Story 1.4: Milestone 3 — Agent Loop, Dispatch, and History

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an attendee,
I want milestone-2 extended with the agent `for` loop, `switch`-based tool dispatch, and `[]Message` history accumulation,
so that I have the complete working agent with all three concepts introduced incrementally and nothing hidden.

## Acceptance Criteria

1. **Given** Ollama running with llama3.2 pulled and a file-editing task **When** `go run ./milestone-3/` is executed **Then** the agent loop runs, dispatches `read_file` and `write_file` tool calls, accumulates message history, and terminates when the model returns a text reply with no tool calls (AD-7)

2. **Given** `milestone-3/main.go` **When** compared line-by-line to `milestone-2/main.go` **Then** every line from milestone-2 is present and unchanged, and the new concept block (loop + switch dispatch + history append) is appended with no deletions or reordering (AD-1) — see **Permitted Delta Exceptions** in Dev Notes for the two governed exceptions

3. **Given** the agent loop in `milestone-3/main.go` **When** tool dispatch is read **Then** it uses `switch msg.ToolCalls[0].Function.Name` with one `case` per tool name — no if/else chain, no function-map (AD-12)

4. **Given** the stop condition in `milestone-3/main.go` **When** it is read **Then** the loop exits when `len(msg.ToolCalls) == 0` — no other stop condition (AD-7)

5. **Given** the agent loop in `milestone-3/main.go` **When** the loop body is read **Then** the `len(msg.ToolCalls) == 0` stop-condition check appears before the `switch` dispatch statement — dispatch is never reached on an empty `ToolCalls` slice (AD-7, AD-12)

6. **Given** the message history in `milestone-3/main.go` **When** it is read **Then** it is a `[]Message` slice appended in the loop, never truncated, never written to disk (AD-6)

7. **Given** `milestone-3/main.go` **When** `go build -o /dev/null ./milestone-3/` is run **Then** it compiles independently without touching any other directory (AD-2, AD-3)

## Tasks / Subtasks

- [x] Task 1: Seed milestone-3/main.go from milestone-2 baseline (AC: 2)
  - [x] Delete `milestone-3/.gitkeep`
  - [x] Create `milestone-3/main.go` as an exact copy of `milestone-2/main.go`

- [x] Task 2: Correct `ToolCallFunction.Arguments` to `map[string]any` (AC: 1, 3) — **CRITICAL, see Permitted Delta Exception #1**
  - [x] Change field from `Arguments string` to `Arguments map[string]any` (json tag `"arguments"` unchanged)
  - [x] This aligns the struct with AD-11 and the real Ollama wire format; without it M3 panics at runtime

- [x] Task 3: Change the prompt to a file-editing task (AC: 1) — **see Permitted Delta Exception #2**
  - [x] Replace `prompt := "What is 2+2?"` with a task that forces a tool call, e.g. `prompt := "Read the file demo.txt, then write its contents reversed to reversed.txt"`
  - [x] A non-file prompt would exit the loop on iteration 1 (no tool_calls) and ACs 1/3/5/6 would never exercise

- [x] Task 4: Seed the message history slice before the loop (AC: 6)
  - [x] Add `messages := []Message{{Role: "user", Content: prompt}}` after the prompt line, before the loop

- [x] Task 5: Wrap the M2 request/response code in the agent `for` loop (AC: 1, 6)
  - [x] Open a `for {` after seeding `messages`
  - [x] Inside the loop, marshal `ChatRequest{... Messages: messages ...}` (was the inline `[]Message{{...}}` literal in M2 — now references the history slice)
  - [x] Keep the M2 `http.Post` → `io.ReadAll` → `json.Unmarshal(... &result)` block unchanged inside the loop body
  - [x] Bind `msg := result.Message` and `messages = append(messages, msg)` (accumulate assistant reply — AD-6)

- [x] Task 6: Add the stop-condition check BEFORE dispatch (AC: 4, 5)
  - [x] `if len(msg.ToolCalls) == 0 { fmt.Println(msg.Content); break }` — this replaces M2's terminal `fmt.Println(result.Message.Content)`
  - [x] Must appear before the `switch` (AC 5)

- [x] Task 7: Add `switch` tool dispatch and append the tool result (AC: 1, 3, 6)
  - [x] `switch msg.ToolCalls[0].Function.Name` with one `case` per tool (`read_file`, `write_file`) — AD-12
  - [x] Each case calls the matching tool's `Run(msg.ToolCalls[0].Function.Arguments)` and captures the output string
  - [x] `messages = append(messages, Message{Role: "tool", Content: output})` so the model sees the result on the next iteration (AD-6)

- [x] Task 8: Verify build, vet, and runtime behavior (AC: 1, 7)
  - [x] `go build -o /dev/null ./milestone-3/` exits 0
  - [x] `go vet ./milestone-3/` exits 0
  - [x] With Ollama running + llama3.2, `go run ./milestone-3/` completes a tool-dispatch loop and produces a visible file change, then prints a final text reply

## Dev Notes

### 🚨 Permitted Delta Exceptions (read before touching AD-1)

AC #2 reads "every line from milestone-2 is present and unchanged." Two lines are **governed exceptions** — they are the only M2 lines this story may modify, and each is mandated by a higher authority than the literal AD-1 wording. Document both in the Change Log so code-review does not revert them (the M1→M2 review *did* revert an ungoverned prompt change — do not let these be confused with that).

**Exception #1 — `ToolCallFunction.Arguments` type (CRITICAL, runtime-breaking if skipped):**

M2 shipped:
```go
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`   // ← WRONG for M3
}
```
M3 MUST be:
```go
type ToolCallFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}
```
Rationale:
- **AD-11 mandates `map[string]any`** (architecture spine, the binding authority — see References). M2 deviated; the deviation was harmless only because `ToolCall` was dead code in M2 (no dispatch existed).
- Ollama `/api/chat` returns `tool_calls[].function.arguments` as a **JSON object**, not a string. With `Arguments string`, `json.Unmarshal` of the response fails with `json: cannot unmarshal object into Go struct field ToolCallFunction.arguments of type string` — M3 crashes on the first tool-calling reply.
- `map[string]any` plugs directly into `Tool.Run(args map[string]any)` (AD-5) — no intermediate `json.Unmarshal` of arguments needed, and no new imports.

> NOTE for facilitator / correct-course: M2 (`milestone-2/main.go:31`) still carries the latent `Arguments string` bug. It is invisible in M2 (dead code) but means M2 and M3 are not struct-identical. Recommend a follow-up to align M2 to `map[string]any` so AD-1 line-identity holds across all milestones. **Out of scope for this story — do not edit M2 here.** Flag it; do not silently fix it.

**Exception #2 — prompt literal:**

`prompt := "What is 2+2?"` → a file-editing task (e.g. `"Read the file demo.txt, then write its contents reversed to reversed.txt"`). Rationale: AC #1/#3/#5/#6 require the loop to actually dispatch `read_file`/`write_file`. "What is 2+2?" returns plain text on iteration 1, the stop condition fires immediately, and dispatch is never exercised — the milestone would not demonstrate its own concept. The prompt change is functionally required by M3, unlike the M2 case where it was reverted (M2 has no dispatch, so no functional reason to change it).

### The Additive Delta (AD-1) — what the M2→M3 diff looks like

M3 = M2 verbatim + the two governed exceptions above + the new concept block. The "new concept block" is the agent loop. Converting M2's straight-line single-shot call into a loop necessarily **wraps and indents** the existing M2 request/response lines — this is permitted (no deletion, no logical reordering; indentation is not reordering). The only structural relocations:
- M2's inline `Messages: []Message{{Role: "user", Content: prompt}}` becomes `Messages: messages` (references the new history slice).
- M2's terminal `fmt.Println(result.Message.Content)` moves into the stop-condition branch (`if len(msg.ToolCalls) == 0 { fmt.Println(msg.Content); break }`).

Everything else from M2 (structs, tools slice, toolDefs loop, `http.Post`, `io.ReadAll`, `json.Unmarshal`) is present and unchanged inside the loop body.

### Canonical M3 main() concept block

Seed history + loop, replacing M2's single-shot tail (the `body :=`/`http.Post`/print sequence). The `tools`, `toolDefs`, and `prompt` setup above the loop is unchanged except the prompt value (Exception #2):

```go
prompt := "Read the file demo.txt, then write its contents reversed to reversed.txt"

messages := []Message{{Role: "user", Content: prompt}}

for {
	body, err := json.Marshal(ChatRequest{
		Model:    "llama3.2",
		Messages: messages,
		Stream:   false,
		Tools:    toolDefs,
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result ChatResponse
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}

	msg := result.Message
	messages = append(messages, msg)

	if len(msg.ToolCalls) == 0 {
		fmt.Println(msg.Content)
		break
	}

	var output string
	switch msg.ToolCalls[0].Function.Name {
	case "read_file":
		output = tools[0].Run(msg.ToolCalls[0].Function.Arguments)
	case "write_file":
		output = tools[1].Run(msg.ToolCalls[0].Function.Arguments)
	}
	messages = append(messages, Message{Role: "tool", Content: output})
}
```

Notes on the block:
- **Stop check before switch** (AD-7, AC 5): the `if len == 0 { break }` precedes the `switch`, so dispatch never touches `ToolCalls[0]` on an empty slice (avoids index-out-of-range panic).
- **Single tool call per turn** (AD-12): milestone code dispatches `ToolCalls[0]` only — one tool per iteration. Keep it simple per spec; do NOT loop over all tool calls.
- **`tools[0]` / `tools[1]` index coupling**: `tools[0]` = `read_file`, `tools[1]` = `write_file` by definition order in the M2 tools slice. The `switch` is the dispatch decision (satisfies AD-12); the index just resolves which `Tool.Run` to call. Keep the tools slice in `read_file`, `write_file` order. (A name→Tool lookup map is NOT permitted as the dispatch mechanism — AD-12 — though the switch itself is the contract; index resolution inside cases is fine.)
- **Tool-result message role**: append the tool output as `Message{Role: "tool", Content: output}` so Ollama feeds it back to the model next iteration. This is what makes the loop converge (AD-6).
- **`defer resp.Body.Close()` inside a `for`**: technically defers stack until function return (loop end), not per-iteration. Acceptable for this teaching binary (short-lived, bounded iterations). Do not add manual `Close()` / restructure — that would be ungoverned scope creep and clutter the pedagogy. Pre-existing pattern from M1/M2.

### Architecture Invariants Binding This Story

| Invariant | Rule |
|---|---|
| AD-1 | Additive-only delta from M2: no deletions, no reordering (except the two governed exceptions above) |
| AD-2 | Standalone — `go build ./milestone-3/` compiles alone, no other dir touched |
| AD-3 | Single file: `main.go` only, no helper files |
| AD-4 | HTTP POST to `http://localhost:11434/api/chat`; no LLM SDK |
| AD-6 | History is `[]Message` accumulated in the loop; never truncated, never written to disk |
| AD-7 | Loop exits ONLY when `len(msg.ToolCalls) == 0`; no max-turns, timeout, or error stop |
| AD-9 | Zero external dependencies; stdlib only |
| AD-11 | Typed Ollama structs; `ToolCallFunction.Arguments` is `map[string]any` |
| AD-12 | Dispatch is `switch msg.ToolCalls[0].Function.Name`, one case per tool; no if/else, no func-map |
| AD-13 | `http.Post(url, "application/json", body)` — no custom http.Client |

### Scope Boundary — DO NOT ADD in M3

- Max-turns / iteration cap / timeout / error-based stop condition (AD-7 forbids; deferred to Phase 7 extension questions)
- Retry/backoff on Ollama calls (AD: `log.Fatal` only)
- Handling multiple tool calls per turn (dispatch `ToolCalls[0]` only — AD-12)
- Path sanitization / file-size caps / permission tuning in `read_file`/`write_file` (deferred items below — pre-existing, out of scope for this story)
- HTTP client timeout (would require custom `http.Client` — violates AD-13)
- New imports — the loop, `switch`, and `append` are language builtins; `map[string]any` Arguments needs no `json.Unmarshal` of args. Import list stays identical to M2.
- Editing `milestone-2/main.go` (the M2 `Arguments string` fix is a separate correct-course item — flag, don't fix)

### Deferred / Pre-existing Items (inherited, NOT this story's job)

These came from M1/M2 reviews and remain deferred (see `deferred-work.md`). `Tool.Run` becomes live code in M3 (dispatch now reaches it), so the first three are now *reachable* — but they remain explicitly deferred for the teaching codebase. Do not address unless ACs require:
- Path traversal in `read_file`/`write_file` (no sanitization) — now reachable in M3, still deferred
- Silent zero-value type assertion `args["path"].(string)` — now reachable, still deferred
- Unbounded `os.ReadFile` (memory) — now reachable, still deferred
- No HTTP timeout (AD-13 forbids the fix), `write_file` hardcoded 0644, schemas never validated — deferred

### Latest Tech Note — Ollama `/api/chat` tool-call wire format

Ollama returns assistant tool calls as:
```json
"message": {
  "role": "assistant",
  "tool_calls": [{"function": {"name": "read_file", "arguments": {"path": "demo.txt"}}}]
}
```
`arguments` is a JSON **object** → Go `map[string]any` (confirms Exception #1). After dispatch, send the result back as a `{"role": "tool", "content": "<output>"}` message in the next request. Requires Ollama 0.3+ and a tool-capable model (llama3.2). [Stack: ARCHITECTURE-SPINE.md#Stack flags llama3.2 tool_calls support as UNVERIFIED — Task 8 manual run is the verification.]

### Project Structure Notes

```
agent-workshop/
  go.mod                    ← unchanged (stdlib only, AD-9)
  milestone-1/main.go       ← unchanged
  milestone-2/main.go       ← unchanged (carries latent Arguments string bug — flag only)
  milestone-3/
    .gitkeep                ← DELETE
    main.go                 ← CREATE (this story)
  starter/main.go           ← unchanged
```

### Build Verification Command

Use `go build -o /dev/null ./milestone-3/` — NOT bare `go build ./milestone-3/`. A binary named `milestone-3` would collide with the `milestone-3/` directory (same issue noted in M1 and M2 stories).

### References

- Agent loop / stop condition: [Source: _bmad-output/planning-artifacts/epics.md#Story-1.4, ARCHITECTURE-SPINE.md#AD-6, #AD-7]
- Switch dispatch: [Source: ARCHITECTURE-SPINE.md#AD-12]
- `Arguments map[string]any` (Exception #1 authority): [Source: ARCHITECTURE-SPINE.md#AD-11]
- Tool struct / Run signature: [Source: ARCHITECTURE-SPINE.md#AD-5]
- Additive delta rule: [Source: ARCHITECTURE-SPINE.md#AD-1]
- M2 current code + scope handoff: [Source: milestone-2/main.go, _bmad-output/implementation-artifacts/1-3-milestone-2-tool-struct-and-tools-slice.md#Scope-Boundary]
- Build collision fix: [Source: 1-3-milestone-2-tool-struct-and-tools-slice.md#Build-Verification-Command]
- Deferred items: [Source: _bmad-output/implementation-artifacts/deferred-work.md]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8[1m] (BMad dev-story workflow)

### Implementation Plan

M3 = M2 verbatim + the two governed AD-1 delta exceptions + the appended agent-loop concept block. No `_test.go` files were authored: AD-3 mandates a single `main.go` and AD-9 forbids any addition beyond stdlib, so a separate test package would itself violate the architecture spine. Validation therefore = `go build -o /dev/null`, `go vet`, structural assertions against each AC (grep), and the Task-8 live runtime run (the spine flags llama3.2 tool_calls as UNVERIFIED — this run is its verification).

### Debug Log References

- `go build -o /dev/null ./milestone-3/` → exit 0
- `go vet ./milestone-3/` → no issues
- Imports lines 3–11 byte-identical to milestone-2 (no new imports — AD-9 held)
- Live run with Ollama + llama3.2:latest, `demo.txt` ("hello agent workshop") seeded at repo root as a throwaway: loop dispatched `read_file` then `write_file`, created `reversed.txt` ("workshop agent hello"), printed a final text reply, exited 0. Throwaway `demo.txt`/`reversed.txt` removed after the run (they are story 1.5 deliverables, not this story's).

### Completion Notes List

- **Exception #1 (CRITICAL)** applied: `ToolCallFunction.Arguments` changed from `string` (M2) to `map[string]any` per AD-11 and the real Ollama wire format. Without it M3 would crash on the first tool-calling reply. json tag `"arguments"` unchanged.
- **Exception #2** applied: prompt changed from `"What is 2+2?"` to `"Read the file demo.txt, then write its contents reversed to reversed.txt"` so the loop actually exercises dispatch (ACs 1/3/5/6).
- AC verification: AC2 line-identity holds (M2 prefix identical; only the two governed exceptions differ); AC3 `switch msg.ToolCalls[0].Function.Name`, one `case` per tool, no if/else, no func-map; AC4 sole stop = `len(msg.ToolCalls) == 0`; AC5 stop-check (L133) precedes `switch` (L139); AC6 `[]Message` appended (L131, L145), never truncated/written to disk; AC7 builds standalone.
- ⚠️ Flag for correct-course (NOT fixed — out of scope per Dev Notes): `milestone-2/main.go:32` still carries the latent `Arguments string` bug. Invisible in M2 (dead code), but means M2/M3 are not struct-identical. Recommend a follow-up to align M2 to `map[string]any`.
- Deferred items (path traversal, silent type assertion, unbounded ReadFile) became *reachable* in M3 (dispatch now hits `Tool.Run`) but remain deferred per Dev Notes — see `deferred-work.md`.

### File List

- `milestone-3/main.go` (new)
- `milestone-3/.gitkeep` (deleted)

### Review Findings

Code review (2026-06-24, 3-layer adversarial: Blind Hunter / Edge Case Hunter / Acceptance Auditor). Acceptance Auditor: **0 AC/invariant violations** — all 7 ACs and AD-1,2,3,4,6,7,9,11,12,13 pass. 0 decision-needed, 0 patch, 1 new defer, ~9 dismissed as spec-governed design or false positives.

- [x] [Review][Defer] Non-2xx HTTP status not checked [milestone-3/main.go:114] — `http.Post` only errors on transport failure; a 4xx/5xx (or 200-with-error) body unmarshals to a zero-value `ChatResponse`, hits `len(ToolCalls)==0`, prints an empty line and exits as if the model finished. New in M3 (the loop relies on response shape). Deferred: a status check is additive error handling beyond AD's `log.Fatal`-only pattern; fixing borders on AD-7 (no error-based stop). Teaching-scope.
- [x] [Review][Defer] Silent zero-value type assertion on `args["path"]`/`args["content"]` [milestone-3/main.go:64,77] — reachable in M3 (dispatch now hits `Tool.Run`). Already tracked in deferred-work.md (from 1-3); not re-added.
- [x] [Review][Defer] Unbounded `os.ReadFile` injected verbatim into next request [milestone-3/main.go:65] — reachable in M3. Already tracked in deferred-work.md (from 1-3); not re-added.

Dismissed (spec-governed / false positive): no-`default` switch + unbounded loop (AD-7 forbids max-turns/error-stop, AD-12 mandates one-case-per-tool); single `ToolCalls[0]` dispatch drops parallel calls (AD-12 mandated); `tools[0/1]` positional coupling (Dev Notes L169 permits); `defer resp.Body.Close()` in loop (Dev Notes L171 permits, pre-existing M1/M2); `ToolCalls[0]` index panic (guarded by AC5 stop-check); missing tool_call_id linkage (Ollama wire format uses role:tool+content, no id — spec L206-215); unbounded `messages` growth (AD-6 never-truncate, AD-7 no-cap); tool-errors-as-strings (teaching design); `log.Fatal` skips defers (consequence of governed error pattern).

## Change Log

- 2026-06-24: Story 1.4 created — Milestone 3 (agent loop, switch dispatch, message history) as additive delta from milestone-2, with two governed AD-1 delta exceptions documented (Arguments→map[string]any per AD-11; prompt→file-editing task).
- 2026-06-24: Story 1.4 implemented — created `milestone-3/main.go` (M2 + agent `for` loop + `switch` dispatch + `[]Message` history). Applied the two governed AD-1 delta exceptions: (1) `ToolCallFunction.Arguments` `string`→`map[string]any` (Exception #1, AD-11, runtime-critical); (2) prompt → file-editing task (Exception #2). Build + vet clean; live Ollama/llama3.2 run dispatched read_file→write_file and converged. Deleted `milestone-3/.gitkeep`. M2 left untouched (latent `Arguments string` flagged for correct-course, not fixed).
