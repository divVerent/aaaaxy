## ![AAAAXY](logo.png)

## Speedrunning

AAAAXY has integrated speedrun timing support. Your time and run
category will be shown in the end credits.

### Run Categories

The following run categories exist:

-   Any%: every run is Any% (but Any% is not shown if 100% is met too).
-   100%: all regular checkpoints have been hit.
-   All Notes: all notes have been read. Required an almost 100% run.
-   All Paths: all connections between the visited checkpoints have been
    used by the player during the run in one or the other direction.
-   All Secrets: all notes in all secret rooms have been read.
-   All Flipped: all checkpoints that were seen were last seen in
    flipped/mirrored state.
-   No Teleports: the checkpoint map was never used to teleport to a
    different checkpoint.
-   No Escape/Backspace/Start: the game menu was never used during the
    run. Implies No Teleports.

Also, All Secrets and No Escape are mutually exclusive, unless you find
a glitch I do not know about yet :)

### Timing

AAAAXY includes an integrated timer with the following rules:

-   Game load time does not count.
-   Timer starts from first action after game load.
-   Timer continues ticking while in the menu.
-   Split times are shown in-game in various places.
-   Timing ends once input is no longer accepted at the end of the game.
-   To show a permanent in-game timer, pass `-show_time` as command line
    argument (optional).
