## ![AAAAXY](logo.png)

## User Guide

### Installing

AAAAXY is released in [binary
form](https://github.com/divVerent/aaaaxy/releases) as a zip file
containing self-contained executables for each supported platform.

So just extract the game executable to a convenient place and run it
from there!

### Input

AAAAXY can be played with a keyboard or any controller good enough for
playing NES games. While the controls do follow usual standards set by
two dimensional games of the past, some details are to be guessed by the
player and experimented with.

The game menu can, of course, be reached using the Escape key or the
Start button.

If your gamepad is not supported yet, you can typically make it work by
passing its definition in
[SDL_GameControllerDB](https://github.com/gabomdq/SDL_GameControllerDB/blob/master/gamecontrollerdb.txt)
format as `-gamepad_override` flag or `SDL_GAMECONTROLLERCONFIG`
environment variable. As an extension, multiple gamepad definitions can
be provided not just separated by newlines but also by semicolons.

### Settings

Press `Esc` or `Start` to get to the game menu which has settings.

### Driver Settings

The following environment variables can be used in case the game does
not run at all at default settings:

-   On Windows:
    -   Rendering by default takes place using DirectX 11.
    -   `EBITENGINE_DIRECTX=version=12` forces rendering using DirectX
        12.
    -   `EBITENGINE_GRAPHICS_LIBRARY=opengl` forces rendering using
        OpenGL 3.2.
-   On Linux:
    -   Rendering by default takes place using OpenGL 3.2.
    -   `EBITENGINE_OPENGL=es` forces rendering using OpenGL ES 3.
    -   `MESA_GL_VERSION_OVERRIDE=3.2` together with
        `MESA_GLSL_VERSION_OVERRIDE=150` may make the game run on some
        OpenGL 2 graphics chips. This configuration is unsupported and
        may break at any time, and it sure will not work on all OpenGL 2
        chips.
-   On macOS:
    -   Rendering by default takes place using Metal.
    -   `EBITENGINE_GRAPHICS_LIBRARY=opengl` forces rendering using
        OpenGL 3.2.

### Video Recording

#### To MP4

To prepare for recording videos, make sure the `ffmpeg` command is
available and in the current search path. On Windows, just [download
it](https://ffmpeg.org/download.html) and place `ffmpeg.exe` right next
to `aaaaxy.exe`.

To then record a video of the game, first record a demo (see below), and
then un the game again, passing the flags
`-dump_media=video.mp4 -demo_play=demo.dem`. This will save a video of
the gameplay under `video.mp4`.

#### Via Raw Files

To record a video of the game, pass the flags
`-dump_video=video.raw -dump_audio=audio.raw`. Then play normally (game
may be slower). When you then exit the game, the console window will
show you a FFmpeg command to turn these files into a finished video!

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

By default AAAAXY stores saved games in the following location:

| Operating System |                                                  Save Games<br>Configuration                                                   |
|:----------------:|:------------------------------------------------------------------------------------------------------------------------------:|
|     Android      | `/sdcard/Android/data/io.github.divVerent.aaaaxy/files/save`<br>`/sdcard/Android/data/io.github.divVerent.aaaaxy/files/config` |
|       iOS        |                        `Library/Application Support/AAAAXY/save`<br>`Library/Preferences/AAAAXY/config`                        |
|  Linux AppImage  |                                         `~/.local/share/AAAAXY`<br>`~/.config/AAAAXY`                                          |
|  Linux FlatPak   |          `~/.var/app/io.github.divverent.aaaaxy/data/AAAAXY`<br>`~/.var/app/io.github.divverent.aaaaxy/config/AAAAXY`          |
|    Linux Snap    |                      `~/snap/aaaaxy/common/.local/share/AAAAXY`<br>`~/snap/aaaaxy/common/.config/AAAAXY`                       |
|   Linux native   |                                         `~/.local/share/AAAAXY`<br>`~/.config/AAAAXY`                                          |
|      macOS       |                        `~/Library/Application Support/AAAAXY`<br>`~/Library/Application Support/AAAAXY`                        |
|       Web        |                                                    `getSave(n)`<br>`get()`                                                     |
|     Windows      |                     `C:\Users\%USERNAME%\Saved Games\AAAAXY`<br>`C:\Users\%USERNAME%\AppData\Local\AAAAXY`                     |
|       Wine       |    `~/.wine/drive_c/users/$USER/Saved Games/AAAAXY`<br>`~/.wine/drive_c/users/$USER/Local Settings/Application Data/AAAAXY`    |

This can be customized by passing the flags `-conig_path` and
`-save_path`, or by passing `-portable`, in which case the files will be
stored in subdirectories named `config` and `save` of the current
directory.

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
the app. For example, to manually backup the primary save slot, run from
an USB-connected PC:

    adb pull /sdcard/Android/data/io.github.divVerent.aaaaxy/files/save/save-0.json aaaaxy-save-0.json

Similarly, to restore it, first launch the AAAAXY app, quit it again,
and then run:

    adb push aaaaxy-save-0.json /sdcard/Android/data/io.github.divVerent.aaaaxy/files/save/save-0.json

On iOS, the file path above can be reached using
[iExplorer](https://macroplant.com/iexplorer) from macOS and Windows,
and using [ifuse](https://github.com/libimobiledevice/ifuse) from Linux
like this:

    mkdir -p ~/mnt
    ifuse --container io.github.divverent.aaaaxy ~/mnt
    cp ~/mnt/Library/Application\ Support/AAAAXY/save/save-0.json aaaaxy-save-0.json
    # or: cp aaaaxy-save-0.json ~/mnt/Library/Application\ Support/AAAAXY/save/save-0.json
    fusermount -u ~/mnt

### Save States

Save states can be switched in `Settings` / `Switch Save State`.
