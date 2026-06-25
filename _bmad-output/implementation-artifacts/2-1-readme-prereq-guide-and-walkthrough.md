# Story 2.1: README Prereq Guide and Walkthrough

---
baseline_commit: cae8ec371799251ef2b4ccf0b02b8535794efcf2
---

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an attendee,
I want clear installation instructions, a setup verification test, and a Tool JSON schema snippet in the README,
so that I can confirm my environment is ready before the session and self-rescue during Part 2 without waiting for the facilitator.

## Acceptance Criteria

1. **Given** a fresh machine **When** the README prereq section is followed step by step **Then** Go and Ollama are installed and both smoke tests succeed: `go build` on the starter file produces no error, and `ollama run llama3.2 "hello"` returns a response

2. **Given** `README.md` **When** the install section is read **Then** it contains OS-specific commands or direct links to canonical install pages for Go and Ollama **And** it contains `ollama pull llama3.2` as a copy-paste command (FR-5)

3. **Given** `README.md` **When** the "verify your setup" section is read **Then** it specifies exactly two checks: (1) `go build` + run on the starter file and (2) `ollama run llama3.2 "hello"` — and states that both must succeed to be session-ready (FR-5)

4. **Given** `README.md` **When** the Tool JSON schema section is read **Then** it contains the complete, copy-pasteable JSON Schema object for the `Parameters` field of a Tool — present in README, not only in FACILITATOR.md (FR-5)

5. **Given** `README.md` **When** the walkthrough section is read **Then** it describes the three-milestone progression so attendees can follow along or catch up independently

## Tasks / Subtasks

- [x] Task 1: Create `README.md` at repo root (AC: 1, 2)
  - [x] Add install section with links to canonical install pages for Go (`https://go.dev/doc/install`) and Ollama (`https://ollama.com/download`)
  - [x] Include `ollama pull llama3.2` as a distinct copy-paste command block (FR-5)
  - [x] Mark as OS-agnostic (links to canonical pages handle per-OS steps — no need to duplicate package manager commands)

- [x] Task 2: Add "Verify your setup" section (AC: 1, 3)
  - [x] Include exactly two checks in order:
    1. `go build -o /dev/null ./starter/` → must compile with no errors (no binary produced)
    2. `ollama run llama3.2 "hello"` → must return a response
  - [x] State explicitly that **both** must succeed to be session-ready
  - [x] Do NOT add a third check or an automated script (deferred to v2 per PRD §5.1 FR-5 Out of Scope)

- [x] Task 3: Add Tool JSON schema snippet (AC: 4)
  - [x] Include the canonical JSON Schema object used as the `Parameters` field value — use the `write_file` schema as the representative example (two properties; shows the full pattern):
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
  - [x] This is the value that goes in the `Parameters json.RawMessage` field of the `Tool` struct — make that context explicit in the README
  - [x] Snippet must be copy-pasteable (formatted, not inline JSON) — this is the self-rescue anchor for Part 2 when attendees are stuck on the Parameters field format

- [x] Task 4: Add three-milestone walkthrough section (AC: 5)
  - [x] Describe each milestone's single new concept in one sentence:
    - Milestone 1: HTTP call to Ollama — sends a prompt, prints the reply
    - Milestone 2: Tool struct + tools slice — adds `read_file` / `write_file` and passes tools in the request
    - Milestone 3: Agent loop + dispatch + history — wraps Milestone 2 in a `for` loop with `switch` dispatch and `[]Message` accumulation
  - [x] Frame as "each milestone = previous milestone + one concept block" (AD-1 additive-delta) so attendees understand how to read the progression
  - [x] Do NOT add code blocks per milestone — FACILITATOR.md carries the detailed talk-track; README walkthrough is summary only

- [x] Task 5: Verify README completeness
  - [x] Confirm `README.md` is at repo root (same level as `go.mod`, `agent-demo`, `FACILITATOR.md`)
  - [x] Run `go build -o /dev/null ./starter/` from repo root → exits 0
  - [x] Confirm `ollama pull llama3.2` command appears verbatim

