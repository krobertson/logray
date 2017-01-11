// Copyright 2012-2013 Apcera Inc. All rights reserved.

package unittest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apcera/logray"
)

func Example() {
	// Setup
	buffer := SetupBuffer()
	defer buffer.DumpToStdout()

	// Unit tests go here.

	// Clear the buffer so nothing ends up being printed.
	buffer.Clear()
}

func ExamplePass() {
	// Setup
	buffer := SetupBuffer()
	defer buffer.DumpToStdout()

	// Log a bunch of stuff.
	logger := logray.New()
	fmt.Println("Expected output.")
	logger.Info("log line 1")
	logger.Info("log line 2")
	logger.Info("log line 3")

	// Clear the buffer so nothing ends up on stdout.
	buffer.Clear()

	// Output: Expected output.
}

func TestLogBufferFields(t *testing.T) {

	logBuffer := SetupBuffer()

	logger := logray.New()
	logger.SetField("TestField1", "Test1")
	logger.SetField("TestField2", "Test2")
	logger.Info("Test Log")

	// Wait for a while and let the logger goroutine run.
	time.Sleep(1 * time.Second)

	if len(logBuffer.buffer) != 1 {
		t.Fatalf("Expected: %d, Got: %d LineData items in the log buffer.", 1, len(logBuffer.buffer))
	}

	if len(logBuffer.fields) != 2 {
		t.Fatalf("Expected: %d, Got: %d fields in the log buffer.", 2, len(logBuffer.fields))
	}

	if fv, ok := logBuffer.fields["TestField1"]; !ok || !strings.EqualFold(fv.(string), "Test1") {
		t.Fatal("Unexpected field found in the log buffer")
	}

	if fv, ok := logBuffer.fields["TestField2"]; !ok || !strings.EqualFold(fv.(string), "Test2") {
		t.Fatal("Unexpected field found in the log buffer")
	}
}
