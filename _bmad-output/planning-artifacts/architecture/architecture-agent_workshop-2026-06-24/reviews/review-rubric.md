---
review-date: 2026-06-24
reviewed-by: Claude Code
verdict: PASS (with minor clarifications required)
---

# Architecture Spine Review: agent_workshop

## Verdict: PASS

Spine cleanly maps PRD capabilities to enforcement rules. Progressive disclosure pattern is well-architected via ADs 1–3. Divergence points for the three milestones are defined and enforceable. However, three clarifications needed before repo implementation.

---

## Findings

**ARCHITECTURE-SPINE.md:92, AD-8** — Architecture hardcodes darwin/arm64 but no conditional build setup defined. If Thomas's machine is Apple Silicon and a contributor is not, binary diverges without detection. Remedy: Either (1) add `GOOS=darwin GOARCH=arm64` as explicit build-command convention in AD-8, or (2) document in FACILITATOR.md how to build for contributor's own arch.

**ARCHITECTURE-SPINE.md:98, AD-9** — Module path asserted as `github.com/tclaudel/agent-workshop` without confirming actual repo name. Verify URL before committing rule to architecture.

**ARCHITECTURE-SPINE.md:116, Stack** — Ollama 0.3+ tool_calls requirement cited from training data; version claim and feature availability not web-verified. Before declaring hard requirement, confirm Ollama 0.3 release date and that tool_calls is stable/documented in that release.

---

## Secondary Issues (non-blocking, low priority)

**ARCHITECTURE-SPINE.md:107, Conventions** — Tool dispatch rule lists "read_file and write_file are the two required tool names." AD-5 (canonical Tool struct shape) allows arbitrary tool names. Clarify: are read/write tools mandatory (i.e., must appear in every milestone-2+), or are they just the example tools used in the demo? If mandatory, add to AD-5 as a rule; if examples only, move naming to Deferred or clarify as "example tools."

---

## Checklist Summary

1. **Divergence points fixed** ✓ — ADs 1–3 cleanly prevent milestone divergence; AD-4–9 govern stack/binary/dispatch.  
2. **Rules enforceable** ~ — Most are (structural/code-review level); AD-4 lacks automation; AD-8 assumes build discipline.  
3. **Deferred scope clean** ✓ — No deferred items create intra-milestone divergence.  
4. **Tech verified-current** ~ — Go 1.22+, llama3.2 OK; Ollama 0.3+ tool_calls needs web check.  
5. **All dimensions owned/deferred** ✓ — Operational envelope (darwin/arm64) assigned to AD-8; session logistics out of scope; logging not applicable to workshop scope.  
6. **All PRD capabilities mapped** ✓ — FR-1 to FR-5 covered; FR-6 (session structure) explicitly deferred to facilitator guide (reasonable); FR-7 depends on FR-3/FR-4; FR-8 out of scope.

---

## Recommendation

Merge with three clarifications: (1) confirm github.com/tclaudel/agent-workshop is the target repo, (2) add explicit build-command pattern to AD-8 or document arch flexibility in FACILITATOR.md, (3) verify Ollama 0.3 tool_calls availability via Ollama changelog/docs before declaring hard requirement.
