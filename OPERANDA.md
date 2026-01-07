# OPERANDA - Current State

## Active Work

Phase 1 complete: minimal RIDE protocol client that connects and executes APL code.

### Current Status

- [x] Created project documentation (CLAUDE.md)
- [x] Set up Latin documentation structure
- [x] Initialize Go module
- [x] Implement RIDE protocol (wire format, handshake)
- [x] Create minimal main.go that connects and runs `1+1`

### Project Structure

```
gritt/
├── main.go           # Entry point, connects and runs 1+1
├── ride/
│   ├── protocol.go   # Wire format (send/recv messages)
│   └── client.go     # Connection, handshake, Execute
├── go.mod
├── CLAUDE.md
├── FACIENDA.md
└── OPERANDA.md
```

### Testing

Start Dyalog:
```bash
RIDE_INIT=SERVE:*:4502 dyalog +s -q
```

Run gritt:
```bash
go run . -addr localhost:4502
```

Or use multiplexer from dyctl.

### Key Findings

- Wire format: 4-byte BE length (includes itself) + "RIDE" + payload
- SERVE mode handshake: Dyalog sends first, we respond
- AppendSessionOutput type 14 = input echo (skip it)
- SetPromptType with type > 0 means ready for input

---

## Previous Issues

(none yet)
