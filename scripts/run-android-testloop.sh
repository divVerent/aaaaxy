#!/bin/sh
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

# Disable zip packing improvements as these do not matter at runtime.
export ADVZIP=:

adb shell pm uninstall -k io.github.divverent.aaaaxy || true
cd AndroidStudioProjects/AAAAXY/
export ANDROID_HOME=$HOME/Android/Sdk
./gradlew assembleRelease
adb install app/build/outputs/apk/release/app-release.apk

adb logcat -c
adb logcat &
logcat_pid=$!

adb shell am start-activity -n io.github.divverent.aaaaxy/.MainActivity -a com.google.intent.action.TEST_LOOP -W

while adb shell pidof io.github.divverent.aaaaxy >/dev/null; do
	sleep 5
done

kill $logcat_pid
wait
