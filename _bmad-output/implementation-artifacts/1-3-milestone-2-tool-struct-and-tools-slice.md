---
baseline_commit: 1-2-milestone-1-http-call-to-ollama
---

# Story 1.3: Milestone 2 — Tool Struct and Tools Slice

Status: done

## Story

As an attendee,
I want milestone-1 extended with a Tool struct, `read_file`/`write_file` implementations, and tools passed in the Ollama request,
so that I can see exactly what "adding tools" means with no other code changed.

## Acceptance Criteria

1. **Given** Ollama is running with llama3.2 pulled and a task that requires file access **When** `go run ./milestone-2/` is executed **Then** the Ollama request body includes a non-empty `tools` array — verified by inspecting the marshalled `ChatRequest` struct before it is sent (not the model response, which is non-deterministic)

2. **Given** `milestone-2/main.go` **When** compared line-by-line to `milestone-1/main.go` **Then** every line from milestone-1 is present and unchanged, and the new concept block (Tool struct + tools slice + ToolDef wrappers) is appended with no deletions or reordering (AD-1)

3. **Given** the `Tool` struct in `milestone-2/main.go` **When** it is read **Then** it has exactly 4 fields in this order: `Name string`, `Description string`, `Parameters json.RawMessage`, `Run func(args map[string]any) string` — no additional fields (AD-5)

4. **Given** `milestone-2/main.go` **When** it is read **Then** `read_file` and `write_file` tools are defined with `ToolDef`/`ToolFunction` wrappers for the Ollama wire format (AD-11) **And** `Tool.Run` is not serialized

5. **Given** `milestone-2/main.go` **When** `go build -o /dev/null ./milestone-2/` is run **Then** it compiles independently without touching any other directory (AD-2, AD-3)

## Tasks / Subtasks

- [x] Task 1: Set up milestone-2/main.go from milestone-1 baseline (AC: 2)
  - [x] Delete `milestone-2/.gitkeep`
  - [x] Create `milestone-2/main.go` starting from exact copy of `milestone-1/main.go`

- [x] Task 2: Add ToolCalls to Message and Tools to ChatRequest (AC: 2, 4)
  - [x] Add `ToolCalls []ToolCall \`json:"tool_calls,omitempty"\`` as third field of Message (after Content)
  - [x] Add `Tools []ToolDef \`json:"tools,omitempty"\`` as fourth field of ChatRequest (after Stream)

- [x] Task 3: Append new M2 type definitions after ChatResponse (AC: 3, 4)
  - [x] Add `ToolCallFunction` struct
  - [x] Add `ToolCall` struct
  - [x] Add `ToolFunction` struct
  - [x] Add `ToolDef` struct
  - [x] Add `Tool` struct with exactly 4 fields in canonical order (AD-5)

- [x] Task 4: Add tools concept block to main() before existing marshal call (AC: 1, 4)
  - [x] Define `tools := []Tool{...}` with read_file and write_file using os.ReadFile/os.WriteFile
  - [x] Build `toolDefs` slice from tools slice with ToolDef/ToolFunction wrappers
  - [x] Add `Tools: toolDefs` to ChatRequest literal

- [x] Task 5: Add "os" to imports (AC: 5)

- [x] Task 6: Verify build and vet (AC: 5)
  - [x] `go build -o /dev/null ./milestone-2/` exits 0
  - [x] `go vet ./milestone-2/` exits 0

### Review Findings

- [x] [Review][Decision] Non-permitted status-code guard block in main() — removed to conform to AD-1. [milestone-2/main.go:119]
- [x] [Review][Patch] Prompt mutated from M1: "What is 2+2?" → "What is an entity in DDD?" — reverted to `"What is 2+2?"`. [milestone-2/main.go:99]
- [x] [Review][Defer] Path traversal in read_file/write_file (no path sanitization) [milestone-2/main.go:64,79] — deferred, pre-existing; Run is dead code in M2, will matter in M3
- [x] [Review][Defer] Silent zero-value type assertion on args["path"]/args["content"] [milestone-2/main.go:64,77] — deferred, pre-existing; dead code in M2
- [x] [Review][Defer] Unbounded file read can exhaust memory [milestone-2/main.go:65] — deferred, pre-existing; dead code in M2
- [x] [Review][Defer] No HTTP timeout on Ollama request [milestone-2/main.go:113] — deferred, pre-existing from M1
- [x] [Review][Defer] json.RawMessage parameter schemas never validated [milestone-2/main.go:62,74] — deferred, pre-existing; dead code in M2
- [x] [Review][Defer] write_file hardcoded 0644 permissions [milestone-2/main.go:79] — deferred, pre-existing; dead code in M2

## Dev Notes

### Additive Delta Rule (AD-1)

M2 = M1 verbatim + minimal struct field tie-ins + appended type definitions + new concept block in main.

**Struct field additions** (internal inserts — only two permitted):
- `Message.ToolCalls` added so M3 (additive from M2) can read tool call responses without further struct modification
- `ChatRequest.Tools` added to carry toolDefs to Ollama wire format

**New type definitions** (appended after ChatResponse, in this order):
ToolCallFunction → ToolCall → ToolFunction → ToolDef → Tool

**New code in main()** (inserted before existing `prompt :=` / before marshal block):
tools slice definition → toolDefs loop → `Tools: toolDefs` in ChatRequest literal

