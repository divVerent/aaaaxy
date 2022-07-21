module github.com/divVerent/aaaaxy

go 1.17

require (
	github.com/Microsoft/go-winio v0.5.2
	github.com/adrg/xdg v0.4.0
	github.com/akavel/rsrc v0.10.2
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/go-cmp v0.5.8
	github.com/google/go-licenses v0.0.0-20211006200916-ceb292363ec8
	github.com/hajimehoshi/ebiten/v2 v2.3.6
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/ncruces/zenity v0.8.9
	golang.org/x/image v0.0.0-20220617043117-41969df76e82
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8
)

require (
	github.com/dchest/jsmin v0.0.0-20220218165748-59f39799265f // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20220712193148-63cf1f4ef61f // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/google/licenseclassifier v0.0.0-20210722185704-3043a050f148 // indirect
	github.com/hajimehoshi/oto/v2 v2.3.0-alpha.4 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jezek/xgb v1.0.1 // indirect
	github.com/jfreymuth/oggvorbis v1.0.3 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/josephspurrier/goversioninfo v1.4.0 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20190725054713-01f96b0aa0cd // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/randall77/makefat v0.0.0-20210315173500-7ddd0e42c844 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/src-d/gcfg v1.4.0 // indirect
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/mobile v0.0.0-20220518205345-8578da9835fd // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

require (
	github.com/otiai10/copy v1.6.0
	github.com/spf13/cobra v0.0.5
	golang.org/x/exp v0.0.0-20220126164734-073fb1339172 // indirect
)

// Workaround for osxcross.
// replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 v2.3.3-0.20220517135915-073c6f4f13b7

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /home/rpolzer/src/ebiten
// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw
