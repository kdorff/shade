// Copyright 2016 Richard Hawkins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ball manages a ball's state

package ball

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/examples/ex1-pong/player"
	"github.com/hurricanerix/shade/sprite"
)

// Ball state
type Ball struct {
	pos    mgl32.Vec3
	Owner  *player.Player
	Sprite *sprite.Context
}

func New(pos, dir mgl32.Vec3, owner *player.Player, s *sprite.Context) *Ball {
	b := Ball{
		pos:    pos,
		Sprite: s,
		Owner:  owner,
	}
	return &b
}

func (b Ball) Pos() mgl32.Vec3 {
	return b.pos
}

func (b *Ball) Update(dt float32, g *[]entity.Entity) {
}

func (b Ball) Draw() {
	b.Sprite.DrawFrame(mgl32.Vec2{0, 0}, b.pos, nil)
}
