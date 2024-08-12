module github.com/divVerent/aaaaxy

go 1.21

toolchain go1.22.3

require (
	github.com/Microsoft/go-winio v0.6.2
	github.com/adrg/xdg v0.5.0
	github.com/akavel/rsrc v0.10.2
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/google/go-cmp v0.6.0
	github.com/google/go-licenses v1.6.1-0.20230903011517-706b9c60edd4
	github.com/hajimehoshi/bitmapfont/v3 v3.1.0
	github.com/hajimehoshi/ebiten/v2 v2.7.8
	github.com/jeandeaual/go-locale v0.0.0-20240223122105-ce5225dcaa49
	github.com/leonelquinteros/gotext v1.6.1
	github.com/lestrrat-go/strftime v1.0.6
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/ncruces/zenity v0.10.13
	github.com/zachomedia/go-bdf v0.0.0-20220611021443-a3af701111be
	golang.org/x/image v0.19.0
	golang.org/x/sys v0.23.0
	golang.org/x/text v0.17.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/jsmin v0.0.0-20220218165748-59f39799265f // indirect
	github.com/ebitengine/gomobile v0.0.0-20240518074828-e86332849895 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.2.0 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-text/typesetting v0.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/licenseclassifier/v2 v2.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/josephspurrier/goversioninfo v1.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/randall77/makefat v0.0.0-20210315173500-7ddd0e42c844 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
)

require (
	github.com/otiai10/copy v1.14.0
	github.com/spf13/cobra v1.8.1
)

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /Users/rpolzer/src/ebiten

// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw

// update-ebitengine-fork.sh changes:
// replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 v2.5.10-with-574925cf7a72deaf73be4c481348a7a44f7b7e19-and-cc247962703eba99eae732876496375191f16cbe-and-58eb1af1eb88915a726f29ffe9279025f98c9be8

// Pin go-locale to work around https://github.com/hajimehoshi/ebiten/issues/2899.
replace github.com/jeandeaual/go-locale => github.com/jeandeaual/go-locale v0.0.0-20220711133428-7de61946b173
