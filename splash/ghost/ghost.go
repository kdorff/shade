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
// Package ghost TODO doc

package ghost

import (
	"math"
	"runtime"

	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Ghost struct {
	Sprite  *sprite.Context
	Rect    *shapes.Rect
	dx      float32
	looking int
	fx      float32
	frame   float32
}

// New TODO doc
func New(group *sprite.Group) (*Ghost, error) {
	// TODO should take a group in as a argument
	c := Ghost{
		looking: 1,
	}

	i, err := sprite.LoadAsset("assets/ghost.png")
	if err != nil {
		return nil, err
	}
	s, err := sprite.New(i, 6, 3)
	if err != nil {
		return nil, err
	}

	c.Sprite = s
	w := float32(c.Sprite.Width) * 2
	r, err := shapes.NewRect(-32, 480.0/2-float32(c.Sprite.Height)/2, w, w)
	if err != nil {
		return nil, err
	}
	c.Rect = r

	c.dx = 0.3
	c.fx = 0.02

	// TODO: this should probably be added outside of player
	if group != nil {
		group.Add(&c)
	}
	return &c, nil
}

// Bind TODO doc
func (c *Ghost) Bind(program uint32) error {
	return c.Sprite.Bind(program)
}

// Update TODO doc
func (c *Ghost) Update(dt float32, g *sprite.Group) {
	c.Rect.X += c.dx * dt
	c.frame += c.fx * dt
	if c.Rect.X >= 400 {
		c.dx = 0
		c.looking = -1
	}
}

// Draw TODO doc
func (c *Ghost) Draw() {
	var x float32 = c.Rect.X
	var y float32 = c.Rect.Y

	top := y + 64.0
	middle := y + 32.0
	bottom := y + 0.0

	left := x + 0.0
	right := x + 32.0

	eyes := 0
	if c.looking == -1 {
		eyes = 0
	} else if c.looking == 0 {
		eyes = 2
	} else {
		eyes = 4
	}

	f := int(math.Mod(float64(int(c.frame)), 3)) * 2

	c.Sprite.DrawFrame(eyes, 0, 1.0, 1.0, left, top, nil, nil)
	c.Sprite.DrawFrame(eyes+1, 0, 1.0, 1.0, right, top, nil, nil)

	c.Sprite.DrawFrame(0, 1, 1.0, 1.0, left, middle, nil, nil)
	c.Sprite.DrawFrame(1, 1, 1.0, 1.0, right, middle, nil, nil)

	c.Sprite.DrawFrame(f, 2, 1.0, 1.0, left, bottom, nil, nil)
	c.Sprite.DrawFrame(f+1, 2, 1.0, 1.0, right, bottom, nil, nil)
}

// Bounds TODO doc
func (c *Ghost) Bounds() chan shapes.Rect {
	ch := make(chan shapes.Rect, 1)
	ch <- *(c.Rect)
	close(ch)
	return ch
}
