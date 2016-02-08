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
// Package camera TODO doc

package camera

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Context struct {
	Pos    mgl32.Vec3
	Width  float32
	Height float32
	Offset mgl32.Vec2
	Top    float32
	Bottom float32
	Left   float32
	Right  float32

	ProjMatrix    mgl32.Mat4
	viewMatrixLoc int32
	ViewMatrix    mgl32.Mat4
}

func New() (*Context, error) {
	c := Context{
		Pos:    mgl32.Vec3{},
		Width:  640,
		Height: 480,
	}
	return &c, nil
}

func (c *Context) Bind(program uint32) {
	var left, right, top, bottom, near, far float32
	right = float32(c.Width)
	top = float32(c.Height)
	near = 0.1
	far = 100.0
	c.ProjMatrix = mgl32.Ortho(left, right, bottom, top, near, far)
	projUniform := gl.GetUniformLocation(program, gl.Str("ProjMatrix\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &c.ProjMatrix[0])

	c.viewMatrixLoc = gl.GetUniformLocation(program, gl.Str("ViewMatrix\x00"))
	c.Move(mgl32.Vec3{})
	gl.UniformMatrix4fv(c.viewMatrixLoc, 1, false, &c.ViewMatrix[0])
}

func (c *Context) Move(pos mgl32.Vec3) {
	pos[0] -= c.Offset[0]
	pos[1] -= c.Offset[1]

	if c.Left != 0 && pos[0] < c.Left {
		pos[0] = c.Left
	} else if c.Right != 0 && pos[0] > c.Right {
		pos[0] = c.Right
	}
	if c.Top != 0 && pos[1] > c.Top {
		pos[1] = c.Top
	} else if c.Bottom != 0 && pos[1] < c.Bottom {
		pos[1] = c.Bottom
	}

	lerp := float32(0.1)
	c.Pos[0] = c.Pos[0] + (pos[0]-c.Pos[0])*lerp
	c.Pos[1] = c.Pos[1] + (pos[1]-c.Pos[1])*lerp

	var eye, center, up mgl32.Vec3
	eye = mgl32.Vec3{c.Pos[0], c.Pos[1], 7.0}
	center = mgl32.Vec3{c.Pos[0], c.Pos[1], -1.0}
	up = mgl32.Vec3{0.0, 1.0, 0.0}
	c.ViewMatrix = mgl32.LookAtV(eye, center, up)

	gl.UniformMatrix4fv(c.viewMatrixLoc, 1, false, &c.ViewMatrix[0])
}

func (c *Context) Update(dt float32) {
}
