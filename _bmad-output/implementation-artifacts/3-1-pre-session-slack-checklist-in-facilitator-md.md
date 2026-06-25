---
baseline_commit: cae8ec371799251ef2b4ccf0b02b8535794efcf2
---

# Story 3.1: Pre-Session Slack Checklist in FACILITATOR.md

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As Thomas,
I want a pre-session checklist in FACILITATOR.md that includes creating a Slack group conversation with registered participants,
so that I have a repeatable process before each session and attendees have a channel for setup questions and follow-up.

## Acceptance Criteria

1. **Given** `FACILITATOR.md` **When** the pre-session checklist section is read **Then** it lists: create a Slack group conversation with all registered participants, send the invite before the session date, and confirm all attendees can access it (FR-8)

2. **Given** the Slack group process in `FACILITATOR.md` **When** it is read **Then** it specifies a direct group conversation scoped per cohort — not a standing public channel, not reused across sessions (FR-8)

3. **Given** Phase 7 of the live session **When** Thomas reaches the end of the session **Then** FACILITATOR.md instructs him to mention the Slack group and confirm attendees have access (FR-8)

## Tasks / Subtasks

- [x] Task 1: Expand the Slack placeholder in FACILITATOR.md pre-session checklist (AC: 1, 2)
  - [x] Remove the `*(Story 3.1 will add detailed sub-steps)*` placeholder annotation from the Slack checklist item
  - [x] Replace the single-line Slack item with a detailed sub-checklist:
    - [x] Create a direct Slack group conversation with all registered participants for this session cohort (one per cohort, not reused across sessions, not a standing public channel)
    - [x] Send a welcome message in the group before the session date — confirm attendees can reply (pre-session setup questions encouraged)
    - [x] Day-of: confirm all attendees show as members before Phase 1 starts
  - [x] Add a clarifying note: the group is scoped per cohort — create a new one for each session; do not reuse a group from a prior session

- [x] Task 2: Verify Phase 7 Slack cue is present and correct (AC: 3)
  - [x] Read FACILITATOR.md Phase 7 section
  - [x] Confirm: cue for Thomas to mention the Slack group and confirm attendees have access is present
  - [x] If present and correct: no change needed (it was implemented in Story 2.2)
  - [x] If missing or incorrect: add/fix the cue: *"Before we wrap — confirm everyone has access to the Slack group I created for this cohort. If you don't see it, come find me after."*

## Dev Notes

### File to modify

`FACILITATOR.md` at repo root — this file was created in Story 2.2 and is complete. This story makes one targeted edit: replace the Slack placeholder in the pre-session checklist with expanded sub-steps.

**Do NOT touch any other section of FACILITATOR.md.** Only the pre-session checklist's Slack item changes.

### Current state of FACILITATOR.md pre-session checklist

The current Slack item (Story 2.2 placeholder, confirmed by reading the file):

```markdown
- [ ] Create a Slack group conversation with all registered participants for this session cohort *(Story 3.1 will add detailed sub-steps)*
```

This must be replaced with an expanded item that satisfies all 3 ACs.

### Target state for the Slack item

Replace the placeholder with this structure (exact wording may be adjusted, but the substance must cover all 3 sub-steps and the per-cohort constraint):

```markdown
- [ ] **Slack group (per-cohort — create new for each session; do not reuse across sessions):**
  - [ ] Create a direct Slack group conversation with all registered participants for this cohort
  - [ ] Send a welcome message before the session date — encourage setup questions; confirm attendees can reply
  - [ ] Day-of: confirm all registered participants are in the group before Phase 1 starts
```

### What must NOT change

- **Phase 7 Slack cue** — already implemented in Story 2.2: *"Before we wrap — confirm everyone has access to the Slack group I created for this cohort. If you don't see it, come find me after."* This satisfies AC 3. Verify it exists; if it does, no edit needed.
- **All other FACILITATOR.md sections** — phases 1–7, pacing risks, copy-paste snippets, ring-structure framing, spot-check script. None of these are in scope.
- **README.md, milestone files, go.mod, agent-demo, demo-task.txt** — out of scope entirely.

### Scope constraint

This story is documentation-only. One targeted edit to one file. No code changes, no new files, no build or test steps.

### FR mapping

- FR-8: Thomas creates a Slack group conversation with all registered participants before each session (pre-session setup questions + post-session follow-up); group mentioned during Phase 7.
- This story implements FR-8 in FACILITATOR.md's pre-session checklist. Story 2.2 already implemented the Phase 7 mention.

### Previous story intelligence (Story 2.2)

- Story 2.2 created `FACILITATOR.md` and intentionally left the Slack checklist item as a placeholder for Story 3.1.
- Story 2.2 Dev Notes explicitly state: "Story 3.1 will add the Slack process detail; this story adds the placeholder item."
- Phase 7 Slack cue was fully implemented in Story 2.2 — confirmed in Story 2.2 completion notes: *"Slack cue: 'Before we wrap — confirm everyone has access to the Slack group I created for this cohort. If you don't see it, come find me after.'"*
- No other items in FACILITATOR.md need updating.

### Project structure notes

```
agent-workshop/
  FACILITATOR.md   ← EDIT (this story — pre-session checklist Slack item only)
  README.md        ← unchanged
  go.mod           ← unchanged
  agent-demo       ← unchanged
  demo-task.txt    ← unchanged
  starter/main.go  ← unchanged
  milestone-1/     ← unchanged
  milestone-2/     ← unchanged
  milestone-3/     ← unchanged
```

