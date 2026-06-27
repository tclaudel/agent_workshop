# Milestone 3 — The Agent Loop

> **New concept:** actually *run* the tools. Wrap the request in a `for` loop, dispatch each tool call, feed the result back, and keep the growing conversation as memory — until the model has no more tools to call.
>
> **Builds on:** [Milestone 2](./milestone-2.md) — same tools, same types. Milestone 2 ended on a cliffhanger: the model asked to call `read_file`, and we printed the request and quit. This milestone picks up that exact dangling call and **runs it** — then feeds the result back so the model can take the next step.

This is the milestone where the program becomes an **agent**: a system that takes an action, observes the result, and decides what to do next — repeatedly — without us scripting each step.

Milestone 2 got one tool *request* out of the model and threw it away. The whole job here is to stop throwing it away: execute the call, hand the model the result, and let it keep going. Do that in a loop and a single ignored request becomes a multi-step agent that reads a file, uppercases it, and writes the result — none of which we scripted by hand.

---

## The core idea: the agent loop

An agent is a loop over a conversation. Each turn, the model either:

- **calls tools** → we run them, append the results, loop again; or
- **answers in plain text** → it's done, print and stop.

```mermaid
flowchart TD
    start([prompt]) --> send[POST /api/chat<br/>with full message history]
    send --> resp[model reply]
    resp --> append[append reply to messages]
    append --> check{ToolCalls<br/>present?}
    check -->|no| done[print Content → exit]
    check -->|yes| run[for each tool call:<br/>look up tool, run it]
    run --> result[append each result<br/>as a 'tool' message]
    result --> send
```

The conversation slice (`messages`) **is** the agent's memory. Every reply and every tool result gets appended, so each new request carries the entire history. That's what lets the model say "I read the file, here are its contents, now I'll write the uppercased version" across multiple turns.

---

## What changed since Milestone 2

Three additions:

1. **The `for` loop** around the request (was a single call).
2. **A dispatch map** `toolsByName` for O(1) lookup of which tool to run.
3. **History accumulation** — `messages` grows every turn instead of being a fixed one-element slice.

(`ToolCallFunction.Arguments` is already `map[string]any` from Milestone 2 — the parsed args drop straight into `Tool.Run`.)

There's also a shared HTTP client with a timeout — a small robustness upgrade now that we make many calls:

```go
var httpClient = &http.Client{Timeout: 60 * time.Second}
```

The two tools (`read_file`, `write_file`) gain explicit argument validation (`if !ok || path == ""`) so a malformed call returns a readable error string instead of misbehaving.

---

## The dispatch map

Milestone 2 had a slice of tools. To *run* the one the model named, we need fast lookup by name:

```go
toolsByName := make(map[string]*Tool, len(tools))
for i := range tools {
	toolsByName[tools[i].Name] = &tools[i]
}
```

> **Note:** the README calls this "switch dispatch" — the implementation uses a **map** instead. Same idea (route a name to its handler), cleaner as the tool list grows.

---

## The loop, line by line

```go
prompt := "Read the file demo.txt, convert its contents to UPPERCASE, then write the result to uppercase.txt"
messages := []Message{{Role: "user", Content: prompt}}

for {
	chatRequest := ChatRequest{
		Model: "llama3.2", Messages: messages, Stream: false, Tools: toolDefs,
	}
	// ... marshal, POST, read, unmarshal into `result` ...
```
Each iteration sends the **entire** `messages` history, not just the latest turn. The model is stateless between calls; the history *is* the state.

```go
	msg := result.Message
	messages = append(messages, msg)   // remember what the model said
```
Append the model's reply to memory before doing anything else.

```go
	if len(msg.ToolCalls) == 0 {
		fmt.Println(msg.Content)
		break                          // no tools requested → the agent is finished
	}
```
The **termination condition**: a reply with no tool calls means the model is talking to *us*, not asking for an action. Print and exit the loop.

```go
	for _, tc := range msg.ToolCalls {
		tool, ok := toolsByName[tc.Function.Name]
		if !ok {
			messages = append(messages, Message{Role: "tool",
				Content: fmt.Sprintf("error: unknown tool %q", tc.Function.Name)})
			continue
		}
		output := tool.Run(tc.Function.Arguments)
		messages = append(messages, Message{Role: "tool", Content: output})
	}
}
```
For each requested call: find the tool, run it, and append its output as a `"tool"` message. Unknown tool? Append an error message — the model can read it and adapt. Then the loop sends everything back and the model decides what's next.

---

## A concrete run: "uppercase demo.txt into uppercase.txt"

```mermaid
sequenceDiagram
    participant U as user prompt
    participant A as agent loop
    participant O as Ollama
    participant FS as filesystem

    U->>A: "read demo.txt, write UPPERCASE to uppercase.txt"
    A->>O: messages=[user]
    O-->>A: tool_call read_file{path:"demo.txt"}
    A->>FS: os.ReadFile("demo.txt")
    FS-->>A: "hello agent workshop"
    A->>O: messages=[user, assistant, tool:"hello agent workshop"]
    O-->>A: tool_call write_file{path:"uppercase.txt", content:"HELLO AGENT WORKSHOP"}
    A->>FS: os.WriteFile("uppercase.txt", "HELLO AGENT WORKSHOP")
    FS-->>A: "ok"
    A->>O: messages=[..., tool:"ok"]
    O-->>A: "Done! I uppercased the text." (no tool calls)
    A->>U: print Content, exit
```

Notice the history growing on every arrow into Ollama. That accumulation is the whole trick.

---

## The role of `Role: "tool"`

Tool results go back with `Role: "tool"`, distinct from `"user"` and `"assistant"`. This tells the model "this text is the *output of a function you asked for*," not a new instruction from the human. Getting the role right is what keeps the conversation coherent.

---

## The catch this leaves (motivates Milestone 4)

This loop trusts the model completely. When it says "Done, I uppercased the text," we believe it and exit. But an LLM is **probabilistic** — it might:

- leave some letters lowercase,
- drop or alter words instead of just changing case,
- write a placeholder like `$(read_file.content)` instead of the real text,
- forget to write the file at all.

This is not hypothetical: `llama3.2` gets this exact task right only about **1 in 3** runs. Nothing here checks. [Milestone 4](./milestone-4.md) adds a deterministic verification step before trusting that final answer — and feeds failures back into *this same loop* so the agent corrects itself.

---

## Run it

`demo.txt` must exist in the project root (where you run the binary).

```bash
go build -o ./milestone-3-bin ./milestone-3/
./milestone-3-bin
```

Expected: JSON logs of each turn, a freshly written `uppercase.txt`, and a final confirmation line. Run it a few times — roughly 2 of 3 runs the model botches the result (a stray lowercase letter, a `$(...)` placeholder, a dropped word) yet still says "Done." That unchecked failure is exactly what Milestone 4 fixes.

---

| | |
|---|---|
| ← Previous | [Milestone 2 — Declaring Tools](../milestone-2/docs.md) |
| → Next | [Milestone 4 — Verify the Output](../milestone-4/docs.md) |
