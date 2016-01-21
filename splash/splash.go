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
	"fmt"
	_ "image/png"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/sprite"
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

	logopath := fmt.Sprintf("%s/src/github.com/hurricanerix/shade/assets/logo.png", os.Getenv("GOPATH"))

	background, err := sprite.Load(logopath, 8, 1)
	if err != nil {
		panic(err)
	}
	background.Bind(screen.Program)

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

		screen.Fill(0.0, 0.0, 0.0)

		// TODO: read from actual window width/height
		var x float32 = (640.0 - 96.0/2.0) / 2.0
		var y float32 = (480.0 - 96.0/2.0) / 2.0

		top := y + 64.0
		middle := y + 32.0
		bottom := y + 0.0

		left := x + 0.0
		center := x + 32.0
		right := x + 64.0

		background.DrawFrame(0, 0, 1.0, 1.0, right, top)
		background.DrawFrame(1, 0, 1.0, 1.0, center, top)
		background.DrawFrame(2, 0, 1.0, 1.0, right, top)

		background.DrawFrame(3, 0, 1.0, 1.0, center, middle)
		background.DrawFrame(4, 0, 1.0, 1.0, right, middle)

		background.DrawFrame(5, 0, 1.0, 1.0, left, bottom)
		background.DrawFrame(6, 0, 1.0, 1.0, center, bottom)
		background.DrawFrame(7, 0, 1.0, 1.0, right, bottom)

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
		if total > 3000.0 {
			running = false
		}
	}
}
