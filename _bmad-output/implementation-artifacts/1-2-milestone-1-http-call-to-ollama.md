---
baseline_commit: 1-1-repository-foundation
---

# Story 1.2: Milestone 1 — HTTP Call to Ollama

Status: done

## Story

As an attendee,
I want a complete working Go program that sends a prompt to Ollama and prints the reply,
so that I have a reference for the first concept: a raw HTTP call to a local LLM.

## Acceptance Criteria

1. **Given** Ollama is running with llama3.2 pulled **When** `go run ./milestone-1/` is executed **Then** the program sends a request to `http://localhost:11434/api/chat` and prints the model's text reply to stdout

2. **Given** `milestone-1/main.go` **When** it is read **Then** all Ollama calls use `http.Post(url, "application/json", body)` — no custom `http.Client`, no `http.NewRequest` (AD-13) **And** request and response bodies are typed Go structs, not `map[string]any` or raw byte templates (AD-11) **And** the request sets `stream: false` **And** there are zero external imports beyond stdlib (AD-4, AD-9)

3. **Given** `milestone-1/main.go` **When** the prompt input mechanism is read **Then** the prompt is a hardcoded string literal in `main()` — no CLI argument parsing, no stdin reading, no flag package

4. **Given** `milestone-1/main.go` **When** `go build -o /dev/null ./milestone-1/` is run **Then** it compiles independently without touching any other directory (AD-2, AD-3)

## Tasks / Subtasks

- [x] Task 1: Remove placeholder, create `milestone-1/main.go` (AC: 1, 2, 3, 4)
  - [x] Delete `milestone-1/.gitkeep`
  - [x] Create `milestone-1/main.go` with `package main` declaration
  - [x] Add stdlib imports only: `bytes`, `encoding/json`, `fmt`, `io`, `log`, `net/http`
  - [x] Define typed structs: `Message`, `ChatRequest`, `ChatResponse` (see Dev Notes for exact shapes)
  - [x] Implement `main()`: hardcoded prompt → marshal → `http.Post` → read body → unmarshal → print content
- [x] Task 2: Verify build and vet (AC: 4)
  - [x] `go build -o /dev/null ./milestone-1/` exits 0 with no errors
  - [x] `go vet ./milestone-1/` exits 0

### Review Findings

- [x] [Review][Patch] No HTTP status code check — silent failure on 4xx/5xx and Ollama error JSON [milestone-1/main.go:41-55] — fixed
- [x] [Review][Defer] No HTTP client timeout [milestone-1/main.go:41] — deferred, violates AD-13
- [x] [Review][Defer] ChatResponse struct missing done/done_reason [milestone-1/main.go:23-25] — deferred, M1 scope boundary
- [x] [Review][Defer] No MaxBytes guard / empty-partial-truncated body edge cases [milestone-1/main.go:47-55] — deferred, pre-existing
- [x] [Review][Defer] Content-Type mismatch accepted silently [milestone-1/main.go:47-55] — deferred, pre-existing

## Dev Notes

### Struct Shapes (AD-11)

M1 uses a minimal `Message` struct — `ToolCalls` is intentionally absent in milestone-1. AD-10 binds to milestone-2/3 only; showing `ToolCalls` in M1 would confuse attendees learning "HTTP call only." M2's concept block will add this field inside the struct as its one permitted internal insert.

```go
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
    Stream   bool      `json:"stream"`
}

type ChatResponse struct {
    Message Message `json:"message"`
}
```

### HTTP Call Pattern (AD-13)

```go
body, err := json.Marshal(ChatRequest{...})
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
fmt.Println(result.Message.Content)
```

### Hardcoded Prompt

Use `"What is 2+2?"` — short, deterministic intent, clearly readable during live session. Set as `Content` in a user `Message` with `Role: "user"`. Model name: `"llama3.2"`.

### Error Handling Convention

`log.Fatal(err)` on all hard errors (Ollama unreachable, marshal/unmarshal failure). No custom error types, no retries (architecture consistency convention).

### Build Verification Command

Use `go build -o /dev/null ./milestone-1/` — NOT bare `go build ./milestone-1/`. When run from repo root, Go defaults the binary output name to `milestone-1`, which collides with the existing `milestone-1/` directory and fails. Learned from story 1.1 debug log (same issue with `./starter/`).

### Scope Boundary

Do NOT add: `Tool` struct, `ToolDef`, `ToolFunction`, `ToolCall`, `ToolCallFunction`, agent loop, `switch` dispatch, message history slice, `[]ToolCall` field on `Message`, CLI arg parsing, `os.Args`, `flag` package. Those belong to stories 1.3 (M2) and 1.4 (M3).

### M2 Compatibility Note

When story 1.3 implements M2, the `Message` struct gains `ToolCalls []ToolCall \`json:"tool_calls,omitempty"\``. This single field insertion inside the existing struct is the one place M2 is not "purely appended" — it is intentional and acceptable under AD-1 (no lines deleted, no reordering). All other M2 additions are appended after existing M1 code.

### Architecture Invariants Binding This Story

| Invariant | Rule |
|---|---|
| AD-2 | Standalone — no cross-milestone imports; `go build ./milestone-1/` compiles alone |
| AD-3 | Single file: `main.go` only — no sub-packages, no helper files |
| AD-4 | HTTP POST to `http://localhost:11434/api/chat`; no LLM SDK |
| AD-9 | Single `go.mod` at repo root; zero `require` directives; stdlib only |
| AD-11 | Typed Go structs for all Ollama request/response — no `map[string]any` |
| AD-13 | `http.Post(url, "application/json", body)` shorthand — no custom `http.Client`, no `http.NewRequest` |

### Project Structure Notes

```
agent-workshop/
  go.mod                      ← unchanged (story 1.1)
  milestone-1/
    main.go                   ← CREATED by this story (replaces .gitkeep)
  milestone-2/
    .gitkeep                  ← unchanged
  milestone-3/
    .gitkeep                  ← unchanged
  starter/
    main.go                   ← unchanged (story 1.1)
```

### References

- Structs and HTTP pattern: [Source: _bmad-output/planning-artifacts/architecture/architecture-agent_workshop-2026-06-24/ARCHITECTURE-SPINE.md#AD-11, AD-13]
- Message struct shape (M2/M3 extension): [Source: ARCHITECTURE-SPINE.md#AD-10]
- Additive delta rule: [Source: ARCHITECTURE-SPINE.md#AD-1]
- Story AC source: [Source: _bmad-output/planning-artifacts/epics.md#Story-1.2]
- Build collision fix: [Source: _bmad-output/implementation-artifacts/1-1-repository-foundation.md#Debug-Log]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6[1m]

### Debug Log References

### Completion Notes List

- Created `milestone-1/main.go`: stdlib-only HTTP POST to Ollama. Typed structs `Message`, `ChatRequest`, `ChatResponse`. Hardcoded prompt `"What is 2+2?"`. All errors via `log.Fatal`. `go build -o /dev/null ./milestone-1/` and `go vet ./milestone-1/` both exit 0.

### File List

- milestone-1/main.go (created)

## Change Log

- 2026-06-24: Story 1.2 implemented — created `milestone-1/main.go` with typed structs and `http.Post` pattern. Deleted `.gitkeep`. Build and vet verified clean.
