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
	"github.com/hurricanerix/transylvania/display"
	"github.com/hurricanerix/transylvania/events"
	"github.com/hurricanerix/transylvania/examples/ex1-platform/block"
	"github.com/hurricanerix/transylvania/examples/ex1-platform/player"
	"github.com/hurricanerix/transylvania/sprite"
	"github.com/hurricanerix/transylvania/time/clock"
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

	background, err := sprite.Load("background.png", 1, 1)
	if err != nil {
		panic(err)
	}
	background.Bind(c.Screen.Program)

	sprites := sprite.NewGroup()
	p, err := player.New(sprites)
	if err != nil {
		panic(err)
	}

	c.Walls = sprite.NewGroup()
	// TODO: should only load image data once.
	//block, err := sprite.Load("block.png", 1)
	//if err != nil {
	//	panic(err)
	//}
	//println(block)

	for x := 0; x < 640; x += 32 {
		for y := 0; y < 480; y += 32 {
			if x == 0 || x == 640-32 || y == 0 || y == 480-32 {
				// NOTE: original python code
				//    wall = pygame.sprite.Sprite(self.walls)
				//    wall.image = block
				//    wall.rect = pygame.rect.Rect((x, y), block.get_size())
				b, err := block.New(c.Walls)
				if err != nil {
					panic(err)
				}
				b.Rect.X = float32(x)
				b.Rect.Y = float32(y)
				b.Rect.Width = float32(b.Image.Width)
				b.Rect.Width = float32(b.Image.Height)
			}
		}
	}
	sprites.Add(c.Walls)

	sprites.Bind(c.Screen.Program)
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
		background.Draw(0, 0)
		sprites.Draw()

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}