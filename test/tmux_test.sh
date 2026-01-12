#!/bin/bash
# tmux-driven TUI test for gritt
# Generates HTML report with snapshots

set -e

SESSION="gritt-test"
GRITT_BIN="${GRITT_BIN:-./gritt}"
TIMEOUT=5

# Report setup
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
REPORT_DIR="test-reports"
REPORT_FILE="$REPORT_DIR/test-$TIMESTAMP.html"
mkdir -p "$REPORT_DIR"

# Test state
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0
SNAPSHOTS=""

cleanup() {
    tmux kill-session -t "$SESSION" 2>/dev/null || true
}
trap cleanup EXIT

# Capture pane and return escaped HTML
capture() {
    tmux capture-pane -t "$SESSION" -p
}

capture_html() {
    capture | sed 's/&/\&amp;/g; s/</\&lt;/g; s/>/\&gt;/g'
}

snapshot() {
    local label="$1"
    local content=$(capture_html)
    SNAPSHOTS+="<div class=\"snapshot\">
<h3>$label</h3>
<pre>$content</pre>
</div>
"
}

wait_for() {
    local pattern="$1"
    local timeout="${2:-$TIMEOUT}"
    local start=$(date +%s)

    while true; do
        if capture | grep -q "$pattern"; then
            return 0
        fi
        local now=$(date +%s)
        if (( now - start > timeout )); then
            return 1
        fi
        sleep 0.2
    done
}

send_keys() {
    tmux send-keys -t "$SESSION" "$@"
}

send_line() {
    send_keys "$1" Enter
}

# Test runner
run_test() {
    local name="$1"
    local check="$2"

    ((TESTS_RUN++))

    if eval "$check"; then
        ((TESTS_PASSED++))
        echo -e "\033[0;32mPASS:\033[0m $name"
        SNAPSHOTS+="<div class=\"result pass\">✓ $name</div>
"
        return 0
    else
        ((TESTS_FAILED++))
        echo -e "\033[0;31mFAIL:\033[0m $name"
        SNAPSHOTS+="<div class=\"result fail\">✗ $name</div>
"
        snapshot "Failed: $name"
        return 1
    fi
}

