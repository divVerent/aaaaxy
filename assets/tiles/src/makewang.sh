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

ifboth() {
	case "$1$2" in
		truetrue)
			echo "0 1"
			;;
		*)
			echo "x"
			;;
	esac
}

orsort() {
	all=$(tr _ '\n' | sed -e 's,en,ne,g; s,es,se,g; s,ws,sw,g; s,wn,nw,g;')
	{
		echo "$all" | grep '^n$'
		echo "$all" | grep '^e$'
		echo "$all" | grep '^s$'
		echo "$all" | grep '^w$'
		echo "$all" | grep '^ne'
		echo "$all" | grep '^se'
		echo "$all" | grep '^sw'
		echo "$all" | grep '^nw'
	} | tr '\n' _ | head -c -1
}

rotations() {
	name=$1
	for rotation in ES SW WN NE SE EN NW WS; do
		lcrotation=$(echo "$rotation" | tr A-Z a-z)
		fullrotation=$lcrotation$(echo "$lcrotation" | tr eswn wnes)
		echo -n "$rotation "
		rotated=$(echo "$name" | tr eswn- "$fullrotation-")
		echo "$rotated" | orsort
		echo
	done
}

i=9
for n in false true; do
	for e in false true; do
		for s in false true; do
			for w in false true; do
				for ne in $(ifboth $n$e); do
					for se in $(ifboth $s$e); do
						for sw in $(ifboth $s$w); do
							for  nw in $(ifboth $n$w); do
								set -- base.png
								name=
								$n && set -- "$@" n.png && name=$name'_n'
								$e && set -- "$@" e.png && name=$name'_e'
								$s && set -- "$@" s.png && name=$name'_s'
								$w && set -- "$@" w.png && name=$name'_w'
								[ $ne != x ] && set -- "$@" ne$ne.png && name=$name"_ne$ne"
								[ $se != x ] && set -- "$@" se$se.png && name=$name"_se$se"
								[ $sw != x ] && set -- "$@" sw$sw.png && name=$name"_sw$sw"
								[ $nw != x ] && set -- "$@" nw$nw.png && name=$name"_wn$nw"
								name=$(printf '%s' "$name" | orsort)
								convert "$@" -compose Over -layers Flatten ../wall_$name.png
								cat <<EOF
								<tile id="$i">
								<image width="16" height="16" source="../tiles/wall_$name.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
EOF
								rotations "$name" | while read -r rot rotname; do
									cat <<EOF
								<property name="img.$rot" type="string" value="../tiles/wall_$rotname.png" />
EOF
								done
								cat <<EOF
								</properties>
								</tile>
EOF
								i=$((i+1))
							done
						done
					done
				done
			done
		done
	done
done


