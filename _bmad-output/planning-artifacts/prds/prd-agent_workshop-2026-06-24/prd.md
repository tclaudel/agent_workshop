---
title: agent_workshop — Build an AI Agent from Scratch (Go + Ollama)
status: final
created: 2026-06-24
updated: 2026-06-24
---

# PRD: agent_workshop

## 0. Document Purpose

This PRD is for Thomas (principal engineer, agentic platform & data) and any co-facilitators or stakeholders involved in planning or running the workshop. It defines scope, requirements, and success criteria for a 1-hour internal engineering workshop that teaches AI agent fundamentals from scratch using Go and Ollama. The primary artifact is a GitHub repository; the live session is facilitated by Thomas using slides (out of scope for this PRD). Downstream work: slide deck design, scheduling, attendee communications.

Inputs used: brainstorming session output (`brainstorm-agent-from-scratch-ollama-golang-2026-06-23/workshop-outline.md`).

---

## 1. Vision

Engineers at Bedrock Streaming interact daily with AI agents — as users of tools like Claude Code, or as builders evaluating LangChain, LlamaIndex, and similar frameworks. Most cannot explain what an agent actually is under the hood: the loop, the tool dispatch, the stop condition. This knowledge gap is benign until it isn't — a debugging session gone opaque, a framework abstraction that papers over a critical assumption, an architecture decision made without basis.

The agent_workshop closes this gap. In 60 minutes, with vanilla Go and a local Ollama instance, every attending engineer builds a working AI agent from scratch — no frameworks, no abstractions, nothing hidden. The session's thesis: demystification is the most durable skill. Once an engineer has implemented the loop themselves, they can read any agent framework's source, debug any agent behavior, and make architecture decisions from first principles rather than cargo-cult confidence.

This workshop is also Thomas's opening statement as principal engineer of the agentic platform & data squads: the team builds on understanding, not abstraction.

---

## 2. Why Now

Thomas joins as principal engineer of the agentic platform & data squads. Running this workshop early — before platform decisions are announced, before a technical direction is set — establishes the team's intellectual posture, surfaces engineers across the company who are interested in agentic work, and seeds the Slack community that will become a lightweight peer network. The signal gathered from the dry run informs whether and how to scale recurring sessions.

---

## 3. Target User

### 3.1 Jobs To Be Done

- Understand how AI agents work internally, not just how to use them
- Build confidence to read and evaluate agent framework source code
- Gain a foundation for architecture decisions about agentic systems
- Connect with other engineers in the company interested in agentic work

### 3.2 Non-Users (v1)

- Non-engineers (PMs, analysts, designers) — session requires writing Go
- Engineers with no Go experience — [ASSUMPTION: basic Go fluency is a prereq; defining "basic" as: can write, build, and run a simple Go program]
- Engineers seeking production-ready agent patterns — this is educational, not production guidance

### 3.3 Key User Journeys

**UJ-1. Marco attends and can finally explain what he just ran.**
Marco, a backend engineer who has used LangChain in a side project but never read the source, joins with his laptop set up (Go + Ollama running, llama3.2 pulled). He watches Thomas run the pre-built binary at minute 0 — a file changes on screen. Over the next 50 minutes he writes three progressive Go programs: an HTTP call to Ollama, a Tool struct with read/write functions, then the agent loop. At minute 57, Thomas runs the binary again. Marco can name every component: the HTTP call, the tool dispatch JSON, the loop breaking on text reply. He compiles his own binary at minute 59. He joins the Slack channel on the way out.
**Edge case:** Marco's Ollama is not running when the session starts. The prereq guide and a 5-minute buffer before minute 0 covers this.

**UJ-2. Thomas runs the session again for a second squad, zero repo changes.**
Three weeks after the dry run, Thomas schedules a second session for a different squad. He opens the same repository, loads the same slides, and runs the same binary demo. Nothing needed updating. The facilitator guide gives him talk-track checkpoints and pacing notes so the session is reproducible without relying on memory.

---

## 4. Glossary

- **Agent** — A Go program that runs an LLM in a loop, dispatches tools based on model output, and stops when the model returns a text reply instead of a tool call.
- **Tool** — A Go function exposed to the model via a struct (`Name`, `Description`, `Parameters`, `Run`). The model requests tools; the Go program executes them.
- **Tool call** — A JSON blob returned by the model requesting execution of a named Tool with specific arguments. The Go program, not the model, decides whether to execute it.
- **Agent loop** — The `for` loop in Go that alternates between calling the model and executing tool calls until a stop condition is met.
- **Stop condition** — The condition under which the agent loop exits: the model returns a plain text reply rather than a tool call.
- **Milestone** — One of three progressive code states in the repo corresponding to workshop Parts 1, 2, and 3.
- **Reference binary** — A pre-compiled Go binary included in the repo, used for the minute-0 hook demo and the minute-57 payoff.
- **Dry run** — The first closed session run with a selected group of engineers before public rollout.

