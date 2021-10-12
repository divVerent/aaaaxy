module github.com/divVerent/aaaaxy

go 1.16

require (
	github.com/adrg/xdg v0.3.4
	github.com/akavel/rsrc v0.10.2
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/go-licenses v0.0.0-20211006200916-ceb292363ec8
	github.com/hajimehoshi/ebiten/v2 v2.3.0-alpha.0.20210928160017-78cdb94552e2
	github.com/hajimehoshi/oto/v2 v2.1.0-alpha.2 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/exp v0.0.0-20211011213208-1d87cf485e27 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/mobile v0.0.0-20210924032853-1c027f395ef7 // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/net v0.0.0-20211011170408-caeb26a5c8c0 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac
	golang.org/x/tools v0.1.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /home/rpolzer/src/ebiten
// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw
