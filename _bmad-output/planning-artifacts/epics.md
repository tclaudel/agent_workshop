---
stepsCompleted: [step-01-validate-prerequisites, step-02-design-epics, step-03-create-stories, step-04-final-validation]
inputDocuments:
  - _bmad-output/planning-artifacts/prds/prd-agent_workshop-2026-06-24/prd.md
  - _bmad-output/planning-artifacts/architecture/architecture-agent_workshop-2026-06-24/ARCHITECTURE-SPINE.md
---

# agent_workshop - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for agent_workshop, decomposing the requirements from the PRD and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

FR-1: Repo contains 3 milestone dirs (`milestone-1/`, `milestone-2/`, `milestone-3/`), each a complete runnable Go program. M1 = HTTP call to Ollama; M2 = M1 + Tool struct + tools slice; M3 = M2 + agent loop + tool dispatch + message history.
FR-2: `starter/main.go` — package declaration, required imports, and an empty `main()` stub; compiles with `go build`; zero logic beyond the stub.
FR-3: Pre-built reference binary committed to repo root as `agent-demo`; runs on Thomas's darwin/arm64 machine; demonstrates observable file changes when given a task (minute-0 hook and minute-57 payoff).
FR-4: `FACILITATOR.md` — per-phase timing targets, key verbal cues, pacing risks (Part 2 Tool JSON schema parameters, Part 3 loop insight "do not rush"), copy-paste snippets; ring-structure framing (minute-0 and minute-57 are same demo, different understanding); Phase 6 framework call-outs (LangChain, LlamaIndex, AutoGPT, Claude Code); self-sufficient for second facilitator with no prior session experience.
FR-5: `README.md` — OS-specific install commands for Go and Ollama (or links to canonical pages), `ollama pull llama3.2`, verify-setup smoke test (go build+run + `ollama run llama3.2 "hello"`), Tool JSON schema copy-paste snippet for `Parameters` field.
FR-6: Session follows 7-phase structure with specified timing: Phase 1 (0–5 min) binary demo, Phase 2 (5–10 min) definitions on board, Phase 3 (10–20 min) M1 checkpoint, Phase 4 (20–35 min) M2 checkpoint, Phase 5 (35–50 min) M3 checkpoint, Phase 6 (50–57 min) binary walkthrough, Phase 7 (57–60 min) attendee build + optional deep-ceiling questions.
FR-7: Session runs with Thomas as sole facilitator; no external services required beyond Ollama running locally on Thomas's machine; all demo artifacts committed to repo.
FR-8: Thomas creates a Slack group conversation with all registered participants before each session (pre-session setup questions + post-session follow-up); group mentioned during Phase 7.

### NonFunctional Requirements

NFR-1: Comprehension rate — ≥3 of 5 spot-checked attendees correctly identify all three components (HTTP call, tool dispatch, loop break) at minute-57 demystification step without prompting. (SM-1)
NFR-2: Build completion rate — ≥70% of attendees successfully run `go build` and execute their binary by minute 60. (SM-2)
NFR-3: Repeatability — same repo used for ≥2 sessions without structural changes required. (SM-4)
NFR-4: FACILITATOR.md self-sufficiency — a second facilitator with no prior session experience can run the session from FACILITATOR.md alone.

### Additional Requirements

