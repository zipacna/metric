// Copyright (c) 2020 Jean Mattes (risingcode.net)
// MIT-License (see LICENSE file)
// This package is a collection of helper functions I wrote while learning Go.
package metric

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

var Started int64

// Takes a snapshot of Unix time in nanoseconds at package init.
// As well as doing init for the logger.
func init() {
	Started = time.Now().UnixNano()

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

// Prints the total runtime of the program (as timestamp difference).
// Usage: call `defer metric.Stop()` at the start of the main function.
func Stop() {
	var stopped = float64(time.Now().UnixNano() - Started)
	var unit10e9 = math.Pow(10.0, 9.0)
	fmt.Printf("ROP: %f sec", stopped/unit10e9)
}

// There is no Ternary operator in Go (if condition true return positive, otherwise negative).
func Ternary(cond bool, positive interface{}, negative interface{}) (result interface{}) {
	if cond {
		return positive
	} else {
		return negative
	}
}

// LogData contains all necessary information to properly handle an error (used as input for CheckError).
// `metric.CheckError(metric.LogData{Err: err1, Tag: "Malformed", UseLogger: true, Severity: "fatal"})`
// `metric.CheckError(metric.LogData{Err: err6, Tag: "Close", UseLogger: true, Severity: "print"})`
// The first example could be the first (or second if counting from zero) error to catch within the control flow,
// the tag "Malformed" could indicate some formatting error; Severity goes from print over fatal to panic.
// The second example could be the sixth (or seventh) error to catch, Close (e.g. closing operation) & prints only.
type LogData struct {
	Err      error
	Tag      string
	Severity string
}

// CheckError sums up the error printing/logging for reuse.
// Call it with the appropriate LogData to simplify error catching and logging.
func CheckError(data LogData) {
	if data.Err != nil {
		var errString = fmt.Sprintf("%s error: %s\n", data.Tag, data.Err)
		switch data.Severity {
		case "print":
			log.Printf(errString)
		case "fatal":
			log.Fatalf("Last error was fatal ... %s", errString)
		case "panic":
			log.Panicf("Last error caused a panic ... %s", errString)
		}
	}
}
