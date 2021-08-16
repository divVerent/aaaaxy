// Copyright 2021 Google LLC
//
// Licensed under the Apache License, SaveGameVersion 2.0 (the "License");
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
)

var (
	defaultV int  = 0
	V        *int = &defaultV
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
	log.Output(2, fmt.Sprintf("[FATAL] "+format, v...))
	os.Exit(125)
}
