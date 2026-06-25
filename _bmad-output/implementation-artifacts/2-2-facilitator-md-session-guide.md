# Story 2.2: FACILITATOR.md Session Guide

---
baseline_commit: cae8ec371799251ef2b4ccf0b02b8535794efcf2
---

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As Thomas (or any future facilitator),
I want a complete facilitator guide with talk-track anchors, per-phase timing, pacing risks, and copy-paste snippets,
so that I can run a reproducible session without relying on memory and a second facilitator can do the same with no prior experience.

## Acceptance Criteria

1. **Given** `FACILITATOR.md` **When** it is read **Then** it contains a section for each of the 7 phases with: timing target (e.g. "0–5 min"), milestone checkpoint where applicable, and the key verbal cue for that phase (FR-4, FR-6)

2. **Given** the Phase 1 talk-track in `FACILITATOR.md` **When** it is read **Then** it explicitly states the ring-structure framing: minute-0 and minute-57 are the same demo run twice — same visual, completely different understanding — and names this as the session's narrative spine (FR-4)

3. **Given** the Phase 6 talk-track in `FACILITATOR.md` **When** it is read **Then** it names LangChain, LlamaIndex, AutoGPT, and Claude Code when making the claim "every agent is this loop" (FR-4)

4. **Given** the pacing risks section in `FACILITATOR.md` **When** it is read **Then** it explicitly calls out both risks: Part 2 (Tool JSON schema `Parameters` field trips attendees — point them to the README snippet) and Part 3 (the loop insight needs time to land — do not rush "there is no magic") (FR-4)

5. **Given** the Phase 7 section in `FACILITATOR.md` **When** it is read **Then** it includes the optional deep-ceiling extension questions (stop condition robustness, tool failure handling, memory persistence, planning step) with framing guidance: "these are directions to explore, not code to write today" (FR-6) **And** it contains a cue for Thomas to mention the per-session Slack group and confirm all attendees have access

6. **Given** `FACILITATOR.md` **When** the pre-session checklist section is read **Then** it exists as a distinct section covering actions Thomas must complete before each session (at minimum: Slack group creation — detailed in Story 3.1)

7. **Given** `FACILITATOR.md` **When** read by a second facilitator with no prior session experience **Then** they can identify without outside help: what to say at each phase, when to check for milestone completion, where pacing risks are, and how to unblock the two most common attendee blockers (NFR-4)

8. **Given** `FACILITATOR.md` **When** copy-paste snippets section is read **Then** it contains at minimum: the Tool JSON schema `Parameters` snippet and the agent loop `for` skeleton that attendees can paste if they fall behind in Part 3

9. **Given** the Phase 6 section in `FACILITATOR.md` **When** it is read **Then** it includes a spot-check script: three specific questions Thomas asks 3–5 attendees at minute-57 to verify they can name all three components (HTTP call, tool dispatch, loop break) without prompting — and a pass threshold of ≥3/5 correct (NFR-1)

10. **Given** the Phase 7 section in `FACILITATOR.md` **When** it is read **Then** it includes a build-completion checkpoint: Thomas scans the room at minute-60 and notes how many attendees successfully ran their binary — target ≥70% (NFR-2)

## Tasks / Subtasks

- [x] Task 1: Create `FACILITATOR.md` at repo root (AC: 1)
  - [x] Add 7 phase sections in order, each with: timing target, milestone checkpoint (where applicable), key verbal cue
  - [x] Phase 1 (0–5 min): binary demo setup — no milestone checkpoint; verbal cue: "watch what this does — we'll come back to it"
  - [x] Phase 2 (5–10 min): definitions on board — no milestone checkpoint; verbal cue: define agent, tool, loop on whiteboard/screen
  - [x] Phase 3 (10–20 min): Milestone 1 checkpoint — `go build ./milestone-1/` + `go run ./milestone-1/`; verbal cue: "HTTP call is all it takes to talk to a model"
  - [x] Phase 4 (20–35 min): Milestone 2 checkpoint — `go build ./milestone-2/`; verbal cue: "tools are just Go functions wrapped in a struct"
  - [x] Phase 5 (35–50 min): Milestone 3 checkpoint — `go build ./milestone-3/` + run; verbal cue: "the loop is the agent"
  - [x] Phase 6 (50–57 min): binary walkthrough — walks through `milestone-3/main.go` on screen; verbal cue: "every agent framework is this loop"
  - [x] Phase 7 (57–60 min): attendee build + optional Q&A

