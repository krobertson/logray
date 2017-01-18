// Copyright 2012-2016 Apcera Inc. All rights reserved.

package logray

import (
	"net/url"
	"testing"
)

const (
	formatString = "%color:class% %classfixed% " +
		"%year%-%month%-%day% %hour%:%minute%:%second%.%nanosecond% " +
		"%tzoffset% %tz% pid=%pid%"
)

func TestLoggerOutputsDownLevel(t *testing.T) {

	ResetDefaultOutput()
	// Create new URI for default stdout output.
	defaultUrlString := "stdout://?format=" + url.QueryEscape(formatString)
	// Append one extra field 'tid' to the url.
	newDefaultUrlString := defaultUrlString + url.QueryEscape(" tid='%field:tid%'")

	// Add default stdout output.
	if err := AddDefaultOutput(defaultUrlString, ALL); err != nil {
		t.Fatal(err)
	}

	logger := New()

	// Update for DEBUG should not 	update the output.
	logger.UpdateOutput(newDefaultUrlString, DEBUG)
	if len(logger.outputs) > 2 {
		t.Fatalf("More than 2 logger outputs defined: %d", len(logger.outputs))
	}
	debugOutput := logger.outputs[0]
	restOutput := logger.outputs[1]

	defaultUrl, err := url.Parse(defaultUrlString)
	if err != nil {
		t.Fatal(err)
	}

	newUrl, err := url.Parse(newDefaultUrlString)
	if err != nil {
		t.Fatal(err)
	}

	// Must have new URI in debug output only.
	if debugOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		debugOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", debugOutput.OutputWrapper)
	}

	if restOutput.OutputWrapper.URL.Scheme != defaultUrl.Scheme ||
		restOutput.OutputWrapper.URL.RawQuery != defaultUrl.RawQuery {
		t.Fatalf("Output: %v", restOutput.OutputWrapper)
	}

	// Now update for ALL.
	logger.UpdateOutput(newDefaultUrlString, ALL)

	if len(logger.outputs) > 2 {
		t.Fatalf("More than 2 logger outputs defined: %d", len(logger.outputs))
	}

	debugOutput = logger.outputs[0]
	restOutput = logger.outputs[1]

	if debugOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		debugOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", debugOutput.OutputWrapper)
	}

	if restOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		restOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", restOutput.OutputWrapper)
	}

}

func TestLoggerOutputsUpLevel(t *testing.T) {

	ResetDefaultOutput()

	// Create new URI for default stdout output.
	defaultUrlString := "stdout://?format=" + url.QueryEscape(formatString)
	// Append one extra field 'tid' to the url.
	newDefaultUrlString := defaultUrlString + url.QueryEscape(" tid='%field:tid%'")

	// Add default stdout output.
	if err := AddDefaultOutput(defaultUrlString, DEBUGPLUS); err != nil {
		t.Fatal(err)
	}

	logger := New()

	// Update for TRACE should not update the output.
	logger.UpdateOutput(newDefaultUrlString, TRACE)
	if len(logger.outputs) > 1 {
		t.Fatalf("More than 1 logger outputs defined: %d", len(logger.outputs))
	}
	debugPlusOutput := logger.outputs[0]

	defaultUrl, err := url.Parse(defaultUrlString)
	if err != nil {
		t.Fatal(err)
	}

	newUrl, err := url.Parse(newDefaultUrlString)
	if err != nil {
		t.Fatal(err)
	}

	// Must have old URI in debug plus output.
	if debugPlusOutput.OutputWrapper.URL.Scheme != defaultUrl.Scheme ||
		debugPlusOutput.OutputWrapper.URL.RawQuery != defaultUrl.RawQuery {
		t.Fatalf("Output: %v", debugPlusOutput.OutputWrapper)
	}

	// Now update for ALL.
	logger.UpdateOutput(newDefaultUrlString, ALL)

	if len(logger.outputs) > 1 {
		t.Fatalf("More than 1 logger outputs defined: %d", len(logger.outputs))
	}

	debugPlusOutput = logger.outputs[0]

	// Must have new URL in the output.
	if debugPlusOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		debugPlusOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", debugPlusOutput.OutputWrapper)
	}
}

func TestLoggerOutputsMultiLevel(t *testing.T) {

	ResetDefaultOutput()

	// Create new URI for default stdout output.
	defaultUrlString := "stdout://?format=" + url.QueryEscape(formatString)
	// Append one extra field 'tid' to the url.
	newDefaultUrlString := defaultUrlString + url.QueryEscape(" tid='%field:tid%'")

	// Add default stdout output.
	if err := AddDefaultOutput(defaultUrlString, DEBUGPLUS); err != nil {
		t.Fatal(err)
	}

	// Add default stdout output.
	if err := AddDefaultOutput(defaultUrlString, TRACE); err != nil {
		t.Fatal(err)
	}

	logger := New()

	// Update for TRACE should not update the output.
	logger.UpdateOutput(newDefaultUrlString, TRACE)
	if len(logger.outputs) > 2 {
		t.Fatalf("More than 2 logger outputs defined: %d", len(logger.outputs))
	}
	debugPlusOutput := logger.outputs[0]
	traceOutput := logger.outputs[1]

	defaultUrl, err := url.Parse(defaultUrlString)
	if err != nil {
		t.Fatal(err)
	}

	newUrl, err := url.Parse(newDefaultUrlString)
	if err != nil {
		t.Fatal(err)
	}

	// Must have old URI in debug plus output.
	if debugPlusOutput.OutputWrapper.URL.Scheme != defaultUrl.Scheme ||
		debugPlusOutput.OutputWrapper.URL.RawQuery != defaultUrl.RawQuery {
		t.Fatalf("Output: %v", debugPlusOutput.OutputWrapper)
	}

	// Must Have new URI in trace output.
	if traceOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		traceOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", debugPlusOutput.OutputWrapper)
	}

	// Now update for DEBUG.
	logger.UpdateOutput(newDefaultUrlString, DEBUG)

	if len(logger.outputs) > 3 {
		t.Fatalf("More than 3 logger outputs defined: %d", len(logger.outputs))
	}

	debugOutput := logger.outputs[0]
	traceOutput = logger.outputs[1]
	restOutput := logger.outputs[2]

	// Must have new URL in the debug output.
	if debugOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		debugOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", debugOutput.OutputWrapper)
	}

	if traceOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		traceOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", traceOutput.OutputWrapper)
	}

	if restOutput.OutputWrapper.URL.Scheme != defaultUrl.Scheme ||
		restOutput.OutputWrapper.URL.RawQuery != defaultUrl.RawQuery {
		t.Fatalf("Output: %v", restOutput.OutputWrapper)
	}

	// Now update for ALL.
	logger.UpdateOutput(newDefaultUrlString, ALL)

	if len(logger.outputs) > 3 {
		t.Fatalf("More than 3 logger outputs defined: %d", len(logger.outputs))
	}

	debugOutput = logger.outputs[0]
	traceOutput = logger.outputs[1]
	restOutput = logger.outputs[2]

	// Must have new URL in the debug output.
	if debugOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		debugOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", debugOutput.OutputWrapper)
	}

	if traceOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		traceOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", traceOutput.OutputWrapper)
	}

	if restOutput.OutputWrapper.URL.Scheme != newUrl.Scheme ||
		restOutput.OutputWrapper.URL.RawQuery != newUrl.RawQuery {
		t.Fatalf("Output: %v", restOutput.OutputWrapper)
	}
}
