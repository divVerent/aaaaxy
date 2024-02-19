module github.com/divVerent/aaaaxy

go 1.19

require (
	github.com/Microsoft/go-winio v0.6.1
	github.com/adrg/xdg v0.4.0
	github.com/akavel/rsrc v0.10.2 // indirect
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/google/go-cmp v0.6.0
	github.com/hajimehoshi/bitmapfont/v3 v3.0.0
	github.com/hajimehoshi/ebiten/v2 v2.7.0-alpha.7.0.20240218172058-1a5832428306
	github.com/jeandeaual/go-locale v0.0.0-20240204043739-672d8d016d9a
	github.com/leonelquinteros/gotext v1.5.2
	github.com/lestrrat-go/strftime v1.0.6
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/ncruces/zenity v0.10.11
	github.com/zachomedia/go-bdf v0.0.0-20220611021443-a3af701111be
	golang.org/x/image v0.15.0
	golang.org/x/sys v0.17.0
	golang.org/x/text v0.14.0
)

require (
	github.com/dchest/jsmin v0.0.0-20220218165748-59f39799265f // indirect
	github.com/ebitengine/gomobile v0.0.0-20240218171544-120934310db5 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.2.0-alpha.4 // indirect
	github.com/ebitengine/purego v0.7.0-alpha // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/josephspurrier/goversioninfo v1.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/randall77/makefat v0.0.0-20210315173500-7ddd0e42c844 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
)

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /Users/rpolzer/src/ebiten

// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw

// update-ebitengine-fork.sh changes:
// replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 v2.5.10-with-574925cf7a72deaf73be4c481348a7a44f7b7e19-and-cc247962703eba99eae732876496375191f16cbe-and-58eb1af1eb88915a726f29ffe9279025f98c9be8
