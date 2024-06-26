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

format=${1:-semver}

if [ -d .git ] && [ x"$AAAAXY_BUILD_USE_VERSION_FILE" != x'true' ]; then
	# Skip documentation-only commits.
	rev=$(git rev-parse HEAD)
	while :; do
		if git rev-parse HEAD^2 >/dev/null 2>&1; then
			# This is a merge commit. Cannot walk up further.
			break
		fi
		parent=$(git rev-parse "$rev"^)
		if ! git diff --quiet "$parent" HEAD -- . ':!docs' ':!io.github.divverent.aaaaxy.metainfo.xml' ':!third_party/SDL_GameControllerDB/assets/input' ':!.gitmoduleversions' ':!.lastreleaseversion' ':!.lastreleaseversioncode' ':!fastlane/metadata/android/en-US/changelogs' ':!XcodeProjects/iOS/aaaaxy.xcodeproj/project.pbxproj'; then
			# Non-doc changes have been found. Do not walk up further.
			break
		fi
		rev=$parent
	done

	gitdesc=${2:-$(git describe --always --long --match 'v*.*' --exclude 'v*.*.*' "$rev")}

	case "$gitdesc" in
		v*.*-*-g*)
			gitcount=${gitdesc%-g*}
			gitcount=${gitcount##*-}
			gittag=${gitdesc%-*-g*}
			commits=${3:-$(($(git log --oneline "$rev" | wc -l)))}  # Is there a better way?
			hash=${gitdesc##*-g}
			date=$(git log -n 1 --pretty=format:%cd --date=format:%Y%m%d "$hash")
			;;
		*)
			echo >&2 "ERROR: Invalid git describe output: $gitdesc."
			echo >&2 "You may provide the version in a file called .lastreleaseversion"
			echo >&2 "and then set AAAAXY_BUILD_USE_VERSION_FILE=true."
			exit 1
			;;
	esac

	case "$gittag" in
		v*.*-*)
			prerelease=-${gittag##*-}
			gitver=${gittag%-*}
			;;
		v*.*)
			prerelease=
			gitver=$gittag
			;;
		*)
			echo >&2 "Invalid version tag: $gitver."
			exit 1
			;;
	esac

	case "$gitver" in
		v*.*)
			major=${gitver%.*}
			major=${major#v}
			minor=${gitver#v*.}
			patch=$gitcount
			;;
		*)
			echo >&2 "Internal error - invalid parsed git version: $gitver."
			exit 1
			;;
	esac
elif [ -f .lastreleaseversion ]; then
	# Re-import from a semver string.
	# Used for packaging source code.
	save_IFS=$IFS
	IFS='+.'
	set -- $(cat .lastreleaseversion)
	IFS=$save_IFS
	case "$#" in
		6)
			major=$1
			minor=$2
			prerelease=
			patch=$3
			date=$4
			commits=$5
			hash=$6
			;;
		7)
			major=$1
			minor=$2
			prerelease=${3#0}
			patch=$4
			date=$5
			commits=$6
			hash=$7
			;;
		*)
			echo >&2 "Internal error - failed to parse .lastreleaseversion file."
			exit 1
			;;
	esac
	echo >&2 "NOTE: version imported from .lastreleaseversion file."
	echo >&2 "NOTE: when building from a git clone, this message should not show up."
else
	echo >&2 "This script must be called from the root of the AAAAXY source code."
	exit 1
fi

# Set of variables here:
# - major
# - minor
# - prerelease
# - patch
# - date (YYYYMMDD)
# - commits
# - hash

case "$prerelease" in
	-alpha)
		prerelease_add=0
		;;
	-beta)
		prerelease_add=10000
		;;
	-rc)
		prerelease_add=20000
		;;
	'')
		prerelease_add=30000
		;;
	*)
		echo >&2 "Invalid prerelease name: $prerelease."
		exit 1
		;;
esac

case "$format" in
	semver)
		case "$prerelease" in
			'')
				printf '%d.%d.%d+%s.%d.%s\n' "$major" "$minor" "$patch" "$date" "$commits" "$hash"
				;;
			-*)
				printf '%d.%d.0%s.%d+%s.%d.%s\n' "$major" "$minor" "$prerelease" "$patch" "$date" "$commits" "$hash"
				;;
			*)
				echo >&2 "Internal error - invalid parsed prerelease version: $prerelease."
				exit 1
				;;
		esac
		;;
	macos|ios)
		case "$prerelease" in
			'')
				printf '%d.%d.%d\n' "$major" "$minor" "$patch"
				;;
			-*)
				printf '%d.%d.0%s%d\n' "$major" "$minor" "${prerelease#-}" "$patch"
				;;
			*)
				echo >&2 "Internal error - invalid parsed prerelease version: $prerelease."
				exit 1
				;;
		esac
		;;
	windows)
		printf '%d.%d.%d.%d\n' \
			"$major" \
			"$minor" \
			"$((patch + prerelease_add))" \
			"$commits"
		;;
	gittag)
		case "$prerelease" in
			'')
				printf 'v%d.%d.%d\n' "$major" "$minor" "$patch"
				;;
			-*)
				printf 'v%d.%d.0%s.%d\n' "$major" "$minor" "$prerelease" "$patch"
				;;
			*)
				echo >&2 "Internal error - invalid parsed prerelease version: $prerelease."
				exit 1
				;;
		esac
		;;
	android)
		# Building an Android version code out of the semver, similar to Windows binary versions.
		# <major><minor:%02d>3<patch:%04d>0
		# Last digit forced to zero to allow F-Droid to extend for
		# split packages.
		printf '%d%02d%05d%d\n' \
			"$major" \
			"$minor" \
			"$((patch + prerelease_add))" \
			"${AAAAXY_BUILD_VERSIONCODE_ADD:-0}"
		;;
	*)
		echo >&2 "Unknown version type: $format."
		;;
esac
