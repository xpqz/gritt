# 0008 - Socket Mode Quote-Quad Input: Unsolved Mystery

*2026-01-24*

## Summary

Quote-quad (`⍞`) input works in gritt TUI and RIDE, but fails with NONCE ERROR in socket mode. Despite extensive investigation, the root cause remains unknown.

## What Works

- **gritt TUI**: Type `⍞`, type input, get result. Works perfectly.
- **RIDE**: Same flow, works perfectly.
- **Socket mode with `⎕`**: Quad input works fine.

## What Fails

```bash
echo "⍞" | nc -U /tmp/apl.sock   # Returns empty (correct - quote-quad doesn't echo)
echo "hello" | nc -U /tmp/apl.sock   # NONCE ERROR
```

## Investigation

### Protocol Analysis

Captured messages from both RIDE (via multiplexer MITM) and gritt socket mode.

**RIDE flow (works):**
```
→ ["Execute",{"trace":0,"text":"      ⍞\n"}]
← ["AppendSessionOutput",{"result":"...","type":14,...}]
← ["SetPromptType",{"type":0}]
← ["SetPromptType",{"type":4}]
→ ["GetAutocomplete",{"line":"","pos":0,"token":0}]  ← UI feature
← ["ReplyGetAutocomplete",{"options":[],...}]
→ ["Execute",{"trace":0,"text":"hello\n"}]
← ["AppendSessionOutput",{"result":"hello\n","type":14,...}]
← ["AppendSessionOutput",{"result":"hello\n","type":2,...}]
← ["SetPromptType",{"type":1}]
```

**Socket mode flow (fails):**
```
→ ["Execute",{"trace":0,"text":"⍞\n"}]
← ["AppendSessionOutput",{"result":"⍞\n","type":14,...}]
← ["SetPromptType",{"type":0}]
← ["SetPromptType",{"type":4}]
→ ["Execute",{"trace":0,"text":"hello\n"}]
← ["HadError",{"error":16,"dmx":0}]   ← NONCE ERROR
```

### What We Tested

1. **JSON key order** - RIDE sends `{"trace":0,"text":...}`, we tried both orders. No difference.

2. **6-space prefix** - RIDE sends `"      ⍞\n"` with leading spaces. We tried with and without. No difference.

3. **GetAutocomplete roundtrip** - RIDE sends this between prompts. We tried adding it. No difference.

4. **Exact JSON format** - Used `SendRaw` to send character-identical JSON. No difference.

### Key Observations

1. **Same RIDE client** - Socket mode uses the same `ride.Client` (same TCP connection) for all socket clients. Connection context should be identical.

2. **TUI has continuous receive loop** - The TUI runs a goroutine that constantly reads from the RIDE connection. Socket mode only reads inside `execCapture`.

3. **Protocol docs incomplete** - The RIDE protocol docs literally say `:red_circle: These modes need explaining with expected behaviour` for prompt types.

## Theories (Unverified)

1. **Hidden connection state** - Dyalog may track state that isn't visible in the protocol messages. RIDE maintains this state through continuous message flow; we don't.

2. **Timing/acknowledgment** - Dyalog may expect some acknowledgment or message flow between SetPromptType and the next Execute. RIDE's UI naturally provides this; we don't.

3. **Execute not valid in type 4** - Despite RIDE sending Execute for quote-quad input, there may be something about HOW or WHEN it sends it that differs.

## Current Status

**Documented as known limitation.** Quote-quad (`⍞`) input does not work in socket mode. Regular expressions and quad (`⎕`) input work fine.

## Next Steps

- Ask Dyalog support for clarification on RIDE protocol for quote-quad input
- Investigate if continuous background receive loop helps
- Check if there's hidden state in the RIDE client connection
