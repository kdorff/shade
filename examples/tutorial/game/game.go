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
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/transylvania/display"
	"github.com/hurricanerix/transylvania/events"
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

	player, err := sprite.Load("player.png")
	if err != nil {
		panic(err)
	}

	running := true
	for running {
		clock.Tick(30)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		if screen.Window.ShouldClose() {
			running = !screen.Window.ShouldClose()
		}
		for _, event := range events.Get() {
			fmt.Println(event)
			if event.Action == glfw.Press && event.Key == glfw.KeyEscape {
				running = false
				event.Window.SetShouldClose(true)
			}
		}

		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)
		screen.Blit(player, 320, 240)
		screen.Flip()

		glfw.PollEvents()
	}
}
