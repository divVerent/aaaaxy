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
	"log"
	"os"
	"runtime/debug"

	"github.com/divVerent/aaaaxy/internal/alert"
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

func Debugf(format string, v ...interface{}) {
	if *V < debugLevel {
		return
	}
	log.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
}

func Infof(format string, v ...interface{}) {
	if *V < infoLevel {
		return
	}
	log.Output(2, fmt.Sprintf("[INFO] "+format, v...))
}

func Warningf(format string, v ...interface{}) {
	if *V < warningLevel {
		return
	}
	log.Output(2, fmt.Sprintf("[WARNING] "+format, v...))
}

func Errorf(format string, v ...interface{}) {
	if *V < errorLevel {
		return
	}
	log.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
}

func Fatalf(format string, v ...interface{}) {
	debug.PrintStack()
	msg := fmt.Sprintf(format, v...)
	log.Output(2, "[FATAL] "+msg)
	if !*Batch {
		alert.Show(msg)
	}
	os.Exit(125)
}
