---
review_type: technology-currency
subject: ARCHITECTURE-SPINE.md
reviewer: Claude Code (caveman-ultra)
date: 2026-06-24
verdict: CRITICAL GAPS — stack versions are asserted from training data, not reality-checked
---

# Technology Version & API Reality-Check Review

## Verdict
**CRITICAL GAPS** — The Stack section (L113–118) asserts Go/Ollama/llama3.2 versions and API shapes without confirmation against current reality. The entire milestone-3 agent design depends on Ollama tool_calls behavior that is not verified.

---

## Findings

### 1. Go 1.22+ — ACCEPTABLE
**Status:** No issue  
**Evidence:** Go 1.22 and later have stable stdlib for `encoding/json` and `net/http`. Zero external dependencies per AD-9 means version doesn't affect compatibility.  
**Risk:** Minimal.

---

### 2. Ollama 0.3+ (tool_calls support) — [UNVERIFIED — web check needed]
**Status:** CRITICAL ASSERTION WITHOUT EVIDENCE  
**Finding:** Line 116 asserts `0.3+` as the minimum version with "tool_calls support required", but provides:
- No release notes or changelog reference confirming tool_calls landed in 0.3
- No date when feature was added
- No API stability guarantee (could have changed in 0.4+, 0.5+, etc.)

**Reality risk:** If Ollama 0.3 did NOT have tool_calls, or if the API shape changed in later releases, all three milestones fail at runtime.  
**Fix needed:** Verify against Ollama changelog / GitHub releases when tool_calls support was introduced, and confirm it remains stable in current release used for workshop.

---

### 3. llama3.2 3B (tool_calls support via `ollama pull`) — [UNVERIFIED — web check needed]
**Status:** CRITICAL ASSUMPTION  
**Finding:** Line 117 specifies "llama3.2 3B (via `ollama pull llama3.2`)" with implicit assumption that:
- This model quantization supports tool use in Ollama
- Ollama's chat API exposes tool_calls for this model
- No mention of which quantization variant (GGUF precision, parameter count)

**Reality risk:** llama3.2 may exist in multiple quantizations; the 3B variant may not support tool use, or tool use may be restricted to larger variants (7B+). Workshop attendees run `ollama pull llama3.2` and receive a model that doesn't support `tool_calls`, breaking milestone-3.  
**Fix needed:** Confirm llama3.2 3B supports tool_calls in Ollama, test locally with `curl` against `/api/chat` endpoint, document exact quantization.

---

### 4. `/api/chat` with `stream: false` — [UNVERIFIED — web check needed]
**Status:** ASSUMED API SHAPE  
**Finding:** Line 106 specifies the HTTP request shape as:
- Endpoint: `/api/chat`
- Parameter: `stream: false`

No evidence that:
- This is the current correct shape for tool-use requests
- Tool calls can be returned with streaming disabled
- Response format matches the Code's tool dispatch logic (switch on `tool_calls[0].function.name`, line 107)

**Reality risk:** Ollama may require `stream: true` for tool use, or tool_calls may be formatted differently in streamed vs. non-streamed responses. The milestone-3 agent loop expects to parse `tool_calls` from a single JSON response; if Ollama requires streaming, the entire loop breaks.  
**Fix needed:** Test locally: issue a tool-use request to Ollama `/api/chat` with `stream: false` and confirm the response contains `tool_calls` at the top level with `[0].function.name` structure.

---

### 5. Tool struct shape: 4 fields, `json.RawMessage` for Parameters — [UNVERIFIED — web check needed]
**Status:** CANONICAL ASSUMPTION (AD-5)  
**Finding:** Lines 65–72 lock the Tool struct shape as:
```go
type Tool struct {
    Name        string
    Description string
    Parameters  json.RawMessage
    Run         func(args map[string]any) string
}
```

This design assumes:
- `Parameters` (as `json.RawMessage`) is the correct field to hold the JSON Schema object sent to Ollama
- Ollama's tool definition format expects exactly these four fields
- The `Run` callback signature matches the dispatch expectation

**Reality risk:** Ollama may expect different field names (e.g., `schema`, `input_schema`, `args_schema`), different field types (e.g., `Parameters` as a struct, not raw JSON), or additional metadata fields. If the struct doesn't match Ollama's expected shape, the HTTP request fails or is rejected.  
**Fix needed:** Test locally: construct a Tool in this shape, marshal to JSON, POST to `/api/chat` with `tools: [Tool]`, and confirm Ollama accepts the request without schema validation errors.

---

### 6. Tool dispatch: `tool_calls[0].function.name` — [UNVERIFIED — web check needed]
**Status:** RESPONSE SHAPE ASSUMPTION  
**Finding:** Line 107 specifies tool dispatch as "Switch on `tool_calls[0].function.name`". This assumes:
- Ollama's response includes a top-level `tool_calls` array
- Each tool call has a `function` object
- The function object has a `name` field (e.g., `"read_file"`, `"write_file"`)

**Reality risk:** Ollama may return tool calls in a different format (e.g., `tools[0].name`, `function_calls[0]`, etc.). If the JSON path doesn't match, the switch statement has no cases and tool dispatch silently fails.  
**Fix needed:** Test locally: invoke a tool and inspect the full JSON response structure. Confirm `tool_calls[0].function.name` is the correct path.

---

## Summary Table

| Item | Status | Risk |
|------|--------|------|
| Go 1.22+ | ✓ Acceptable | Low |
| Ollama 0.3+ | [UNVERIFIED] | **CRITICAL** |
| llama3.2 3B tools | [UNVERIFIED] | **CRITICAL** |
| `/api/chat` shape | [UNVERIFIED] | **CRITICAL** |
| Tool struct shape | [UNVERIFIED] | **CRITICAL** |
| `tool_calls[0].function.name` | [UNVERIFIED] | **CRITICAL** |

---

## Recommendation

Before finalizing the architecture spine:

1. **Test Ollama tool use locally** against the current version (0.5+?). Verify:
   - Ollama version and release date of tool_calls feature
   - llama3.2 3B model supports tool use (test with a simple tool request)
   - Request/response shapes match the code's expectations

2. **Document verified versions** in Stack section with evidence (e.g., "Ollama 0.4+ (tool_calls GA in 0.4.0, 2024-MM-DD)")

3. **Add API contract section** to architecture spine: document the exact JSON request shape, response shape for tool_calls, and tool definition schema — with example curl commands for verification.

4. **Lock the Tool struct to tested shape** once confirmed against Ollama API.

---

**Impact:** If these unverified assumptions are wrong, the first attendee running milestone-3 will encounter a runtime failure (tool_calls not returned, struct mismatch, or response parsing error). This blocks the core learning objective.
