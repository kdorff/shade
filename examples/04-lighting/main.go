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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/light"
	"github.com/hurricanerix/shade/sprite"
)

const windowWidth = 640
const windowHeight = 480

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	screen, err := display.SetMode("03-lighting", windowWidth, windowHeight)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}
	ambientColor := mgl32.Vec4{0.2, 0.2, 0.2, 1.0}

	face, err := loadSprite("color.png", "normal.png", 1, 1)
	if err != nil {
		panic(err)
	}
	face.Bind(screen.Program)

	light := light.Positional{
		Pos:   mgl32.Vec3{0.5, 0.5, 1.0},
		Color: mgl32.Vec4{0.8, 0.8, 1.0, 1.0},
		Power: 1000,
	}

	for running := true; running; {
		screen.Fill(0.0, 0.0, 0.0)

		// TODO move this somewhere else (maybe a Clear method of display
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// TODO refector events to be cleaner
		if screen.Window.ShouldClose() {
			running = !screen.Window.ShouldClose()
		}

		for _, event := range events.Get() {
			if event.KeyEvent && event.Action == glfw.Press && event.Key == glfw.KeyEscape {
				running = false
				event.Window.SetShouldClose(true)
			}
			if !event.KeyEvent {
				light.Pos[0] = event.X
				light.Pos[1] = float32(windowHeight) - event.Y
			}
		}
		face.DrawFrame(0, 0, 1.0, 1.0, windowWidth/2-float32(face.Width)/2, windowHeight/2-float32(face.Height)/2, nil, nil, &ambientColor, &light)

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}

}

func loadSprite(colorPath, normalPath string, framesWide, framesHigh int) (*sprite.Context, error) {
	c, err := sprite.Load(colorPath)
	if err != nil {
		return nil, err
	}

	n, err := sprite.Load(normalPath)
	if err != nil {
		return nil, err
	}

	s, err := sprite.New(c, n, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}
