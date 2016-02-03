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
	"fmt"
	"math"
	"runtime"

	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Ball TODO doc
type Ball struct {
	Sprite *sprite.Context
	Rect   *shapes.Rect
	dx     float32
	dy     float32
}

// New TODO doc
func New(x, y, speed, angle float32, s *sprite.Context, group *sprite.Group) (*Ball, error) {
	// TODO should take a group in as a argument
	b := Ball{
		Sprite: s,
	}

	b.dx = float32(math.Cos(float64(angle))) * speed
	b.dy = float32(math.Sin(float64(angle))) * speed

	rect, err := shapes.NewRect(x, y, float32(b.Sprite.Width), float32(b.Sprite.Height))
	if err != nil {
		return &b, fmt.Errorf("could create rect: %v", err)
	}
	b.Rect = rect

	// TODO: this should probably be added outside of ball
	group.Add(&b)
	return &b, nil
}

// Bind TODO doc
func (b *Ball) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

// Update TODO doc
func (b *Ball) Update(dt float32, g *sprite.Group) {
	lastR := shapes.Rect{b.Rect.X, b.Rect.Y, b.Rect.Width, b.Rect.Height}

	b.Rect.X += b.dx * dt
	b.Rect.Y += b.dy * dt

	newR := b.Rect

	switchDx := false
	switchDy := false

	for _, cell := range sprite.Collide(b, g, false) {
		for cb := range cell.Bounds() {
			if lastR.Left() <= cb.Right() && lastR.Right() >= cb.Left() {
				switchDx = true
			}
			if lastR.Bottom() <= cb.Top() && lastR.Top() >= cb.Bottom() {
				switchDy = true
			}
		}
	}
	if switchDx {
		newR.X = lastR.X
		b.dx *= -1
	}
	if switchDy {
		newR.Y = lastR.Y
		b.dy *= -1
	}
}

// Draw TODO doc
func (b *Ball) Draw() {
	b.Sprite.Draw(b.Rect.X, b.Rect.Y)
}

// Bounds TODO doc
func (b *Ball) Bounds() chan shapes.Rect {
	ch := make(chan shapes.Rect, 1)
	ch <- *(b.Rect)
	close(ch)
	return ch
}