No M1 lines are deleted or reordered.

### Canonical Struct Shapes

```go
// AD-10 with ToolCalls added for M2/M3
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// M2 ChatRequest with Tools
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []ToolDef `json:"tools,omitempty"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

// -- M2 CONCEPT BLOCK: new type definitions appended here --

type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ToolDef struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// AD-5: exactly 4 fields, exact order
type Tool struct {
	Name        string
	Description string
	Parameters  json.RawMessage
	Run         func(args map[string]any) string
}
```

### Tools Concept Block in main()

```go
tools := []Tool{
	{
		Name:        "read_file",
		Description: "Read the contents of a file",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}},"required":["path"]}`),
		Run: func(args map[string]any) string {
			path, _ := args["path"].(string)
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			return string(data)
		},
	},
	{
		Name:        "write_file",
		Description: "Write content to a file",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"content":{"type":"string"}},"required":["path","content"]}`),
		Run: func(args map[string]any) string {
			path, _ := args["path"].(string)
			content, _ := args["content"].(string)
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			return "ok"
		},
	},
}

toolDefs := make([]ToolDef, len(tools))
for i, t := range tools {
	toolDefs[i] = ToolDef{
		Type: "function",
		Function: ToolFunction{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		},
	}
}
```

Then the existing M1 marshal block gains `Tools: toolDefs`:

```go
body, err := json.Marshal(ChatRequest{
	Model: "llama3.2",
	Messages: []Message{
		{Role: "user", Content: prompt},
	},
	Stream: false,
	Tools:  toolDefs,
})
```

### Tool.Run Serialization Note (AC 4)

`encoding/json` skips `func` fields silently — no json tag needed; `Run` won't appear in marshalled output. Satisfies "Tool.Run is not serialized."

### Import Change

M1 imports: `bytes`, `encoding/json`, `fmt`, `io`, `log`, `net/http`
M2 adds: `"os"` — required for `os.ReadFile`, `os.WriteFile`

### Scope Boundary — DO NOT ADD in M2

- Agent loop (`for` loop)
- `switch` dispatch on ToolCalls
- `[]Message` history accumulation
- JSON unmarshal of `ToolCall.Function.Arguments`
- Any tool invocation / dispatch code
- CLI arg parsing, stdin reading

All of the above belong to Story 1.4 (M3).

### Build Verification Command

Use `go build -o /dev/null ./milestone-2/` — NOT bare `go build ./milestone-2/`. Binary named `milestone-2` would collide with the `milestone-2/` directory (same issue as M1 story).

### Architecture Invariants Binding This Story

| Invariant | Rule |
|---|---|
| AD-1 | Additive-only delta from M1: no deletions, no reordering |
| AD-2 | Standalone — `go build ./milestone-2/` compiles alone |
| AD-3 | Single file: `main.go` only, no helper files |
| AD-4 | HTTP POST to `http://localhost:11434/api/chat`; no LLM SDK |
| AD-5 | Tool struct: exactly 4 fields in order (Name, Description, Parameters, Run) |
| AD-9 | Zero external dependencies; stdlib only |
| AD-11 | Typed Go structs for all Ollama wire types (ToolCallFunction, ToolCall, ToolFunction, ToolDef) |
| AD-13 | `http.Post(url, "application/json", body)` — no custom http.Client |

### Project Structure Notes

```
agent-workshop/
  go.mod                    ← unchanged
  milestone-1/
    main.go                 ← unchanged (story 1.2)
  milestone-2/
    .gitkeep                ← DELETE
    main.go                 ← CREATE (this story)
  milestone-3/
    .gitkeep                ← unchanged
  starter/
    main.go                 ← unchanged (story 1.1)
```

### References

- Tool struct shape: [Source: _bmad-output/planning-artifacts/epics.md#AD-5]
- Message with ToolCalls: [Source: epics.md#AD-10, _bmad-output/implementation-artifacts/1-2-milestone-1-http-call-to-ollama.md#M2-Compatibility-Note]
- Ollama wire structs: [Source: epics.md#AD-11]
- Additive delta rule: [Source: epics.md#AD-1]
- Build collision fix: [Source: 1-2-milestone-1-http-call-to-ollama.md#Build-Verification-Command]
- M1 current code: [Source: milestone-1/main.go]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6[1m]

### Debug Log References

### Completion Notes List

- Created `milestone-2/main.go` as additive delta from milestone-1: added `ToolCalls` to `Message`, `Tools` to `ChatRequest`, appended 5 new structs (ToolCallFunction, ToolCall, ToolFunction, ToolDef, Tool with exactly 4 fields per AD-5), added tools concept block in main() with read_file/write_file implementations, added "os" import.
- Verified: `go build -o /dev/null ./milestone-2/` exits 0, `go vet ./milestone-2/` exits 0.
- Verified: marshalled ChatRequest contains non-empty `tools` array (2 entries: read_file, write_file).
- No M1 lines deleted or reordered (AD-1 satisfied). Single file, stdlib only (AD-3, AD-9).

### File List

- milestone-2/.gitkeep (deleted)
- milestone-2/main.go (created)

## Change Log

- 2026-06-24: Implemented Story 1.3 — created milestone-2/main.go extending milestone-1 with Tool struct, ToolDef/ToolFunction wire types, read_file/write_file tool definitions, and tools passed in Ollama ChatRequest.
