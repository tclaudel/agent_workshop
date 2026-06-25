# Deferred Work

## Deferred from: code review of 1-1-repository-foundation (2026-06-24)

- `go 1.22` without `toolchain` directive in `go.mod:3` — Go 1.21+ introduced the `toolchain` directive; omitting it means `go mod tidy` on a newer toolchain may silently rewrite go.mod, causing non-reproducible builds across team members using different Go versions. Not a spec violation for this story; worth addressing in a future go.mod maintenance pass.

## Deferred from: code review of 1-2-milestone-1-http-call-to-ollama (2026-06-24)

- No HTTP client timeout — violates AD-13 to fix (requires custom http.Client)
- ChatResponse struct missing done/done_reason — M1 scope boundary
- No MaxBytes guard / empty-partial-truncated body edge cases — MINOR
- Content-Type mismatch accepted silently — MINOR

## Deferred from: code review of 1-3-milestone-2-tool-struct-and-tools-slice (2026-06-24)

- Path traversal in read_file/write_file — no path sanitization; Tool.Run dead code in M2, will matter in M3 when dispatch is added
- Silent zero-value type assertion on args["path"]/args["content"] — missing-key returns "" silently; dead code in M2
- Unbounded file read (os.ReadFile) can exhaust memory — no size cap; dead code in M2
- No HTTP timeout on Ollama request — pre-existing from M1, same as above
- json.RawMessage parameter schemas never parsed/validated — typo would produce malformed schema; dead code path in M2
- write_file hardcoded 0644 permissions — silently overwrites existing file perms; dead code in M2

## Deferred from: code review of 1-4-milestone-3-agent-loop-dispatch-and-history (2026-06-24)

- Non-2xx HTTP status not checked (milestone-3/main.go:114) — `http.Post` only errors on transport failure; a 4xx/5xx or 200-with-error body unmarshals to a zero-value ChatResponse, hits `len(ToolCalls)==0`, and the loop prints an empty line and exits as if the model finished normally. NEW in M3 (loop depends on response shape). A status check is additive error handling beyond AD's `log.Fatal`-only pattern and borders on AD-7 (no error-based stop) — teaching-scope, deferred.
- (Reconfirmed reachable in M3, already tracked from 1-3 — not duplicated above): silent zero-value type assertion on args["path"]/["content"]; unbounded os.ReadFile; path traversal in read_file/write_file. Tool.Run is now live code (dispatch reaches it), so these are reachable rather than dead-code; they remain deferred per the story's Dev Notes.

## Deferred from: code review of 1-5-reference-binary-and-demo-task (2026-06-25)

- CWD-relative paths in binary — `demo.txt`/`reversed.txt` resolved relative to CWD; `./agent-demo` must be run from project root or `read_file` errors silently and LLM may fabricate content (milestone-3/main.go:99) — frozen code
- demo.txt trailing newline propagates to reversed output — `os.ReadFile` preserves the `\n`; LLM reversal puts it first in `reversed.txt`; non-deterministic anyway per story decision (milestone-3/main.go:65) — frozen code
- `reversed.txt` silently overwritten on repeat runs — no existence check; acceptable for demo but can obscure re-run behavior (milestone-3/main.go:79) — frozen code
- `defer resp.Body.Close()` inside for-loop — defers accumulate until function return, leaking response body file descriptors across loop iterations (milestone-3/main.go:118) — frozen code
- Only `ToolCalls[0]` dispatched — extra tool calls in a multi-call LLM response silently dropped; valid Ollama behavior not handled (milestone-3/main.go:139-143) — frozen code
- `_bmad-output/` directory not in .gitignore — planning artifacts untracked; project-wide decision needed on whether to commit or ignore them

## Deferred from: code review of 2-1-readme-prereq-guide-and-walkthrough (2026-06-25)

- `starter/` directory not explained in README — attendees see `go build -o /dev/null ./starter/` but README never says what `starter/` is (participant's working copy vs. template vs. reference); FACILITATOR.md (Story 2.2) is the appropriate place for workshop flow context
- No clone/repo-acquisition instruction — README assumes a local copy of the repo exists; standard workshop onboarding step typically handled outside README (workshop invite, GitHub link, or facilitator setup)

## Deferred from: code review of 2-2-facilitator-md-session-guide (2026-06-25)

- Slack group creation has no pre-session verification step — first check is Phase 7 verbal cue; remediation impossible mid-session. Story 3.1 scope.
- Build-completion checkpoint at minute 60 = session end — no time to act on scan results within the session window; inherent to the 60-minute format design.

## Deferred from: code review of 3-1-pre-session-slack-checklist-in-facilitator-md (2026-06-26)

Story 3.1 changes clean (all ACs satisfied). All findings below are pre-existing in FACILITATOR.md from Story 2.2.

- Snippet 2 (Copy-Paste): only `ToolCalls[0]` dispatched — multi-tool response silently drops remaining calls (also tracked in 1-5 deferred as frozen code)
- Snippet 2 (Copy-Paste): `output` variable undeclared — snippet won't compile as-is when pasted
- Snippet 2 (Copy-Paste): magic-index dispatch (`tools[0]`/`tools[1]`) contradicts switch-on-name teaching and breaks if tool slice order changes
- Phase 6 spot-check: 7-minute Phase 6 window is too tight for re-run + 3-component walkthrough + 3–5 person verbal spot-check
- Phase 1: no cleanup of `reversed.txt` after live Phase 1 run — Phase 6 re-run and attendee builds in Phase 5 may produce inconsistent behavior
- Phase 4 checkpoint: `go build -o /dev/null` fails on Windows — no documented fallback
- Pacing Risk 1: no `read_file` schema snippet — only points to `milestone-2/main.go`; no rescue anchor if that file is incomplete
- Phase 1: `demo-task.txt` vs `demo.txt` distinction mentioned once as parenthetical, never reinforced
- Pre-session checklist: no item verifies `demo-task.txt` and `demo.txt` exist before session start
- Pre-session checklist: Ollama verification step has no recovery/fix instructions (only verifies, no remediation path)
- Pre-session checklist: no Go installation or minimum version verification step
- Phase 1: no fallback if `agent-demo` fails live in front of room (Ollama down, missing binary, missing file)
- Phase 3: no fallback for attendee connection-refused errors (Ollama host/port mismatch)
- Phase 4: `Arguments` type change (string→map[string]any) flagged in transition note but no fix hint for attendees with their own code
- Phase 6 spot-check: pass threshold undefined when fewer than 3 attendees present
- Phase 7: 3-minute window (57–60) insufficient for cold Milestone 3 build; below-70% threshold has no in-session recovery path
- Phase 7: no recovery documented if attendee is not in Slack group at wrap time
- Phase 7: no room-size guidance — visual build-completion scan unreliable for groups larger than ~15
- General: no zero-attendee scenario documented
