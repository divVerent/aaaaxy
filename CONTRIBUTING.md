# How to Contribute

We'd love to accept your patches and contributions to this project.
There are just a few small guidelines you need to follow.

## Code reviews

All submissions, including submissions by project members, require
review. We use GitHub pull requests for this purpose. Consult [GitHub
Help](https://help.github.com/articles/about-pull-requests/) for more
information on using pull requests.

## Translating

Go to <https://www.transifex.com/aaaaxy/aaaaxy> to translate AAAAXY to
your language! You can always copy the translation files from Transifex
to `assets/locales/<language>/` to try them out locally.

Some automated verification of translation files is done at startup to
check that text fits in boxes and format strings work properly - see the
log output of the game in the terminal for details.

### Inflections

Right now, some strings are used in multiple places and thus will need
different inflection in languages like Latin and Russian; before
translating to languages that require this, this may need to be resolved
in code first.

This in particular applies to proper names expanded via `{{BigCity}}`
into both `Welcome to %s` and `%s Road Rage`. If your target language
does not inflect proper names, this simplifies things a lot.

### Reordering

Sometimes arguments in a format string need to be reordered when
translating.

This can be achieved using [Go's reordering
syntax](https://pkg.go.dev/fmt#hdr-Explicit_argument_indexes). In
particular, the following things can be done:

-   `%s walks towards %s`
-   `%[1]s walks towards %[2]s` - same but with explicit indexes
-   `%[2]s is where %[1]s walks towards` - reordered

### Trying It Out

If you compiled the game from source code, you can quickly try out your
downloaded translations from Transifex using
`scripts/try-from-transifex.sh`. When using the game from a binary
download from GitHub's "Releases" section, you can try it out as follows
instead (here, `xx` stands for the short name of the language you are
translating to):

1.  Unpack the release zip file, as usual.
2.  Open a command prompt in the directory containing the game.
3.  Run `./aaaaxy -dump_embedded_assets=data`
4.  If the language isn't in the game yet:
    1.  Edit `data/locales/LINGUAS` to include the `xx`.
    2.  Create a subdirectory in `data/locales` named `xx`.
5.  Put your translated `game.po` and `level.po` files in the
    subdirectory within `data/locales/xx`.
6.  Run `./aaaaxy -cheat_replace_embedded_assets=data -language=xx` to
    run the game with the modified data.

## Community Guidelines

This project follows [Google's Open Source Community
Guidelines](https://opensource.google/conduct/).
