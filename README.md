# ![AAAAXY](https://divVerent.github.io/aaaaxy/logo.png)

AAAAXY is a nonlinear 2D puzzle platformer taking place in impossible
spaces.

Although your general goal is reaching the surprising end of the game,
you are encouraged to set your own goals while playing. Exploration will
be rewarded, and secrets await you!

So jump and run around, and enjoy losing your sense of orientation in
this World of Wicked Weirdness. Find out what Van Vlijmen will make you
do. Pick a path, get inside a Klein Bottle, recognize some memes, and by
all means: don't look up.

And beware of a minor amount of trolling.

To reach the end, a new player will take about 4 to 6 hours, a full
playthrough can be finished in about 1 hour and the end can be reached
in about 15 minutes.

The game is available for the following platforms:

| Platform     | Downloads                                                                                                                                                                                                           |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Android      | [F-Droid](https://f-droid.org/en/packages/io.github.divverent.aaaaxy/), [Google Play](https://play.google.com/store/apps/details?id=io.github.divverent.aaaaxy)                                                     |
| HTML5 (slow) | [netcup](https://rm.cloudns.org/aaaaxy/current/aaaaxy.html)                                                                                                                                                         |
| iOS          | [App Store](https://apps.apple.com/us/app/aaaaxy/id6447063878)                                                                                                                                                      |
| Linux        | [Flathub](https://flathub.org/apps/details/io.github.divverent.aaaaxy), [GitHub](https://github.com/divVerent/aaaaxy/releases), [Itch](https://divverent.itch.io/aaaaxy), [Snap Store](https://snapcraft.io/aaaaxy) |
| macOS        | [GitHub](https://github.com/divVerent/aaaaxy/releases), [Itch](https://divverent.itch.io/aaaaxy)                                                                                                                    |
| Windows      | [GitHub](https://github.com/divVerent/aaaaxy/releases), [Itch](https://divverent.itch.io/aaaaxy)                                                                                                                    |

Available languages: Chinese (Simplified), English, German, Latin,
Portuguese and Ukrainian.

## Screenshots

[![shot1](https://divVerent.github.io/aaaaxy/screenshots/shot1.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot1.png)
[![shot2](https://divVerent.github.io/aaaaxy/screenshots/shot2.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot2.png)
[![shot3](https://divVerent.github.io/aaaaxy/screenshots/shot3.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot3.png)
[![shot4](https://divVerent.github.io/aaaaxy/screenshots/shot4.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot4.png)
[![shot5](https://divVerent.github.io/aaaaxy/screenshots/shot5.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot5.png)
[![shot6](https://divVerent.github.io/aaaaxy/screenshots/shot6.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot6.png)
[![shot7](https://divVerent.github.io/aaaaxy/screenshots/shot7.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot7.png)
[![shot8](https://divVerent.github.io/aaaaxy/screenshots/shot8.jpg)](https://divVerent.github.io/aaaaxy/screenshots/shot8.png)

## More Resources

This documentation is intended for developers; user-centric
documentation is available on the [game's
website](https://divverent.github.io/aaaaxy/).

## Notable Libraries

This game is based on the following libraries:

-   [Ebitengine](https://github.com/hajimehoshi/ebiten) for low level
    graphics and input
-   [Oto](https://github.com/hajimehoshi/oto) for sound.
-   [tmx](https://github.com/fardog/tmx) for parsing
    [Tiled](https://www.mapeditor.org/) tile maps.

## Compiling

This game has been successfully compiled on Linux, FreeBSD and Windows,
and will likely compile just fine on other unixoid systems as well. For
other platforms some minor porting may be required so save games and
settings can be retained; in particular see `vfs/state_*.go`.

To build the game for yourself, install `git`, `golang`, `graphviz`,
`imagemagick` and `pandoc`, and then run:

``` sh
git clone https://github.com/divVerent/aaaaxy
cd aaaaxy
git submodule update --init --remote
make
```

NOTE: On FreeBSD systems, use `gmake` instead of `make` to compile.

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

## Editing

### In a Source Checkout

Use Tiled on the included `aaaaxy.tiled-project`. The map is in
`assets/maps/level.tmx`.

If checkpoints were changed, run `make assets-update` to regenerate the
checkpoint map.

`make run` will use the modified data.

### On a Release Binary

Run the game with `-dump_embedded_assets=/path/to/folder/for/editing`.
Download the `mappingsupport` pack from the release and extract it to
that same directory. Tell Tiled to use the included \`objecttypes.xml.

Running the game with
`-cheat_replace_embedded_assets=/path/to/folder/for/editing` will then
use the modified data.

Changing checkpoints is not supported this way.

## Mathematical Notes

This game does not take place in the Euclidean space you're used to -
instead, you are experiencing the universal cover of a massively twisted
space, which feels like a space with a locally-Euclidean topology and
geometry, which however violates some of their axioms globally (e.g. the
shortest path from one point to another *modulo the equivalence
relation* is not necessarily a line, and parallels are not always
unique).

For some added confusion, the space is occasionally reconfigured by
turning portals on/off at runtime, further violating expectations.

This is rather similar to seamless portals, but yields stronger
immersion and is generally an interesting approach I wanted to try out.
In particular, gravity behaves consistently across portals, and objects
are entirely glitch-free around them. The player is only ever visible
once. On the other hand, this approach can not sensibly support
multi-player games; a more traditional portal based engine would be more
appropriate there, where the same object may be seen multiple times on
the screen. Sadly that approach would appear rather confusing in places
where this would cause the player themselves to show up multiple times
at once, which is why I did not go with it.

In 3D games with transparent portals/warpzones, immersion is usually
achieved by treating each portal as a dynamic texture surface which
shows a view out of a camera projected to the other side. An open source
implementation of this can e.g. be found in
[Xonotic](https://www.xonotic.org). The reason why this works is that
only those parts of an object are shown that have a line of sight to the
player - as expected in a first person game. Implementing a third person
view that way is already a bit more tricky but usually one can work
around the view origin mismatching but being close to the player origin;
however what is usually not possible with transparent portals in a
top-down or other 2D-ish view, as conflicting (e.g. self-overlapping)
geometry may need to be rendered to the screen at the same time.

There are however two approaches to solve this:

-   Design levels so that conflicting geometry is never on screen.
    -   In other words, when transparently teleporting in order to move
        the player past a portal, a screen-sized environment of the
        source position must always match a screen-sized environment of
        the destination position.
    -   This approach is simple and very immersive and has already been
        used in the original Super Mario Bros. game on the NES.
    -   It however is not very flexible as any impossible geometry has
        to be rather large and behave fully Euclidean on every
        screen-sized environment around positions the player can visit.
-   Hide anything that has no line of sight to the player.
    -   This actually matches the approach used in first person 3D games
    -   With this approach, the game needs to be consistent with
        Euclidean geometry only in small environments around each
        object.
        -   In this implementation, the consistency requirement is that
            an 1-tile environment around every portal must match, and
            that the same "screen tile" cannot be reached by a 1-tile
            environment around a line of sight through two different
            sets of portals at the same time.
        -   As this game demonstrates, this can yield rather interesting
            while still obvious non-Euclidean topologic properties.
    -   This is the approach has been explored in this game as well -
        but very likely for the first time in a two dimensional game.

## License

This project is released under the [Apache 2.0 License](LICENSE).

## Disclaimer

This is not an officially supported Google product.
