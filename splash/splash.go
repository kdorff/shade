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
// Package splash TODO doc

package splash

import (
	_ "image/png"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/splash/ghost"
	"github.com/hurricanerix/shade/time/clock"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func Main(screen *display.Context) {
	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	font, err := fonts.SimpleASCII()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)
	msg := "Shade SDK"
	color := mgl32.Vec4{1.0, 1.0, 1.0, 1.0}
	textSize := float32(4.0)
	_, h := font.SizeText(textSize, textSize, msg)

	g, err := ghost.New(nil)
	if err != nil {
		panic(err)
	}
	g.Bind(screen.Program)

	total := float32(0.0)
	for running := true; running; {
		dt := clock.Tick(30)
		total += dt

		// TODO move this somewhere else (maybe a Clear method of display
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// TODO refector events to be cleaner
		if screen.Window.ShouldClose() {
			running = !screen.Window.ShouldClose()
		}

		for _, event := range events.Get() {
			if event.Action == glfw.Press && event.Key == glfw.KeyEscape {
				running = false
			}
		}

		g.Update(dt, nil)

		screen.Fill(0.0, 0.0, 0.0)

		font.DrawText(screen.Width/5, screen.Height/2+h/2, textSize, textSize, &color, msg)
		g.Draw()

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
		if total > 3000.0 {
			running = false
		}
	}
}
