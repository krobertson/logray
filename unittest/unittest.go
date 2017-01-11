// Copyright 2012-2014 Apcera Inc. All rights reserved.

// A simple package for dealing with logging module output during the
// unit testing cycle.
package unittest

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/apcera/logray"
)

const (
	// This format string is used to replace the time stamp with a fixed number
	// representing the number of log lines contained in the buffer. We fake
	// this a bit by using epoch and mutating the real time stamp.
	defaultFormatStr = "%color:class%[%hour%:%minute%:%second%.%nanosecond% " +
		"%class% category='%category%' context='%context%'] %message%" +
		"%color:default%"
)

// Interface that matches both testing.T and testing.B.
type Logger interface {
	Failed() bool
}

// This will be used as a Output object as a logging destination. It is
// configured to capture all logs and store them in memory, flushing them
// only if the unit test fails.
type LogBuffer struct {
	// This will store a copy of the LineData objects that the logging library
	// generates in order to commit them later if necessary.
	buffer []*logray.LineData

	// The last read item in the buffer. This is used with NewLines()
	lastRead int

	// Mutex used to protect internal data structures.
	mutex sync.Mutex

	// fields contains all log line fields to be logged.
	fields map[string]interface{}
}

// Tracker implementation of io.Writer to let this object receive logs.
func (l *LogBuffer) Write(ld *logray.LineData) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if len(ld.Fields) > 0 {
		if l.fields == nil {
			l.fields = make(map[string]interface{}, len(ld.Fields))
		}
		for k, v := range ld.Fields {
			l.fields[k] = v
		}
	}
	l.buffer = append(l.buffer, ld)
	return nil
}

// Required to implement the Output interface.
func (l *LogBuffer) Flush() error { return nil }

// Clears the list of lines in the buffer.
func (l *LogBuffer) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.buffer = nil
	l.lastRead = 0
}

// Gets all received lines from the tracker.
func (l *LogBuffer) Lines() []string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	buffer := bytes.NewBuffer(nil)
	// Ignore err, there is nothing we can do with it anyway.
	output, _ := logray.NewIOWriterOutput(buffer, defaultFormatStr, "auto")
	for _, ld := range l.buffer {
		// Ignore errors.. there is nothing we can do about it anyway.
		output.Write(ld)
	}
	return strings.Split(string(buffer.Bytes()), "\n")
}

func (l *LogBuffer) NewLines() []string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	buffer := bytes.NewBuffer(nil)
	// Ignore err, there is nothing we can do with it anyway.
	output, _ := logray.NewIOWriterOutput(buffer, defaultFormatStr, "auto")
	for i, ld := range l.buffer {
		if i < l.lastRead {
			continue
		}
		// Ignore errors.. there is nothing we can do about it anyway.
		output.Write(ld)
	}
	l.lastRead = len(l.buffer)
	lines := strings.Split(string(buffer.Bytes()), "\n")
	// Remove empty last line.
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// Dumps all the lines in the buffer to stdout.
func (l *LogBuffer) DumpToStdout() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	formatStr := defaultFormatStr
	// Add format string part for new fields added to the buffer.
	if len(l.fields) > 0 {
		if ind := strings.Index(formatStr, "]"); ind != -1 {
			prefix := formatStr[:ind]
			fieldsOutput := make([]string, 0)
			for k, v := range l.fields {
				if v != nil {
					fieldsOutput = append(fieldsOutput, fmt.Sprintf(" %s='%%field:%s%%'", k, k))
				}
			}
			middle := strings.Join(fieldsOutput, " ")
			suffix := formatStr[ind:]
			formatStr = fmt.Sprintf("%s%s%s", prefix, middle, suffix)
		}
	}

	output, err := logray.NewIOWriterOutput(os.Stdout, formatStr, "auto")
	if err != nil {
		fmt.Printf("Error from logray.NewIOWriterOutput: %s\n", err)
		return
	}
	for _, ld := range l.buffer {
		if err := output.Write(ld); err != nil {
			fmt.Printf("Error in Write(): %s\n", err)
		}
	}
}

func (l *LogBuffer) DumpToFile(path string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var err error
	var w io.WriteCloser

	if w, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log file %q, outputting on STDERR", path)
		w = os.Stderr
	}

	var output logray.Output
	if output, err = logray.NewIOWriterOutput(w, defaultFormatStr, "auto"); err != nil {
		fmt.Printf("Error from logray.NewIOWriterOutput: %s\n", err)
		return
	}

	for _, ld := range l.buffer {
		if err := output.Write(ld); err != nil {
			fmt.Fprintf(os.Stderr, "Error in Write(): %s\n", err)
		}
	}
}

// Dumps the logs in the log buffer if the test failed, otherwise
// clears them for the next test.
func (l *LogBuffer) FinishTest(t Logger) {
	if t.Failed() {
		l.DumpToStdout()
	}
	l.Clear()
}

// Sets up everything needed to unit test against a LogBuffer object.
// The returned object will use numerical counters rather than dates in
// order to make the output predictable.
func SetupBuffer() *LogBuffer {
	b := new(LogBuffer)
	newOutput := func(u *url.URL) (logray.Output, error) {
		return b, nil
	}
	logray.ResetCachedOutputs()
	logray.ResetDefaultOutput()
	logray.AddNewOutputFunc("logbuffer", newOutput)
	logray.AddDefaultOutput("logbuffer://", logray.ALL)
	return b
}
