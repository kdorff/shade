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
	"bufio"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/camera"
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
}

// New TODO doc
func New(screen *display.Context) (Context, error) {
	return Context{
		Screen: screen,
	}, nil
}

type Scene struct {
	Sprites *sprite.Group
	Player  *player.Player
	Walls   *sprite.Group
}

// Main TODO doc
func (c *Context) Main(screen *display.Context) {
	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Offset = mgl32.Vec2{200, 100}
	cam.Top = 64 * 6.5 // TODO: should be 64 x 14
	cam.Right = 64 * 54
	cam.Left = 1
	cam.Bind(c.Screen.Program)

	scene, err := loadMap("map.data")
	if err != nil {
		panic(err)
	}

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	scene.Sprites.Bind(screen.Program)

	for running := true; running; {

		screen.Fill(0, 0, 0)

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
			scene.Player.HandleEvent(event, dt/1000.0)
		}

		cam.Move(mgl32.Vec3{scene.Player.Rect.X, scene.Player.Rect.Y, 0.0})

		scene.Sprites.Update(dt/1000.0, scene.Walls)

		effect := sprite.Effects{
			Scale:          mgl32.Vec3{1.0, 1.0, 1.0},
			EnableLighting: true,
			AmbientColor:   mgl32.Vec4{0.3, 0.3, 0.3, 1.0},
			Light:          *scene.Player.Light}

		scene.Sprites.Draw(&effect)

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}

// sprites, player, collidable
func loadMap(path string) (*Scene, error) {
	scene := Scene{}

	scene.Sprites = sprite.NewGroup()
	scene.Walls = sprite.NewGroup()

	playerSprite, err := loadSpriteAsset("assets/gopher128x128.png", "assets/gopher128x128.normal.png", 1, 1)
	if err != nil {
		return &scene, err
	}
	blockSprite, err := loadSpriteAsset("assets/block64x64.png", "assets/block64x64.normal.png", 1, 1)
	if err != nil {
		return &scene, err
	}

	f, err := os.Open(path)
	if err != nil {
		return &scene, err
	}

	s := bufio.NewScanner(f)
	count := 0
	lines := []string{}
	for s.Scan() {
		count += 1
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return &scene, err
	}

	x := float32(0)
	y := float32(0)
	for i := count - 1; i >= 0; i -= 1 {
		for _, c := range lines[i] {
			switch c {
			case '#':
				_, err := block.New(float32(x), float32(y), blockSprite, scene.Walls)
				if err != nil {
					panic(err)
				}
			case 'S':
				p, err := player.New(x, y, playerSprite, scene.Sprites)
				if err != nil {
					panic(err)
				}
				scene.Player = p
			}
			x += float32(blockSprite.Width)
		}
		x = 0
		y += float32(blockSprite.Height)
	}
	scene.Sprites.Add(scene.Walls)

	return &scene, nil
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
