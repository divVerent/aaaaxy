# Releasing

AAAAXY releases are published in the following places:

  - [GitHub Releases](https://github.com/divVerent/aaaaxy/releases)
      - [Analytics](https://github.com/divVerent/aaaaxy/graphs/traffic)
  - AppImage
      - Primary copy on GitHub Releases; seems mirrored in lots of
        places.
  - [Snap](https://snapcraft.io/aaaaxy)
      - [Analytics](https://snapcraft.io/aaaaxy/metrics)
  - [Flathub](https://flathub.org/apps/details/io.github.divverent.aaaaxy)
      - [Analytics](https://klausenbusk.github.io/flathub-stats/#ref=io.github.divverent.aaaaxy&interval=infinity&downloadType=installs%2Bupdates)
      - To see active users per release, run `sh
        scripts/flathub-stats.sh`.
  - [Itch](https://divverent.itch.io/aaaaxy)
      - [Analytics](https://itch.io/game/summary/1199736)

# Binaries

The official binary release is built using `sh
scripts/binary-release.sh`. This will also print instructions about
releasing to AppImage, Snap, Flathub and Itch.

# Versioning

We follow semantic versioning, but in an adapted form to suppport
speedrunning leaderboards.

Thus, the version components are incremented as follows:

  - Level version must be incremented on changes that break most save
    games in ways that are not easy to "repair".
      - In particular removing or renaming checkpoint internal names
        causes this; the game will not start if the current checkpoint
        is not on the map.
      - However, changing a checkpoint's "text" property is OK and does
        not require a level version bump.
      - Generally such changes are strongly discouraged.
      - This also requires bumping the major version. Set the level
        version to the new major version then.
  - Major version must be incremented on changes that likely break
    existing speedruns.
      - In particular, slowing down a section required for any speedrun
        categories, including 100%, requires a major version bump.
      - Exception: when some major cheese (grossly unintended skip) is
        fixed, a minor version bump is sufficient.
      - Exception: slowing down "All Secrets" speedruns only requires a
        minor version bump.
  - Minor version must be incremented on changes that likely add faster
    speedruns.
      - In particular, simplifying a section or making it optional
        requires a minor version bump.
  - Patch level must be incremented in any other case.
  - We use `*-alpha.*` pre-release sub-versions whenever major or minor
    have been bumped; this only gets turned into a real new
    `<major>.<minor>.0` version when actually releasing.
      - When patch level is bumped, switching to a pre-release version
        is not necessary.

Note that I say &ldquo;likely&rdquo;; if a section is not required for
any known/published speedrun, slowing it down may not require a major
version bump.

## Automation

Version is partially automated using git.

The version is built based on the closest git `v<major>.<minor>-alpha`
or `v<major>.<minor>` tag. Binary releases will be tagged
`v<major>.<minor>.<patch>` or `v<major>.<minor>.0-alpha<patch>`; any
commit on git counts as a source-only release with a patchlevel and
needs not be released as a binary.

## Conversion to Windows Version Numbers

The Windows scheme differs a little from semantic versioning by using a
`major.minor.revision.buildnumber` scheme. We convert as follows:

  - `major` maps to `major`.
  - `minor` maps to `minor`.
  - `patch` maps to `revision+N` where N is `0` for alpha, `10000` for
    beta, `20000` for rc and `30000` for finished versions.
  - `buildnumber` is always the total number of commits in the
    repository.
