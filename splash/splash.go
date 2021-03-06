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
	"github.com/hurricanerix/shade/camera"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/splash/ghost"
	"github.com/hurricanerix/shade/sprite"
	"github.com/hurricanerix/shade/time/clock"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func Main(screen *display.Context) {
	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(screen.Program)

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	font, err := loadFont()
	if err != nil {
		panic(err)
	}
	font.Sprite.Bind(screen.Program)
	msg := "Shade SDK"

	g := ghost.New()
	g.Bind(screen.Program)

	total := float32(0.0)
	for running := true; running; {
		dt := clock.Tick(30)
		total += dt

		screen.Fill(0.0, 0.0, 0.0)

		// TODO move this somewhere else (maybe a Clear method of display
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, event := range events.Get() {
			if event.Type == events.KeyUp && event.Key == glfw.KeyEscape {
				// Send window close event
				screen.Close()
			}
			if event.Type == events.WindowClose {
				// Handle window close
				running = false
			}
		}

		g.Update(dt, nil)

		effect := sprite.Effects{
			Scale:          mgl32.Vec3{1.0, 1.0, 1.0},
			Tint:           mgl32.Vec4{1.0, 1.0, 1.0, 0.0},
			EnableLighting: true,
			AmbientColor:   g.AmbientColor,
			Light:          *g.Light}
		_, h := font.SizeText(&effect, msg)

		pos := mgl32.Vec3{
			screen.Width / 5,
			screen.Height/2 + h/2,
			0}

		font.DrawText(pos, &effect, msg)
		g.Draw()

		screen.Flip()
		events.Poll()

		if total > 3000.0 {
			running = false
		}
	}
}

func loadFont() (*fonts.Context, error) {
	c, err := sprite.LoadAsset("assets/splash-font.png")
	if err != nil {
		return nil, err
	}

	n, err := sprite.LoadAsset("assets/splash-font.normal.png")
	if err != nil {
		return nil, err
	}

	s, err := sprite.New(c, n, 32, 3)
	if err != nil {
		return nil, err
	}

	m := make(map[int32]mgl32.Vec2, s.Width*s.Height)
	for y := 0; y < 3; y++ {
		for x := 0; x < 32; x++ {
			m[int32((y+1)*32+x)] = mgl32.Vec2{float32(x), float32(y)}
		}
	}

	u := mgl32.Vec2{31, 1}

	f, err := fonts.New(s, m, u)
	if err != nil {
		return nil, err
	}
	return f, nil
}
