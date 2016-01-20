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
// Package fonts TODO doc

package fonts

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hurricanerix/transylvania/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type location struct {
	x int
	y int
}

type Context struct {
	Image      *sprite.Context
	LocMap     map[int32]location
	UnknownLoc location
}

func New() (*Context, error) {
	c := Context{
		LocMap:     make(map[int32]location, 96),
		UnknownLoc: location{y: 1, x: 31},
	}
	//for y := 32; y < 32*3; y += 32 {
	for y := 0; y < 3; y++ {
		for x := 0; x < 32; x++ {
			c.LocMap[int32((y+1)*32+x)] = location{y: y, x: x}
		}
	}

	path := fmt.Sprintf("%s/src/github.com/hurricanerix/transylvania/assets/font.png", os.Getenv("GOPATH"))
	i, err := sprite.Load(path, 32, 3)
	if err != nil {
		return &c, err
	}
	c.Image = i

	return &c, nil
}

func (c *Context) Bind(program uint32) {
	c.Image.Bind(program)
}

func (c Context) DrawText(x, y, sx, sy float32, msg string) {
	cx := x
	cy := y
	for _, r := range msg {
		if l, ok := c.LocMap[r]; ok {
			c.Image.DrawFrame(l.x, l.y, sx, sy, cx, cy)
			cx += float32(c.Image.Width) * sx
		} else if r == 10 {
			cx = x
			cy -= float32(c.Image.Height) * sy
		} else {
			c.Image.DrawFrame(c.UnknownLoc.x, c.UnknownLoc.y, sx, sy, cx, cy)
			cy += float32(c.Image.Width) * sx
		}
	}
}