- [x] Task 2: Add ring-structure framing in Phase 1 (AC: 2)
  - [x] State explicitly: minute-0 and minute-57 are the same `agent-demo` binary run with the same `demo-task.txt` task
  - [x] Name this the session's narrative spine: "same visual, completely different understanding"
  - [x] Include the run command: `./agent-demo < demo-task.txt` (or however Thomas runs it — confirm this is the actual invocation from the binary behavior in Story 1.5)
  - [x] Frame for facilitator: run the demo, let it complete, say nothing about how it works — the explanation comes at minute-57

- [x] Task 3: Add Phase 6 framework call-outs (AC: 3)
  - [x] When making the "every agent is this loop" claim, name: LangChain, LlamaIndex, AutoGPT, and Claude Code
  - [x] One sentence framing: "LangChain, LlamaIndex, AutoGPT, Claude Code — they're all this loop with more layers on top"
  - [x] Walk through `milestone-3/main.go` on screen: identify the three components (HTTP call lines, switch dispatch, loop break condition)

- [x] Task 4: Add Phase 6 spot-check script (AC: 9, NFR-1)
  - [x] Three specific questions Thomas asks 3–5 attendees verbatim:
    1. "What does the agent call to get a model response?" → expected: an HTTP POST / http.Post call
    2. "What makes the agent decide which tool to run?" → expected: switch on tool name / tool_calls
    3. "What makes the loop stop?" → expected: no tool_calls in the response / len(msg.ToolCalls) == 0
  - [x] Pass threshold: ≥3 of 5 spot-checked attendees answer all three correctly without prompting (NFR-1)
  - [x] If threshold not met: extend Phase 7 Q&A, defer deep-ceiling questions

- [x] Task 5: Add pacing risks section (AC: 4)
  - [x] Risk 1 — Part 2 (Phase 4), Tool JSON schema `Parameters` field:
    - Symptom: attendees stuck on what value to put in `Parameters json.RawMessage`
    - Response: point them to README.md "Tool JSON Schema" section — the write_file example is copy-pasteable
    - Do NOT re-explain the JSON Schema format verbally; the README does it better
  - [x] Risk 2 — Part 3 (Phase 5), loop insight:
    - Symptom: attendees moving on before the "for loop = agent" insight lands
    - Response: pause after completing Milestone 3; say "there is no magic" explicitly; let the silence sit for 5–10 seconds
    - Verbal cue: "This is it. There's nothing else hidden. The for loop is the agent."

- [x] Task 6: Add pre-session checklist section (AC: 6)
  - [x] Distinct section titled "Pre-Session Checklist" near the top of the file (before Phase 1)
  - [x] Minimum items:
    - [x] Create a Slack group conversation with all registered participants (Story 3.1 will expand this)
    - [x] Verify `agent-demo` binary is executable: `ls -la agent-demo` → shows executable bit
    - [x] Verify Ollama is running and llama3.2 is pulled: `ollama run llama3.2 "hello"` → returns response
    - [x] Run the demo dry: `./agent-demo < demo-task.txt` (or actual invocation) — confirm it completes in under 2 minutes
  - [x] Note: Story 3.1 will add the Slack process detail; this story adds the placeholder item

- [x] Task 7: Add copy-paste snippets section (AC: 8)
  - [x] Section titled "Copy-Paste Snippets" — place near the end or as an appendix
  - [x] Snippet 1 — Tool JSON schema `Parameters` value (same as README):
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
    Label: "Value for `Parameters json.RawMessage` field in Tool struct (write_file example)"
  - [x] Snippet 2 — Agent loop `for` skeleton for attendees who fall behind in Phase 5:
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
    Label: "Milestone 3 agent loop — paste into main() after toolDefs slice"
  - [x] Do NOT include full milestone code blocks — snippets are rescue anchors only

