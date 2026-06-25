---
title: Reconciliation — workshop-outline.md vs prd.md
created: 2026-06-24
source-input: brainstorm-agent-from-scratch-ollama-golang-2026-06-23/workshop-outline.md
source-prd: prd-agent_workshop-2026-06-24/prd.md
---

# Reconciliation: Workshop Outline vs PRD

## Method

Line-by-line comparison of qualitative content (tone, voice, framing, pacing notes, facilitator talk-tracks, optional extensions) present in the workshop outline against what was captured in the PRD's FRs, vision, and facilitator guide requirements. The PRD is evaluated for silent drops — content that existed in the input but left no trace in the output.

---

## GAP 1 — The "same moment" hook/payoff framing is lost

**In the outline (Minute 0–5):**
> "The minute-one hook and the final takeaway are the same moment. Attendees will watch this again at minute 60 and understand every mechanism underneath it."

This is not a structural requirement — it is a *deliberate narrative device*. The session is designed as a ring structure: the opening demo and the closing demo are identical visual events, but the attendee's understanding of them is completely different. The PRD captures the two demo moments as separate checklist items (FR-3: reference binary; FR-6: Phase 1 and Phase 6) but nowhere articulates that they are intentionally the *same moment* serving a rhetorical purpose. The `FACILITATOR.md` FR (FR-4) says "talk-track anchors" but doesn't name or protect this specific device.

**Risk:** A second facilitator reading `FACILITATOR.md` without the outline might treat Phase 6 as a review step rather than a culminating payoff. The structural symmetry — and why it matters for attendee experience — is absent from any FR.

**Where it should appear:** FR-4 (FACILITATOR.md) should explicitly require that the guide articulates the ring structure and the verbal cue that connects minute 0 to minute 57.

---

## GAP 2 — "There is no magic" as a design principle for the loop section

**In the outline (Minute 35–50):**
> "The 'intelligence' is the model deciding what to do next given the full conversation history. The loop is just Go. The stop condition is just a `break`. There is no magic."

And the pacing note:
> "Part 3 (the Loop) is where the insight lands. Do not rush through it. The 'there is no magic' moment needs to breathe."

The PRD captures the loop's technical requirements correctly (FR-1: milestone-3 with loop, dispatch switch, message history). It also captures the pacing risk for Part 2 (Tool JSON schema). But the pacing note for Part 3 — the explicit instruction that the loop section is the *emotional and intellectual payoff* of the session and must not be rushed — is entirely absent from FR-4 and from FR-6.

**Risk:** The PRD's SM-C1 counter-metric ("do not optimize for session speed") is related but too abstract. It does not tell a facilitator *which phase* carries the most insight weight and therefore must be protected from time pressure. The outline is explicit: Part 3 is the moment. The PRD does not say this.

**Where it should appear:** FR-4 (FACILITATOR.md) should call out Part 3 as the critical insight phase, not just Part 2 as the technical pacing risk.

---

## GAP 3 — The "after today" transfer goal: framework auditing

**In the outline (Minute 5–10):**
> "After today: open any agent framework's source and see through it immediately."

And later (Minute 50–57):
> "Every agent you will ever read — LangChain, LlamaIndex, AutoGPT, Claude Code — is this loop. The frameworks add error handling, retries, memory backends, multi-model routing. But the skeleton is what you just wrote."

The PRD's vision section captures the demystification framing well. But the *specific, concrete transfer promise* — that the attendee will be able to audit LangChain, LlamaIndex, AutoGPT, Claude Code source code after this session — appears only weakly in the vision ("read any agent framework's source, debug any agent behavior"). It does not appear at all in the FR-6 phase structure for the demystification payoff (Phase 6), nor in the facilitator guide requirement (FR-4).

The outline treats the naming of specific frameworks (LangChain, LlamaIndex, AutoGPT, Claude Code) as a deliberate verbal move that lands the transfer claim concretely. The PRD abstracts it to generic language.

**Risk:** A facilitator will not know to name those specific frameworks unless it is in `FACILITATOR.md`. The specificity is what makes the claim feel credible and earned. Generic "frameworks" does not land the same way.

**Where it should appear:** FR-4 should require the FACILITATOR.md to include the specific frameworks to name during Phase 6 as part of the talk-track anchor for that phase.

---

## GAP 4 — The optional "deep ceiling" extension questions

