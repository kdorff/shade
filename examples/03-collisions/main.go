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
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/shade/camera"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/examples/03-collisions/ball"
	"github.com/hurricanerix/shade/examples/03-collisions/block"
	"github.com/hurricanerix/shade/sprite"
	"github.com/hurricanerix/shade/time/clock"
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

	cam, err := camera.New()
	if err != nil {
		panic(err)
	}
	cam.Bind(screen.Program)

	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	objects := []entity.Entity{}
	balls := []entity.Entity{}
	walls := []entity.Entity{}

	blockSprite, err := loadSprite("assets/block32x32.png", "", 2, 1)
	if err != nil {
		panic(err)
	}
	blockSprite.Bind(screen.Program)

	for x := 0; float32(x) < screen.Width; x += 32 {
		for y := 0; float32(y) < screen.Height; y += 32 {
			if x == 0 || x == 640-32 || y == 0 || y == 480-32 {
				_, err := block.New(float32(x), float32(y), blockSprite, &walls)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	_, err = block.New(float32(blockSprite.Width)*4, float32(blockSprite.Height)*4, blockSprite, &walls)
	if err != nil {
		panic(err)
	}
	_, err = block.New(float32(blockSprite.Width)*4, windowHeight-float32(blockSprite.Height)*5, blockSprite, &walls)
	if err != nil {
		panic(err)
	}
	_, err = block.New(windowWidth-float32(blockSprite.Width)*5, float32(blockSprite.Height)*4, blockSprite, &walls)
	if err != nil {
		panic(err)
	}
	_, err = block.New(windowWidth-float32(blockSprite.Width)*5, windowHeight-float32(blockSprite.Height)*5, blockSprite, &walls)
	if err != nil {
		panic(err)
	}

	for _, w := range walls {
		objects = append(objects, w)
	}

	ballSprite, err := loadSprite("assets/ball.png", "", 1, 1)
	if err != nil {
		panic(err)
	}
	ballSprite.Bind(screen.Program)

	//rand.Seed(1)
	rand.Seed(time.Now().Unix())

	b := addBall(screen.Width/2, screen.Height/2, ballSprite, &balls)
	objects = append(objects, b)

	//	sprites.Bind(screen.Program)
	for running := true; running; {
		dt := clock.Tick(30)

		screen.Fill(0.0, 0.0, 0.0)

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

			if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeySpace {
				addBall(screen.Width/2, screen.Height/2, ballSprite, &balls)
			}
		}

		for _, b := range balls {
			bs := b.(*ball.Ball)
			bs.Update(dt/1000.0, objects)
			bs.Draw()
		}
		for _, o := range objects {
			switch o.Type() {
			case "ball":
				o.(*ball.Ball).Draw()
			case "block":
				o.(*block.Block).Draw()
			}
		}

		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}

func addBall(x, y float32, s *sprite.Context, g *[]entity.Entity) *ball.Ball {
	speed := float32(rand.Intn(500) + 200)
	angle := float32(rand.Intn(360))
	b, err := ball.New(x, y, speed, angle, s, g)
	if err != nil {
		panic(err)
	}
	return b
}

func loadSprite(colorName, normalName string, framesWide, framesHigh int) (*sprite.Context, error) {
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
