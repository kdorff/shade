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
	"github.com/hurricanerix/shade/camera"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Config TODO doc
type Config struct {
	DevMode bool
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
func (c *Context) Main(screen *display.Context, config Config) {
	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Offset = mgl32.Vec2{200, 100}
	cam.TopStop = 64 * 6.5 // TODO: should be 64 x 14
	cam.RightStop = 64 * 54
	cam.LeftStop = 1
	cam.Bind(c.Screen.Program)

	//clock, err := clock.New()
	//if err != nil {
	//panic(err)
	//}

	font, err := fonts.SimpleASCII()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)

	for running := true; running; {

		screen.Fill(0, 0, 0)

		//dt := clock.Tick(30)

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

		if config.DevMode {
			deveff := sprite.Effects{
				EnableLighting: false,
				Scale:          mgl32.Vec3{2.0, 2.0, 1.0},
				Tint:           mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
			}
			msg := "Dev Mode!\n"
			font.DrawText(mgl32.Vec3{cam.Left + 20, cam.Top - 40, 0}, &deveff, msg)
		}
		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}

func loadSpriteAsset(colorName, normalName string, framesWide, framesHigh int) (*sprite.Context, error) {
	c, err := sprite.LoadAsset(colorName)
	if err != nil {
		return nil, err
	}
	n, err := sprite.LoadAsset(normalName)
	if err != nil {
		return nil, err
	}
	s, err := sprite.New(c, n, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}
