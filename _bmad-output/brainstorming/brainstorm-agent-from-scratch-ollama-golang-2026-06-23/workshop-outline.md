# Workshop Outline: Build an AI Agent from Scratch with Ollama + Go

**Duration:** 60 minutes  
**Stack:** Vanilla Go (no external deps beyond Ollama client), Ollama running locally  
**Promise:** By the end, you will understand exactly how an agent works — and you will have built one you can launch from your own compiled binary.

---

## MINUTE 0–5 — The Hook: Watch It Change a File

**Facilitator action:** Run the pre-built reference binary. Give it a live coding task.  
*Example: "Add a function that reverses a string to main.go"*

The agent calls the model, decides to use a write-file tool, and the file changes on screen.

**Say aloud:**
> "That's it. That's the whole magic. By the end of this hour, you'll know exactly why that worked — every line of it."

**Why this matters:** The minute-one hook and the final takeaway are the same moment. Attendees will watch this again at minute 60 and understand every mechanism underneath it.

---

## MINUTE 5–10 — Context: What We're Building and Why From Scratch

**Talk track (5 min, no code yet):**

- You've probably used LangChain, LlamaIndex, or a similar framework. That's fine. This workshop is about demystification, not replacement.
- Frameworks evolve fast and carry supply-chain risk. Understanding internals lets you audit them, debug them, and not blindly trust them.
- Hard constraints for today: vanilla Go, one hour, nothing hidden. Constraints are the point — they make every mechanism visible.
- After today: open any agent framework's source and see through it immediately.

**Core mental model to establish:**
> Agent = LLM in a loop + tools + a stop condition

Write this on the board. It stays up for the whole session.

---

## MINUTE 10–20 — Part 1: The LLM Call (10 min)

**What attendees write:**

1. A `main.go` that sends a user prompt to Ollama via HTTP and prints the response.
2. No abstraction yet — raw `http.Post`, read the JSON, print the content.

**Key teaching moment:**  
Strip away the mystery of "talking to an LLM." It is an HTTP call with JSON in and JSON out. Show the raw request and response bodies.

**Checkpoint:** Attendees have a working Go program that prompts Ollama and reads a reply.

---

## MINUTE 20–35 — Part 2: Tools (15 min)

**Concept to introduce (2 min):**  
The model cannot act on the world by itself. Tools are Go functions the model is *allowed to call*. You decide which tools exist. You decide what they do.

**What attendees write:**

1. A `Tool` struct: `Name`, `Description`, `Parameters` (JSON schema-style), and a `Run` function.
2. Two concrete tools:
   - `read_file(path string) string`
   - `write_file(path string, content string)`
3. A `tools` slice passed as part of the Ollama request so the model knows what's available.

**Key teaching moment:**  
The model does not execute tools. It outputs a *request* to call a tool (a JSON blob). Your Go code reads that request and decides whether to call it. You are always in control.

**Checkpoint:** Attendees can see the model responding with a tool-call JSON when asked to read or write a file.

---

## MINUTE 35–50 — Part 3: The Loop (15 min)

**Concept to introduce (2 min):**  
One tool call is not an agent. An agent loops: call model → if tool requested, run tool → feed result back → call model again → repeat until model says it's done.

**What attendees write:**

1. A `for` loop wrapping the model call.
2. A switch on the response: tool call → execute → append result to message history → continue; text reply → print and break.
3. The message history slice that accumulates context across turns.

**Key teaching moment:**  
The "intelligence" is the model deciding what to do next given the full conversation history. The loop is just Go. The stop condition is just a `break`. There is no magic.

**Checkpoint:** Give the agent a task that requires two steps (e.g., read a file, then write a modified version). Watch it chain the calls without any hand-holding.

---

## MINUTE 50–57 — The Demystification Payoff

**Facilitator action:** Return to the minute-one demo. Run the reference binary again with the same task.

**Walk through it line by line:**
- The HTTP call to Ollama — that's Part 1.
- The tool definitions the model received — that's Part 2.
- The loop that kept running until the file was written — that's Part 3.

**Say aloud:**
> "Every agent you will ever read — LangChain, LlamaIndex, AutoGPT, Claude Code — is this loop. The frameworks add error handling, retries, memory backends, multi-model routing. But the skeleton is what you just wrote."

---

## MINUTE 57–60 — Run Your Own Binary

**Facilitator action:** `go build -o my-agent .`

Every attendee compiles their own binary and runs it with a live task.

**The ending moment:**  
Same file-change behavior as minute zero. Same visual. Completely different understanding of what happened.

**Optional deep ceiling (if time and interest):**  
Point attendees toward the questions that make agents better — not more code to write today, but directions to explore:
- How does the agent decide *when* it has enough information to stop?
- What happens when a tool fails mid-loop?
- How would you add memory that persists across sessions?
- How would you give it a planning step before acting?

These are architecture decisions. They now have the foundation to reason about them.

---

## Facilitator Notes

**Materials needed:**
- Ollama running locally on each machine with a capable model (e.g., `llama3`, `qwen2.5-coder`, or similar)
- Reference binary pre-compiled and ready for the minute-zero demo
- Starter file with package declaration and imports only (attendees write the rest)

**Pacing risks:**
- Part 2 (Tools) is where attendees slow down — the JSON schema for tool parameters trips people up. Have a copy-paste snippet ready.
- Part 3 (the Loop) is where the insight lands. Do not rush through it. The "there is no magic" moment needs to breathe.

**Framing to repeat throughout:**
> "This is not for production. This is for understanding. Once you understand it, you can use any framework confidently — or skip it entirely."
