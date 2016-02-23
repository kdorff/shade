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
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/light"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Player struct {
	pos      mgl32.Vec3
	Shape    *shapes.Shape
	Sprite   *sprite.Context
	Light    *light.Positional
	Facing   float32
	resting  bool
	dy       float32
	leftKey  bool
	rightKey bool
	jumpKey  bool
}

// New TODO doc
func New(x, y float32, s *sprite.Context) *Player {
	// TODO should take a group in as a argument
	p := Player{
		pos:    mgl32.Vec3{x, y, 1.0},
		Shape:  shapes.NewRect(32, 96, 0, 96),
		Sprite: s,
		Facing: 2,
	}
	light := light.Positional{
		Pos:   mgl32.Vec3{p.pos[0], float32(s.Height), 50.0},
		Color: mgl32.Vec4{0.7, 0.7, 1.0, 1.0},
		Power: 10000,
	}
	p.Light = &light
	return &p
}

func (p Player) Bounds() shapes.Shape {
	return *p.Shape
}

func (p Player) Pos() mgl32.Vec3 {
	return p.pos
}

// HandleEvent TODO doc
func (p *Player) HandleEvent(event events.Event, dt float32) {
	// TODO: move this to SDK to handle things like holding Left & Right at the same time correctly

	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeyLeft {
		p.leftKey = true
	}
	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeyRight {
		p.rightKey = true
	}
	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeySpace {
		p.jumpKey = true
	}
	if event.Action == glfw.Release && event.Key == glfw.KeyLeft {
		p.leftKey = false
	}
	if event.Action == glfw.Release && event.Key == glfw.KeyRight {
		p.rightKey = false
	}
	if event.Action == glfw.Release && event.Key == glfw.KeySpace {
		p.jumpKey = false
	}
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Sprite.Bind(program)
}

// Update TODO doc
func (p *Player) Update(dt float32, group *[]entity.Entity) {
	lastPos := mgl32.Vec3{p.pos[0], p.pos[1], p.pos[2]}

	if p.leftKey {
		p.pos[0] -= 300.0 * dt
		p.Light.Pos[0] = p.pos[0]
		p.Facing = 1
	}
	if p.rightKey {
		p.pos[0] += 300.0 * dt
		p.Facing = 2
		p.Light.Pos[0] = p.pos[0] + float32(p.Sprite.Width)
	}
	if p.resting && p.jumpKey {
		p.dy = 1500.0
	}
	p.dy = float32(math.Min(float64(1500.0), float64(p.dy-40.0)))

	p.pos[1] += p.dy * dt

	newPos := &p.pos
	p.resting = false

	if p.pos[1] < 127 {
		p.resting = true
		p.pos[1] = 128
		p.dy = 0.0
	}

	var cgroup []entity.Collider
	for i := range *group {
		if c, ok := (*group)[i].(entity.Collider); ok {
			cgroup = append(cgroup, c)
		}
	}
	for _, c := range entity.Collide(p, &cgroup, false) {
		pos := c.Hit.Pos()

		if c.Dir[0] > 0.5 {
			newPos[0] = lastPos[0]
		} else if c.Dir[0] < -0.5 {
			newPos[0] = lastPos[0]
		}

		if c.Dir[1] > 0.5 {
			// Hit top of tile
			newPos[1] = pos[1]
		} else if c.Dir[1] < -0.5 {
			// Hit bottom of tile
			p.resting = true
			newPos[1] = pos[1] + 64 + 1
			p.dy = 0.0
		}
	}
	p.Light.Pos[1] = p.pos[1] + float32(p.Sprite.Height)

}

// Draw TODO doc
func (p *Player) Draw() {
	//e *sprite.Effects) {
	//p.Sprite.DrawFrame(mgl32.Vec2{1, p.Facing}, p.Pos, e)
	p.Sprite.DrawFrame(mgl32.Vec2{1, p.Facing}, p.pos, nil)
}