### References

- Story 3.1 ACs: [Source: _bmad-output/planning-artifacts/epics.md#Story-3.1]
- FR-8: [Source: _bmad-output/planning-artifacts/epics.md#FR-8]
- Story 2.2 Slack placeholder: [Source: _bmad-output/implementation-artifacts/2-2-facilitator-md-session-guide.md#Task-6]
- Phase 7 Slack cue (already implemented): [Source: FACILITATOR.md#Phase-7]
- Epic 3 context: [Source: _bmad-output/planning-artifacts/epics.md#Epic-3]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6[1m]

### Debug Log References

None.

### Completion Notes List

- Replaced single-line Slack placeholder in FACILITATOR.md pre-session checklist with 3-item sub-checklist covering: create group, send welcome message, day-of confirmation. Per-cohort constraint made explicit.
- Verified Phase 7 Slack cue already present (Story 2.2 implementation) — no change required.
- Documentation-only story. No code changes.

### File List

- FACILITATOR.md

### Review Findings

Acceptance Auditor: ✅ all 3 ACs satisfied. No violations. Story 3.1 changes are clean.

All findings below are pre-existing issues in FACILITATOR.md from Story 2.2 — not introduced by Story 3.1.

- [x] [Review][Defer] Snippet 2: only ToolCalls[0] dispatched — multi-tool response silently drops remaining calls [FACILITATOR.md#snippet-2] — deferred, pre-existing (also tracked in 1-5 deferred)
- [x] [Review][Defer] Snippet 2: `output` variable undeclared — snippet won't compile as-is [FACILITATOR.md#snippet-2] — deferred, pre-existing
- [x] [Review][Defer] Snippet 2: magic-index dispatch (tools[0]/tools[1]) contradicts switch-on-name teaching and breaks if tool order changes [FACILITATOR.md#snippet-2] — deferred, pre-existing
- [x] [Review][Defer] Phase 6 spot-check timing unrealistic — 7 min for re-run + 3-component walkthrough + 3-5 person verbal spot-check [FACILITATOR.md#phase-6] — deferred, pre-existing
- [x] [Review][Defer] No cleanup of reversed.txt after Phase 1 live run — Phase 6 re-run and attendee builds may produce inconsistent behavior [FACILITATOR.md#phase-1] — deferred, pre-existing
- [x] [Review][Defer] Milestone 2 checkpoint `go build -o /dev/null` fails on Windows [FACILITATOR.md#phase-4] — deferred, pre-existing
- [x] [Review][Defer] Pacing Risk 1: no read_file schema snippet — only points to milestone-2/main.go with no rescue anchor [FACILITATOR.md#pacing-risks] — deferred, pre-existing
- [x] [Review][Defer] demo-task.txt vs demo.txt distinction introduced once as parenthetical, never reinforced [FACILITATOR.md#phase-1] — deferred, pre-existing
- [x] [Review][Defer] Pre-session: no checklist item verifies demo-task.txt and demo.txt exist before session [FACILITATOR.md#pre-session-checklist] — deferred, pre-existing
- [x] [Review][Defer] Pre-session: Ollama check has no recovery/fix instructions if it fails [FACILITATOR.md#pre-session-checklist] — deferred, pre-existing
- [x] [Review][Defer] Pre-session: no Go installation or version verification step [FACILITATOR.md#pre-session-checklist] — deferred, pre-existing
- [x] [Review][Defer] Phase 1: no fallback documented if agent-demo fails live in front of room [FACILITATOR.md#phase-1] — deferred, pre-existing
- [x] [Review][Defer] Phase 3: no fallback for attendee connection-refused errors (wrong host/port for Ollama) [FACILITATOR.md#phase-3] — deferred, pre-existing
- [x] [Review][Defer] Phase 4: Arguments type change (string→map[string]any) flagged but no hint on fix for attendees evolving their own code [FACILITATOR.md#phase-4] — deferred, pre-existing
- [x] [Review][Defer] Phase 6 spot-check: pass threshold undefined when fewer than 3 attendees present [FACILITATOR.md#phase-6] — deferred, pre-existing
- [x] [Review][Defer] Phase 7: 3-minute window (57–60) insufficient for cold Milestone 3 build; below-70% threshold has no real in-session recovery path [FACILITATOR.md#phase-7] — deferred, pre-existing (also tracked in 2-2 deferred)
- [x] [Review][Defer] Phase 7: no documented recovery if attendee is not in Slack group at wrap time [FACILITATOR.md#phase-7] — deferred, pre-existing
- [x] [Review][Defer] No room-size guidance — visual build-completion scan unreliable for groups >15 [FACILITATOR.md#phase-7] — deferred, pre-existing
- [x] [Review][Defer] No zero-attendee scenario documented [FACILITATOR.md] — deferred, pre-existing

## Change Log

- 2026-06-26: Expanded Slack pre-session checklist item in FACILITATOR.md — replaced placeholder with 3-step sub-checklist (create group, send welcome, day-of confirm) and per-cohort scoping note. Verified Phase 7 Slack cue present — no change needed.
- 2026-06-26: Code review complete — 0 patch, 0 decision-needed, 19 deferred (all pre-existing from Story 2.2). Story 3.1 changes clean.
