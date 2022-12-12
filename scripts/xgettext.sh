#!/bin/sh
# Copyright 2022 Google LLC
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

xgettext --its=scripts/tmx.its --from-code=utf-8 -F \
	-o assets/locales/level.pot assets/maps/level.tmx
go run github.com/leonelquinteros/gotext/cli/xgotext \
	-default game \
	-in internal/ \
	-out assets/locales/

LF='
'
all_languages=
bad_languages=

for domain in level game; do
	for f in assets/locales/*/"$domain".po; do
		language=${f%/*}
		language=${language##*/}
		msgmerge -U "$f" assets/locales/"$domain".pot
		total=$(grep -c '^#:' "$f")
		untranslated=$(msgattrib --untranslated "$f" | grep -c '^#:')
		fuzzy=$(msgattrib --only-fuzzy "$f" | grep -c '^#:')
		score=$(((total - untranslated - fuzzy) * 100 / total))
		echo "$f: $score%: $untranslated/$total untranslated, $fuzzy/$total fuzzy"
		all_languages="$all_languages$language$LF"
		if [ $score -lt 90 ]; then
			bad_languages="$bad_languages$language$LF"
		fi
	done
done

good_languages=$(
	{
		echo "$all_languages" | sort -u
		echo "$bad_languages"
	} | sort | uniq -u
)
echo "Good languages:" $good_languages
