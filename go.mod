module github.com/divVerent/aaaaxy

go 1.19

require (
	github.com/Microsoft/go-winio v0.6.1
	github.com/adrg/xdg v0.4.0
	github.com/akavel/rsrc v0.10.2
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/google/go-cmp v0.6.0
	github.com/google/go-licenses v0.0.0-20211006200916-ceb292363ec8
	github.com/hajimehoshi/bitmapfont/v3 v3.0.0
	github.com/hajimehoshi/ebiten/v2 v2.6.2
	github.com/jeandeaual/go-locale v0.0.0-20220711133428-7de61946b173
	github.com/leonelquinteros/gotext v1.5.2
	github.com/lestrrat-go/strftime v1.0.6
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/ncruces/zenity v0.10.10
	github.com/zachomedia/go-bdf v0.0.0-20220611021443-a3af701111be
	golang.org/x/image v0.13.0
	golang.org/x/mobile v0.0.0-20231006135142-2b44d11868fe
	golang.org/x/sys v0.13.0
	golang.org/x/text v0.13.0
)

require (
	github.com/dchest/jsmin v0.0.0-20220218165748-59f39799265f // indirect
	github.com/ebitengine/purego v0.4.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/google/licenseclassifier v0.0.0-20210722185704-3043a050f148 // indirect
	github.com/hajimehoshi/oto/v2 v2.5.0-alpha.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
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
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

require (
	github.com/otiai10/copy v1.6.0
	github.com/spf13/cobra v0.0.5
)

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /Users/rpolzer/src/ebiten

// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw

// update-ebitengine-fork.sh changes:
replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 v2.5.10-with-574925cf7a72deaf73be4c481348a7a44f7b7e19-and-cc247962703eba99eae732876496375191f16cbe-and-58eb1af1eb88915a726f29ffe9279025f98c9be8