- AD-1 (Additive-only milestone delta): `milestone-N+1/main.go` = `milestone-N/main.go` + exactly one new contiguous concept block. No lines deleted, no reordering of existing code between milestones.
- AD-2 (Milestone isolation): Each milestone is a standalone Go program; no cross-milestone imports; `go build ./milestone-N/` compiles without touching any other directory.
- AD-3 (Single-file per milestone): Each milestone contains exactly one file: `main.go`. No sub-packages, no `internal/`, no helper files.
- AD-4 (Ollama HTTP API only): All model calls use `http.Post` to `http://localhost:11434/api/chat`; no LLM SDK; no framework import; stdlib only.
- AD-5 (Canonical Tool struct): 4 fields in order — `Name string`, `Description string`, `Parameters json.RawMessage`, `Run func(args map[string]any) string`. No additional fields.
- AD-6 (Message history): `[]Message` accumulated in agent loop; never truncated; never written to disk; cleared on process exit.
- AD-7 (Stop condition): Agent loop exits when model response contains no `tool_calls` field. No other stop condition in milestone code.
- AD-8 (Reference binary): Compiled from `milestone-3/main.go` for `GOOS=darwin GOARCH=arm64`; committed to repo root as `agent-demo`.
- AD-9 (Single go.mod): One `go.mod` at repo root; module path `github.com/tclaudel/agent-workshop`; zero external dependencies; stdlib only.
- AD-10 (Canonical Message struct): Fields in order — `Role string`, `Content string`, `ToolCalls []ToolCall \`json:"tool_calls,omitempty"\``.
- AD-11 (Typed Ollama structs): Typed Go structs for all request/response bodies: `ChatRequest`, `ToolDef`, `ToolFunction`, `ChatResponse`, `ToolCall`, `ToolCallFunction`. No `map[string]any`, no raw byte templates.
- AD-12 (Switch dispatch): Tool dispatch inside agent loop uses `switch msg.ToolCalls[0].Function.Name`; one `case` per tool. No if/else chain, no function-map dispatch.
- AD-13 (http.Post shorthand): All Ollama calls use `http.Post(url, "application/json", body)`; no custom `http.Client`, no `http.NewRequest`, no timeout config in milestone code.

### UX Design Requirements

N/A — workshop has no UI component.

### FR Coverage Map

FR-1: Epic 1 — 3 milestone dirs, each a complete runnable Go program
FR-2: Epic 1 — starter/main.go (package + imports only)
FR-3: Epic 1 — agent-demo binary committed at repo root
FR-4: Epic 2 — FACILITATOR.md with talk-track and pacing
FR-5: Epic 2 — README prereq guide + Tool JSON schema snippet
FR-6: Epic 2 — 7-phase session structure documented in FACILITATOR.md
FR-7: Epic 2 — solo facilitation requirements (no external services) documented
FR-8: Epic 3 — per-session Slack group conversation

## Epic List

### Epic 1: Workshop Codebase
Thomas and attendees have a complete, runnable workshop codebase they can clone and build from. Creates go.mod, starter/main.go, milestone-1 through milestone-3 (each a standalone Go program with additive concept blocks), the agent-demo reference binary, and demo-task.txt.
**FRs covered:** FR-1, FR-2, FR-3
**Arch invariants:** AD-1 through AD-13

### Epic 2: Workshop Documentation
Any facilitator can run a reproducible session without prior experience; attendees can self-rescue during setup and Part 2 Tool JSON schema. Creates README.md (prereq guide, smoke test, Tool JSON schema snippet) and FACILITATOR.md (7-phase talk-track, timing, ring-structure framing, pacing risks, Phase 6 framework call-outs).
**FRs covered:** FR-4, FR-5, FR-6, FR-7
**NFRs addressed:** NFR-3, NFR-4

### Epic 3: Post-Workshop Community
Attendees have a Slack channel for pre-session setup help and post-session follow-up; Thomas has a documented process for creating it per session.
**FRs covered:** FR-8

---

## Epic 1: Workshop Codebase

Thomas and attendees have a complete, runnable workshop codebase they can clone and build from. Creates go.mod, starter/main.go, milestone-1 through milestone-3 (each a standalone Go program with additive concept blocks), the agent-demo reference binary, and demo-task.txt.

### Story 1.1: Repository Foundation

As Thomas,
I want an initialized Go module with the correct directory skeleton and a starter file,
So that the repo is clonable and buildable before any milestone code exists.

**Acceptance Criteria:**

**Given** the repo is cloned
**When** `go build ./starter/` is run
**Then** it compiles with no errors and produces no output (no implementation code)

**Given** the repo root
**When** the directory structure is listed
**Then** `go.mod`, `starter/`, `milestone-1/`, `milestone-2/`, `milestone-3/` all exist

**Given** the `go.mod` file
**When** it is read
**Then** the module path is `github.com/tclaudel/agent-workshop`, Go version is 1.22+, and there are zero `require` directives (stdlib only — AD-9)

**Given** `starter/main.go`
**When** it is read
**Then** it contains only the `package main` declaration, import statements, and an empty `main()` function body — no other logic (FR-2)

### Story 1.2: Milestone 1 — HTTP Call to Ollama

As an attendee,
I want a complete working Go program that sends a prompt to Ollama and prints the reply,
So that I have a reference for the first concept: a raw HTTP call to a local LLM.

**Acceptance Criteria:**

**Given** Ollama is running with llama3.2 pulled
**When** `go run ./milestone-1/` is executed with a prompt
**Then** the program sends a request to `http://localhost:11434/api/chat` and prints the model's text reply to stdout

**Given** `milestone-1/main.go`
**When** it is read
**Then** all Ollama calls use `http.Post(url, "application/json", body)` — no custom `http.Client`, no `http.NewRequest` (AD-13)
**And** request and response bodies are typed Go structs, not `map[string]any` or raw byte templates (AD-11)
**And** the request sets `stream: false`
**And** there are zero external imports beyond stdlib (AD-4, AD-9)

**Given** `milestone-1/main.go`
**When** the prompt input mechanism is read
**Then** the prompt is a hardcoded string literal in `main()` — no CLI argument parsing, no stdin reading, no flag package (pedagogically simplest for live session; attendees see the full intent at a glance)

**Given** `milestone-1/main.go`
**When** `go build ./milestone-1/` is run
**Then** it compiles independently without touching any other directory (AD-2, AD-3)

### Story 1.3: Milestone 2 — Tool Struct and Tools Slice

As an attendee,
I want milestone-1 extended with a Tool struct, `read_file`/`write_file` implementations, and tools passed in the Ollama request,
So that I can see exactly what "adding tools" means with no other code changed.

**Acceptance Criteria:**

**Given** Ollama is running with llama3.2 pulled and a task that requires file access
**When** `go run ./milestone-2/` is executed
**Then** the Ollama request body includes a non-empty `tools` array — verified by inspecting the marshalled `ChatRequest` struct before it is sent (not the model response, which is non-deterministic)

**Given** `milestone-2/main.go`
**When** compared line-by-line to `milestone-1/main.go`
**Then** every line from milestone-1 is present and unchanged, and the new concept block (Tool struct + tools slice + ToolDef wrappers) is appended with no deletions or reordering (AD-1)

**Given** the `Tool` struct in `milestone-2/main.go`
**When** it is read
**Then** it has exactly 4 fields in this order: `Name string`, `Description string`, `Parameters json.RawMessage`, `Run func(args map[string]any) string` — no additional fields (AD-5)

**Given** `milestone-2/main.go`
**When** it is read
**Then** `read_file` and `write_file` tools are defined with `ToolDef`/`ToolFunction` wrappers for the Ollama wire format (AD-11)
**And** `Tool.Run` is not serialized

**Given** `milestone-2/main.go`
**When** `go build ./milestone-2/` is run
**Then** it compiles independently without touching any other directory (AD-2, AD-3)

### Story 1.4: Milestone 3 — Agent Loop, Dispatch, and History

As an attendee,
I want milestone-2 extended with the agent `for` loop, `switch`-based tool dispatch, and `[]Message` history accumulation,
So that I have the complete working agent with all three concepts introduced incrementally and nothing hidden.

**Acceptance Criteria:**

**Given** Ollama running with llama3.2 pulled and a file-editing task
**When** `go run ./milestone-3/` is executed
**Then** the agent loop runs, dispatches `read_file` and `write_file` tool calls, accumulates message history, and terminates when the model returns a text reply with no tool calls (AD-7)

**Given** `milestone-3/main.go`
**When** compared line-by-line to `milestone-2/main.go`
**Then** every line from milestone-2 is present and unchanged, and the new concept block (loop + switch dispatch + history append) is appended with no deletions or reordering (AD-1)

**Given** the agent loop in `milestone-3/main.go`
**When** tool dispatch is read
**Then** it uses `switch msg.ToolCalls[0].Function.Name` with one `case` per tool name — no if/else chain, no function-map (AD-12)

**Given** the stop condition in `milestone-3/main.go`
**When** it is read
**Then** the loop exits when `len(msg.ToolCalls) == 0` — no other stop condition (AD-7)

**Given** the agent loop in `milestone-3/main.go`
**When** the loop body is read
**Then** the `len(msg.ToolCalls) == 0` stop-condition check appears before the `switch` dispatch statement — dispatch is never reached on an empty `ToolCalls` slice (AD-7, AD-12)

**Given** the message history in `milestone-3/main.go`
**When** it is read
**Then** it is a `[]Message` slice appended in the loop, never truncated, never written to disk (AD-6)

**Given** `milestone-3/main.go`
**When** `go build ./milestone-3/` is run
**Then** it compiles independently without touching any other directory (AD-2, AD-3)

### Story 1.5: Reference Binary and Demo Task

As Thomas,
I want a pre-built `agent-demo` binary and a `demo-task.txt` committed to the repo root,
So that I can run the minute-0 hook and minute-57 demystification demos without compiling during the session.

**Acceptance Criteria:**

**Given** Thomas's darwin/arm64 machine with Ollama running and llama3.2 pulled
**When** `./agent-demo` is run with the contents of `demo-task.txt` as the task
**Then** the agent executes, makes at least one tool call, and produces an observable file change on screen within a few loop iterations

**Given** the repo root
**When** `ls -la agent-demo` is run
**Then** the binary exists, is executable, and `file agent-demo` reports `Mach-O 64-bit executable arm64` (AD-8)

**Given** `demo-task.txt`
**When** it is read
**Then** it contains a short, clear file-editing task (e.g. "add a function that reverses a string to a file named demo-output.go") that produces a visible file change within a few agent loop iterations

**Given** the `agent-demo` binary
**When** `go version -m agent-demo` is run
**Then** the module path is `github.com/tclaudel/agent-workshop` and the build path references `milestone-3/main.go` — confirming the binary was compiled from the correct source (AD-8, FR-3)

---

## Epic 2: Workshop Documentation

Any facilitator can run a reproducible session without prior experience; attendees can self-rescue during setup and Part 2 Tool JSON schema. Creates README.md (prereq guide, smoke test, Tool JSON schema snippet) and FACILITATOR.md (7-phase talk-track, timing, ring-structure framing, pacing risks, Phase 6 framework call-outs).

### Story 2.1: README Prereq Guide and Walkthrough

As an attendee,
I want clear installation instructions, a setup verification test, and a Tool JSON schema snippet in the README,
So that I can confirm my environment is ready before the session and self-rescue during Part 2 without waiting for the facilitator.

**Acceptance Criteria:**

**Given** a fresh machine
**When** the README prereq section is followed step by step
**Then** Go and Ollama are installed and both smoke tests succeed: `go build` on the starter file produces no error, and `ollama run llama3.2 "hello"` returns a response

**Given** `README.md`
**When** the install section is read
**Then** it contains OS-specific commands or direct links to canonical install pages for Go and Ollama
**And** it contains `ollama pull llama3.2` as a copy-paste command (FR-5)

**Given** `README.md`
**When** the "verify your setup" section is read
**Then** it specifies exactly two checks: (1) `go build` + run on the starter file and (2) `ollama run llama3.2 "hello"` — and states that both must succeed to be session-ready (FR-5)

**Given** `README.md`
**When** the Tool JSON schema section is read
**Then** it contains the complete, copy-pasteable JSON Schema object for the `Parameters` field of a Tool — present in README, not only in FACILITATOR.md (FR-5)

**Given** `README.md`
**When** the walkthrough section is read
**Then** it describes the three-milestone progression so attendees can follow along or catch up independently

### Story 2.2: FACILITATOR.md Session Guide

