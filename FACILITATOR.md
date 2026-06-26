# FACILITATOR.md — Agent Workshop Session Guide

This document is the complete run-sheet for the agent workshop. A second facilitator with no prior session experience should be able to run the session end-to-end using only this file.

---

## Pre-Session Checklist

Complete these **before** attendees arrive:

- [ ] **Slack group (per-cohort — create new for each session; do not reuse across sessions):**
  - [ ] Create a direct Slack group conversation with all registered participants for this cohort
  - [ ] Send a welcome message before the session date — encourage setup questions; confirm attendees can reply
  - [ ] Day-of: confirm all registered participants are in the group before Phase 1 starts
- [ ] Build `agent-demo` from source: `go build -o agent-demo ./milestone-3/` → binary appears in repo root
- [ ] Verify Ollama is running and `llama3.2` is pulled: `ollama run llama3.2 "hello"` → returns a response
- [ ] Run the demo dry: `./agent-demo < demo-task.txt` — confirm it completes in under 2 minutes and produces `reversed.txt`
- [ ] Delete `reversed.txt` after dry run: `rm -f reversed.txt`

---

## Session Phases

### Phase 1 — Binary Demo (0–5 min)

**Milestone checkpoint:** None — do not reference milestones yet.

**What to do:** Run `agent-demo` cold, in front of the room. Say nothing about how it works.

```bash
./agent-demo < demo-task.txt
```

Watch: the agent uses its `read_file` tool to read `demo.txt`, then writes `reversed.txt`. (`demo-task.txt` is the task prompt piped via stdin — `demo.txt` is the file the agent reads and reverses.) Files appear. The loop runs. Text prints.

**Verbal cue:** *"Watch what this does — we'll come back to it."*

Then move directly to the whiteboard. Do not explain the binary. The explanation comes at minute 57.

> **Ring-structure narrative spine:** This is minute 0 of a two-run structure. Minute 57 (Phase 6) runs the same binary with the same `demo-task.txt` — same visual, completely different understanding. The session's entire narrative arc is built around this ring: plant the hook now, close it in Phase 6. Your job here is to create curiosity, not explain anything.

---

### Phase 2 — Definitions on Board (5–10 min)

**Milestone checkpoint:** None.

**What to do:** Write these three definitions on the whiteboard or screen. Say them aloud as you write.

1. **Agent** — a program that calls a model in a loop
2. **Tool** — a function the model can request
3. **Loop** — runs until the model stops requesting tools

**Verbal cue:** *"Agent. Tool. Loop. Three words. By the end of today you'll have built all three."*

---

### Phase 3 — Milestone 1: HTTP Call to Ollama (10–20 min)

**Milestone checkpoint:**

```bash
go build -o ./milestone-1-bin ./milestone-1/
./milestone-1-bin
```

Expected: program sends a chat request to Ollama and prints the model response.

**What to do:** Walk attendees through `milestone-1/main.go`. Show the `http.Post` call to `http://localhost:11434/api/chat`. Show the JSON marshal/unmarshal. Let them build and run.

**Verbal cue:** *"An HTTP call is all it takes to talk to a model. That's the first component."*

---

### Phase 4 — Milestone 2: Tool Struct and Tools Slice (20–35 min)

**Milestone checkpoint:**

```bash
go build -o /dev/null ./milestone-2/
```

Expected: compiles cleanly. (No run needed — milestone-2 exercises the struct definitions.)

**What to do:** Walk attendees through the `Tool` struct and `toolDefs` slice. Show how `json.RawMessage` holds the JSON schema. Let them build.

**Verbal cue:** *"Tools are just Go functions wrapped in a struct. The model never calls your function directly — it asks, you run it."*

> **Transition note for Milestone 3:** The `Arguments` field changes type from `string` (Milestone 2) to `map[string]any` (Milestone 3). Attendees evolving their own code will need to update that field in their `ToolCall` struct.

