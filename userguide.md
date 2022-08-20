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

  - Especially on laptops with 4K displays connected, graphics
    performance may be poor. As a workaround, press `Esc` to switch to
    the game menu, then in the settings switch to a lower graphics
    setting. <https://github.com/hajimehoshi/ebiten/issues/1772> tracks
    a fix for this issue.

### Video Recording

#### To MP4

To prepare for recording videos, make sure the `ffmpeg` command is
available and in the current search path. On Windows, just [download
it](https://ffmpeg.org/download.html) and place `ffmpeg.exe` right next
to `aaaaxy.exe`.

To then record a video of the game, first record a demo (see below), and
then un the game again, passing the flags `-dump_media=video.mp4
-demo_play=demo.dem`. This will save a video of the gameplay under
`video.mp4`.

#### Via Raw Files

To record a video of the game, pass the flags `-dump_video=video.raw
-dump_audio=audio.raw`. Then play normally (game may be slower). When
you then exit the game, the console window will show you a FFmpeg
command to turn these files into a finished video\!

Video recording can be sped up by using `make FASTER_VIDEO_DUMPING=true`
when compiling to enable multithreading, and by passing
`-dump_video_fps_divisor=2` to limit the video to SFR (30fps).

### Demo Recording

To record a demo of the game, pass the flags `-demo_record=demo.dem`.

The resulting `demo.dem` file can be played back with only this exact
same version of the game by passing `-demo_play=demo.dem`; however the
above section on video recording can be used to turn the demo into a
video.

Note that demo playback during video recording is never at realtime. You
don't want any duplicate frames, right?

### Data Storage

AAAAXY stores saved games in the following location:

| Operating System |                                                  Save Games<br>Configuration                                                   |
| :--------------: | :----------------------------------------------------------------------------------------------------------------------------: |
|     Android      | `/sdcard/Android/data/io.github.divVerent.aaaaxy/files/save`<br>`/sdcard/Android/data/io.github.divVerent.aaaaxy/files/config` |
|  Linux AppImage  |                                         `~/.local/share/AAAAXY`<br>`~/.config/AAAAXY`                                          |
|  Linux FlatPak   |          `~/.var/app/io.github.divverent.aaaaxy/data/AAAAXY`<br>`~/.var/app/io.github.divverent.aaaaxy/config/AAAAXY`          |
|    Linux Snap    |                      `~/snap/aaaaxy/common/.local/share/AAAAXY`<br>`~/snap/aaaaxy/common/.config/AAAAXY`                       |
|   Linux native   |                                         `~/.local/share/AAAAXY`<br>`~/.config/AAAAXY`                                          |
|      macOS       |                            `~/Library/Application Support/AAAAXY`<br>`~/Library/Preferences/AAAAXY`                            |
|       Web        |                                                    `getSave(n)`<br>`get()`                                                     |
|     Windows      |                     `C:\Users\%USERNAME%\Saved Games\AAAAXY`<br>`C:\Users\%USERNAME%\AppData\Local\AAAAXY`                     |
|       Wine       |    `~/.wine/drive_c/users/$USER/Saved Games/AAAAXY`<br>`~/.wine/drive_c/users/$USER/Local Settings/Application Data/AAAAXY`    |

It is recommended to back up these files.

WARNING: Do not edit the save game files. If needed, cheats are provided
as command line options.

To edit a config setting on the web, type into the developer console
something like:

    setConf({show_fps: true, vsync: false});

To do this on Itch's web player, select the `index.html` subframe that
is hosted on `hwcdn.net` first.

On Android, the file path above can only be reached via `adb` (USB
Debugging). However, Android One backup of the savegames is enabled in
the app.

### Save States

Save states can be switched in `Settings` / `Switch Save State`.
