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

# Skip documentation-only commits.
rev=$(git rev-parse HEAD)
while :; do
	if git rev-parse HEAD^2 >/dev/null 2>&1; then
		# This is a merge commit. Cannot walk up further.
		break
	fi
	parent=$(git rev-parse "$rev"^)
	if ! git diff --quiet "$parent" HEAD -- . ':!docs'; then
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
		echo >&2 "Assuming you're building from a tarball and building fake version info."
		gitcount=0
		gittag=v0.0
		commits=0
		hash=unknown
		date=$(date +%Y%m%d)
		;;
esac

case "$gittag" in
	v*.*-*)
		prerelease=-${gittag##*-}
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
			*)
				echo >&2 "Invalid prerelease name: $prerelease."
				exit 1
				;;
		esac
		gitver=${gittag%-*}
		;;
	v*.*)
		prerelease=
		prerelease_add=30000
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

case "$format" in
	semver)
		case "$prerelease" in
			'')
				echo "$major.$minor.$patch+$date.$commits.$hash"
				;;
			-*)
				echo "$major.$minor.0$prerelease.$patch+$date.$commits.$hash"
				;;
			*)
				echo >&2 "Internal error - invalid parsed prerelease version: $prerelease."
				exit 1
				;;
		esac
		;;
	macos)
		case "$prerelease" in
			'')
				echo "$major.$minor.$patch"
				;;
			-*)
				echo "$major.$minor.0${prerelease#-}$patch"
				;;
			*)
				echo >&2 "Internal error - invalid parsed prerelease version: $prerelease."
				exit 1
				;;
		esac
		;;
	windows)
		echo "$major.$minor.$((patch + prerelease_add)).$commits"
		;;
	gittag)
		case "$prerelease" in
			'')
				echo "v$major.$minor.$patch"
				;;
			-*)
				echo "v$major.$minor.0$prerelease.$patch"
				;;
			*)
				echo >&2 "Internal error - invalid parsed prerelease version: $prerelease."
				exit 1
				;;
		esac
		;;
esac
