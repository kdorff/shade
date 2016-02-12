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
// Package player TODO doc

package player

import (
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Ball TODO doc
type Player struct {
	Pos       mgl32.Vec3
	Sprites   []*sprite.Context
	Shapes    []*shapes.Shape
	Collision bool
	With      string
	current   int
}

// New TODO doc
func New(x, y float32, sprites []*sprite.Context, shps []*shapes.Shape, group *[]entity.Entity) (*Player, error) {
	// TODO should take a group in as a argument
	b := Player{
		Pos:     mgl32.Vec3{x, y, 1.0},
		Sprites: sprites,
		Shapes:  shps,
		current: 1,
	}

	// TODO: this should probably be added outside of ball
	if group != nil {
		*group = append(*group, &b)
	}
	return &b, nil
}

func (p Player) Type() string {
	return "ball"
}

func (p Player) Label() string {
	return ""
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Sprites[p.current].Bind(program)
}

func (p Player) Bounds() *shapes.Shape {
	return p.Shapes[p.current]
}

func (p Player) Pos2() *mgl32.Vec3 {
	return &p.Pos
}

// Update TODO doc
func (p *Player) Update(dt float32, g []entity.Entity) {
	p.Collision = false
	p.With = ""
	for _, e := range *sprite.Collide(p, &g, false) {
		p.Collision = true
		p.With = e.Type()
	}
}

// Draw TODO doc
func (p Player) Draw() {
	p.Sprites[p.current].Draw(p.Pos, nil)
}

func (p *Player) NextShape() {
	p.current += 1
	if p.current >= len(p.Sprites) {
		p.current = 0
	}
}
