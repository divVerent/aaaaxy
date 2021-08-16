# Releasing

AAAAXY releases are published in the following places:

  - GitHub Releases

TODO(divVerent): add more; provide links.

# Versioning

We follow semantic versioning, but in an adapted form to suppport
speedrunning leaderboards.

Thus, the version components are incremented as follows:

  - Major version must be incremented on changes that likely break
    existing speedruns.
      - In particular, slowing down a section required for any speedrun
        categories, including 100%, requires a major version bump.
  - Minor version must be incremented on changes that likely add faster
    speedruns.
      - In particular, simplifying a section or making it optional
        requires a minor version bump.
  - Patch level must be incremented in any other case.

Note that I say &ldquo;likely&rdquo;; if a section is not required for
any known/published speedrun, slowing it down may not require a major
version bump.
