# How to Contribute

We'd love to accept your patches and contributions to this project.
There are just a few small guidelines you need to follow.

## Code reviews

All submissions, including submissions by project members, require
review. We use GitHub pull requests for this purpose. Consult [GitHub
Help](https://help.github.com/articles/about-pull-requests/) for more
information on using pull requests.

### Translating

Go to <https://www.transifex.com/aaaaxy/aaaaxy> to translate AAAAXY to
your language\! You can always copy the translation files from Transifex
to `assets/locales/<language>/` to try them out locally.

Some automated verification of translation files is done at startup to
check that text fits in boxes and format strings work properly - see the
log output of the game in the terminal for details.

A script to quickly try out your downloaded translations from Transifex
is provided in `scripts/try-from-transifex.sh`.

Right now, some strings are used in multiple places and thus will need
different inflection in languages like Latin and Russian; before
translating to languages that require this, this may need to be resolved
in code first.

This in particular applies to proper names expanded via `{{BigCity}}`
into both `Welcome to %s` and `%s Road Rage`. If your target language
does not inflect proper names, this simplifies things a lot.

## Community Guidelines

This project follows [Google's Open Source Community
Guidelines](https://opensource.google/conduct/).