- [x] Task 8: Add Phase 7 deep-ceiling questions + Slack cue + build checkpoint (AC: 5, 10, NFR-2)
  - [x] Deep-ceiling extension questions (framing: "these are directions to explore, not code to write today"):
    1. Stop condition robustness: what happens if the model never stops calling tools? (hint: max-turns counter)
    2. Tool failure handling: what if `read_file` returns an error string? Does the model retry? Should the agent retry?
    3. Memory persistence: what if you want the agent to remember context across process restarts?
    4. Planning step: what if you sent a "plan first, then act" system message before the loop?
  - [x] Slack cue: "Before we wrap — confirm everyone has access to the Slack group I created for this cohort. If you don't see it, come find me." (FR-8, Story 2.2 AC 5)
  - [x] Build-completion checkpoint: at minute-60, Thomas scans the room. Target: ≥70% of attendees have successfully run `go build ./milestone-3/` and executed their binary (NFR-2). Note the count for repeatability tracking (NFR-3).

- [x] Task 9: Verify FACILITATOR.md self-sufficiency (AC: 7)
  - [x] Read through the whole document as a second facilitator would
  - [x] Confirm each question is answerable without outside help:
    - What to say at each phase? → verbal cues present for all 7 phases ✓
    - When to check for milestone completion? → phases 3, 4, 5 have explicit checkpoint commands ✓
    - Where are the pacing risks? → pacing risks section names both blockers ✓
    - How to unblock Part 2 attendees stuck on Parameters? → pacing risks section → README ✓
    - How to unblock Part 3 attendees who missed the loop insight? → pacing risks section verbal cue ✓

### Review Findings

- [x] [Review][Patch] Arguments type change milestone-2→milestone-3 — add Phase 4→5 transition note: "`Arguments` changes from `string` to `map[string]any` in Milestone 3 — update your struct" [FACILITATOR.md:Phase 4]
- [x] [Review][Patch] agent-demo binary not committed — remove binary from repo; add build-from-source step to pre-session checklist so any platform can build it [FACILITATOR.md:Pre-Session Checklist]
- [x] [Review][Patch] go build ./milestone-N/ fails — directory names conflict with binary output names; all three milestones fail to build from project root [FACILITATOR.md:58-61, 75-76, 92-96]
- [x] [Review][Patch] Phase 5 run path unreachable and stdin ignored — `./milestone-3/milestone-3` only exists if built from inside subdirectory (never explained); milestone-3 ignores `< demo-task.txt` (hardcodes prompt) [FACILITATOR.md:93-95]
- [x] [Review][Patch] demo-task.txt vs demo.txt ambiguity in Phase 1 — "Watch: the agent reads `demo.txt`" accurate but second facilitator may confuse with demo-task.txt (stdin); clarify distinction [FACILITATOR.md:30]
- [x] [Review][Patch] Phase 6 timing impossible — 7 minutes (50–57) must hold code walkthrough + second binary run + spot-check of 3–5 attendees; ring closure at minute 57 = Phase 7 start with no allocation guidance [FACILITATOR.md:114-154]
- [x] [Review][Patch] Pacing Risk 1 reference gap — README has only write_file schema; read_file schema absent; attendees needing read_file schema find nothing at the pointer [FACILITATOR.md:191-197]
- [x] [Review][Patch] Spot-check threshold undefined for <5 attendees — "≥3 of 5" breaks with 3 or 4 people; add proportional fallback (e.g. ≥60%) [FACILITATOR.md:151]
- [x] [Review][Patch] Pacing Risks "Part 2/Part 3" don't match phase numbering — phases numbered 1–7; "Part 2/Part 3" create cross-reference confusion for second facilitator [FACILITATOR.md:187, 201]
- [x] [Review][Patch] Pre-session dry run creates reversed.txt with no cleanup step — file persists into Phase 1 live demo; add delete step to checklist [FACILITATOR.md:14]
- [x] [Review][Patch] AD-12 gap — Phase 6 "Switch dispatch" label missing verbal cue "switch on tool name" [FACILITATOR.md:122]
- [x] [Review][Defer] Slack group creation has no pre-session verification step — deferred, pre-existing; Story 3.1 scope
- [x] [Review][Defer] Build-completion checkpoint at minute 60 = session end, no time to act — deferred, pre-existing; inherent to 60-min session design

## Dev Notes

### File to create

`FACILITATOR.md` at repo root — this file does **not** exist yet (confirmed: no `FACILITATOR.md` in repo root at story creation time; only `README.md`, `go.mod`, `agent-demo`, `demo-task.txt` exist).

### Session structure (exact timing)

