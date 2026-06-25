---
reviewer: claude-sonnet-4-6[1m]
reviewed: 2026-06-24
prd: prd.md
verdict: needs-fixes
---

# PRD Quality Review — agent_workshop

## Verdict: NEEDS-FIXES

The PRD is well-scoped and clearly written for an internal workshop. Structure, vision, non-goals, and assumptions are strong. Three issues block a clean pass; two more are worth fixing before implementation begins.

---

## Findings

### [HIGH] FR-6 phase count contradicts section header — §5.2

**Location:** §5.2 FR-6 description line vs. consequences list.

The description says "5-phase structure" but the consequences list enumerates seven distinct phases (Phase 1 through Phase 7). This is a direct internal contradiction. A second facilitator reading `FACILITATOR.md` (which this FR drives) will receive ambiguous guidance on how many phases to plan for.

**Fix:** Change "5-phase structure" in the FR-6 description line to "7-phase structure" to match the enumerated consequences.

---

### [HIGH] Glossary missing "Milestone" coverage gap — §4 vs FRs

**Location:** §4 Glossary entry for "Milestone" vs. FR-1 consequences.

The glossary defines "Milestone" as "one of three progressive code states." FR-1 then specifies those as `milestone-1/`, `milestone-2/`, `milestone-3/` directories. However, FR-6 references numbered phases (Phase 1–7) that do not map 1-to-1 to the three milestones. The glossary does not define how milestones relate to phases, and neither section cross-references the other. An implementer building the repo structure could mis-align the directory layout with the session flow.

**Fix:** Add a sentence to the Milestone glossary entry: "Milestone 1 corresponds to Phase 3, Milestone 2 to Phase 4, Milestone 3 to Phase 5." Alternatively, add a cross-reference in FR-1 to FR-6.

---

### [MEDIUM] UJ-2 lacks a named persona — §3.3

**Location:** §3.3 UJ-2.

UJ-1 (Marco) is concrete: named persona, specific background, beat-by-beat narrative. UJ-2 describes Thomas re-running the session, but Thomas is the facilitator, not a distinct user type. The journey has no named second-facilitator persona, which is the actual scenario this UJ is supposed to validate (per §9 Open Questions and SM-4). The journey is functional but tests the wrong actor — it validates Thomas's repeatability, not the hand-off scenario the PRD calls out as a risk.

**Fix:** Either (a) rename UJ-2 to make it explicitly a "Thomas repeats the session" journey and add a separate UJ-3 with a named second facilitator persona (e.g., "Léa, a senior engineer who was not at the dry run, takes over a session"), or (b) reframe UJ-2 around the second-facilitator scenario since that is the higher-risk case. Option (a) is preferred given SM-4 explicitly validates repeatability by a different facilitator.

---

### [MEDIUM] SM-1 comprehension target is not binary-measurable — §8

**Location:** §8 SM-1.

SM-1 describes a spot-check of 3–5 attendees, with the target that "attendees can identify the HTTP call, tool dispatch, and loop break." There is no pass/fail threshold stated (e.g., "4 of 5 attendees can name all three components"). As written, one attendee naming one component would technically satisfy the metric. For an internal workshop this is low stakes, but if SM-4 (repeatability) depends on consistent session quality, a vague SM-1 weakens the feedback loop.

**Fix:** Add a threshold: e.g., "Target: ≥4 of 5 spot-checked attendees can name all three components unprompted."

---

### [LOW] "Basic Go fluency" assumption is underdefined for enforcement — §3.2 / §10

**Location:** §3.2 Non-Users and §10 Assumptions Index.

The assumption defines "basic Go fluency" as "can write, build, and run a simple Go program." This is correctly flagged as an assumption, but no mechanism is described for how this prereq is communicated or verified before registration. The non-users section lists this as a disqualifier, yet §5.1 FR-5 only covers the setup guide (Go install, Ollama, llama3.2 pull) — it does not include a Go fluency self-check. An attendee who has Go installed but has never written a Go program would pass the FR-5 verify step and arrive unqualified.

**Fix:** Add a line to FR-5 consequences: "Prereq section includes a Go fluency self-check (e.g., 'If you have never written a Go function, this session requires a 30-minute Go tour first — link')." This keeps the assumption explicit and actionable rather than passive.

---

## Criteria Summary

| Criterion | Result |
|---|---|
| Essential sections present and scoped for internal workshop | Pass |
| FRs are capabilities not implementations, testable, actor-scoped | Pass |
| Glossary covers all domain nouns in FRs and UJs | Partial — milestone/phase cross-reference missing |
| Success metrics have targets and cross-reference FRs | Partial — SM-1 threshold unquantified |
| Assumptions explicit, not buried | Pass |
| Non-goals specific enough to prevent scope creep | Pass |
| No contradictions between sections | Fail — FR-6 phase count contradicts description |
| UJs concrete with named personas and specific beats | Partial — UJ-2 wrong actor |
| PRD length appropriate for stakes (~5–8 pages internal) | Pass — tight, well-scoped |
