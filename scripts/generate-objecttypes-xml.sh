#!/bin/sh

cat <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<objecttypes>
EOF

grep ' \[INFO\] default for ' "$@" |\
cut -d ' ' -f 6- |\
LC_ALL=C sort -u |\
while read -r type xml; do
	type=${type#\"}
	type=${type%\":}
	if [ x"$type" != x"$curtype" ]; then
		if [ -n "$curtype" ]; then
			echo " </objecttype>"
		fi
		case "$type" in
			_TileMod)             color=0000ff ;;
			Animation)            color=ffffff ;;
			AppearBlock)          color=00aa00 ;;
			Checkpoint)           color=008000 ;;
			CheckpointTarget)     color=008000 ;;
			CoverSprite)          color=ffffff ;;
			CreditsTarget)        color=ff00ff ;;
			DelayTarget)          color=000000 ;;
			DisappearBlock)       color=00aa00 ;;
			ExitButton)           color=ffffff ;;
			FadeTarget)           color=ff00ff ;;
			ForceField)           color=ff00ff ;;
			Give)                 color=ffff00 ;;
			Goal)                 color=ffff00 ;;
			JumpPad)              color=ff00ff ;;
			LogicalGate)          color=000000 ;;
			MovableSprite)        color=ffffff ;;
			MovingAnimation)      color=ffffff ;;
			OneWay)               color=0000ff ;;
			Player)               color=008000 ;;
			PrintToConsoleTarget) color=000000 ;;
			QuestionBlock)        color=000000 ;;
			RespawnPlayer)        color=ff0000 ;;
			Riser)                color=000080 ;;
			RiserFsck)            color=ffff80 ;;
			SequenceCollector)    color=00ff00 ;;
			SequenceTarget)       color=00ff00 ;;
			SetState)             color=ff0000 ;;
			SetStateTarget)       color=ff0000 ;;
			SoundTarget)          color=0000ff ;;
			SpawnCounter)         color=ff0000 ;;
			Sprite)               color=ffffff ;;
			StopTimerTarget)      color=ff00ff ;;
			Switch)               color=ff0000 ;;
			SwitchableAnimation)  color=ffffff ;;
			SwitchableJumpPad)    color=ff00ff ;;
			SwitchableSprite)     color=ffffff ;;
			SwitchableText)       color=ffffff ;;
			SwitchMusic)          color=00ff00 ;;
			SwitchMusicTarget)    color=00ff00 ;;
			Text)                 color=ffffff ;;
			TnihSign)             color=ffff00 ;;
			VVVVVV)               color=00ff00 ;;
			WarpZone)             color=ff0000 ;;
			ZoomTarget)           color=ff00ff ;;
			*) echo >&2 "Add type: $type"; exit 1 ;;
		esac
		echo " <objecttype name=\"$type\" color=\"#$color\">"
		curtype=$type
	fi
	echo "  $xml"
done

cat <<EOF
 </objecttype>
</objecttypes>
EOF
