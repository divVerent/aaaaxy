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

xgettext --its=scripts/tmx.its --from-code=utf-8 -F --no-location \
	-o - assets/maps/level.tmx |\
	sed -e 's/^#. #:/#:/g' \
	> assets/locales/level.pot
go run github.com/leonelquinteros/gotext/cli/xgotext \
	-default game_raw \
	-in internal/ \
	-out assets/locales/

# Passing in game_raw.pot in two different ways so that
# messages that just have a comment are not emitted.
# This just outputs any obsolete strings in game.pot.comments.
cat assets/locales/game_raw.pot |\
msgcomm \
	--less-than=2 \
	--more-than=0 \
	-s \
	assets/locales/game.pot.comments \
	assets/locales/game_raw.pot \
	-

# Passing in game_raw.pot in two different ways so that
# messages that just have a comment are not emitted.
cat assets/locales/game_raw.pot |\
msgcomm \
	--less-than=999999 \
	--more-than=1 \
	-s \
	assets/locales/game.pot.comments \
	assets/locales/game_raw.pot \
	- |\
	sed -e '/^#: / s/:[1-9][0-9]*//g' \
	> assets/locales/game.pot

LF='
'
all_linguas=
bad_linguas=

for d in assets/locales/*/; do
	language=${d%/}
	language=${language##*/}
	# Go's x/text/language always uses dashes as separator.
	lingua=$(echo "$language" | tr _ -)
	all_linguas="$all_linguas$lingua$LF"
	for domain in level game; do
		f=assets/locales/"$language"/"$domain".po
		if ! [ -f "$f" ]; then
			echo "$f: not found"
			bad_linguas="$bad_linguas$lingua$LF"
			continue
		fi
		msgmerge -o "$f.new" "$f" assets/locales/"$domain".pot
		total=$(grep -c '^#:' "$f.new")
		untranslated=$(msgattrib --untranslated "$f.new" | grep -c '^#:')
		fuzzy=$(msgattrib --only-fuzzy "$f.new" | grep -c '^#:')
		score=$(((total - untranslated - fuzzy) * 100 / total))
		echo "$f: $score%: $untranslated/$total untranslated, $fuzzy/$total fuzzy"
		if [ $score -lt 90 ]; then
			bad_linguas="$bad_linguas$lingua$LF"
		fi
	done
done

good_lingua=$(
	{
		echo "$all_linguas"
		echo "$bad_linguas"
		echo "$bad_linguas"
	} | sort | uniq -u
)
echo "Good languages:" $good_lingua

echo "$good_lingua" > assets/locales/LINGUAS

make
languages=$(
	printf "'en'"
	xvfb-run ./aaaaxy -dump_languages |\
	grep . |\
	while IFS=- read -r lang variant; do
		case "$lang-$variant" in
			be@tarask-)
				res=b+be+Latn
				;;
			zh-Hans)
				res=b+zh+Hans
				;;
			zh-Hant)
				res=b+zh+Hant
				;;
			*)
				res=$lang${variant:+-r$variant}
				;;
		esac
		printf ", '%s'" "$res"
	done
)
sed -i -e "s/resConfigs .*/resConfigs $languages/" \
	./AndroidStudioProjects/AAAAXY/app/build.gradle
