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
gitdesc=${2:-$(git describe --always --long --match 'v*')}
commits=${3:-$(($(git log --oneline | wc -l)))}  # Is there a better way?

hash=${gitdesc##*-g}
date=$(git log -n 1 --pretty=format:%cd --date=format:%Y%m%d "$hash")

case "$gitdesc" in
	v*.*-*-g*)
		gitcount=${gitdesc%-g*}
		gitcount=${gitcount##*-}
		gittag=${gitdesc%-*-g*}
		;;
	*)
		echo >&2 "Invalid git describe output: $gitdesc."
		exit 1
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
	windows)
		echo "$major.$minor.$((patch + prerelease_add)).$commits"
		;;
esac
