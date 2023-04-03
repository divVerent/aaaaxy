// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
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

//go:build ios
// +build ios

package fun

import (
	"time"

	m "github.com/divVerent/aaaaxy/internal/math"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/NSCalendar.h>
#import <Foundation/NSDate.h>
#import <Foundation/Foundation.h>

int time_zone_hours() {
	NSCalendar *cal = [NSCalendar currentCalendar];
	NSDate *now = [NSDate now];
	NSDate *jan1 = [cal
		dateWithEra: [cal component:NSCalendarUnitEra fromDate: now]
		year: [cal component:NSCalendarUnitYear fromDate: now]
		month: 1 day: 1 hour: 0 minute: 0 second: 0 nanosecond: 0];
	NSInteger secs = [[cal timeZone] secondsFromGMTForDate: jan1];
	if (secs < 0) {
		// Make sure to still round down.
		return ~((~secs) / 3600);
	} else {
		return secs / 3600;
	}
}
*/
import "C"

func init() {
	SetTimeZoneHours(C.time_zone_hours())
}
