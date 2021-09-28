#!/bin/sh
# Copyright 2021 Google LLC
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

: ${GO:=go}

cat <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<component type="desktop-application">
  <id>io.github.divverent.aaaaxy-$($GO env GOOS)-$($GO env GOARCH)</id>
  <name>AAAAXY</name>
  <summary>A nonlinear puzzle platformer taking place in non-Euclidean geometry.</summary>
  <metadata_license>*** TODO ***</metadata_license>
  <project_license>Apache-2.0</project_license>
  <requires>
    <control>keyboard</control>
    <control>gamepad</control>
  </requires>
  <description>
    <p>
      Although your general goal is reaching the surprising end of the game, you are encouraged to set your own goals while playing. Exploration will be rewarded, and secrets await you!
    </p>
    <p>
      So jump and run around, and enjoy losing your sense of orientation in this World of Wicked Weirdness. Find out what Van Vlijmen will make you do. Pick a path, get inside a Klein Bottle, recognize some memes, and by all means: don&apos;t look up.
    </p>
  </description>
  <launchable type="desktop-id">aaaaxy-$($GO env GOOS)-$($GO env GOARCH)</launchable>
  <screenshots>
    <screenshot type="default">
      <image>https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot1.png</image>
    </screenshot>
    <screenshot>
      <image>https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot5.png</image>
    </screenshot>
    <screenshot>
      <image>https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot8.png</image>
    </screenshot>
  </screenshots>
</component>
EOF
