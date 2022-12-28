#!/bin/sh
#
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

if [ $# -ne 2 ]; then
	echo >&2 'First, go to https://www.transifex.com/aaaaxy/aaaaxy,'
	echo >&2 'your language, and "Download for use" both game.pot and'
	echo >&2 'level.pot translations.'
	echo >&2
	echo >&2 'Then, call this script as follows:'
	echo >&2 "$0 download-folder language"
	exit 1
fi

dir=$1
lang=$2

prefix=for_use_aaaaxy_assets-locales

mkdir -p assets/locales/test
for domain in game level; do
	cp -v \
		"$dir"/"$prefix"-"$domain"-pot--main_"$lang".po \
		assets/locales/test/"$domain".po
done

cp -v ~/.config/AAAAXY/config.json ~/.config/AAAAXY/config.json.save
make run ARGS='-language=test'
cp -v ~/.config/AAAAXY/config.json.save ~/.config/AAAAXY/config.json

echo >&2
echo >&2 'Please look for lines marked ERROR in above log.'
