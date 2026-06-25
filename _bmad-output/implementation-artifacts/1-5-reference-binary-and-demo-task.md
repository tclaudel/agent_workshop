---
baseline_commit: 1-4-milestone-3-agent-loop-dispatch-and-history
source_binary: milestone-3/main.go
build_target: GOOS=darwin GOARCH=arm64
---

# Story 1.5: Reference Binary and Demo Task

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As Thomas,
I want a pre-built `agent-demo` binary and a `demo-task.txt` committed to the repo root,
so that I can run the minute-0 hook and minute-57 demystification demos without compiling during the session.

## Acceptance Criteria

1. **Given** Thomas's darwin/arm64 machine with Ollama running and llama3.2 pulled **When** `./agent-demo` is run **Then** the agent executes, makes at least one tool call, and produces an observable file change on screen (`reversed.txt` appears) within a few loop iterations

2. **Given** the repo root **When** `ls -la agent-demo` is run **Then** the binary exists, is executable, and `file agent-demo` reports `Mach-O 64-bit executable arm64` (AD-8)

3. **Given** `demo-task.txt` **When** it is read **Then** it contains a short, clear file-editing task that mirrors the binary's behaviour and produces a visible file change within a few agent loop iterations — see **Decision #1** in Dev Notes for why this MUST match `milestone-3/main.go`'s hardcoded prompt, not an arbitrary task

4. **Given** the `agent-demo` binary **When** `go version -m agent-demo` is run **Then** the `mod` line shows module path `github.com/tclaudel/agent-workshop` and the `path` line is `github.com/tclaudel/agent-workshop/milestone-3` — confirming the binary was compiled from `milestone-3/main.go` (AD-8, FR-3)