generate_report() {
    local status_class="pass"
    local status_text="All tests passed"
    if (( TESTS_FAILED > 0 )); then
        status_class="fail"
        status_text="$TESTS_FAILED test(s) failed"
    fi

    cat > "$REPORT_FILE" << EOF
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>gritt test report - $TIMESTAMP</title>
    <style>
        body {
            font-family: 'SF Mono', 'Menlo', 'Monaco', 'Cascadia Code', 'Consolas', 'DejaVu Sans Mono', -apple-system, monospace;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #1a1a2e;
            color: #eee;
        }
        h1 { color: #00d9ff; }
        h2 { color: #888; border-bottom: 1px solid #333; padding-bottom: 10px; }
        h3 { color: #aaa; margin: 10px 0 5px 0; }
        .summary {
            background: #252540;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
        }
        .summary.pass { border-left: 4px solid #00ff88; }
        .summary.fail { border-left: 4px solid #ff4444; }
        .stats { display: flex; gap: 30px; margin-top: 15px; }
        .stat { text-align: center; }
        .stat-value { font-size: 2em; font-weight: bold; }
        .stat-label { color: #888; font-size: 0.9em; }
        .stat-value.pass { color: #00ff88; }
        .stat-value.fail { color: #ff4444; }
        .result {
            padding: 8px 15px;
            margin: 5px 0;
            border-radius: 4px;
        }
        .result.pass { background: #1a3d2a; color: #00ff88; }
        .result.fail { background: #3d1a1a; color: #ff4444; }
        .snapshot {
            margin: 20px 0;
            background: #252540;
            border-radius: 8px;
            overflow: hidden;
        }
        .snapshot h3 {
            background: #1a1a2e;
            margin: 0;
            padding: 10px 15px;
        }
        .snapshot pre {
            margin: 0;
            padding: 15px;
            overflow-x: auto;
            font-size: 14px;
            line-height: 1.2;
            background: #0a0a15;
            color: #00d9ff;
            font-family: 'SF Mono', 'Menlo', 'Monaco', 'Cascadia Code', 'Consolas', 'DejaVu Sans Mono', monospace;
        }
        .timestamp { color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>gritt test report</h1>
    <p class="timestamp">$TIMESTAMP</p>

    <div class="summary $status_class">
        <strong>$status_text</strong>
        <div class="stats">
            <div class="stat">
                <div class="stat-value">$TESTS_RUN</div>
                <div class="stat-label">Total</div>
            </div>
            <div class="stat">
                <div class="stat-value pass">$TESTS_PASSED</div>
                <div class="stat-label">Passed</div>
            </div>
            <div class="stat">
                <div class="stat-value fail">$TESTS_FAILED</div>
                <div class="stat-label">Failed</div>
            </div>
        </div>
    </div>

    <h2>Test Progress</h2>
    $SNAPSHOTS
</body>
</html>
EOF
    echo ""
    echo "Report saved to: $REPORT_FILE"
}

# --- Main ---

echo "=== gritt tmux test ==="
echo ""

# Build
echo "Building gritt..."
go build -o "$GRITT_BIN" . || { echo "Build failed"; exit 1; }

# Check Dyalog
if ! nc -z localhost 4502 2>/dev/null; then
    echo "Error: Dyalog not running on port 4502"
    echo "Start with: RIDE_INIT=SERVE:*:4502 dyalog +s -q"
    exit 1
fi

# Start gritt
echo "Starting gritt..."
cleanup
tmux new-session -d -s "$SESSION" -x 100 -y 40 "$GRITT_BIN"
sleep 1

snapshot "Initial state"

# Test 1: Initial render
run_test "Initial render shows title" "capture | grep -q 'gritt'"

# Test 2: F12 toggles debug
send_keys F12
sleep 0.3
snapshot "After F12 (debug pane open)"
run_test "F12 opens debug pane" "capture | grep -q 'debug'"

# Test 3: Focus indicator
run_test "Focused pane has double border" "capture | grep -q '╔'"

# Test 4: Esc closes
send_keys Escape
sleep 0.3
snapshot "After Esc (debug pane closed)"
run_test "Esc closes debug pane" "! capture | grep -q '╔.*debug'"

# Test 5: F12 reopens
send_keys F12
sleep 0.3
run_test "F12 reopens debug pane" "capture | grep -q 'debug'"
send_keys Escape
sleep 0.2

# Test 6: Execute 1+1
send_line "1+1"
sleep 1
snapshot "After executing 1+1"
run_test "Execute 1+1 returns 2" "capture | grep -E '^│?2 *│?$'"

# Test 7: Execute iota
send_line "⍳5"
sleep 1
snapshot "After executing ⍳5"
run_test "Execute ⍳5 returns sequence" "capture | grep '1 2 3 4 5'"

# Test 8: Navigate up and edit
send_keys Up Up Up Up  # Go to 1+1 line
sleep 0.3
send_keys End
send_keys BSpace  # Delete the second 1
send_keys "2"     # Make it 1+2
sleep 0.3
snapshot "After editing 1+1 to 1+2"
send_keys Enter
sleep 1
snapshot "After executing edited line"
run_test "Edit and re-execute works" "capture | grep -E '^│?3 *│?$'"

# Test 9: Debug pane shows protocol messages
send_keys F12
sleep 0.3
snapshot "Debug pane with protocol log"
run_test "Debug pane shows Execute messages" "capture | grep -q 'Execute'"

# Final snapshot
send_keys Escape
sleep 0.2
snapshot "Final state"

# Generate report
generate_report

echo ""
if (( TESTS_FAILED == 0 )); then
    echo "=== All $TESTS_PASSED tests passed ==="
else
    echo "=== $TESTS_FAILED of $TESTS_RUN tests failed ==="
    exit 1
fi
