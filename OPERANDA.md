# OPERANDA - Current State

## Just Fixed (Critical Bugs)

**Bug #1: Escape in edit mode was closing window instead of returning to tracer**
- Root cause: Global `ClosePane` handler (bound to Escape) was intercepting the key BEFORE the pane's `HandleKey`
- The global handler in `tui.go` always called `closeEditor` for tracer panes without checking edit mode
- Fix: Added check in `ClosePane` handler - if tracer is in edit mode, `break` to let pane handle it
- File: `tui.go:415-418`

**Bug #2: Breakpoints were being cleared when leaving tracer**
- Root cause: `ToggleStop` only set `Modified=true` for non-debugger windows
- This meant tracer breakpoints weren't saved on close (since `Modified=false` skips `saveEditor`)
- Fix: Always set `Modified=true` when breakpoints change, regardless of window type
- File: `editor.go:130-143`

**Also fixed: Code didn't compile**
- `matchKey` helper was being called but never defined
- Added `matchKey(r rune, configKey string) bool` method to `EditorPane`
- File: `editor_pane.go:549-555`

---

## Active Work: Tracer & Connection Resilience

Phase 4 tracer work complete. Connection resilience improved with automatic window restoration.

### Completed Earlier (this session)

**Breakpoints:**
- Toggle breakpoint: `C-] b` (on current line in editor/tracer)
- Visual indicator: red `●` in gutter
- `SetLineAttributes` sent immediately (breakpoints work without save)
- Breakpoints saved with function via `SaveChanges`

**Tracer Controls (single-key in tracer mode):**
| Key | Action | RIDE Message |
|-----|--------|--------------|
| Enter / n | Step over (next line) | RunCurrentLine |
| i | Step into | StepInto |
| o | Step out | ContinueTrace |
| c | Continue | Continue |
| r | Resume all threads | RestartThreads |
| p | Previous (trace backward) | TraceBackward |
| f | Forward (skip line) | TraceForward |
| e | Enter edit mode | (local) |
| Esc | Close/pop frame | CloseWindow |

**Edit Mode in Tracer:**
- Tracer windows (`Debugger=true`) are read-only by default
- Press `e` to enter edit mode → title changes from `[tracer]` to `[edit]`
- Can edit function code while debugging
- Press `Esc` in edit mode → saves changes, returns to `[tracer]` mode
- Press `Esc` in tracer mode → closes/pops frame

**Automatic Window Restoration:**
- On connect/reconnect, gritt sends `GetWindowLayout`
- Dyalog responds with `OpenWindow` for any orphaned editors/tracers
- Windows automatically restored (no manual intervention needed)
- `close-all-windows` command still available via command palette for manual clearing

**Command Palette Fixes:**
- Added scrolling support
- Fixed rendering issues

**Bug Fixes:**
- Tracer now properly blocks text insertion (checks `Debugger` flag)
- Breakpoints persist across editor close (set `Modified=true`)
- Command palette mouse clicks account for scroll offset

### Key Files Changed

| File | Changes |
|------|---------|
| `editor.go` | `ToggleStop` always sets `Modified` (breakpoints persist) |
| `editor_pane.go` | Tracer mode key handling, edit mode toggle, `InTracerMode()` |
| `tui.go` | `GetWindowLayout` on connect, `closeAllWindows`, tracer controls |
| `command_palette.go` | Scrolling support, rendering fixes |
| `tui_test.go` | Tracer mode tests, APLcart filter test |
| `cmd/explore/main.go` | Protocol exploration tool |

### Pending

1. **Tracer status bar**: Show tracer keys at bottom when tracer focused
2. **Configurable tracer keys**: Move hardcoded keys (n, i, o, etc.) to config

### Testing

```bash
# Run full test suite (54 tests)
go test -v -run TestTUI

# Manual testing with protocol log
./gritt -log debug.log

# Protocol exploration
go run ./cmd/explore/
```

### Project Structure

```
gritt/
├── main.go              # Entry point, CLI flags, color detection
├── tui.go               # bubbletea TUI - Model, Update, View
├── apl                  # Shell script for ephemeral Dyalog
├── pane.go              # Floating pane system, cellbuf compositor
├── editor.go            # EditorWindow struct (Stop, Modified, etc.)
├── editor_pane.go       # Editor/tracer pane (tracer mode, edit mode)
├── stack_pane.go        # Stack frame list pane
├── debug_pane.go        # Debug log pane
├── keys_pane.go         # Key mappings pane
├── command_palette.go   # Searchable command list pane (with scrolling)
├── apl_symbols.go       # Backtick map and symbol definitions
├── symbol_search.go     # Symbol search pane
├── aplcart.go           # APLcart integration
├── keys.go              # KeyMap struct definition
├── config.go            # Config loading (with embedded default)
├── gritt.default.json   # Default key bindings (embedded at build)
├── tui_test.go          # TUI tests (54 tests)
├── uitest/              # Test framework (tmux, HTML reports)
├── cmd/explore/         # Protocol exploration tool
├── ride/
│   ├── protocol.go      # Wire format
│   ├── client.go        # Connection, handshake
│   └── logger.go        # Protocol logging
├── test-reports/        # Generated HTML test reports
└── adnotata/            # Design notes and exploration
```

---

## Next Session

1. Add tracer-specific status bar (show tracer keys when focused)
2. Make tracer keys configurable via gritt.json
