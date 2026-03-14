# How to Contribute

We'd love to accept your patches and contributions to this project.
There are just a few small guidelines you need to follow.

## Code reviews

All submissions, including submissions by project members, require
review, and must be fully understood by the submitter before
contributing. The submitter must be available for review comments. We
use GitHub pull requests for this purpose. Consult [GitHub
Help](https://help.github.com/articles/about-pull-requests/) for more
information on using pull requests.

## Licensing

All contributions must be licensed under the [Apache 2.0
License](LICENSE).

To ensure proper licensing, as an author you agree to the following
terms:

- You are holding sufficient rights to the work to publish it under the
  [Apache 2.0 License](LICENSE). This usually requires you to be the
  sole copyright holder, or to obtain permission by all other
  contributors to the change; in the latter case, you must mention all
  other contributors and include evidence that these other contributors
  have given said permission.

- You are submitting your change according to the [Apache 2.0
  License](LICENSE) with no additional terms attached.

## AI Policy

Contributions by a LLM or other AI entity still must follow the above
requirements, just like human contributions.

In particular this means:

- The model or entity that has been used must be named.

- It must not have been trained or prompted on works with incompatible
  licenses. Neither must the output for any other reason be derivative
  of works with incompatible licenses.

- All authors of works that went into the training or prompting, or that
  the submission is otherwise derivative of, must be named and the
  respective licenses be stated.

- Note that an excessive number of contributors likely makes reviewing
  accuracy of the licensing information infeasible and usually leads to
  rejection.

- The human submitter certifies they have fully understood their
  submission and can respond to comments or questions about it.

## Translating

Go to <https://app.transifex.com/aaaaxy/aaaaxy> to translate AAAAXY to
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

- `%s walks towards %s`
- `%[1]s walks towards %[2]s` - same but with explicit indexes
- `%[2]s is where %[1]s walks towards` - reordered

### Trying It Out

If you compiled the game from source code, you can quickly try out your
downloaded translations from Transifex using
`scripts/try-from-transifex.sh`. When using the game from a binary
download from GitHub's "Releases" section, you can try it out as follows
instead:

1.  Unpack the release zip file, as usual.
2.  Put your translated `game.po` and `level.po` files - renamed to
    those exact file names - in the same directory as the game
    executable.
3.  Run `./aaaaxy -language=.` to run the game with the modified data
    (or simply run the game normally and select "user provided" in the
    language selector, which will be second open from the beginning of
    the list).

## Community Guidelines

This project follows [Google's Open Source Community
Guidelines](https://opensource.google/conduct/).
