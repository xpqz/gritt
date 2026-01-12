# Editor Prerequisites

Before implementing editor windows (Phase 3), these foundations need to be in place.

## 1. Key Mappings (config.json)

All keys defined in configuration file, not hardcoded.

- Reference: Dyalog Unix TTY docs for key codes
- Load from `config.json` at startup
- No hot-reload needed initially
- **Ctrl+?** (Ctrl+Shift+/) to show current key mappings in a floating pane

## 2. Floating Panes

All new panes are floating, resizable, moveable. Replaces current fixed-split approach (F12 debug pane becomes floating).

- Mouse support: drag to move, drag edges to resize
- Keyboard support: shortcuts TBD in discussion
- Visual focus indicator: border style or title highlight for active pane

**Open questions:**
- Z-order: click to raise? Explicit bring-to-front key?
- Minimum size constraints?
- Snap-to-edges or free positioning?
- How to handle pane overlapping session area?

## 3. Ctrl+C Confirmation

Don't kill immediately. Show confirmation screen.

- Also need to sort out Ctrl+C for copy (context-dependent)
- Ctrl+C in input = interrupt? copy? confirm quit?
- Need clear mental model for users

## 4. Clipboard Support

Essential for editors.

- Copy: Ctrl+C (when selection exists? or always in non-input context?)
- Paste: Ctrl+V
- Cut: Ctrl+X
- Selection model needed (shift+arrows? mouse drag?)

## 5. Command Syntax (⍝«)

`⍝«command args` pattern for gritt-specific commands.

Known commands:
- `⍝«log save` - save session history

Future possibilities:
- `⍝«config reload`
- `⍝«connect host:port`
- `⍝«theme dark`

## 6. Connection Resilience

- Visual indicator of connection state (connected/disconnected/reconnecting)
- **Preserve client-side state on disconnect** - don't trash history like RIDE does
- Ability to reconnect and continue with existing session buffer
- Graceful handling of interpreter crashes

## 7. Editor Undo/Redo

Deferred to later. Session buffer is append-only, so "undo" = re-execute previous commands.

For actual editor windows (editing functions), will need proper undo/redo stack - but that's part of editor implementation, not a prerequisite.

---

## Implementation Order (Proposed)

1. Key mappings config (foundation for everything else)
2. Floating pane system (needed for key mappings display, debug pane, editors)
3. Ctrl+C / clipboard (sort out key conflicts)
4. Connection resilience (before editors, so we don't lose edited functions)
5. ⍝« command parsing (can be incremental)

Then: Editor windows (Phase 3)
