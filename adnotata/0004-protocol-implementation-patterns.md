# Protocol Implementation Patterns

What we've learned implementing RIDE protocol features.

## Adding a New Message Handler

1. **Capture real traffic first** - Use dyctl multiplexer or `-log` to see actual message flow
2. **Check ~/dev/ride** - JS RIDE is the source of truth for edge cases
3. **Add case in tui.go's handleRIDEMessage()** - Parse args, update state, log
4. **Test with real Dyalog** - Protocol quirks only show up in practice

## Message Flow Gotchas

### Request-Response Timing

Some messages require waiting for a response before sending the next:

```
BAD:  SaveChanges → CloseWindow  (CloseWindow ignored while save pending)
GOOD: SaveChanges → ReplySaveChanges → CloseWindow
```

**Pattern**: Add a `Pending*` flag, send follow-up in response handler.

### Integer Booleans

Dyalog sends `0`/`1` integers, not JSON booleans:

```go
// Wrong
if debugger, ok := args["debugger"].(bool); ok { ... }

// Right
if debugger, ok := args["debugger"].(float64); ok {
    w.Debugger = debugger != 0
}
```

### Token-Based Window Management

Each editor/tracer window has a unique token. Track them:

```go
editors map[int]*EditorWindow  // token → window state
```

Dyalog reuses tokens after CloseWindow, so clean up properly.

## Tracer vs Editor

Same `OpenWindow` message, distinguished by `debugger` field:
- `debugger: 0` → regular editor (create new pane)
- `debugger: 1` → tracer window (part of call stack)

Tracer windows come in batches when errors occur (one per stack frame).

## Testing Protocol Features

1. Write test that exercises the feature
2. Run with `-log test-reports/protocol.log`
3. Check log for unexpected message ordering
4. Add assertions that verify state changes (not just sleeps)

## Reference

- `~/dev/dyctl/CLAUDE.md` - Protocol details
- `~/dev/ride/` - JS implementation (source of truth)
- `adnotata/0003-debugging-protocol.md` - Debugging with logs