---

## 5. Features

### 5.1 Workshop Repository

The repo is the primary product artifact. It lives on Thomas's personal GitHub. Attendees clone it before the session and use it throughout.

**Description:** The repo provides three things: (1) progressive milestone code mirroring the workshop structure, (2) a self-sufficient walkthrough guide in the README, and (3) facilitator materials so any future session is reproducible without memory or slides. It is the single source of truth for attendees and facilitator alike.

**Functional Requirements:**

#### FR-1: Progressive milestone code
Repo contains three milestone directories (or clearly-named files) corresponding to workshop Parts 1, 2, and 3. Each milestone is a complete, runnable Go program.

**Consequences (testable):**
- `milestone-1/` contains a working Go program that sends a prompt to Ollama via HTTP and prints the reply.
- `milestone-2/` extends milestone-1 with a `Tool` struct, `read_file` and `write_file` implementations, and the tools slice passed in the Ollama request.
- `milestone-3/` extends milestone-2 with the agent loop, tool dispatch switch, and message history accumulation.

#### FR-2: Starter file
Repo contains a `starter/main.go` with only the package declaration and required imports. Attendees build from this during the session.

**Consequences (testable):**
- `starter/main.go` compiles with `go build` (no logic yet, but valid Go).
- File contains no implementation code beyond the package line and imports.

#### FR-3: Reference binary
Repo contains a pre-built reference binary for the minute-0 hook demo and the minute-57 demystification payoff. Thomas builds this binary for his own machine; no multi-architecture distribution required.

**Consequences (testable):**
- Reference binary is committed to the repo and runs on Thomas's facilitation machine.
- Running the binary with a task like "add a function that reverses a string to main.go" causes a file change observable on screen.

#### FR-4: Facilitator guide
Repo contains a `FACILITATOR.md` with talk-track anchors, per-phase timing, pacing risks, and copy-paste snippets. A second facilitator with no prior session experience must be able to run the session from this file alone.

**Consequences (testable):**
- `FACILITATOR.md` has a section per phase with timing target and key verbal cue.
- Phase 1 talk-track includes the ring-structure framing: minute-0 and minute-57 are the same demo, run twice — same visual, completely different understanding. This is the session's narrative spine.
- Phase 6 talk-track names specific frameworks (LangChain, LlamaIndex, AutoGPT, Claude Code) when making the claim "every agent is this loop."
- Pacing risks call out both: Part 2 (Tool JSON schema parameters trips attendees) and Part 3 (the loop insight needs to breathe — do not rush the "there is no magic" moment).

#### FR-5: Prereq setup guide
Repo README includes exact commands to install Go, install Ollama, and pull llama3.2. It also includes the Tool JSON schema copy-paste snippet so attendees can self-rescue during Part 2 without waiting for the facilitator.

**Consequences (testable):**
- Commands are OS-specific or link to canonical install pages for Go and Ollama.
- A copy-pasteable `ollama pull llama3.2` command is present.
- Prereq section includes a "verify your setup" step with a Go build+run smoke test and an `ollama run llama3.2 "hello"` check — attendees who cannot complete both are not ready for the session.
- Tool JSON schema snippet for the `Parameters` field is present in the README (not only in FACILITATOR.md).

**Out of Scope:**
- Automated prereq verification script (v2)
- Docker-based fallback environment

### 5.2 Live Workshop Session

The 60-minute facilitated session using the repo and slides.

**Description:** Thomas facilitates live, following the 7-phase structure. Attendees write code at their own machines. The session has one success criterion: attendee compiles their own agent by minute 60. The repo and facilitator guide enforce the structure.

**Functional Requirements:**

#### FR-6: Session phase structure
Session follows the defined 7-phase structure in order.

**Consequences (testable):**
- Phase 1 (0–5 min): Reference binary demo runs live.
- Phase 2 (5–10 min): "Agent = LLM + loop + tools + stop condition" is written on the board and stays visible.
- Phase 3 (10–20 min): Attendees reach Milestone 1 checkpoint (HTTP call to Ollama works).
- Phase 4 (20–35 min): Attendees reach Milestone 2 checkpoint (model returns tool-call JSON).
- Phase 5 (35–50 min): Attendees reach Milestone 3 checkpoint (two-step chained tool call completes).
- Phase 6 (50–57 min): Reference binary re-run, walked line by line.
- Phase 7 (57–60 min): Attendees run `go build -o my-agent .` and execute their binary with a live task. Facilitator may offer the deep-ceiling extension questions if time and interest allow (see FR-6 notes).

