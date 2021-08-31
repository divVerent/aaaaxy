// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trigger

import (
	"fmt"
	go_image "image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// ForceField, when hit by the player, sends the player away from it.
type ForceField struct {
	mixins.NonSolidTouchable
	World  *engine.World
	Entity *engine.Entity

	AlphaMod  float64
	AnimFrame int
	Active    bool

	TouchedFrame int
	ShockSound   *sound.Sound
	SourceImg    *ebiten.Image
}

const (
	forceFieldStrength = 1536 * constants.SubPixelScale / engine.GameTPS
	ffFadeFrames       = 12
	ffActiveThreshold  = 4
	ffAlphaBrownian    = 0.5
	ffAlphaMin         = 0.5
	ffAlphaMax         = 1.0
)

func (f *ForceField) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	f.NonSolidTouchable.Init(w, e)
	f.World = w
	f.Entity = e
	w.SetOpaque(e, false)
	w.SetSolid(e, true)
	w.SetZIndex(e, constants.ForceFieldZ)

	// TODO: change to a dedicated sound.
	var err error
	f.ShockSound, err = sound.Load("forcefield.ogg")
	if err != nil {
		return fmt.Errorf("could not load jump sound: %v", err)
	}

	f.SourceImg, err = image.Load("sprites", "forcefield.png")
	if err != nil {
		return fmt.Errorf("failed to load forcefield sprite: %v", err)
	}

	// Force fields always spawn active.
	f.Active = true
	f.AnimFrame = ffFadeFrames
	f.AlphaMod = (ffAlphaMin + ffAlphaMax) / 2

	return nil
}

func (f *ForceField) Despawn() {}

func (f *ForceField) Update() {
	// Respond to button.
	active := true
	playerAbilities := f.World.Player.Impl.(interfaces.Abilityer)
	if playerAbilities.HasAbility("control") {
		playerButtons := f.World.Player.Impl.(interfaces.ActionPresseder)
		if playerButtons.ActionPressed() {
			active = false
		}
	}

	// Brownian motion.
	f.AlphaMod = f.AlphaMod + (2*rand.Float64()-1)*ffAlphaBrownian
	if f.AlphaMod < ffAlphaMin {
		f.AlphaMod = ffAlphaMin
	}
	if f.AlphaMod > ffAlphaMax {
		f.AlphaMod = ffAlphaMax
	}

	if active {
		f.AnimFrame++
	} else {
		f.AnimFrame--
	}
	if f.AnimFrame <= 0 {
		f.Entity.Alpha = 0
		f.AnimFrame = 0
	} else if f.AnimFrame >= ffFadeFrames {
		f.Entity.Alpha = f.AlphaMod
		f.AnimFrame = ffFadeFrames
	} else {
		alpha := float64(f.AnimFrame) / float64(ffFadeFrames)
		f.Entity.Alpha = f.AlphaMod * alpha
	}
	f.Active = f.AnimFrame >= ffActiveThreshold
	f.World.SetSolid(f.Entity, f.Active)

	// Set image to random subsection.
	// Need to take orientation into account to do that.
	gotW, gotH := f.SourceImg.Size()
	wantW, wantH := f.Entity.Rect.Size.DX, f.Entity.Rect.Size.DY
	if f.Entity.Orientation.Right.DX == 0 {
		wantW, wantH = wantH, wantW
	}
	xOffset, yOffset := rand.Intn(gotW-wantW+1), rand.Intn(gotH-wantH+1)
	f.Entity.Image = f.SourceImg.SubImage(go_image.Rectangle{
		Min: go_image.Point{
			X: xOffset,
			Y: yOffset,
		},
		Max: go_image.Point{
			X: xOffset + wantW,
			Y: yOffset + wantH,
		},
	}).(*ebiten.Image)

	// Regular updating.
	f.NonSolidTouchable.Update()
	if f.TouchedFrame > 0 {
		f.TouchedFrame--
	}
}

func (f *ForceField) Touch(other *engine.Entity) {
	// Maybe it has been turned off?
	if !f.Active {
		return
	}
	// Do we actually touch the player?
	if other != f.World.Player {
		return
	}
	p := other.Impl.(interfaces.Physics)
	// Require player to leave before jumping the player again.
	prevTouchedFrame := f.TouchedFrame
	f.TouchedFrame = 2
	if prevTouchedFrame > 0 {
		return
	}
	// TODO: compute jump vector.
	// Should point away from jumppad, and 45 degrees at the corners.
	// So, look at center-center vector.
	cc := other.Rect.Center().Delta(f.Entity.Rect.Center())
	if cc.IsZero() {
		// Can't jump if overlapping at center. Oops. Shouldn't happen though.
		log.Errorf("Forcefield: refusing to jump player due to full overlap.")
		return
	}
	// Scale so if touching precisely the corner, it is 45 degrees.
	cc.DX *= (other.Rect.Size.DY + f.Entity.Rect.Size.DY)
	cc.DY *= (other.Rect.Size.DX + f.Entity.Rect.Size.DX)
	// Scale to strength.
	away := cc.WithLength(forceFieldStrength)
	// Perform the jump.
	p.SetVelocityForJump(away)
	f.ShockSound.Play()
}

func init() {
	engine.RegisterEntityType(&ForceField{})
}