**In the outline (Minute 57–60):**
> "Optional deep ceiling (if time and interest): Point attendees toward the questions that make agents better — not more code to write today, but directions to explore:
> - How does the agent decide when it has enough information to stop?
> - What happens when a tool fails mid-loop?
> - How would you add memory that persists across sessions?
> - How would you give it a planning step before acting?
> These are architecture decisions. They now have the foundation to reason about them."

This content — the four specific extension questions and the framing that they are *architecture decisions, not code to write* — is entirely absent from the PRD. There is no FR that captures it, no mention in FR-4 (FACILITATOR.md), and no mention in FR-6 (phase structure). The "Optional deep ceiling" is treated as if it doesn't exist.

**Risk:** The four questions are not arbitrary; they map to the exact gaps the workshop deliberately leaves open (stop condition robustness, error handling, memory persistence, planning). They give attendees a structured off-ramp toward further learning without requiring the workshop to cover them. A facilitator without the outline will not know to offer this, and the session ends abruptly at "run go build" with no forward horizon for attendees.

**Where it should appear:** FR-4 or FR-6 should include the optional extension questions as a Phase 7 element when time and energy allow. This is low-cost to add and high-value for attendee experience.

---

## GAP 5 — The copy-paste snippet for Tool JSON schema is under-specified

**In the outline (Facilitator Notes / Pacing risks):**
> "Part 2 (Tools) is where attendees slow down — the JSON schema for tool parameters trips people up. Have a copy-paste snippet ready."

The PRD captures this pacing risk and names it in FR-4:
> "specifically: the Tool JSON schema parameters snippet flagged as a pacing risk"

However, FR-4 only says the snippet must be *present* in FACILITATOR.md. It does not require the snippet to be *ready in the starter file or a separate resource attendees can access directly*. The outline's intent ("have a copy-paste snippet ready") implies it should be immediately accessible to attendees, not just to the facilitator.

**Risk:** If the snippet is only in `FACILITATOR.md` (a facilitator-facing file), attendees who fall behind cannot self-rescue. The outline implies the snippet should be attendee-accessible — likely in the repo's README, the starter code as a comment, or a dedicated `snippets/` reference file.

**Where it should appear:** FR-2 (starter file) or a new FR should specify that the Tool JSON schema snippet is accessible to attendees directly, not only to the facilitator via `FACILITATOR.md`.

---

## Minor Drops (below gap threshold, noted for completeness)

- **Model flexibility language:** The outline lists `llama3`, `qwen2.5-coder`, "or similar" as acceptable models. FR-5 specifies only `llama3.2`. The outline's flexibility framing ("a capable model") is lost — the PRD is more prescriptive than the input intended.

- **"Constraints are the point" framing:** The outline explicitly calls out that the constraints (vanilla Go, one hour, no deps) are pedagogically intentional: "Constraints are the point — they make every mechanism visible." This framing does not appear in the vision or anywhere in the PRD. It is the philosophical justification for the stack choice, not just a scope decision.

- **"Not for production" framing:** The outline says to repeat throughout: "This is not for production. This is for understanding." The PRD captures the non-production scope in non-goals but does not require the facilitator to verbally repeat this framing during the session. It is a tone anchor, not just a scope boundary.

---

## Summary Table

| Gap | Content Dropped | Severity | Recommended Fix |
|-----|----------------|----------|-----------------|
| 1 | Ring structure: minute-0 and minute-57 are the same moment | High | FR-4: require FACILITATOR.md to name and protect this device |
| 2 | Part 3 loop is the emotional/intellectual payoff; must not be rushed | High | FR-4: add pacing note for Part 3, not just Part 2 |
| 3 | Specific framework names (LangChain, LlamaIndex, AutoGPT, Claude Code) in Phase 6 talk-track | Medium | FR-4: include specific names in Phase 6 verbal cue |
| 4 | Optional "deep ceiling" — four extension questions + architecture framing | Medium | FR-4 or FR-6: add as optional Phase 7 element |
| 5 | Tool JSON schema snippet must be attendee-accessible, not only facilitator-facing | Medium | FR-2 or new FR: place snippet where attendees can access it directly |
| — | Model flexibility ("capable model, e.g. llama3, qwen2.5-coder") vs hard llama3.2 | Low | FR-5: soften to "llama3.2 or equivalent capable model" |
| — | "Constraints are the point" philosophical framing | Low | Vision: add one sentence on why vanilla Go + one hour is the stack |
| — | "Not for production" verbal anchor to repeat throughout session | Low | FR-4: add as recurring verbal cue across all phases |
