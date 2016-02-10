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

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/light"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Ghost struct {
	Pos          mgl32.Vec3
	Sprite       *sprite.Context
	Light        *light.Positional
	AmbientColor mgl32.Vec4
	dx           float32
	looking      int
	fx           float32
	frame        float32
	dl           float32
}

// New TODO doc
func New(group *[]entity.Entity) (*Ghost, error) {
	// TODO should take a group in as a argument
	i, err := sprite.LoadAsset("assets/ghost.png")
	if err != nil {
		return nil, err
	}
	s, err := sprite.New(i, nil, 6, 3)
	if err != nil {
		return nil, err
	}

	c := Ghost{
		Pos:          mgl32.Vec3{-32, 480.0/2 - float32(s.Height)/2, 1.0},
		Sprite:       s,
		looking:      1,
		AmbientColor: mgl32.Vec4{0.2, 0.2, 0.2, 1.0},
	}

	light := light.Positional{
		Pos:   mgl32.Vec3{float32(s.Width) / 2, float32(s.Height) / 2, 100.0},
		Color: mgl32.Vec4{0.7, 0.7, 1.0, 1.0},
		Power: 10000,
	}
	c.Light = &light

	c.dx = 0.3
	c.fx = 0.02

	// TODO: this should probably be added outside of player
	if group != nil {
		*group = append(*group, &c)
	}
	return &c, nil
}

func (g Ghost) Type() string {
	return "ghost"
}

func (g Ghost) Label() string {
	return ""
}

// Bind TODO doc
func (c *Ghost) Bind(program uint32) error {
	return c.Sprite.Bind(program)
}

// Update TODO doc
func (c *Ghost) Update(dt float32, g *[]entity.Entity) {
	c.Pos[0] += c.dx * dt
	c.frame += c.fx * dt
	if c.Pos[0] >= 400 {
		c.dx = 0
		c.looking = -1
		c.dl = 0.01
	}

	if c.AmbientColor[0] <= 0.5 {
		c.AmbientColor[0] += c.dl
		c.AmbientColor[1] += c.dl
		c.AmbientColor[2] += c.dl
	}
	c.Light.Pos[0] = c.Pos[0] + float32(c.Sprite.Width)*2
	c.Light.Pos[1] = c.Pos[1] + float32(c.Sprite.Height)*2
}

// Draw TODO doc
func (c *Ghost) Draw(e *sprite.Effects) {
	var x float32 = c.Pos[0]
	var y float32 = c.Pos[1]

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

	c.Sprite.DrawFrame(mgl32.Vec2{float32(eyes), 0}, mgl32.Vec3{left, top, 0.0}, nil)
	c.Sprite.DrawFrame(mgl32.Vec2{float32(eyes) + 1, 0}, mgl32.Vec3{right, top, 0.0}, nil)

	c.Sprite.DrawFrame(mgl32.Vec2{0, 1}, mgl32.Vec3{left, middle, 0.0}, nil)
	c.Sprite.DrawFrame(mgl32.Vec2{1, 1}, mgl32.Vec3{right, middle, 0.0}, nil)

	c.Sprite.DrawFrame(mgl32.Vec2{float32(f), 2}, mgl32.Vec3{left, bottom, 0.0}, nil)
	c.Sprite.DrawFrame(mgl32.Vec2{float32(f) + 1, 2}, mgl32.Vec3{right, bottom, 0.0}, nil)
}
