# Intent: Build-an-Agent Workshop (Ollama + Vanilla Go)

## Core Promise

Demystify how AI agents work. Attendees leave able to read any agent codebase — LangChain, Claude Code, anything — and see through it. No more black box.

## Format

- Duration: 1 hour, hands-on
- Audience: developers (may have used frameworks, not built from scratch)
- Output: attendees run their own compiled Go binary — a working CLI agent

## Hard Constraints

- Vanilla Go only — no external deps beyond the Ollama client
- No advanced Go abstractions
- Not production-grade (explicit; constraints are the teaching tool, not a liability)
- Single track: no Python alternative, no multi-LLM scope

## Narrative Arc

**Minute-one hook** — attendee gives the agent a coding task and watches it change a file on disk. Instant "I made an agent" moment.

**Middle** — walk through why it worked: the agent loop, tool calls, LLM interaction over Ollama. Every part is visible because there are no layers.

**Ending payoff** — attendee compiles and runs their own binary, uses it like a CLI tool (Claude Code-style UX: simple surface, deep ceiling). The minute-one demo is now fully understood.

## Key Insight (from devil's advocate stress-test)

The from-scratch constraint *is* the value. Frameworks exist; this workshop isn't about production code. It's about reading any framework's source afterward and understanding it immediately. Vanilla Go + 1h + "not for prod" make every mechanism visible — nothing is hidden behind convenience.

## Agent Design Principle to Embed

Simple surface, deep ceiling — easy to start (understand the loop), endless depth (redesign the architecture, add planning, memory, multi-step reasoning). Workshop structure mirrors the agent design: one path, two depths, no explicit fork.

## Out of Scope

- Terminal UI / full Claude Code-style CLI experience (idea surfaced, not chosen)
- Multi-LLM support
- Python parallel track
- Advanced agent features (memory, planning loops) as workshop deliverables — mention as "deep ceiling" only