> **See [Pacing Risks](#pacing-risks) → Risk 1** if attendees get stuck on the `Parameters` field value.

---

### Phase 5 — Milestone 3: Agent Loop, Dispatch, and History (35–50 min)

**Milestone checkpoint:**

```bash
go build -o ./milestone-3-bin ./milestone-3/
./milestone-3-bin
```

Expected: the agent reads `demo.txt`, writes `reversed.txt`, stops when no more tool calls.

**What to do:** Walk attendees through the `for` loop in `milestone-3/main.go`. Show the stop condition (`len(msg.ToolCalls) == 0`) **before** the switch dispatch — this order is critical. Show how `messages` accumulates history. Let them build and run.

**Verbal cue:** *"The loop is the agent. The for loop is the whole thing."*

After the binary runs — **pause**. Say:

> *"This is it. There's nothing else hidden. The for loop is the agent."*

Let that sit for 5–10 seconds. Do not move on immediately.

> **See [Pacing Risks](#pacing-risks) → Risk 2** if the loop insight doesn't land.

---

### Phase 6 — Binary Walkthrough (50–57 min)

**Milestone checkpoint:** None (code reading, not building).

**What to do:**

1. **Run `agent-demo` again** (ring closure — do this first):

```bash
./agent-demo < demo-task.txt
```

**Verbal cue:** *"Same binary as minute 0. Same task file. Different eyes."*

2. Open `milestone-3/main.go` on screen. Walk through the three components:
   - **HTTP call lines** — `http.Post` to `http://localhost:11434/api/chat`
   - **Switch on tool name** — `switch msg.ToolCalls[0].Function.Name` (say: *"switch on tool name"*)
   - **Loop break condition** — `if len(msg.ToolCalls) == 0 { break }`

**Verbal cue:** *"LangChain, LlamaIndex, AutoGPT, Claude Code — they're all this loop with more layers on top. Every agent framework is this loop."*

---

#### Phase 6 Spot-Check Script

**Target: minute 54–55** (before Phase 7 begins at minute 57). Ask **3–5 attendees** these three questions verbatim. Do not prompt or hint.

1. *"What does the agent call to get a model response?"*
   → Expected answer: an HTTP POST / `http.Post` call

2. *"What makes the agent decide which tool to run?"*
   → Expected answer: switch on tool name / `tool_calls` field

3. *"What makes the loop stop?"*
   → Expected answer: no tool calls in the response / `len(msg.ToolCalls) == 0`

**Pass threshold:** ≥3 of 5 spot-checked attendees answer all three correctly without prompting (or ≥60% if fewer than 5 are present).

**If threshold not met:** extend Phase 7 Q&A time. Defer the deep-ceiling extension questions. Spend the time reinforcing the three components before wrapping.

---

### Phase 7 — Attendee Build + Optional Q&A (57–60 min)

**Milestone checkpoint:** Build-completion scan at minute 60.

**What to do:** Give attendees time to run their own binary. Circulate. Answer questions.

**Slack cue:** *"Before we wrap — confirm everyone has access to the Slack group I created for this cohort. If you don't see it, come find me after."*

**Build-completion checkpoint (minute 60):** Scan the room. Count how many attendees have successfully run `go build -o ./milestone-3-bin ./milestone-3/` and executed their binary.
- **Target:** ≥70% of attendees built and ran successfully.
- Note the count for repeatability tracking (e.g., "8/10 built successfully").
- If below 70%, note blockers for post-session debrief.

---

#### Deep-Ceiling Extension Questions

*Use only if time allows and spot-check threshold was met. Frame as:*

> *"These are directions to explore, not code to write today."*

1. **Stop condition robustness:** What happens if the model never stops calling tools? *(hint: add a max-turns counter)*
2. **Tool failure handling:** What if `read_file` returns an error string? Does the model retry? Should the agent retry?
3. **Memory persistence:** What if you want the agent to remember context across process restarts?
4. **Planning step:** What if you sent a "plan first, then act" system message before the loop?

---

### Phase 8 *(optional)* — Milestone 4: Test the Output (overflow / advanced attendees)

**When to run:** Only if the spot-check threshold was met *and* time remains, or for attendees who finished Milestone 3 early and are already fluent. This is not part of the core 60-minute arc — it is an extension for the deep-ceiling crowd.

**Concept:** An LLM is probabilistic — the model can return text that is the wrong length or never write the file at all. Milestone 4 adds **output validation fed back into the loop**: when the model says it's done, the agent runs a deterministic check first; on failure it sends the failure back as a new turn and the agent self-corrects. A `maxTurns` cap stops an infinite retry cycle.

**Milestone checkpoint:**

```bash
go test ./milestone-4/
go build -o ./milestone-4-bin ./milestone-4/
./milestone-4-bin
```

Expected: `go test` passes 3 checks with no Ollama running (the checks are pure functions). The binary reads `demo.txt`, writes `reversed.txt`, validates it, and stops once validation passes.

**What to do:** Walk attendees through three things in `milestone-4/main.go`:

- `checkReversal` — the invariant we assert on the result (the reversed text must keep the same length). *"This is the test. The model made a request to reverse; this is us checking it didn't drop or add characters."*
- The validation branch where `len(msg.ToolCalls) == 0` used to just `break` — now it verifies before trusting the model, and on failure appends a `user` message and loops again.
- `maxTurns` — *"the self-correction loop needs a stop, or a stubborn model loops forever."*

**Verbal cue:** *"The loop is the agent. The tests are what make the agent trustworthy. Every production agent has both."*

> **Tie-back:** This is the natural home for Deep-Ceiling Extension Question 1 (stop-condition robustness) and 2 (tool failure handling) — Milestone 4 is those questions turned into running code.

---

## Pacing Risks

### Risk 1 — Milestone 2 (Phase 4): `Parameters json.RawMessage` field

**Symptom:** Attendees stare at the `Tool` struct and don't know what JSON value to put in the `Parameters` field.

**Response:** Point them to the README.md "Tool JSON Schema" section — the `write_file` schema is already there as a copy-paste block.

```
→ README.md → "Tool JSON Schema" section → write_file example
```

**Do NOT re-explain the JSON Schema format verbally.** The README snippet is clearer and self-contained. Just point and move on.

If the attendee needs the `read_file` schema (path only, no content): point them to `milestone-2/main.go` tool definitions — the `read_file` schema is defined there.

---

### Risk 2 — Milestone 3 (Phase 5): Loop insight doesn't land

**Symptom:** Attendees finish Milestone 3 and move on without the "aha" — they built the loop but don't feel that the for loop *is* the agent.

**Response:** After the binary runs, stop. Do not move to the next thing.

Say: *"This is it. There's nothing else hidden. The for loop is the agent."*

Then wait 5–10 seconds in silence. Let it land. The pause is the point.

---

## Copy-Paste Snippets

These are rescue anchors for attendees who fall behind. Do not use them as teaching material — they are for unblocking, not explaining.

---

### Snippet 1 — Tool JSON Schema `Parameters` value

**Label:** Value for `` `Parameters json.RawMessage` `` field in `Tool` struct (`write_file` example)

```json
{
  "type": "object",
  "properties": {
    "path":    { "type": "string" },
    "content": { "type": "string" }
  },
  "required": ["path", "content"]
}
```

---

### Snippet 2 — Milestone 3 Agent Loop Skeleton

**Label:** Milestone 3 agent loop — paste into `main()` after `toolDefs` slice

```go
for {
    // marshal ChatRequest{Model, Messages, Stream:false, Tools:toolDefs}
    // http.Post to "http://localhost:11434/api/chat"
    // unmarshal ChatResponse, get msg := result.Message
    messages = append(messages, msg)
    if len(msg.ToolCalls) == 0 {
        fmt.Println(msg.Content)
        break
    }
    switch msg.ToolCalls[0].Function.Name {
    case "read_file":
        output = tools[0].Run(msg.ToolCalls[0].Function.Arguments)
    case "write_file":
        output = tools[1].Run(msg.ToolCalls[0].Function.Arguments)
    }
    messages = append(messages, Message{Role: "tool", Content: output})
}
```

> **Critical order:** The `len(msg.ToolCalls) == 0` stop-condition check appears **before** the `switch` dispatch. Preserve this order.
