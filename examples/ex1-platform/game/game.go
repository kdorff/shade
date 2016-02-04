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
// Package game manages the main game loop.

package game

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/examples/ex1-platform/block"
	"github.com/hurricanerix/shade/examples/ex1-platform/player"
	"github.com/hurricanerix/shade/sprite"
	"github.com/hurricanerix/shade/time/clock"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Context TODO doc
type Context struct {
	Screen *display.Context
	Player *player.Player
	Walls  *sprite.Group
}

// New TODO doc
func New(screen *display.Context) (Context, error) {
	return Context{
		Screen: screen,
	}, nil
}

// Main TODO doc
func (c *Context) Main(screen *display.Context) {
	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	background, err := loadSprite("background.png", 1, 1)
	if err != nil {
		panic(err)
	}
	background.Bind(c.Screen.Program)

	sprites := sprite.NewGroup()
	c.Walls = sprite.NewGroup()

	blockSprite, err := loadSprite("block.png", 1, 1)
	if err != nil {
		panic(err)
	}
	blockSprite.Bind(screen.Program)

	for x := 0; float32(x) < screen.Width; x += 32 {
		for y := 0; float32(y) < screen.Height; y += 32 {
			if x == 0 || x == 640-32 || y == 0 || y == 480-32 {
				_, err := block.New(float32(x), float32(y), blockSprite, c.Walls)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	sprites.Add(c.Walls)

	playerSprite, err := loadSprite("player.png", 1, 1)
	if err != nil {
		panic(err)
	}
	playerSprite.Bind(screen.Program)
	p, err := player.New(float32(screen.Width)/2, float32(screen.Height)/2, playerSprite, sprites)
	if err != nil {
		panic(err)
	}

	for running := true; running; {
		dt := clock.Tick(30)

		// TODO move this somewhere else (maybe a Clear method of display
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// TODO refector events to be cleaner
		if screen.Window.ShouldClose() {
			running = !screen.Window.ShouldClose()
		}

		for _, event := range events.Get() {
			if event.Action == glfw.Press && event.Key == glfw.KeyEscape {
				running = false
				event.Window.SetShouldClose(true)
			}
			p.HandleEvent(event, dt/1000.0)
		}

		sprites.Update(dt/1000.0, c.Walls)
		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)
		background.Draw(mgl32.Vec3{0, 0, 0}, nil)
		sprites.Draw()

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}

func loadSprite(path string, framesWide, framesHigh int) (*sprite.Context, error) {
	i, err := sprite.Load(path)
	if err != nil {
		return nil, err
	}
	s, err := sprite.New(i, nil, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}
