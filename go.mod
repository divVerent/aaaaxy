module github.com/divVerent/aaaaxy

go 1.24.0

toolchain go1.24.6

require (
	github.com/Microsoft/go-winio v0.6.2
	github.com/adrg/xdg v0.5.3
	github.com/akavel/rsrc v0.10.2
	github.com/fardog/tmx v0.0.0-20210504210836-02c45f261672
	github.com/google/go-cmp v0.7.0
	github.com/google/go-licenses v1.6.1-0.20230903011517-706b9c60edd4
	github.com/hajimehoshi/bitmapfont/v3 v3.3.0
	github.com/hajimehoshi/ebiten/v2 v2.9.6
	github.com/jeandeaual/go-locale v0.0.0-20250612000132-0ef82f21eade
	github.com/leonelquinteros/gotext v1.7.2
	github.com/lestrrat-go/strftime v1.1.1
	github.com/lucasb-eyer/go-colorful v1.3.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/ncruces/zenity v0.10.14
	github.com/zachomedia/go-bdf v0.0.0-20220611021443-a3af701111be
	golang.org/x/image v0.34.0
	golang.org/x/sys v0.39.0
	golang.org/x/text v0.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/jsmin v1.0.0 // indirect
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.4.0 // indirect
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/licenseclassifier/v2 v2.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jezek/xgb v1.2.0 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/josephspurrier/goversioninfo v1.5.0 // indirect
	github.com/otiai10/mint v1.6.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/randall77/makefat v0.0.0-20210315173500-7ddd0e42c844 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
)

require (
	github.com/otiai10/copy v1.14.1
	github.com/spf13/cobra v1.10.2
)

// Use when playing around with extended/fixed tmx versions.
// replace github.com/fardog/tmx => github.com/divVerent/tmx v0.0.0-20210504110059-b8d75006ad02

// For debugging:
// replace github.com/hajimehoshi/ebiten/v2 => /Users/rpolzer/src/ebiten

// replace github.com/go-gl/glfw/v3.3/glfw => /home/rpolzer/src/go-gl-glfw/v3.3/glfw

// Pin go-locale to work around https://github.com/hajimehoshi/ebiten/issues/2899.
// replace github.com/jeandeaual/go-locale => github.com/jeandeaual/go-locale v0.0.0-20220711133428-7de61946b173

// update-ebitengine-fork.sh changes:
// replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 v2.7.9-with-99ffe09b63e0d906cc1f502c24f4d2325e6cc09d
