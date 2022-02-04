// Copyright 2021 Google LLC
//
// Licensed under the Apache Livense, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/divVerent/aaaaxy/internal/alert"
	"github.com/divVerent/aaaaxy/internal/atexit"
)

var (
	defaultV     int   = 0
	V            *int  = &defaultV
	defaultBatch bool  = false
	Batch        *bool = &defaultBatch
)

const (
	debugLevel   = 1
	infoLevel    = 0
	warningLevel = -1
	errorLevel   = -2
)

// Do not allow exploiting log parsers or the terminal. #log4j
var (
	newline            = "\n"
	newlineReplacement = "\n\t"
	// Note: explicitly allow the special Unicode characters that we actually use.
	// Do not allow anything else, as it could be nasty stuff like overrides or invisible stuff.
	specialRE          = regexp.MustCompile(`[^\n\t\040-\176µτπö¾©]|\\`)
	specialReplacement = "?"
)

func logSprintf(format string, v ...interface{}) string {
	return specialRE.ReplaceAllStringFunc(strings.ReplaceAll(fmt.Sprintf(format, v...), newline, newlineReplacement), func(s string) string {
		repl := ""
		for _, ch := range []byte(s) {
			repl += fmt.Sprintf("\\%03o", ch)
		}
		return repl
	})
}

func Debugf(format string, v ...interface{}) {
	if *V < debugLevel {
		return
	}
	log.Output(2, logSprintf("[DEBUG] "+format, v...))
}

func Infof(format string, v ...interface{}) {
	if *V < infoLevel {
		return
	}
	log.Output(2, logSprintf("[INFO] "+format, v...))
}

func Warningf(format string, v ...interface{}) {
	if *V < warningLevel {
		return
	}
	log.Output(2, logSprintf("[WARNING] "+format, v...))
}

func Errorf(format string, v ...interface{}) {
	if *V < errorLevel {
		return
	}
	log.Output(2, logSprintf("[ERROR] "+format, v...))
}

func TraceErrorf(format string, v ...interface{}) {
	if *V < errorLevel {
		return
	}
	debug.PrintStack()
	log.Output(2, logSprintf("[ERROR] "+format, v...))
}

func Fatalf(format string, v ...interface{}) {
	debug.PrintStack()
	msg := logSprintf(format, v...)
	log.Output(2, "[FATAL] "+msg)
	CloseLogFile()
	if !*Batch {
		alert.Show(msg)
	}
	atexit.Finish()
	os.Exit(125)
}

var (
	logFiles []io.Closer
)

func AddLogFile(file string) {
	wr, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Errorf("failed to open log file: %v", err)
		return
	}
	logFiles = append(logFiles, wr)
	log.SetOutput(io.MultiWriter(log.Writer(), wr))
}

func CloseLogFile() {
	for _, wr := range logFiles {
		err := wr.Close()
		if err != nil {
			Errorf("failed to close log file: %v", err)
		}
	}
	logFiles = nil
}
