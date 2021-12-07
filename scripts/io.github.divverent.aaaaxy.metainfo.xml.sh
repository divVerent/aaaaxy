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

# Note: for AppImage we tag the metainfo by the architecture as these are
# downloaded and installed manually and thus having multiple arch versions
# active at the same time is conceivable; for FlatPak we don't do this as
# FlatPak uses a package manager that ensures there is only one arch installed
# at a time anyway.

sed -e "
	s,io\\.github\\.divverent\\.aaaaxy,io.github.divverent.aaaaxy_$($GO env GOARCH),;
	s,aaaaxy\\.desktop,aaaaxy-$($GO env GOOS)-$($GO env GOARCH).desktop,;
" < io.github.divverent.aaaaxy.metainfo.xml
