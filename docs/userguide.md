## ![AAAAXY](logo.png)

## User Guide

### Installing

AAAAXY is released in [binary
form](https://github.com/divVerent/aaaaxy/releases) as a zip file
containing self-contained executables for each supported platform.

So just extract the game executable to a convenient place and run it
from there\!

### Input

AAAAXY can be played with a keyboard or any controller good enough for
playing NES games. While the controls do follow usual standards set by
two dimensional games of the past, some details are to be guessed by the
player and experimented with.

The game menu can, of course, be reached using the Escape key or the
Start button.

If your gamepad is not supported yet, you can typically make it work by
passing its definition in
[SDL\_GameControllerDB](https://github.com/gabomdq/SDL_GameControllerDB/blob/master/gamecontrollerdb.txt)
format as `-gamepad_override` flag or `SDL_GAMECONTROLLERCONFIG`
environment variable. As an extension, multiple gamepad definitions can
be provided not just separated by newlines but also by semicolons.

### Settings

Press `Esc` or `Start` to get to the game menu which has settings.

### Known Issues

  - Especially on laptops with 5K displays connected, graphics
    performance may be poor. As a workaround, press `Esc` to switch to
    the game menu, then in the settings switch to a lower graphics
    setting. <https://github.com/hajimehoshi/ebiten/issues/1772> tracks
    a fix for this issue.
  - On some Linux systems, fullscreen mode uses the wrong scaling
    factor. As a workaround, press `F` to switch to windowed mode, then
    maximize the window.

### Video Recording

To record a video of the game, pass the flags `-dump_video=video.raw
-dump_audio=audio.raw`. Then play normally (game may be slower). When
you then exit the game, the console window will show you a FFmpeg
command to turn these files into a finished video\!

Video recording can be sped up by using `make FASTER_VIDEO_DUMPING=true`
when compiling to enable multithreading, and by passing
`-dump_video_fps_divisor=2` to limit the video to SFR (30fps).

### Data Storage

AAAAXY stores saved games in the following location:

  - Windows: `Saved Games/AAAAXY`
  - Linux: `~/.local/share/AAAAXY`

Configuration is stored as follows:

  - Windows: `AppData/Local/AAAAXY`
  - Linux: `~/.config/AAAAXY`

It is recommended to back up these files.

### Save States

TODO(divVerent): Fix.

At the moment there is no menu to select save states; to play a
particular save state, please run the game binary with a command line
argument such as `-save_state=1` to select the save state with index 1.
The default save state can be switched to by passing `-save_state=0`.

Save states record progress at checkpoints only and the game
automatically saves when hitting one.

WARNING: Do not edit the save game files. If needed, cheats are provided
as command line options.
