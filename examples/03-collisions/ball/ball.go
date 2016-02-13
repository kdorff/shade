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
// Package ball TODO doc

package ball

import (
	"math"
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
type Ball struct {
	Pos    mgl32.Vec3
	Sprite *sprite.Context
	Shape  *shapes.Shape
	dx     float32
	dy     float32
}

// New TODO doc
func New(x, y, speed, angle float32, s *sprite.Context, group *[]entity.Entity) (*Ball, error) {
	// TODO should take a group in as a argument
	b := Ball{
		Pos:    mgl32.Vec3{x, y, 1.0},
		Sprite: s,
		Shape:  shapes.NewCircle(mgl32.Vec2{float32(s.Width) / 2, float32(s.Width) / 2}, float32(s.Width)/2),
	}

	b.dx = float32(math.Cos(float64(angle))) * speed
	b.dy = float32(math.Sin(float64(angle))) * speed

	// TODO: this should probably be added outside of ball
	*group = append(*group, &b)
	return &b, nil
}

func (b Ball) Type() string {
	return "ball"
}

func (b Ball) Label() string {
	return ""
}

// Bind TODO doc
func (b *Ball) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

func (b Ball) Bounds() *shapes.Shape {
	return b.Shape
}

func (b Ball) Pos2() *mgl32.Vec3 {
	return &b.Pos
}

// Update TODO doc
func (b *Ball) Update(dt float32, g []entity.Entity) {
	lastPos := mgl32.Vec3{b.Pos[0], b.Pos[1], b.Pos[2]}
	switchDx := false
	switchDy := false

	b.Pos[0] += b.dx * dt
	b.Pos[1] += b.dy * dt

	for _, c := range *sprite.Collide(b, &g, false) {
		eb := c.Entity.Bounds()
		if eb == nil {
			continue
		}
		ep := c.Entity.Pos2()
		if ep == nil {
			continue
		}

		if math.Abs(float64(c.Dir[0])) > math.Abs(float64(c.Dir[1])) {
			switchDx = true
			if eb.Type == "circle" {
				c.Entity.(*Ball).dx *= -1
			}
		} else if math.Abs(float64(c.Dir[1])) > math.Abs(float64(c.Dir[0])) {
			switchDy = true
			if eb.Type == "circle" {
				c.Entity.(*Ball).dx *= -1
				c.Entity.(*Ball).dy *= -1
			}
		} else {
			switchDx = true
			switchDy = true
			if eb.Type == "circle" {
				c.Entity.(*Ball).dx *= -1
				c.Entity.(*Ball).dy *= -1
			}
		}
	}
	if switchDx {
		b.Pos[0] = lastPos[0]
		b.dx *= -1
	}
	if switchDy {
		b.Pos[1] = lastPos[1]
		b.dy *= -1
	}
}

// Draw TODO doc
func (b Ball) Draw() {
	b.Sprite.Draw(b.Pos, nil)
}