### Review Findings

- [x] [Review][Decision] Windows compatibility — dismissed; workshop is macOS/Linux only (user confirmed 2026-06-25)
- [x] [Review][Patch] Missing Go minimum version — go.mod requires Go 1.22 but README does not state a minimum version; attendees with older installs get cryptic errors [README.md]
- [x] [Review][Patch] "Part 2" naming inconsistency — README.md line 45 uses "Part 2 (Milestone 2)" but the rest of the document uses only "Milestone N" terminology [README.md:45]
- [x] [Review][Patch] No working directory instruction — Verify section uses relative path `./starter/` with no "run from project root" note; attendees who paste the command from their home directory get a go.mod-not-found error [README.md:24-41]
- [x] [Review][Patch] Ollama Linux daemon not mentioned — on Linux, Ollama service must be started manually (`systemctl start ollama`); the README gives no hint and the troubleshooting fallback does not cover this [README.md:34-41]
- [x] [Review][Defer] `starter/` directory not explained [README.md] — deferred, pre-existing; FACILITATOR.md (Story 2.2) is the right place for workshop flow context
- [x] [Review][Defer] No clone/repo-acquisition instruction [README.md] — deferred, pre-existing; standard workshop onboarding step typically handled outside README

## Dev Notes

### File to create

`README.md` at repo root — this file does **not** exist yet (confirmed: no `README.md` in repo root at story creation time).

### Canonical install links (do not hard-code package manager commands)

- Go: `https://go.dev/doc/install` — covers macOS, Linux, Windows
- Ollama: `https://ollama.com/download` — covers macOS, Linux, Windows

Reason: hard-coding per-OS commands (`brew install go`, `apt install golang`, etc.) diverges quickly from canonical upstream instructions. Links are stable and self-maintaining.

### Exact smoke test commands

