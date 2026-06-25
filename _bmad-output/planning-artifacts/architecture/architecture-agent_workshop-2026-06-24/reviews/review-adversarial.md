---
title: Adversarial Compatibility Review — agent_workshop ARCHITECTURE-SPINE.md
purpose: Find holes where two independent implementors obey all 9 ADs yet produce incompatible code
reviewer: Claude Code adversarial analysis
date: 2026-06-24
status: findings
---

# Adversarial Code Incompatibility Review

## Verdict

INCOMPATIBLE. Two competent Go engineers, each following every AD to the letter, will produce code that **compiles separately but diverges in observable structure** in ways that break the additive-diff property (AD-1) when attendees compare their code side-by-side or try to combine snippets during the live session.

---

## Top 5 Compatibility Holes

### HOLE 1: Message struct shape unspecified
**Risk Tier:** High  
**Affects:** Milestone-3 compatibility (AD-1, AD-6)  
**Problem:**

The spine mandates storing message history (AD-6: `[]Message` accumulated in the loop, never truncated) but never specifies the `Message` struct fields. The AD only says history is `[]Message` and is non-persistent.

**Adversarial Build A:**
```go
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```
(mirrors Ollama's native message shape)

**Adversarial Build B:**
```go
type Message struct {
    Author  string `json:"role"`
    Text    string `json:"content"`
    Tools   []ToolCall `json:"tool_calls,omitempty"`
}
```
(mirrors the Ollama response envelope with nested tool_calls)

**Why incompatible:** Both follow AD-6 (unbounded slice, never cleared). Both compile. But when Milestone-3 code builds the request body, Build A appends `Message{Role: "assistant", Content: response.Message.Content}` while Build B appends `Message{Author: "assistant", Text: response.Message.Content, Tools: toolCalls}`. An attendee trying to merge the two or compare line-diffs will see a completely different struct definition line 50–60, which violates AD-1 (additive only, no restructuring).

---

### HOLE 2: Tool request/response struct shape unspecified  
**Risk Tier:** High  
**Affects:** Milestone-2 and Milestone-3 (AD-4, AD-5, AD-7)  
**Problem:**

The spine specifies the `Tool` struct (AD-5: four fields, Name/Description/Parameters/Run). But the struct for the Ollama request body shape and the response shape (where tool_calls live) is completely unspecified. Ollama's HTTP API returns nested structures like:

```json
{
  "message": {
    "role": "assistant",
    "content": "...",
    "tool_calls": [
      {
        "function": {
          "name": "read_file",
          "arguments": "{...}"
        }
      }
    ]
  }
}
```

**Adversarial Build A:** Models Ollama's schema exactly inline:
```go
type ChatRequest struct {
    Model    string        `json:"model"`
    Messages []Message     `json:"messages"`
    Tools    []ChatTool    `json:"tools"`
    Stream   bool          `json:"stream"`
}
type ChatTool struct {
    Type     string        `json:"type"`
    Function ToolSpec      `json:"function"`
}
type ToolSpec struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Parameters  json.RawMessage `json:"parameters"`
}
type ChatResponse struct {
    Message Message `json:"message"`
}
type Message struct {
    Role      string     `json:"role"`
    Content   string     `json:"content"`
    ToolCalls []ToolCall `json:"tool_calls"`
}
type ToolCall struct {
    Function struct {
        Name      string `json:"name"`
        Arguments string `json:"arguments"`
    } `json:"function"`
}
```

**Adversarial Build B:** Flattens and renames for clarity:
```go
type OllamaRequest struct {
    Model   string   `json:"model"`
    Msgs    []Msg    `json:"messages"`
    Tools   []struct {
        Type  string `json:"type"`
        Func  struct {
            Name  string          `json:"name"`
            Desc  string          `json:"description"`
            Param json.RawMessage `json:"parameters"`
        } `json:"function"`
    } `json:"tools"`
    Stream bool `json:"stream"`
}
type OllamaReply struct {
    Msg struct {
        Role string
        Text string
        Calls []struct {
            Fn struct {
                FuncName string
                Args     string
            }
        }
    }
}
```

**Why incompatible:** Both follow AD-4 (HTTP only), both use `json.RawMessage` for Parameters (AD-5 consequence). Both compile. Both marshal/unmarshal correctly. But the struct names, field names, and nesting depth differ wildly. When an attendee opens Build A and Build B side-by-side at minute 30, the request marshalling code is completely different — not additive, not a diff, not something they can reason about as "the same thing with one new line added." This breaks the pedagogical invariant (AD-1).

---

### HOLE 3: Ollama request marshalling pattern unspecified
**Risk Tier:** Medium  
**Affects:** Milestone-1, Milestone-2, Milestone-3 (AD-4)  
**Problem:**

AD-4 says "All model calls are `http.Post` to `http://localhost:11434/api/chat`" and "request body marshalled from a Go struct, not a raw string" (Consistency Convention, line 106). But it does not mandate:
- Whether the request struct is marshalled inline or assigned to a variable first.
- What the variable is named (e.g., `req`, `request`, `chatReq`, `payload`).
- Whether it's marshalled on the same line as `json.Marshal` or split across lines.

**Adversarial Build A:**
```go
payload := map[string]interface{}{
    "model": "llama3.2",
    "messages": history,
    "tools": tools,
    "stream": false,
}
body, _ := json.Marshal(payload)
resp, _ := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
```

**Adversarial Build B:**
```go
type ApiRequest struct {
    Model    string        `json:"model"`
    Messages []Message     `json:"messages"`
    Tools    []Tool        `json:"tools"`
    Stream   bool          `json:"stream"`
}
req := ApiRequest{Model: "llama3.2", Messages: history, Tools: tools, Stream: false}
b, _ := json.Marshal(req)
resp, _ := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(b))
```

**Why incompatible:** Build A uses a bare `map[string]interface{}` (valid JSON marshalling per stdlib); Build B uses a typed struct. Both POST the same bytes to Ollama. Both follow AD-4. But an attendee's code will look *completely* different from the reference. If the workshop uses Build B as the reference, attendees who write Build A (especially those new to Go) will struggle during minute 30 when the facilitator points to "this line" — because the structure is unrecognizable. Also violates AD-1 additive property: adding the loop (M2→M3) requires different request shapes, not just appending lines.

---

### HOLE 4: Tool dispatch switch statement structure unspecified
**Risk Tier:** Medium  
**Affects:** Milestone-3 (AD-7)  
**Problem:**

AD-7 says "The agent loop exits when the model's response contains no `tool_calls` field" and Consistency Convention (line 107) says "Switch on `tool_calls[0].function.name`" for the two required tool names (`read_file` and `write_file`). But the spine does not specify:
- Whether dispatch is a switch on string, or if/else, or a map[string]func, or a range loop.
- Whether there's a default case for unknown tools.
- Whether there's an explicit check for empty `tool_calls` slice before indexing `[0]`.
- Whether the loop continues or breaks after tool execution.

**Adversarial Build A:**
```go
for {
    // ... call model, unmarshal response
    if len(response.Message.ToolCalls) == 0 {
        break
    }
    switch response.Message.ToolCalls[0].Function.Name {
    case "read_file":
        // ...
    case "write_file":
        // ...
    default:
        log.Fatal("unknown tool")
    }
}
```

**Adversarial Build B:**
```go
for {
    // ... call model, unmarshal response
    calls := response.Message.ToolCalls
    if calls == nil || len(calls) == 0 {
        break
    }
    name := calls[0].Function.Name
    if name == "read_file" {
        // ...
    } else if name == "write_file" {
        // ...
    }
    // implicitly continues loop
}
```

**Why incompatible:** Both follow AD-7 (stop on no tool_calls). Both exit the loop identically. But the *shape* of the dispatch code is completely different — switch vs if/else. An attendee comparing Build A to Build B will see a structurally different loop, which undermines the pedagogical claim that "the loop is simple — here it is" (FR-6 Phase 5 checkpoint). This violates AD-1 again: you cannot get from one to the other by appending; you must restructure.

---

### HOLE 5: HTTP client setup (global vs local, timeout, error handling) unspecified
**Risk Tier:** Medium  
**Affects:** Milestone-1, Milestone-2, Milestone-3 (AD-4)  
**Problem:**

AD-4 and the Consistency Convention say to use `http.Post`, but do not specify:
- Whether the client is the global `http.DefaultClient` or a custom `&http.Client{}`.
- Whether a timeout is set on the client or not (if timeout, what value? — PRD defers stop-condition robustness, line 178).
- Whether error handling is `if err != nil { log.Fatal(err) }` or a bare `if err != nil { log.Fatal() }`.
- Whether response body is closed immediately or left open.

**Adversarial Build A:**
```go
resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
data, _ := io.ReadAll(resp.Body)
```

**Adversarial Build B:**
```go
client := &http.Client{Timeout: 30 * time.Second}
resp, err := client.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
if err != nil {
    log.Fatal("request failed")
}
defer resp.Body.Close()
data, err := io.ReadAll(resp.Body)
if err != nil {
    log.Fatal("read failed")
}
```

**Why incompatible:** Build A is more compact; Build B is more defensive. Both follow AD-4 and work. But the import lists differ (Build B needs `time` package; Build A does not). The code structure differs at Milestone-1, making the reference binary's source (compiled from M3 which includes this pattern across all three milestones) unrecognizable to attendees who followed the other pattern. Violates AD-1 and breaks the "read the reference source and see your code" moment at minute 57.

---

## Secondary Gaps (not blocking, but unhelpful)

- **Model name string literal placement:** Spine does not say whether `"llama3.2"` is a const at the top, a magic string in the request marshalling, or passed as a function parameter. This affects code clarity but not compilation.
- **JSON unmarshalling error handling:** AD-4 says `json.Marshal` but does not specify `json.Unmarshal` error path (log.Fatal, panic, return error?). Both are valid per AD-4; inconsistency breaks attendee comprehension.
- **Imports ordering and grouping:** No convention given. Two implementations may have different import blocks, breaking the visual diff at milestone boundaries.

---

## Impact on Workshop Invariants

| Invariant | Hole Num | Impact |
| --- | --- | --- |
| AD-1 (Additive only) | 1, 2, 3, 4, 5 | High: Different struct/switch shapes mean adding milestone code is not a linear diff; attendees cannot follow the "one new block" narrative. |
| AD-6 (Message history) | 1 | High: If Message struct differs between builds, the history accumulation pattern is unrecognizable. |
| AD-7 (Stop condition) | 4 | High: Switch vs if/else dispatch changes the visible structure of "the loop" despite identical behavior. |
| AD-4 (HTTP only) | 3, 5 | Medium: Different marshalling and client patterns both work but confuse attendees during side-by-side code comparison. |

---

## Remediation Priorities

**Critical (implement before dry run):**
1. **Specify `Message` struct shape exactly** — use Ollama's native shape or define once in Consistency Conventions.
2. **Specify request/response struct hierarchy** — give explicit field names and nesting levels for the full Ollama request/response envelope.
3. **Mandate dispatch pattern** — if switch, provide the exact switch statement skeleton in the Consistency Conventions section; if if/else, mandate that instead.

**High (implement before first attendee session):**
4. **HTTP client setup convention** — specify whether timeout is required, whether custom client or default, whether error handling prints message or just logs.
5. **Marshalling pattern** — decide between bare map or typed struct, and provide the code skeleton in Consistency Conventions.

**Nice-to-have (before session if time permits):**
6. **Import ordering** — add a section to Consistency Conventions specifying import block order (stdlib only, grouped if multiple groups).
7. **Model name placement** — recommend const or magic string, and provide the line in the Consistency Conventions section.

---

## Conclusion

The ARCHITECTURE-SPINE.md successfully constrains *behavior* (what the code does) but under-constrains *structure* (how it looks). This is fatal for an educational workshop where the narrative is "here is the same code, evolved once, line by line." Two implementors following every rule will produce code that works identically but looks completely different — sufficient to shatter attendee confidence during the "compare your code to the reference" moment.

**Recommendation:** Convert the Consistency Conventions section from prose to code skeletons (full struct definitions, full dispatch switch, full marshalling pattern). Attendees will follow those skeletons and arrive at mutually recognizable code.

