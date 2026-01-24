# 0007 - Character Input (⍞) Causes NONCE ERROR

*2026-01-24*

## Problem

When Dyalog is waiting for `⍞` (character input), sending input via socket causes NONCE ERROR.

## Reproduction

```bash
gritt -addr localhost:14502 -sock /tmp/apl.sock &

# This works (⎕ input)
echo "⎕" | nc -U /tmp/apl.sock
# → ⎕:
echo "42" | nc -U /tmp/apl.sock
# → 42

# This fails (⍞ input)
echo "⍞" | nc -U /tmp/apl.sock
# → (waiting for character input)
echo "hello" | nc -U /tmp/apl.sock
# → NONCE ERROR
```

Also seen in STARMAP's DISPLAY function which uses `⍞` for "Insert fine plotting element and press ENTER" prompt.

## Analysis

`⎕` and `⍞` are different:
- `⎕` = evaluated input - Dyalog expects an APL expression, evaluates it
- `⍞` = character input - Dyalog expects raw characters, no evaluation

Currently gritt always sends `Execute`:

```go
client.Send("Execute", map[string]any{
    "text":  expr + "\n",
    "trace": 0,
})
```

This works for `⎕` because `Execute` evaluates the expression. But for `⍞`, Dyalog wants raw character data, not an expression to execute.

The NONCE ERROR is likely Dyalog saying "I asked for characters, you sent an Execute command - that operation isn't valid in this context."

## Suggested Fix

The RIDE protocol likely has different message types for different input modes. When Dyalog sends `SetPromptType`:
- Type 0 (or similar) = waiting for `⎕` → respond with `Execute`
- Type 4 (or similar) = waiting for `⍞` → respond with different message (maybe `ReplyInput`?)

gritt should:
1. Track the current prompt type from `SetPromptType` messages
2. When sending input, use `Execute` for `⎕` prompts, different message for `⍞` prompts

## References

- Check RIDE protocol documentation for input message types
- Look at how the RIDE JS client handles `⍞` input
- SetPromptType type values and their meanings
