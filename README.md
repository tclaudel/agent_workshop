# Agent Workshop

A hands-on Go workshop where you build a minimal LLM agent from scratch — one concept at a time.

## Prerequisites

### Install Go

Follow the official installer for your OS: <https://go.dev/doc/install>

Requires Go 1.22 or later — verify with `go version` after installing.

### Install Ollama

[Ollama](https://ollama.com) is a local model runtime — it downloads and serves open-weight LLMs on your machine so you can call them over HTTP without any cloud account or API key. The workshop uses it as the model backend. Full documentation: <https://ollama.com/docs>

Follow the official installer for your OS: <https://ollama.com/download>

### Pull the model

```bash
ollama pull llama3.2
```

**Linux users:** start the Ollama service before the session if it is not already running:

```bash
sudo systemctl start ollama
```

## Verify Your Setup

Run both checks from the project root (the directory containing `go.mod`). **Both must succeed to be session-ready.**

**Check 1 — Go:** build the starter file (no binary output, just a compilation check)

```bash
go build -o /dev/null ./starter/
```

Expected: exits with no output and no errors.

**Check 2 — Ollama:** run the model with a simple prompt

```bash
ollama run llama3.2 "hello"
```

Expected: the model returns a response.

If either check fails, re-read the install steps above or ask the facilitator before the session starts.

## Milestone Walkthrough

The workshop is structured as three additive milestones. Each milestone is a standalone Go program (`milestone-N/main.go`) that builds on the previous one by introducing exactly one new concept block.

| Milestone | Directory | New concept |
|-----------|-----------|-------------|
| 1 | `milestone-1/` | HTTP call to Ollama — sends a prompt, prints the reply |
| 2 | `milestone-2/` | Tool struct + tools slice — adds `read_file` / `write_file`, passes tools in the request; the model now *asks* to call one (we surface it but don't run it) |
| 3 | `milestone-3/` | Agent loop + dispatch + history — wraps Milestone 2 in a `for` loop with `switch` dispatch and `[]Message` accumulation |
| 4 *(optional)* | `milestone-4/` | Output validation + self-correction — a deterministic check on the LLM's result, fed back into the loop so the agent retries |

Each milestone = previous milestone + one concept block. If you fall behind, locate the milestone file for the current step and read it to catch up independently.

## Optional Milestone 4 — Don't Trust the Model, Test It

For attendees already fluent with the agent loop. An LLM is **probabilistic**: "uppercase this text" is a request, not a guarantee. `llama3.2` gets this exact task right only about **1 in 3** runs — it may leave a letter lowercase, alter a word, write a `$(...)` placeholder, or never write the file. Milestone 4 adds the missing piece — **verify the output, and feed failures back into the loop**.

It introduces one concept block over Milestone 3:

- `checkUppercase(original, output)` — a deterministic check on the result (the output must equal the original with every letter uppercased). Returns every failure it finds.
- `validateUppercase(src, dst)` — reads both files and runs the checks.
- A loop that, when the model says it's done, **validates before trusting it**. On failure it sends the failure message back as a new user turn — the agent self-corrects. A `maxTurns` cap prevents an infinite "fix it / still wrong" cycle.

The checks live in plain functions, so you can unit-test them with **no Ollama and no network**:

```bash
go test ./milestone-4/
```

Build and run the full agent the same way as Milestone 3:

```bash
go build -o ./milestone-4-bin ./milestone-4/
./milestone-4-bin
```

**Takeaway:** real agents pair the loop with tests on the output. Verification is what turns a probabilistic model into a reliable system.
