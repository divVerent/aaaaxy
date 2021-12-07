# ![AAAAXY](docs/logo.png)

AAAAXY is a nonlinear puzzle platformer taking place in non-Euclidean
geometry.

Although your general goal is reaching the surprising end of the game,
you are encouraged to set your own goals while playing. Exploration will
be rewarded, and secrets await you\!

So jump and run around, and enjoy losing your sense of orientation in
this World of Wicked Weirdness. Find out what Van Vlijmen will make you
do. Pick a path, get inside a Klein Bottle, recognize some memes, and by
all means: don't look up.

[![Get it from the Snap
Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/aaaaxy)

[Or from
Flathub](https://flathub.org/apps/details/io.github.divverent.aaaaxy)

[Or from Itch](https://divverent.itch.io/aaaaxy)

## Screenshots

[![shot1](docs/screenshots/shot1.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot1.png)
[![shot2](docs/screenshots/shot2.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot2.png)
[![shot3](docs/screenshots/shot3.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot3.png)
[![shot4](docs/screenshots/shot4.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot4.png)
[![shot5](docs/screenshots/shot5.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot5.png)
[![shot6](docs/screenshots/shot6.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot6.png)
[![shot7](docs/screenshots/shot7.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot7.png)
[![shot8](docs/screenshots/shot8.jpg)](https://raw.githubusercontent.com/divVerent/aaaaxy/main/docs/screenshots/shot8.png)

## More Resources

This documentation is intended for developers; user-centric
documentation is available on the [game's
website](https://divverent.github.io/aaaaxy/) or [in this
repository](docs/index.md).

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
make
```

NOTE: On FreeBSD systems, use `gmake` instead of `make` to compile.

To update and rebuild, run:

``` sh
cd aaaaxy
git pull
make
```

You can also immediately compile and run the game using:

``` sh
make run
```

## Mathematical Notes

This game does not take place in the Euclidean space you're used to -
instead, you are experiencing the universal cover of a massively twisted
space. This is rather similar to seamless portals, but yields stronger
immersion and is generally an interesting approach I wanted to try out.
In particular, gravity behaves consistently across portals, and objects
are entirely glitch-free around them. On the other hand, this approach
can not support multi-player games; a more traditional portal based
engine would be more appropriate there, where the same object may be
seen multiple times on the screen. Sadly that approach would appear
rather confusing in places where this would cause the player themselves
to show up multiple times at once, which is why I did not go with it.

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

  - Design levels so that conflicting geometry is never on screen.
      - In other words, when transparently teleporting in order to move
        the player past a portal, a screen-sized environment of the
        source position must always match a screen-sized environment of
        the destination position.
      - This approach is simple and very immersive and has already been
        used in the original Super Mario Bros. game on the NES.
      - It however is not very flexible as any non-Euclidean geometry
        has to be rather large and behave fully Euclidean on every
        screen-sized environment around positions the player can visit.
  - Hide anything that has no line of sight to the player.
      - This actually matches the approach used in first person 3D games
      - With this approach, the game needs to be consistent with
        Euclidean geometry only in small environments around each
        object.
          - In this implementation, the consistency requirement is that
            an 1-tile environment around every portal must match, and
            that the same "screen tile" cannot be reached by a 1-tile
            environment around a line of sight through two different
            sets of portals at the same time.
          - As this game demonstrates, this can yield rather interesting
            while still obvious non-Euclidean geometry.
      - This is the approach has been explored in this game as well -
        but very likely for the first time in a two dimensional game.

## License

This project is released under the [Apache 2.0 License](LICENSE).

## Disclaimer

This is not an officially supported Google product.
