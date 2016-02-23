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
	"fmt"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Ball TODO doc
type Player struct {
	pos       mgl32.Vec3
	Sprites   []sprite.Context
	Font      fonts.Context
	Shapes    []shapes.Shape
	Collision *entity.Collision
	With      string
	current   int
}

// New TODO doc
func New(x, y float32, sprites []sprite.Context, shapes []shapes.Shape, font fonts.Context) Player {
	b := Player{
		pos:     mgl32.Vec3{x, y, 1.0},
		Sprites: sprites,
		Font:    font,
		Shapes:  shapes,
		current: 1,
	}
	return b
}

func (p Player) Pos() mgl32.Vec3 {
	return p.pos
}

// SetPos of player
func (p *Player) SetPos(pos mgl32.Vec3) {
	p.pos = pos
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Sprites[p.current].Bind(program)
}

func (p Player) Bounds() shapes.Shape {
	return p.Shapes[p.current]
}

// Update TODO doc
func (p *Player) Update(dt float32, g []entity.Collider) {
	p.Collision = nil
	p.With = ""
	for _, c := range entity.Collide(*p, g) {
		p.Collision = &c
	}
}

// Draw TODO doc
func (p Player) Draw() {
	msg := fmt.Sprintf("(%.0f,%.0f)\n", p.pos[0], p.pos[1])
	if p.Collision == nil {
		msg += fmt.Sprintf("Collision: nil\n")
	} else {
		c, ok := p.Collision.Hit.(entity.Entity)
		if ok {
			msg += fmt.Sprintf("Collision: {\n")
			msg += fmt.Sprintf("  Type: %T\n", c)
			msg += fmt.Sprintf("  Dir: (%.1f,%.1f,%.1f)\n", p.Collision.Dir[0], p.Collision.Dir[1], p.Collision.Dir[2])
			msg += fmt.Sprintf("}\n")
		}
	}

	if p.Shapes[p.current].Type == "rect" {
		msg += fmt.Sprintf("Data: [\n")
		msg += fmt.Sprintf("  Left: %.0f\n", p.Shapes[p.current].Data[0])
		msg += fmt.Sprintf("  Right: %.0f\n", p.Shapes[p.current].Data[1])
		msg += fmt.Sprintf("  Top: %.0f\n", p.Shapes[p.current].Data[2])
		msg += fmt.Sprintf("  Bottom: %.0f\n", p.Shapes[p.current].Data[3])
		msg += fmt.Sprintf("]\n")
	} else {
		msg += fmt.Sprintf("Data: [\n")
		msg += fmt.Sprintf("  Center: (%.0f, %.0f)\n", p.Shapes[p.current].Data[0], p.Shapes[p.current].Data[1])
		msg += fmt.Sprintf("  Radius: %.0f\n", p.Shapes[p.current].Data[2])
		msg += fmt.Sprintf("]\n")
	}
	efx := sprite.Effects{
		Scale: mgl32.Vec3{2.0, 2.0, 1.0},
	}
	p.Font.DrawText(mgl32.Vec3{p.pos[0], p.pos[1] - 16, 0}, &efx, msg)
	p.Sprites[p.current].Draw(p.pos, nil)
}

func (p *Player) NextShape() {
	p.current += 1
	if p.current >= len(p.Sprites) {
		p.current = 0
	}
}
