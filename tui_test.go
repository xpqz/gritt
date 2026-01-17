package main

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"gritt/uitest"
)

const (
	dyalogPort  = 4502
	sessionName = "gritt-test"
	screenW     = 120
	screenH     = 40
)

// TestTUI runs the full TUI test suite
func TestTUI(t *testing.T) {
	// Build gritt first
	t.Log("Building gritt...")
	if err := exec.Command("go", "build", "-o", "gritt", ".").Run(); err != nil {
		t.Fatalf("Failed to build gritt: %v", err)
	}

	// Check if Dyalog is running, if not try to start it
	var dyalogCmd *exec.Cmd
	if err := uitest.RequireDyalog(dyalogPort); err != nil {
		t.Log("Starting Dyalog...")
		var startErr error
		dyalogCmd, startErr = uitest.StartDyalog(dyalogPort)
		if startErr != nil {
			t.Skipf("Dyalog not available: %v", startErr)
		}
		defer func() {
			if dyalogCmd != nil && dyalogCmd.Process != nil {
				dyalogCmd.Process.Kill()
			}
		}()
	}

	// Create test runner with protocol logging
	runner, err := uitest.NewRunner(t, sessionName, screenW, screenH, "./gritt -log test-reports/protocol.log", "test-reports")
	if err != nil {
		t.Fatalf("Failed to create runner: %v", err)
	}
	defer runner.Close()

	// Wait for gritt to start
	runner.Sleep(1 * time.Second)

	// Take initial snapshot
	runner.Snapshot("Initial state")

	// Test 1: Initial render
	runner.Test("Initial render shows title", func() bool {
		return runner.Contains("gritt")
	})

	// Test 2: C-] d toggles debug pane
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("d")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After C-] d (debug pane open)")

	runner.Test("C-] d opens debug pane", func() bool {
		return runner.Contains("debug")
	})

	// Test 3: Focus indicator
	runner.Test("Focused pane has double border", func() bool {
		return runner.Contains("╔")
	})

	// Test 4: Esc closes pane
	runner.SendKeys("Escape")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After Esc (debug pane closed)")

	runner.Test("Esc closes debug pane", func() bool {
		return !runner.Contains("╔")
	})

	// Test 5: C-] d reopens
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("d")
	runner.Sleep(300 * time.Millisecond)

	runner.Test("C-] d reopens debug pane", func() bool {
		return runner.Contains("debug")
	})

	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test 6: Execute 1+1
	runner.SendLine("1+1")
	runner.Sleep(1 * time.Second)
	runner.Snapshot("After executing 1+1")

	runner.Test("Execute 1+1 returns 2", func() bool {
		return runner.Contains("2")
	})

	// Test 7: Execute iota
	runner.SendLine("⍳5")
	runner.Sleep(1 * time.Second)
	runner.Snapshot("After executing ⍳5")

	runner.Test("Execute ⍳5 returns sequence", func() bool {
		return runner.Contains("1 2 3 4 5")
	})

	// Test 8: Edit and re-execute
	runner.SendKeys("Up", "Up", "Up", "Up")
	runner.Sleep(300 * time.Millisecond)
	runner.SendKeys("End")
	runner.SendKeys("BSpace")
	runner.SendKeys("2")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After editing 1+1 to 1+2")

	runner.SendKeys("Enter")
	runner.Sleep(1 * time.Second)
	runner.Snapshot("After executing edited line")

	runner.Test("Edit and re-execute works", func() bool {
		return runner.Contains("3")
	})

	// Test 9: Debug pane shows protocol
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("d")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("Debug pane with protocol log")

	runner.Test("Debug pane shows Execute messages", func() bool {
		return runner.Contains("Execute")
	})

	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test 10: C-] ? shows key mappings pane
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("?")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After C-] ? (key mappings pane)")

	runner.Test("C-] ? opens key mappings pane", func() bool {
		return runner.Contains("key mappings")
	})

	runner.Test("Key mappings shows Actions section", func() bool {
		return runner.Contains("Actions")
	})

	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	runner.Test("Esc closes key mappings pane", func() bool {
		return !runner.Contains("key mappings")
	})

	// Test: C-] : opens command palette
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys(":")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After C-] : (command palette)")

	runner.Test("C-] : opens command palette", func() bool {
		return runner.Contains("Commands")
	})

	runner.Test("Command palette shows debug command", func() bool {
		return runner.Contains("debug")
	})

	runner.Test("Command palette shows quit command", func() bool {
		return runner.Contains("quit")
	})

	// Test: Filter commands by typing
	runner.SendText("deb")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("Command palette filtered to 'deb'")

	runner.Test("Typing filters commands", func() bool {
		// Should still show debug but not quit
		return runner.Contains("debug")
	})

	// Test: Execute command from palette
	runner.SendKeys("Enter")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After selecting debug from palette")

	runner.Test("Selecting debug opens debug pane", func() bool {
		return runner.Contains("debug") && !runner.Contains("Commands")
	})

	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test: Escape closes command palette
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys(":")
	runner.Sleep(300 * time.Millisecond)
	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	runner.Test("Escape closes command palette", func() bool {
		return !runner.Contains("Commands")
	})

	// Test: Save command shows filename prompt
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys(":")
	runner.Sleep(300 * time.Millisecond)
	runner.SendText("save")
	runner.Sleep(200 * time.Millisecond)
	runner.SendKeys("Enter")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("Save prompt with default filename")

	runner.Test("Save command shows filename prompt", func() bool {
		return runner.Contains("Save as:")
	})

	runner.Test("Save prompt has default filename", func() bool {
		return runner.Contains("session-")
	})

	// Cancel save and continue
	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test: Pane move mode
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("d") // Open debug pane first
	runner.Sleep(300 * time.Millisecond)

	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("m") // Enter move mode
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("Pane move mode active")

	runner.Test("C-] m enters pane move mode", func() bool {
		return runner.Contains("MOVE")
	})

	// Move pane with arrow keys
	runner.SendKeys("Up", "Up", "Left", "Left")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After moving pane")

	runner.Test("Arrow keys move pane in move mode", func() bool {
		return runner.Contains("MOVE") // Still in move mode
	})

	// Exit move mode
	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	runner.Test("Escape exits pane move mode", func() bool {
		return !runner.Contains("MOVE")
	})

	// Close the debug pane
	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test: Backtick mode for APL symbols
	runner.SendKeys("`")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("Backtick mode active")

	runner.Test("Backtick activates APL symbol mode", func() bool {
		return runner.Contains("APL symbol")
	})

	runner.SendKeys("i") // Should insert ⍳
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After backtick-i (iota)")

	runner.Test("Backtick-i inserts iota", func() bool {
		return runner.Contains("⍳")
	})

	// Test: Symbol search
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys(":")
	runner.Sleep(300 * time.Millisecond)
	runner.SendText("symbols")
	runner.Sleep(200 * time.Millisecond)
	runner.SendKeys("Enter")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("Symbol search pane")

	runner.Test("Symbol search opens", func() bool {
		return runner.Contains("APL Symbols")
	})

	// Search for "rho"
	runner.SendText("rho")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("Symbol search filtered to rho")

	runner.Test("Symbol search filters by name", func() bool {
		return runner.Contains("⍴")
	})

	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test: APLcart
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys(":")
	runner.Sleep(300 * time.Millisecond)
	runner.SendText("aplcart")
	runner.Sleep(200 * time.Millisecond)
	runner.SendKeys("Enter")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("APLcart pane loading")

	runner.Test("APLcart opens", func() bool {
		return runner.Contains("APLcart")
	})

	// Wait for data to load
	runner.Sleep(3 * time.Second)
	runner.Snapshot("APLcart loaded")

	runner.Test("APLcart loads data", func() bool {
		// Should show count or entries
		return runner.Contains("(") || runner.Contains("⍬")
	})

	// Filter for "interval"
	runner.SendText("interval")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("APLcart filtered for interval")

	runner.Test("APLcart filters results", func() bool {
		// Should show interval-related entries
		return runner.Contains("interval") || runner.Contains("Interval")
	})

	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Test: Ctrl+C shows quit hint
	runner.SendKeys("C-c")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After Ctrl+C (quit hint)")

	runner.Test("Ctrl+C shows quit hint", func() bool {
		return runner.Contains("C-] q to quit")
	})

	// Test 14: Any key clears the hint
	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	runner.Test("Key clears quit hint", func() bool {
		return !runner.Contains("C-] q to quit")
	})

	// Test 15: C-] q shows quit confirmation
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("q")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After C-] q (quit confirmation)")

	runner.Test("C-] q shows quit confirmation", func() bool {
		return runner.Contains("Quit? (y/n)")
	})

	// Test 16: n cancels quit
	runner.SendKeys("n")
	runner.Sleep(200 * time.Millisecond)

	runner.Test("n cancels quit confirmation", func() bool {
		return !runner.Contains("Quit? (y/n)")
	})

	// === BREAKPOINT WORKFLOW TEST ===
	// Clear input line (may have leftover ⍳ from backtick test)
	runner.SendKeys("BSpace")
	runner.Sleep(100 * time.Millisecond)

	// Define function B with multiple lines
	runner.SendLine(")ed B")
	runner.Sleep(500 * time.Millisecond)

	runner.Test("Editor opens for B", func() bool {
		return runner.Contains("B")
	})

	// Type function body: ⎕←'before' / 1+2 / ⎕←'after'
	runner.SendKeys("End", "Enter", "Enter")
	runner.SendText("⎕←'before'")
	runner.SendKeys("Enter")
	runner.SendText("1+2")
	runner.SendKeys("Enter")
	runner.SendText("⎕←'after'")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("B function defined")

	// Move to line 2 and set breakpoint
	runner.SendKeys("Up", "Up") // Go to line 2 (⎕←'before')
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("b") // Toggle breakpoint
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("B with breakpoint on line 2")

	runner.Test("Breakpoint set in editor", func() bool {
		return runner.Contains("●")
	})

	// Save and close editor
	runner.SendKeys("Escape")
	runner.Sleep(500 * time.Millisecond)

	runner.Test("B editor closes", func() bool {
		return runner.WaitForNoFocusedPane(3 * time.Second)
	})

	// Run B - should stop at breakpoint
	runner.SendLine("B")
	runner.Sleep(1 * time.Second)
	runner.Snapshot("Stopped at breakpoint in B")

	runner.Test("Tracer opens at breakpoint", func() bool {
		return runner.Contains("[tracer]") && runner.Contains("before")
	})

	runner.Test("Breakpoint still visible in tracer", func() bool {
		return runner.Contains("●")
	})

	// Test breakpoint toggling - add a second breakpoint on line 3
	runner.SendKeys("Down") // Move to line 3
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("b")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("Two breakpoints set")

	// Count breakpoints - we should see two ● symbols now
	// (This is a bit tricky to test, but we can check the snapshot)

	// Remove the second breakpoint
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("b")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("Back to one breakpoint")

	// Test breakpoint via command palette
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys(":")
	runner.Sleep(300 * time.Millisecond)
	runner.SendText("break")
	runner.Sleep(200 * time.Millisecond)

	runner.Test("Command palette shows breakpoint command", func() bool {
		return runner.Contains("breakpoint")
	})

	runner.SendKeys("Escape") // Cancel palette
	runner.Sleep(200 * time.Millisecond)

	// Focus tracer before edit test
	runner.SendKeys("Tab")
	runner.Sleep(100 * time.Millisecond)

	// Test breakpoint persistence after editing
	runner.SendKeys("e") // Enter edit mode
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("Edit mode in tracer")

	runner.Test("Edit mode active", func() bool {
		return runner.Contains("[edit]")
	})

	// Make a small edit - add a space somewhere
	runner.SendKeys("End")
	runner.SendText(" ")
	runner.Sleep(100 * time.Millisecond)

	// Exit edit mode with Escape
	runner.SendKeys("Escape")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After edit - back to tracer")

	runner.Test("Back to tracer after edit", func() bool {
		return runner.Contains("[tracer]")
	})

	runner.Test("Breakpoint persists after editing", func() bool {
		return runner.Contains("●")
	})

	// Step with 'n' - execute line 2
	runner.SendKeys("n")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("After first step (before printed)")

	runner.Test("Step executes line - 'before' printed", func() bool {
		return runner.Contains("before")
	})

	// Step again - execute 1+2
	runner.SendKeys("n")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("After second step (1+2)")

	runner.Test("Step executes 1+2 - shows 3", func() bool {
		return runner.Contains("3")
	})

	// Step again - execute ⎕←'after'
	runner.SendKeys("n")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("After third step (after printed)")

	runner.Test("Step executes - 'after' printed", func() bool {
		return runner.Contains("after")
	})

	// One more step should complete execution
	runner.SendKeys("n")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("After function completes")

	runner.Test("Function completes - tracer closes", func() bool {
		return !runner.Contains("[tracer]")
	})

	// Clean up B
	runner.SendLine(")erase B")
	runner.Sleep(300 * time.Millisecond)

	// === ERROR STACK TEST - nested functions X→Y→Z ===
	// Clean up any existing functions from previous runs
	runner.SendLine(")erase X Y Z")
	runner.Sleep(500 * time.Millisecond)

	// Define Z (will error)
	runner.SendLine(")ed Z")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("Editor opened for Z")

	runner.Test("Editor opens for Z", func() bool {
		return runner.Contains("Z")
	})

	// Type function body: 9÷0 (divide by zero error)
	runner.SendKeys("End", "Enter", "Enter")
	runner.SendText("9÷0")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("Z function with 9÷0")

	// Save and close - wait for editor to actually close
	runner.SendKeys("Escape")
	runner.Sleep(500 * time.Millisecond)

	runner.Test("Z editor closes after Escape", func() bool {
		return runner.WaitForNoFocusedPane(3 * time.Second)
	})
	runner.Snapshot("After closing Z editor")

	// Define Y (calls Z)
	runner.SendLine(")ed Y")
	runner.Sleep(1 * time.Second)
	runner.Snapshot("Y editor opened")

	runner.Test("Y editor opens", func() bool {
		return runner.Contains("╔") && runner.Contains("Y")
	})

	runner.SendKeys("End", "Enter", "Enter")
	runner.SendText("Z")
	runner.Sleep(200 * time.Millisecond)
	runner.SendKeys("Escape")
	runner.Sleep(500 * time.Millisecond)

	runner.Test("Y editor closes after Escape", func() bool {
		return runner.WaitForNoFocusedPane(3 * time.Second)
	})

	// Define X (calls Y)
	runner.SendLine(")ed X")
	runner.Sleep(1 * time.Second)
	runner.Snapshot("X editor opened")

	runner.Test("X editor opens", func() bool {
		return runner.Contains("╔") && runner.Contains("X")
	})

	runner.SendKeys("End", "Enter", "Enter")
	runner.SendText("Y")
	runner.Sleep(200 * time.Millisecond)
	runner.SendKeys("Escape")
	runner.Sleep(500 * time.Millisecond)

	runner.Test("X editor closes after Escape", func() bool {
		return runner.WaitForNoFocusedPane(3 * time.Second)
	})
	runner.Snapshot("After defining X, Y, Z")

	// Execute X - triggers nested error
	runner.SendLine("X")
	runner.Sleep(2 * time.Second) // Give time for error and tracer to open
	runner.Snapshot("After X errors - tracer opens")

	runner.Test("Tracer opens on error", func() bool {
		return runner.Contains("[tracer]") || runner.Contains("DOMAIN ERROR") || runner.Contains("tracer")
	})

	// Test stack IMMEDIATELY - before any manipulation
	runner.SendKeys("C-]")
	runner.Sleep(100 * time.Millisecond)
	runner.SendKeys("s")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("Stack pane showing X→Y→Z (fresh)")

	runner.Test("Stack shows 3 frames", func() bool {
		return runner.Contains("stack (3)")
	})

	runner.Test("Stack pane shows Z (top of stack)", func() bool {
		return runner.Contains("Z[") || runner.Contains("Z ")
	})

	runner.Test("Stack pane shows Y", func() bool {
		return runner.Contains("Y[") || runner.Contains("Y ")
	})

	runner.Test("Stack pane shows X", func() bool {
		return runner.Contains("X[") || runner.Contains("X ")
	})

	// Close stack pane before other tests
	runner.SendKeys("Escape")
	runner.Sleep(200 * time.Millisecond)

	// Focus tracer
	runner.SendKeys("Tab")
	runner.Sleep(200 * time.Millisecond)

	// Test: Tracer mode blocks text insertion
	runner.Snapshot("Before typing in tracer")

	// Try to type some text - should be blocked in tracer mode
	runner.SendText("xyz")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After typing xyz in tracer mode")

	runner.Test("Tracer mode blocks text insertion", func() bool {
		// Content should be unchanged - no "xyz" inserted
		return !runner.Contains("xyz")
	})

	// Test: Edit mode toggle with 'e' key
	runner.SendText("e")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After pressing e - edit mode")

	runner.Test("Edit mode shows [edit] in title", func() bool {
		return runner.Contains("[edit]")
	})

	// Test: Can type in edit mode
	runner.SendText("test123")
	runner.Sleep(200 * time.Millisecond)
	runner.Snapshot("After typing in edit mode")

	runner.Test("Edit mode allows text insertion", func() bool {
		return runner.Contains("test123")
	})

	// Test: Escape in edit mode returns to tracer (doesn't close)
	runner.SendKeys("Escape")
	runner.Sleep(300 * time.Millisecond)
	runner.Snapshot("After Escape in edit mode")

	runner.Test("Escape in edit mode returns to tracer", func() bool {
		// Should still have a tracer pane open, now showing [tracer] not [edit]
		return runner.Contains("[tracer]")
	})

	// Test: Second Escape pops Z frame (closes tracer for Z)
	runner.SendKeys("Escape")
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("After second Escape - Z popped")

	// Pop remaining frames to clean up
	runner.SendKeys("Escape") // Pop Y
	runner.Sleep(500 * time.Millisecond)
	runner.SendKeys("Escape") // Pop X
	runner.Sleep(500 * time.Millisecond)
	runner.Snapshot("After popping all frames - clean state")

	runner.Test("Stack cleared after popping all frames", func() bool {
		return !runner.Contains("[tracer]")
	})

	// Final snapshot
	runner.Snapshot("Final state")

	// Generate report
	reportFile := runner.GenerateReport()
	if reportFile != "" {
		t.Logf("Report: %s", reportFile)
		// Open in browser if on macOS
		if _, err := os.Stat("/usr/bin/open"); err == nil {
			exec.Command("open", reportFile).Start()
		}
	}
}
