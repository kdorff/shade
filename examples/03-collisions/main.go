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
// Package app manages the main game loop.

package main

import (
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/transylvania/display"
	"github.com/hurricanerix/transylvania/events"
	"github.com/hurricanerix/transylvania/examples/03-collisions/ball"
	"github.com/hurricanerix/transylvania/examples/ex1-platform/block"
	"github.com/hurricanerix/transylvania/sprite"
	"github.com/hurricanerix/transylvania/time/clock"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	screen, err := display.SetMode("03-collisions", windowWidth, windowHeight)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	sprites := sprite.NewGroup()

	walls := sprite.NewGroup()
	for x := 0; float32(x) < screen.Width; x += 32 {
		for y := 0; float32(y) < screen.Height; y += 32 {
			if x == 0 || x == 640-32 || y == 0 || y == 480-32 {
				b, err := block.New(walls)
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
	sprites.Add(walls)

	b, err := ball.New(sprites)
	if err != nil {
		panic(err)
	}
	b.Bind(screen.Program)

	sprites.Bind(screen.Program)
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
		}

		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)

		sprites.Update(dt/1000.0, walls)
		sprites.Draw()

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}

}