| Phase | Timing | Activity | Milestone checkpoint |
|-------|--------|----------|---------------------|
| 1 | 0–5 min | Binary demo (ring-structure hook) | None — run `agent-demo` only |
| 2 | 5–10 min | Definitions on board | None |
| 3 | 10–20 min | Milestone 1 | `go build ./milestone-1/` + `go run ./milestone-1/` |
| 4 | 20–35 min | Milestone 2 | `go build ./milestone-2/` |
| 5 | 35–50 min | Milestone 3 | `go build ./milestone-3/` + run |
| 6 | 50–57 min | Binary walkthrough | None (code reading, not building) |
| 7 | 57–60 min | Attendee build + Q&A | Build-completion scan at minute-60 |

Source: FR-6, epics.md#FR-6

### Ring-structure narrative spine

The session's narrative arc is built around two identical demo runs:

- **Minute 0** (Phase 1): Run `agent-demo` cold. Attendees see the agent execute — files appear, the loop runs, text prints. They don't know how it works. Say nothing about the internals.
- **Minute 57** (Phase 6): Run `agent-demo` again. Attendees now understand the HTTP call, the tool dispatch, and the loop. Same visual. Completely different understanding.

The facilitator's job in Phase 1 is to plant the hook ("watch this — we'll explain every line later") and do nothing else. The payoff in Phase 6 closes the ring.

### The two most common attendee blockers (NFR-4)

1. **Part 2 — `Parameters json.RawMessage` format**: Attendees stare at the Tool struct and don't know what JSON to put in the Parameters field. **Response**: point to README.md "Tool JSON Schema" section — the write_file schema is already there as a copy-paste block. Do not re-explain verbally; the README snippet is clearer.

2. **Part 3 — loop insight doesn't land**: Attendees finish Milestone 3 but don't feel the "aha" that the for loop is the whole agent. **Response**: after the binary runs in Phase 5, pause. Say "There is no magic. The for loop is the agent." Let it sit. Don't move on for 5–10 seconds.

### Agent loop `for` skeleton (exact reference from milestone-3/main.go)

The loop that actually ships (from `milestone-3/main.go:103–127`):

```go
for {
    body, err := json.Marshal(ChatRequest{
        Model:    "llama3.2",
        Messages: messages,
        Stream:   false,
        Tools:    toolDefs,
    })
    if err != nil { log.Fatal(err) }

    resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
    if err != nil { log.Fatal(err) }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil { log.Fatal(err) }

    var result ChatResponse
    if err := json.Unmarshal(data, &result); err != nil { log.Fatal(err) }

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

The FACILITATOR.md snippet should be a **trimmed skeleton** (with comments replacing the HTTP boilerplate) so it's useful as a rescue anchor without being the full implementation. Full code is in `milestone-3/main.go`.

**Critical order constraint (AD-7, AD-12):** The `len(msg.ToolCalls) == 0` stop-condition check MUST appear before the `switch` dispatch. This is enforced by AD-7 and is the single most important structural thing to preserve in any loop skeleton or walkthrough.

### Tool JSON schema snippet (exact values from milestone-2/main.go)

The canonical schemas from `milestone-2/main.go`:
- `read_file`: `{"type":"object","properties":{"path":{"type":"string"}},"required":["path"]}`
- `write_file`: `{"type":"object","properties":{"path":{"type":"string"},"content":{"type":"string"}},"required":["path","content"]}`

Use `write_file` as the example (shows multi-property pattern — better teaching example). Pretty-print it in FACILITATOR.md (same format as in README.md).

### What the pre-session checklist placeholder needs

Story 3.1 will expand the Slack process detail. This story only needs to establish the section with a minimum viable item list. The Slack item should read: "Create a Slack group conversation with all registered participants for this session cohort." — Story 3.1 will add the detailed sub-steps.

### Scope boundaries — DO NOT do in this story

- **Do NOT create or modify any milestone files** — Epic 1 is frozen
- **Do NOT modify README.md** — Story 2.1 is done; README is complete
- **Do NOT implement Story 3.1 Slack detail** — Story 3.1 covers the Slack process; this story only adds the placeholder checklist item
- **Do NOT add slide deck content or links** — out of scope per PRD §6
- **Do NOT add a Docker fallback section** — explicitly out of scope
- **Do NOT add per-milestone code blocks** beyond the two required snippets (loop skeleton, Parameters schema)

### Architecture invariants binding this story

FACILITATOR.md is documentation only — no structural invariants from ARCHITECTURE-SPINE.md apply directly. However, the content must accurately reflect the architecture:

| Invariant | Impact on FACILITATOR.md content |
|-----------|----------------------------------|
| AD-1 | Milestone walkthrough in phase sections must describe additive delta, not rewrites |
| AD-4 | Phase 6 walkthrough uses `http.Post` — not a client library |
| AD-7 | Stop-condition explanation: "loop exits when model returns no tool_calls" — no other condition |
| AD-9 | Build commands are `go build ./milestone-N/` with no `go get` or module installs |
| AD-12 | Switch dispatch — facilitator says "switch on tool name" not "map dispatch" |

### Previous story intelligence (Story 2.1)

- `README.md` is complete at repo root. It contains: install links, `ollama pull llama3.2`, two-step verify section, Tool JSON schema snippet (write_file), three-milestone walkthrough table.
- The README "Tool JSON Schema" section is the canonical self-rescue anchor for Part 2 — FACILITATOR.md should reference it by name when describing the Part 2 pacing risk.
- Story 2.1 dev notes: `go build ./starter/` was adapted to `go build -o /dev/null ./starter/` due to directory naming conflict — this affects the pre-session checklist verify step (if included) but NOT the phase checkpoints (phases 3–5 use `./milestone-N/`).
- No issues with milestone-1, milestone-2, milestone-3 builds — all compile independently.

### Project structure notes

```
agent-workshop/
  FACILITATOR.md   ← CREATE (this story)
  README.md        ← unchanged (Story 2.1 done)
  go.mod           ← unchanged
  agent-demo       ← unchanged (binary, darwin/arm64)
  demo-task.txt    ← unchanged
  starter/main.go  ← unchanged
  milestone-1/main.go ← unchanged
  milestone-2/main.go ← unchanged
  milestone-3/main.go ← unchanged
