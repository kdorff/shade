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

// Package player manages a player's state

package player

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/sprite"
)

// Player state
type Player struct {
	pos    mgl32.Vec3
	Score  int
	Sprite *sprite.Context
}

func New(x, y float32, s *sprite.Context) *Player {
	p := Player{
		pos:    mgl32.Vec3{x, y, 0.0},
		Sprite: s,
	}
	return &p
}

func (p Player) Pos() mgl32.Vec3 {
	return p.pos
}

func (p *Player) Update(dt float32, group *[]entity.Entity) {
}

func (p Player) Draw() {
	posX := p.pos[0]
	posY := p.pos[1]
	p.Sprite.DrawFrame(mgl32.Vec2{0, 0}, mgl32.Vec3{posX, posY, 0}, nil)
	p.Sprite.DrawFrame(mgl32.Vec2{0, 1}, mgl32.Vec3{posX, posY - 8, 0}, nil)
	p.Sprite.DrawFrame(mgl32.Vec2{0, 2}, mgl32.Vec3{posX, posY - 16, 0}, nil)
}
