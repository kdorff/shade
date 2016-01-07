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

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/transylvania/events"
	"github.com/hurricanerix/transylvania/rect"
	"github.com/hurricanerix/transylvania/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Player struct {
	Image *sprite.Context
	Rect  *rect.Rect
}

// New TODO doc
func New(group *sprite.Group) (*Player, error) {
	// TODO should take a group in as a argument
	p := Player{}

	player, err := sprite.Load("player.png")
	if err != nil {
		return &p, fmt.Errorf("could not load player: %v", err)
	}
	p.Image = player

	rect, err := rect.New(320, 240, 0, 0)
	if err != nil {
		return &p, fmt.Errorf("could create rect: %v", err)
	}
	p.Rect = rect

	// TODO: this should probably be added outside of player
	group.Add(&p)
	return &p, nil
}

// HandleEvent TODO doc
func (p *Player) HandleEvent(event events.Event, dt float32) {
	// TODO: move this to SDK to handle things like holding Left & Right at the same time correctly
	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeyLeft {
		p.Rect.X -= 300.0 * dt
	}
	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeyRight {
		p.Rect.X += 300.0 * dt
	}
	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeyUp {
		p.Rect.Y += 300.0 * dt
	}
	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeyDown {
		p.Rect.Y -= 300.0 * dt
	}
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Image.Bind(program)
}

// Update TODO doc
func (p *Player) Update(dt float32) {
	// TODO: Myabe handeling events should be done here, and not in a seperate "HandleEvents" func?
}

// Draw TODO doc
func (p *Player) Draw() {
	p.Image.Draw(p.Rect.X, p.Rect.Y)
}