```

### References

- Story 2.2 ACs + FR-4, FR-6: [Source: _bmad-output/planning-artifacts/epics.md#Story-2.2]
- 7-phase session structure: [Source: epics.md#FR-6]
- Ring-structure framing: [Source: epics.md#FR-4]
- Framework call-outs (LangChain, LlamaIndex, AutoGPT, Claude Code): [Source: epics.md#FR-4]
- Pacing risks: [Source: epics.md#FR-4]
- Deep-ceiling extension questions: [Source: epics.md#FR-6]
- Slack group cue: [Source: epics.md#Story-2.2 AC 5, FR-8]
- Build-completion checkpoint (NFR-2): [Source: epics.md#NFR-2]
- Spot-check script (NFR-1): [Source: epics.md#NFR-1]
- Self-sufficiency requirement (NFR-4): [Source: epics.md#NFR-4]
- Agent loop skeleton: [Source: milestone-3/main.go:103–127]
- Tool struct Parameters schemas: [Source: milestone-2/main.go]
- Stop-condition constraint (AD-7): [Source: architecture/ARCHITECTURE-SPINE.md#AD-7]
- Switch dispatch (AD-12): [Source: architecture/ARCHITECTURE-SPINE.md#AD-12]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6[1m]

### Debug Log References

None.

### Completion Notes List

- Created `FACILITATOR.md` at repo root with all 9 tasks implemented in a single pass.
- 7 phase sections (0–60 min) with timing, milestone checkpoints (phases 3–5), and verbal cues.
- Ring-structure narrative spine: Phase 1 plants hook at minute-0, Phase 6 closes ring at minute-57 with same `./agent-demo < demo-task.txt` invocation.
- Phase 6 spot-check script: 3 questions, ≥3/5 pass threshold, defer deep-ceiling if not met.
- Pacing Risks section: Part 2 Parameters blocker → README.md; Part 3 loop insight → 5–10s pause + verbal cue.
- Copy-paste snippets: `write_file` JSON schema and trimmed loop skeleton with critical stop-before-switch order preserved.
- Phase 7: deep-ceiling extension questions, Slack cue, build-completion checkpoint (≥70% at minute-60).
- Pre-session checklist: 4 items including Slack placeholder (Story 3.1 expands), binary verify, Ollama verify, dry run.
- Self-sufficiency verified: all 5 second-facilitator questions answerable without outside help.
- No milestone files, README, or Story 3.1 Slack detail touched (scope boundary respected).

### File List

- FACILITATOR.md (created)