5. **Given** the repo root **When** `cat demo.txt` is run **Then** a committed `demo.txt` input file exists with short, visible content — required because the binary's hardcoded prompt reads `demo.txt`; without it the tool loop errors and AC #1's observable change is non-deterministic (see **Decision #2**)

## Tasks / Subtasks

- [x] Task 1: Seed the `demo.txt` input file at repo root (AC: 5) — **see Decision #2**
  - [x] Create `demo.txt` with short, visibly-reversible text, e.g. `hello agent workshop` (single line, no trailing logic)
  - [x] This file is READ by the binary's hardcoded prompt; it must exist and be committed so `./agent-demo` works out-of-the-box at session time

- [x] Task 2: Write `demo-task.txt` to mirror the binary's hardcoded prompt (AC: 3) — **see Decision #1**
  - [x] Read the exact prompt literal from `milestone-3/main.go:99`
  - [x] Set `demo-task.txt` content to that same task text (verbatim): `Read the file demo.txt, then write its contents reversed to reversed.txt`
  - [x] Do NOT use the AC-3 epic example ("add a function that reverses a string...") — it describes a different task than the binary performs and would break demo coherence (NFR-1)

- [x] Task 3: Compile the reference binary for darwin/arm64 (AC: 2, 4) — **DO NOT edit milestone-3/main.go**
  - [x] From repo root: `GOOS=darwin GOARCH=arm64 go build -o agent-demo ./milestone-3/`
  - [x] Confirm no source change to `milestone-3/main.go` (it is frozen — story 1.4 is done; AD-8 binds the binary to it as-is)

- [x] Task 4: Verify the binary metadata (AC: 2, 4)
  - [x] `file agent-demo` → must report `Mach-O 64-bit executable arm64`
  - [x] `ls -la agent-demo` → file exists and has the executable bit (go build sets it; no `chmod` needed)
  - [x] `go version -m agent-demo` → `mod` line `github.com/tclaudel/agent-workshop`, `path` line `github.com/tclaudel/agent-workshop/milestone-3`

- [x] Task 5: Live demo verification (AC: 1)
  - [x] Ensure Ollama is running and `llama3.2` is pulled (`ollama list`)
  - [x] Ensure `demo.txt` (Task 1) is present at repo root
  - [x] Run `./agent-demo` → the loop dispatches `read_file` then `write_file`, `reversed.txt` is created at repo root, and a final text reply prints
  - [x] Confirm `reversed.txt` is the OUTPUT and is created live (observable). Do NOT commit `reversed.txt` — see Decision #3. Remove it after the verification run so the committed repo state has `demo.txt` + binary but no pre-existing `reversed.txt`.

- [x] Task 6: Confirm commit hygiene (AC: 2, 5)
  - [x] `agent-demo`, `demo-task.txt`, `demo.txt` are NOT matched by `.gitignore` (current `.gitignore` only ignores `_bmad/` and `bmad-*` skill dirs — all three deliverables are committable)
  - [x] Add `reversed.txt` to `.gitignore` (generated demo output, must not be committed)

## Dev Notes

### 🚨 Decision #1 — `demo-task.txt` MUST mirror the binary's hardcoded prompt (demo coherence)

The binary is compiled from `milestone-3/main.go`, whose prompt is a **hardcoded string literal** (`milestone-3/main.go:99`):

```go
prompt := "Read the file demo.txt, then write its contents reversed to reversed.txt"
```

`milestone-3/main.go` has **no CLI argument parsing, no stdin reading, no flag package** (story 1.2 design decision, preserved additively through M3 per AD-1). Therefore the binary does **not** read `demo-task.txt` at runtime. AC #1's "run with the contents of `demo-task.txt` as the task" cannot mean piping the file in — it means: **`demo-task.txt` is the human-readable narration of what the binary already does.** For the demo to be coherent (Thomas reads the task aloud → runs binary → the named files change), `demo-task.txt`'s text MUST match the hardcoded prompt.

The epic AC #3 parenthetical ("e.g. add a function that reverses a string to a file named `demo-output.go`") is **illustrative only** and describes a *different* task than the binary performs. Using it would mean Thomas narrates one task while the binary does another — directly damaging the minute-57 comprehension check (NFR-1: attendees identify all three components). **Override the example; mirror the real prompt verbatim.**

> Do NOT "fix" this by adding arg/stdin parsing to `milestone-3/main.go` to make it read `demo-task.txt`. That would (a) edit a done, code-reviewed milestone, (b) violate AD-1 additive-delta and the story-1.2 no-CLI decision, and (c) break AC #4 (`go version -m` must still reference `milestone-3/main.go`). The binary is frozen to M3 source.

### 🚨 Decision #2 — `demo.txt` input file is a required deliverable of THIS story

The hardcoded prompt **reads `demo.txt`** before writing `reversed.txt`. If `demo.txt` does not exist at repo root when `./agent-demo` runs, `read_file` returns `error: open demo.txt: no such file or directory`, the model's subsequent behaviour is non-deterministic, and AC #1's "observable file change" is not reliably produced.

Story 1.4's Dev Agent Record explicitly flagged this:
> "Throwaway `demo.txt` ... seeded at repo root ... loop dispatched `read_file` then `write_file`, created `reversed.txt` ... Throwaway `demo.txt`/`reversed.txt` removed after the run (**they are story 1.5 deliverables, not this story's**)."

So `demo.txt` is owned here. It is not listed verbatim in the epic AC table, but it is a hard prerequisite for AC #1 — added as AC #5. Commit it with short, clearly-reversible content (e.g. `hello agent workshop`).

### 🚨 Decision #3 — `reversed.txt` is generated output, NOT committed

`reversed.txt` is what the binary WRITES — it is the observable payoff of the demo. It must be **created live** during the session (and during Task-5 verification) to be observable. Do not commit it: a pre-existing `reversed.txt` removes the "a new file just appeared" moment. Add it to `.gitignore`. After Task-5 verification, delete the generated `reversed.txt` so the committed repo ships `demo.txt` + `agent-demo` + `demo-task.txt` only.

> Reversal-accuracy caveat (facilitator note, not a bug): the byte-reversal is performed by llama3.2, not by deterministic code. The reversed content may be imperfect. The demo's payoff is the **observable file write + the tool-dispatch loop**, NOT character-perfect reversal. Do not treat imperfect reversal as a failure.

### Build command and binary-name collision

From repo root:
```bash
GOOS=darwin GOARCH=arm64 go build -o agent-demo ./milestone-3/
```
- `GOOS`/`GOARCH` are explicit per AD-8. Host is already darwin/arm64 (`go1.24.1 darwin/arm64`), so they are redundant but required by the invariant — keep them for reproducibility/documentation.
- The output name `agent-demo` does **not** collide with any directory (there is no `agent-demo/` dir), so `-o agent-demo` is safe. (Contrast the milestone stories, which used `-o /dev/null` to avoid a binary-vs-directory name clash. That issue does not apply here.)
- `go build` produces an executable with the executable bit set — no `chmod +x` needed.

### What `go version -m agent-demo` reports (AC #4)

```
agent-demo: go1.24.1
        path    github.com/tclaudel/agent-workshop/milestone-3
        mod     github.com/tclaudel/agent-workshop      (devel)
        ...
```
- `path` = the main-package import path = `…/milestone-3` → satisfies "references milestone-3" (AD-8, FR-3).
- `mod` = module path = `github.com/tclaudel/agent-workshop` → satisfies AC #4 module-path check.
- Do NOT use `-trimpath` — it strips local path info but the embedded `path`/`mod` build metadata (which AC #4 inspects) is preserved regardless; plain `go build` is correct and simplest. No reason to add build flags.

### Architecture Invariants Binding This Story

| Invariant | Rule (as it applies here) |
|---|---|
| AD-8 | Binary compiled from `milestone-3/main.go` for `GOOS=darwin GOARCH=arm64`; committed to repo root as `agent-demo`; no multi-arch |
| AD-9 | One `go.mod` at root, module `github.com/tclaudel/agent-workshop`, stdlib only — unchanged by this story |
| FR-3 | Pre-built reference binary at repo root; runs on Thomas's darwin/arm64 machine; observable file change for minute-0 hook / minute-57 payoff |
| AD-1/2/3 | `milestone-3/main.go` is NOT touched — binary is built from it as-is |

### Scope Boundary — DO NOT do in this story

- **Do NOT edit `milestone-3/main.go`** (or any milestone). It is done and code-reviewed. The binary is frozen to it (AD-8). Changing the prompt, adding arg/stdin parsing, or "improving" the demo task all violate this.
- **Do NOT change the M3 prompt to the AC-3 example task.** Mirror the actual hardcoded prompt in `demo-task.txt` instead (Decision #1).
- **Do NOT commit `reversed.txt`** (Decision #3).
- **Do NOT add a Makefile/build script, multi-arch builds, or a build target to the justfile** — the justfile is BMad-only; a single documented `go build` line is the deliverable mechanism. Multi-arch is explicitly deferred (ARCHITECTURE-SPINE.md#Deferred).
- **Do NOT address deferred items** (path traversal in `read_file`/`write_file`, unbounded `os.ReadFile`, no HTTP timeout, non-2xx status check). They live in `deferred-work.md` and are out of scope.
- This is the **last story in Epic 1**. After it, `epic-1-retrospective` is optional.

### Environment (verified at story-creation time)

- `go version` → `go1.24.1 darwin/arm64` (≥1.22 per AD-9 Stack table) ✓
- `uname -m` → `arm64` ✓ (matches AD-8 target)
- `ollama list` → `llama3.2:latest` present ✓ (required for Task-5 live run)
- Repo root currently has **no** `agent-demo`, `demo-task.txt`, or `demo.txt` — all three are created by this story.
- `.gitignore` ignores only `_bmad/` and `bmad-*` skill dirs — none of this story's deliverables are accidentally ignored.

### Project Structure Notes

```
agent-workshop/
  go.mod                 ← unchanged (AD-9)
  agent-demo             ← CREATE (this story) — darwin/arm64 binary from milestone-3/main.go
  demo-task.txt          ← CREATE (this story) — narration; mirrors M3 hardcoded prompt
  demo.txt               ← CREATE (this story) — input file the binary reads (Decision #2)
  reversed.txt           ← GENERATED at runtime — DO NOT commit; add to .gitignore (Decision #3)
  .gitignore             ← UPDATE — add reversed.txt
  milestone-3/main.go    ← UNCHANGED (frozen; binary source — AD-8)
  milestone-1/main.go    ← unchanged
  milestone-2/main.go    ← unchanged
  starter/main.go        ← unchanged
```

### Previous Story Intelligence (1.4)

- M3 is done; its hardcoded prompt (`milestone-3/main.go:99`) is the single source of truth for what the binary does — this story's `demo-task.txt` and `demo.txt` are built around it.
- 1.4's live run confirmed the loop dispatches `read_file` → `write_file` and converges with `demo.txt` seeded → strong evidence Task-5 will succeed with Ollama + llama3.2.
- 1.4 flagged a latent `Arguments string` bug in `milestone-2/main.go:31` (M2/M3 not struct-identical). **Not this story's concern** — it is a separate correct-course item; do not touch M2.

### References

- Reference binary source/target/name: [Source: _bmad-output/planning-artifacts/architecture/architecture-agent_workshop-2026-06-24/ARCHITECTURE-SPINE.md#AD-8]
- Single go.mod / module path: [Source: ARCHITECTURE-SPINE.md#AD-9]
- Structural seed (demo-task.txt + agent-demo at root): [Source: ARCHITECTURE-SPINE.md#Structural Seed]
- Multi-arch deferred: [Source: ARCHITECTURE-SPINE.md#Deferred]
- Story 1.5 ACs + FR-3: [Source: _bmad-output/planning-artifacts/epics.md#Story-1.5, #FR-3]
- Hardcoded prompt + no-CLI decision: [Source: milestone-3/main.go:99, _bmad-output/implementation-artifacts/1-2-milestone-1-http-call-to-ollama.md]
- demo.txt ownership handoff from 1.4: [Source: _bmad-output/implementation-artifacts/1-4-milestone-3-agent-loop-dispatch-and-history.md#Debug-Log-References]
- Deferred items (out of scope): [Source: _bmad-output/implementation-artifacts/deferred-work.md]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6[1m]

### Implementation Plan

1. Seed `demo.txt` with `hello agent workshop` (single line, visibly reversible)
2. Write `demo-task.txt` verbatim from `milestone-3/main.go:99` prompt literal
3. Compile: `GOOS=darwin GOARCH=arm64 go build -o agent-demo ./milestone-3/`
4. Verify binary metadata (file type, executable bit, go version -m)
5. Live run: `./agent-demo` → confirmed read_file + write_file dispatched, `reversed.txt` created with `workshop agent hello`
6. Remove `reversed.txt`; add to `.gitignore`; confirm deliverables not ignored

### Debug Log References

- milestone-3/main.go:99: prompt literal confirmed as `"Read the file demo.txt, then write its contents reversed to reversed.txt"`
- Live run output: `The contents of demo.txt were written to reversed.txt as follows: workshop agent hello` — AC#1 satisfied
- `file agent-demo` → `Mach-O 64-bit executable arm64` ✓
- `go version -m agent-demo` → path `github.com/tclaudel/agent-workshop/milestone-3`, mod `github.com/tclaudel/agent-workshop (devel)` ✓
- `git check-ignore` confirms: `agent-demo`, `demo-task.txt`, `demo.txt` not ignored; `reversed.txt` ignored ✓

### Completion Notes List

- ✅ `demo.txt` created at repo root with `hello agent workshop`
- ✅ `demo-task.txt` created verbatim from M3 hardcoded prompt (mirrors binary behaviour per Decision #1)
- ✅ `agent-demo` binary compiled from `milestone-3/main.go` for darwin/arm64 — no source changes to milestone-3
- ✅ All AC#2/4 binary metadata checks pass
- ✅ AC#1 live run confirmed: loop dispatched read_file → write_file, `reversed.txt` created and removed post-verification
- ✅ `reversed.txt` added to `.gitignore` (Decision #3)
- ✅ No regressions — no existing code touched; milestone files untouched

### File List

- `demo.txt` — created (repo root; input file the binary reads)
- `demo-task.txt` — created (repo root; human-readable narration, mirrors M3 hardcoded prompt)
- `agent-demo` — created (repo root; darwin/arm64 binary compiled from milestone-3/main.go)
- `.gitignore` — updated (added `reversed.txt` entry)

## Change Log

- 2026-06-25: Story 1.5 created — reference `agent-demo` binary (darwin/arm64 from milestone-3/main.go), `demo-task.txt` (mirrors M3 hardcoded prompt — Decision #1), and `demo.txt` input file (required prerequisite for the observable demo, added as AC #5 — Decision #2). milestone-3/main.go frozen; reversed.txt is generated output and not committed (Decision #3). Comprehensive context engine analysis completed - comprehensive developer guide created.
- 2026-06-25: Story 1.5 implemented — all 6 tasks complete. `demo.txt`, `demo-task.txt`, `agent-demo` created; `.gitignore` updated with `reversed.txt`. Live run confirmed AC#1 (read_file→write_file dispatched, reversed.txt created). All binary metadata ACs pass. Status → review.

### Review Findings

- [x] [Review][Patch] Deliverables untracked — all 4 files exist on disk but no git commit exists; run `git add agent-demo demo.txt demo-task.txt .gitignore && git commit` [repo root]
- [x] [Review][Patch] .gitignore uses bare `reversed.txt` — does not cover subdirectory variants (e.g. if LLM writes to a non-root path); use `**/reversed.txt` [.gitignore:8]
- [x] [Review][Defer] CWD-relative paths in binary — `demo.txt`/`reversed.txt` resolved relative to CWD; binary must be invoked from project root or `read_file` errors silently [milestone-3/main.go:99] — deferred, pre-existing frozen code
- [x] [Review][Defer] demo.txt trailing newline causes reversed output to start with `\n` — `os.ReadFile` preserves the `\n`, LLM reversal is non-deterministic anyway [milestone-3/main.go:65] — deferred, pre-existing frozen code
- [x] [Review][Defer] `reversed.txt` silently overwritten on repeat runs — no existence check or warning [milestone-3/main.go:79] — deferred, pre-existing frozen code (see also 1-3 deferred)
- [x] [Review][Defer] `defer resp.Body.Close()` inside for-loop — defers accumulate until function return, leaking file descriptors across iterations [milestone-3/main.go:118] — deferred, pre-existing frozen code
- [x] [Review][Defer] Only `ToolCalls[0]` dispatched — extra tool calls in a multi-call LLM response are silently dropped [milestone-3/main.go:139-143] — deferred, pre-existing frozen code
- [x] [Review][Defer] `_bmad-output/` directory not covered by .gitignore — planning artifacts are untracked; .gitignore new in this story but coverage is a project-wide decision [.gitignore] — deferred, out of scope for story 1.5