**Notes:**
- Milestone-to-phase mapping: Milestone 1 = Phase 3 checkpoint; Milestone 2 = Phase 4 checkpoint; Milestone 3 = Phase 5 checkpoint.
- Deep-ceiling extension (optional, Phase 7 tail): facilitator may surface four architecture questions if time remains — stop condition robustness, tool failure handling, memory persistence across sessions, adding a planning step. Framing: "these are directions to explore, not code to write today." Attendees now have the foundation to reason about them.

#### FR-7: Solo facilitation
Session runs with Thomas as sole facilitator, no additional infra beyond laptop + projector. Realizes UJ-2.

**Consequences (testable):**
- No external services required beyond Ollama running locally on Thomas's machine.
- All demo artifacts (binary, demo task file) are committed to the repo.

### 5.3 Post-Workshop Community

**Description:** Thomas creates a per-session Slack group conversation with participants before each session. The group serves two purposes: pre-session (attendees can ask prereq and setup questions) and post-session (continued discussion, experiments, follow-up). No public channel — direct group conversation scoped per cohort.

**Functional Requirements:**

#### FR-8: Per-session Slack conversation
Thomas creates a Slack group conversation with all registered participants before each session.

**Consequences (testable):**
- Group exists before the session so attendees can ask setup/prereq questions in advance.
- Thomas mentions the group during Phase 7 to confirm attendees have access for post-session follow-up.

---

## 6. Non-Goals (Explicit)

- **No async/self-paced format** — no video recording, no LMS, no written-only course. The live coding dynamic is the product.
- **No internal company hosting** — repo lives on Thomas's personal GitHub, not corporate GitLab or internal infra.
- **No non-Go implementations** — Python, TypeScript variants are out of scope for v1.
- **No production patterns** — error handling, retries, persistent memory, multi-agent orchestration are explicitly out of scope and should not be introduced during the session.
- **No multi-model support** — llama3.2 via Ollama only. Cloud API alternatives are not supported.
- **No attendance tracking or LMS integration** — informal headcount only.
- **No slide ownership** — slides are Thomas's responsibility and are not a repo artifact.

---

## 7. MVP Scope

### 7.1 In Scope

- GitHub repo with 3 milestone files, starter file, reference binary (Thomas's machine), README walkthrough, prereq guide, and `FACILITATOR.md` (self-sufficient for a second facilitator)
- 1 dry run with selected engineers (~2026-07-08)
- 1+ public live sessions
- Per-session Slack group conversation with participants

### 7.2 Out of Scope for MVP

- Automated prereq verification script — deferred to v2 if pacing issues surface in dry run
- Multi-architecture binary distribution — Thomas builds for his own machine only
- Slide deck in the repo — Thomas owns separately
- Video recording of a session — [NOTE FOR PM: high-value if Thomas wants async reach; revisit after first public session]
- Workshop variants (90-min deep-dive, 30-min lightning) — deferred until demand is confirmed

---

## 8. Success Metrics

**Primary**

- **SM-1**: Comprehension rate — at the demystification step (minute 57), facilitator spot-checks 3–5 attendees on "what just happened." Target: ≥3 of 5 spot-checked attendees correctly identify all three components (HTTP call, tool dispatch, loop break) without prompting. Validates FR-6.
- **SM-2**: Build completion rate — % of attendees who successfully run `go build` and execute their binary by minute 60. Target: ≥70%. Validates FR-1, FR-2, FR-6.

**Secondary**

- **SM-3**: Post-session engagement — at least one follow-up question or discussion thread per session in the Slack group conversation. Validates FR-8.
- **SM-4**: Repeatability — Thomas (or another facilitator) runs ≥2 sessions from the same repo without structural changes. Validates FR-4, FR-7. Realizes UJ-2.

**Counter-metrics (do not optimize)**

- **SM-C1**: Do not optimize for session speed. Sessions that skip Parts 1–3 in favor of lecture-and-demo sacrifice the build-completion moment that makes comprehension stick. A faster session is a worse session.

---

## 9. Open Questions

1. Who is the second facilitator? Thomas will identify and onboard them after the first public session. `FACILITATOR.md` must be self-sufficient before that point.

---

## 10. Assumptions Index

- **§3.2** — Basic Go fluency defined as: can write, build, and run a simple Go program. If attendees lack this, a prereq module is needed before the session.
- **§5.3 FR-8** — Post-session Slack group is created per-session by Thomas, not a standing public channel.
