#!/usr/bin/perl
# Copyright 2025 Google LLC
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

use strict;
use warnings;

# 2025/02/16 21:28:22.176348 [INFO] entity "MovableSprite" property "required_orientation" type "string" default ""

my %ents;

while (<>) {
	/\[INFO\] entity "(\w+)" property "(\w+)" type "(\w+)"(?: default "(.*)")?$/
		or next;
	my $p = \$ents{$1}{$2};
	if (defined $$p) {
		if ($$p->[0] ne $3) {
			warn "Mismatching info for entity \"$1\" property \"$2\": $$p->[0]($$p->[1]) vs $3($4)\n";
		}
		if (defined $$p->[1] && defined $4 && $$p->[1] ne $4) {
			warn "Mismatching info for entity \"$1\" property \"$2\": $$p->[0]($$p->[1]) vs $3($4)\n";
		}
		$$p->[1] = $4
			if not defined $$p->[1];
	} else {
		$$p = [$3, $4];
	}
}

my %colors = (
	Animation => "#ffffffff",
	AppearBlock => "#ff00aa00",
	CenterPrintTarget => "#ff0000ff",
	Checkpoint => "#ff008000",
	CoverSprite => "#ffffffff",
	CreditsTarget => "#ffff00ff",
	DelayTarget => "#ff000000",
	DisappearBlock => "#ff00aa00",
	ExitButton => "#ffffffff",
	FadeTarget => "#ffff00ff",
	ForceField => "#ffff00ff",
	Give => "#ffffff00",
	Goal => "#ffffff00",
	JumpPad => "#ffff00ff",
	LogicalGate => "#ff000000",
	MovableSprite => "#ffffffff",
	MovingAnimation => "#ffffffff",
	OneWay => "#ff0000ff",
	Player => "#ff008000",
	PrintToConsoleTarget => "#ff000000",
	QuestionBlock => "#ff000000",
	RespawnPlayer => "#ffff0000",
	Riser => "#ff000080",
	RiserFsck => "#ffffff80",
	SequenceCollector => "#ff00ff00",
	SequenceTarget => "#ff00ff00",
	SetState => "#ffff0000",
	SoundTarget => "#ff0000ff",
	SpawnCounter => "#ffff0000",
	Sprite => "#ffffffff",
	StopTimerTarget => "#ffff00ff",
	Switch => "#ffff0000",
	SwitchMusic => "#ff00ff00",
	SwitchMusicTarget => "#ff00ff00",
	SwitchableAnimation => "#ffffffff",
	SwitchableJumpPad => "#ffff00ff",
	SwitchableSprite => "#ffffffff",
	SwitchableText => "#ffffffff",
	Text => "#ffffffff",
	TnihSign => "#ffffff00",
	VVVVVV => "#ff00ff00",
	WarpZone => "#ffff0000",
	ZoomTarget => "#ffff00ff",
	_TileMod => "#ff0000ff",
);

my $id = 0;
my $json = '';
for my $entity (sort keys %ents) {
	++$id;
	$json .= <<EOF;
        {
            "color": "$colors{$entity}",
            "id": $id,
            "members": [
EOF
	for my $property (sort keys %{$ents{$entity}}) {
		my ($type, $default) = @{$ents{$entity}{$property}};
		my $default_json = <<EOF;
,
                    "value":\040
EOF
		chomp $default_json;
		if (!defined $default) {
			$default_json = '';
		} elsif ($type eq 'bool') {
			$default_json .= $default;
		} elsif ($type eq 'color') {
			$default_json .= "\"$default\"";
		} elsif ($type eq 'float') {
			$default_json .= $default;
		} elsif ($type eq 'int') {
			$default_json .= $default;
		} elsif ($type eq 'string') {
			$default_json .= "\"$default\"";
		} else {
			die "Unsupported type: $type (value: $default)";
		}
		$json .= <<EOF;
                {
                    "name": "$property",
                    "type": "$type"$default_json
                },
EOF
	}
	$json .= <<EOF;
            ],
            "name": "$entity",
            "type": "class",
            "useAs": [
                "object",
                "tile"
            ]
        },
EOF
}
$json =~ s{,(\s*[\]\}])}{$1}g;
print $json;