These are the **verbatim** commands to include in the verify section (AC #3):

```bash
# Check 1: Go — build the starter file (no binary output, just compilation check)
go build ./starter/

# Check 2: Ollama — run the model with a simple prompt
ollama run llama3.2 "hello"
```

`go build ./starter/` satisfies AC #1 ("go build on the starter file") and matches the exact command from Story 1.1's AC. It uses `./starter/` directory syntax consistent with every other milestone build command in the repo.

### Tool JSON schema snippet — exact value from milestone-2/main.go

The canonical `Parameters` schemas committed to the repo are (from `milestone-2/main.go:62,75`):

`read_file`:
```json
{"type":"object","properties":{"path":{"type":"string"}},"required":["path"]}
```

`write_file`:
```json
{"type":"object","properties":{"path":{"type":"string"},"content":{"type":"string"}},"required":["path","content"]}
```

Use **`write_file`** as the README example — it has both a `path` and a `content` property, making it the better teaching example of the schema pattern (shows multi-property, required array). The README snippet should be pretty-printed (not compact inline JSON) for readability and copy-paste ease.

Context line to include in README: _"This is the value you pass to the `Parameters json.RawMessage` field of the `Tool` struct."_

### Three-milestone progression (exact content for walkthrough section)

From AD-1 (additive-only delta) and the epics:

| Milestone | Single new concept | What it proves |
|---|---|---|
| `milestone-1/` | HTTP call to Ollama | The model is just an HTTP endpoint |
| `milestone-2/` | Tool struct + tools slice | Tools are Go functions wrapped in a struct |
| `milestone-3/` | Agent loop + dispatch + history | The "agent" is a `for` loop with a `switch` |

Each milestone file is `milestone-N/main.go` — a standalone Go program that compiles independently (`go build ./milestone-N/`).

### Scope boundaries — DO NOT do in this story

- **Do NOT create `FACILITATOR.md`** — that is Story 2.2's deliverable
- **Do NOT add an automated prereq script** — explicitly deferred to v2 (PRD §5.1 FR-5 Out of Scope)
- **Do NOT add per-milestone code blocks** in the walkthrough — keep it to one-sentence descriptions; the milestone files themselves are the reference
- **Do NOT mention the `agent-demo` binary in the prereq/verify section** — the binary is a facilitator artifact (FR-3), not an attendee prereq; Thomas runs it, attendees observe it
- **Do NOT link to slide deck** — out of scope per PRD §6
- **Do NOT add Docker fallback** — explicitly out of scope (PRD §5.1 FR-5 Out of Scope)
- **Do NOT touch any milestone files** — Epic 1 is done; all milestone code is frozen

### Architecture invariants binding this story

This story creates documentation only. No milestone code is touched. The following apply as informational constraints:

| Invariant | Relevance |
|---|---|
| AD-2 / AD-3 | Each milestone is independent; build command is `go build ./milestone-N/` — verify section uses this exact form |
| AD-4 / AD-9 | Stdlib only, no external deps — README should not imply any `go get` or `go mod tidy` steps are needed |
| FR-5 | README must contain: OS-specific/links install, `ollama pull llama3.2`, two-step verify, Tool JSON schema snippet |

### Previous story intelligence (Epic 1 summary)

All of Epic 1 is done. Key facts affecting README content:

- `starter/main.go` exists and compiles with `go build ./starter/` — use this exact command in verify section
- `milestone-1/`, `milestone-2/`, `milestone-3/` all exist and compile independently — build command pattern is `go build ./milestone-N/`
- `agent-demo` binary is compiled from `milestone-3/main.go` for `GOOS=darwin GOARCH=arm64` — facilitator-only artifact; not in attendee prereq path
- `demo.txt` exists at repo root (binary reads it) — no attendee interaction needed
- `go.mod` at root, module `github.com/tclaudel/agent-workshop`, Go 1.22+, zero external deps
- Tool struct Parameters schemas are in `milestone-2/main.go:62,75` — use `write_file` schema as README example

### Project structure notes

```
agent-workshop/
  README.md     ← CREATE (this story)
  go.mod        ← unchanged
  agent-demo    ← unchanged (binary)
  demo-task.txt ← unchanged
  demo.txt      ← unchanged
  FACILITATOR.md ← does not exist yet (Story 2.2)
  starter/main.go ← unchanged
  milestone-1/main.go ← unchanged
  milestone-2/main.go ← unchanged
  milestone-3/main.go ← unchanged
```

### References

- Story 2.1 ACs + FR-5: [Source: _bmad-output/planning-artifacts/epics.md#Story-2.1, #FR-5]
- Verify commands (go build ./starter/): [Source: epics.md#Story-1.1 AC]
- Tool struct Parameters schemas: [Source: milestone-2/main.go:62,75]
- Additive-delta milestone structure: [Source: ARCHITECTURE-SPINE.md#AD-1]
- Standalone build commands: [Source: ARCHITECTURE-SPINE.md#AD-2]
- No external deps: [Source: ARCHITECTURE-SPINE.md#AD-9]
- Automated script deferred: [Source: _bmad-output/planning-artifacts/prds/prd-agent_workshop-2026-06-24/prd.md#FR-5 Out of Scope]
- Docker fallback out of scope: [Source: prd.md#FR-5 Out of Scope]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6[1m]

### Debug Log References

- `go build ./starter/` fails at repo root: `starter/` directory exists, Go cannot write binary named `starter` over an existing directory. Fixed by using `go build -o /dev/null ./starter/` — compiles with no output, exits 0. Story spec assumed no binary would be produced but `starter/main.go` has `package main` + `func main()`, so Go always attempts binary output without `-o`.

### Completion Notes List

- Created `README.md` at repo root with: install links (go.dev, ollama.com), `ollama pull llama3.2` command, two-step verify section (Go + Ollama), Tool JSON schema snippet (write_file, pretty-printed), three-milestone walkthrough table.
- Adapted verify command from `go build ./starter/` to `go build -o /dev/null ./starter/` — same compilation check, avoids directory naming conflict. All ACs satisfied.
- No milestone files touched; documentation only.

### File List

- README.md

## Change Log

- 2026-06-25: Created README.md — install section, verify section, Tool JSON schema snippet, three-milestone walkthrough (Story 2.1)
