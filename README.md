# AAAAXY

AAAAXY is a nonlinear puzzle platformer in a geometrically impossible
world.

Jump and run around, collect notes and find the surprising ending of the
game\! Fastest completion of the game wins\!

Enjoy losing your sense of orientation\!

## Input

AAAAXY can be played with a keyboard or any game pad having at least the
NES buttons. The exact controls are to be guessed by the player.

If your gamepad is not working, you may need to put its definition in
the `SDL_GAMECONTROLLERCONFIG` environment variable.

## Installing

AAAAXY is released in binary form as a zip file containing
self-contained executables for each supported platform.

So just extract the game executable to a convenient place and run it
from there\!

## Compiling

To build the game for yourself, run:

``` sh
git clone https://github.com/divVerent/aaaaxy
cd aaaaxy
git submodule update --init --remote
make
```

To update and rebuild, run:

``` sh
cd aaaaxy
git pull
git submodule update --remote
make
```

You can also immediately compile and run the game using:

``` sh
make run
```

## Video Recording

To record a video of the game, pass the flags `-dump_video=video.raw
-dump_audio=audio.raw`. Then play normally. When you then exit the game,
the console window will show you a FFmpeg command to turn these files
into a finished video\!

Video recording can be sped up by using `make FASTER_VIDEO_DUMPING=true`
when compiling.

## Data Storage

AAAAXY stores saved games in the following location:

  - Windows: `Saved Games/AAAAXY`
  - Linux: `~/.local/share/AAAAXY`

Configuration is stored as follows:

  - Windows: `AppData/Local/AAAAXY`
  - Linux: `~/.config/AAAAXY`

It is recommended to back up these files.

## Save States

TODO(divVerent): Fix.

At the moment there is no menu to select save states; to play a
particular save state, please run the game binary with a command line
argument such as `-save_state=1` to select the save state with index 1.
The default save state can be switched to by passing `-save_state=0`.

Save states record progress at checkpoints only and the game
automatically saves when hitting one.

## Mathematical Notes

This game does not take place in the Euclidean space you're used to -
instead, you are experiencing the universal cover of a massively twisted
space. This is rather similar to seamless portals, but yields stronger
immersion and is generally an interesting approach I wanted to try out.
In particular, gravity behaves consistently across portals, and objects
are entirely glitch-free around them.

## License

This project is released under the [Apache 2.0 License](LICENSE).

## Disclaimer

This is not an officially supported Google product.
