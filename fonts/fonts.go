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

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type Location struct {
	X int
	Y int
}

// Context TODO doc
type Context struct {
	Sprite     *sprite.Context
	LocMap     map[int32]Location
	UnknownLoc Location
}

// New TODO doc
func New(s *sprite.Context, m map[int32]Location, u Location) (*Context, error) {
	c := Context{
		Sprite:     s,
		LocMap:     m,
		UnknownLoc: u,
	}
	return &c, nil
}

func SimpleASCII() (*Context, error) {
	path := fmt.Sprintf("%s/src/github.com/hurricanerix/shade/assets/font.png", os.Getenv("GOPATH"))
	i, err := sprite.Load(path)
	if err != nil {
		return nil, err
	}

	s, err := sprite.New(i, nil, 32, 3)
	if err != nil {
		return nil, err
	}

	m := make(map[int32]Location, s.Width*s.Height)
	for y := 0; y < 3; y++ {
		for x := 0; x < 32; x++ {
			m[int32((y+1)*32+x)] = Location{Y: y, X: x}
		}
	}

	u := Location{Y: 1, X: 31}

	f, err := New(s, m, u)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Bind TODO doc
func (c *Context) Bind(program uint32) {
	c.Sprite.Bind(program)
}

// DrawText TODO doc
func (c Context) DrawText(x, y, sx, sy float32, color *mgl32.Vec4, msg string) {
	cx := x
	cy := y
	addColor := mgl32.Vec4{}
	subColor := mgl32.Vec4{}

	if color != nil {
		addColor[0] = color[0]
		addColor[1] = color[1]
		addColor[2] = color[2]
		subColor[3] = 1.0 - color[3]
	}

	for _, r := range msg {
		if l, ok := c.LocMap[r]; ok {
			c.Sprite.DrawFrame(l.X, l.Y, sx, sy, cx, cy, &addColor, &subColor, nil, nil)
			cx += float32(c.Sprite.Width) * sx
		} else if r == 10 {
			cx = x
			cy -= float32(c.Sprite.Height) * sy
		} else {
			c.Sprite.DrawFrame(c.UnknownLoc.X, c.UnknownLoc.Y, sx, sy, cx, cy, &addColor, &subColor, nil, nil)
			cy += float32(c.Sprite.Width) * sx
		}
	}
}

// SizeText TODO doc
func (c Context) SizeText(sx, sy float32, msg string) (float32, float32) {
	var lx float32 = 0.0
	var cx float32 = 0.0
	var cy float32 = float32(c.Sprite.Height) * sy
	for _, r := range msg {
		if r == 10 { // code for newline
			cx = 0
			cy += float32(c.Sprite.Height) * sy
		} else {
			cx += float32(c.Sprite.Width) * sx
		}
		if cx > lx {
			lx = cx
		}
	}
	return lx, cy
}