As Thomas (or any future facilitator),
I want a complete facilitator guide with talk-track anchors, per-phase timing, pacing risks, and copy-paste snippets,
So that I can run a reproducible session without relying on memory and a second facilitator can do the same with no prior experience.

**Acceptance Criteria:**

**Given** `FACILITATOR.md`
**When** it is read
**Then** it contains a section for each of the 7 phases with: timing target (e.g. "0–5 min"), milestone checkpoint where applicable, and the key verbal cue for that phase (FR-4, FR-6)

**Given** the Phase 1 talk-track in `FACILITATOR.md`
**When** it is read
**Then** it explicitly states the ring-structure framing: minute-0 and minute-57 are the same demo run twice — same visual, completely different understanding — and names this as the session's narrative spine (FR-4)

**Given** the Phase 6 talk-track in `FACILITATOR.md`
**When** it is read
**Then** it names LangChain, LlamaIndex, AutoGPT, and Claude Code when making the claim "every agent is this loop" (FR-4)

**Given** the pacing risks section in `FACILITATOR.md`
**When** it is read
**Then** it explicitly calls out both risks: Part 2 (Tool JSON schema `Parameters` field trips attendees — point them to the README snippet) and Part 3 (the loop insight needs time to land — do not rush "there is no magic") (FR-4)

**Given** the Phase 7 section in `FACILITATOR.md`
**When** it is read
**Then** it includes the optional deep-ceiling extension questions (stop condition robustness, tool failure handling, memory persistence, planning step) with framing guidance: "these are directions to explore, not code to write today" (FR-6)
**And** it contains a cue for Thomas to mention the per-session Slack group and confirm all attendees have access

**Given** `FACILITATOR.md`
**When** the pre-session checklist section is read
**Then** it exists as a distinct section covering actions Thomas must complete before each session (at minimum: Slack group creation — detailed in Story 3.1)

**Given** `FACILITATOR.md`
**When** read by a second facilitator with no prior session experience
**Then** they can identify without outside help: what to say at each phase, when to check for milestone completion, where pacing risks are, and how to unblock the two most common attendee blockers (NFR-4)

**Given** `FACILITATOR.md`
**When** copy-paste snippets section is read
**Then** it contains at minimum: the Tool JSON schema `Parameters` snippet and the agent loop `for` skeleton that attendees can paste if they fall behind in Part 3

**Given** the Phase 6 section in `FACILITATOR.md`
**When** it is read
**Then** it includes a spot-check script: three specific questions Thomas asks 3–5 attendees at minute-57 to verify they can name all three components (HTTP call, tool dispatch, loop break) without prompting — and a pass threshold of ≥3/5 correct (NFR-1)

**Given** the Phase 7 section in `FACILITATOR.md`
**When** it is read
**Then** it includes a build-completion checkpoint: Thomas scans the room at minute-60 and notes how many attendees successfully ran their binary — target ≥70% (NFR-2)

---

## Epic 3: Post-Workshop Community

Attendees have a Slack channel for pre-session setup help and post-session follow-up; Thomas has a documented process for creating it per session.

### Story 3.1: Pre-Session Slack Checklist in FACILITATOR.md

As Thomas,
I want a pre-session checklist in FACILITATOR.md that includes creating a Slack group conversation with registered participants,
So that I have a repeatable process before each session and attendees have a channel for setup questions and follow-up.

**Acceptance Criteria:**

**Given** `FACILITATOR.md`
**When** the pre-session checklist section is read
**Then** it lists: create a Slack group conversation with all registered participants, send the invite before the session date, and confirm all attendees can access it (FR-8)

**Given** the Slack group process in `FACILITATOR.md`
**When** it is read
**Then** it specifies a direct group conversation scoped per cohort — not a standing public channel, not reused across sessions (FR-8)

**Given** Phase 7 of the live session
**When** Thomas reaches the end of the session
**Then** FACILITATOR.md instructs him to mention the Slack group and confirm attendees have access — cross-reference to Story 2.2 Phase 7 AC (FR-8)
