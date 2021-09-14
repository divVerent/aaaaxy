module github.com/divVerent/aaaaxy

go 1.16

require (
	github.com/adrg/xdg v0.3.4
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20210727001814-0db043d8d5be // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/ebiten/v2 v2.2.0-alpha.14.0.20210914160304-923c84a3d6ca
	github.com/hajimehoshi/oto/v2 v2.1.0-alpha.0.20210912073017-18657977e3dc // indirect
	github.com/jezek/xgb v0.0.0-20210312150743-0e0f116e1240 // indirect
	github.com/jfreymuth/oggvorbis v1.0.3 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2
	golang.org/x/exp v0.0.0-20210910231120-3d0173ecaa1e // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/mobile v0.0.0-20210902104108-5d9a33257ab5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /home/rpolzer/src/ebiten
// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw
