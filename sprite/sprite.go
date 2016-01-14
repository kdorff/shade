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
// Package sprite manages images

package sprite

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // register PNG decode
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/transylvania/shapes"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Sprite TODO doc
type Sprite interface {
	// TODO rename this to something more interfacer
	Bind(uint32) error
	Update(float32, *Group)
	Draw()
	Bounds() shapes.Rect
}

// Context TODO doc
type Context struct {
	Image       image.Image
	Width       int
	Height      int
	frames      int
	vao         uint32
	vbo         uint32
	texLoc      uint32
	model       mgl32.Mat4
	modelMatrix int32
	tex         mgl32.Mat3
	texMatrix   int32
}

// Bind TODO doc
func (c *Context) Bind(program uint32) error {
	gl.UseProgram(program)

	colorMap := gl.GetUniformLocation(program, gl.Str("ColorMap\x00"))
	gl.Uniform1i(colorMap, 0)

	c.modelMatrix = gl.GetUniformLocation(program, gl.Str("ModelMatrix\x00"))
	gl.UniformMatrix4fv(c.modelMatrix, 1, false, &c.model[0])

	c.texMatrix = gl.GetUniformLocation(program, gl.Str("TexMatrix\x00"))
	gl.UniformMatrix3fv(c.texMatrix, 1, false, &c.tex[0])

	// TODO: These prob don't need to be re-created every time.
	if c.vao == 0 {
		gl.GenVertexArrays(1, &c.vao)
		gl.BindVertexArray(c.vao)
	}

	if c.vbo == 0 {
		gl.GenBuffers(1, &c.vbo)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	mcVertex := uint32(gl.GetAttribLocation(program, gl.Str("MCVertex\x00")))
	gl.EnableVertexAttribArray(mcVertex)
	gl.VertexAttribPointer(mcVertex, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoord0 := uint32(gl.GetAttribLocation(program, gl.Str("TexCoord0\x00")))
	gl.EnableVertexAttribArray(texCoord0)
	gl.VertexAttribPointer(texCoord0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return nil
}

// Draw TODO doc
func (c *Context) Draw(x, y float32) {
	c.DrawFrame(0, x, y)
}

func (c *Context) DrawFrame(f int, x, y float32) {
	c.model = mgl32.Ident4()
	c.model = c.model.Mul4(mgl32.Translate3D(float32(c.Width)/2.0, float32(c.Height)/2.0, 0.0))
	c.model = c.model.Mul4(mgl32.Translate3D(x, y, 0.0))
	c.model = c.model.Mul4(mgl32.Scale3D(float32(c.Width), float32(c.Height), 0.0))
	gl.UniformMatrix4fv(c.modelMatrix, 1, false, &c.model[0])

	c.tex = mgl32.Ident3()
	c.tex = c.tex.Mul3(mgl32.Scale2D(1.0/float32(c.frames), 1.0))
	c.tex = c.tex.Mul3(mgl32.Translate2D(float32(f), 1.0))
	gl.UniformMatrix3fv(c.texMatrix, 1, false, &c.tex[0])

	gl.BindVertexArray(c.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.texLoc)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

// Update TODO doc
func (c *Context) Update(dt float32) {
}

// Bounds TODO doc
func (c *Context) Bounds() shapes.Rect {
	return shapes.Rect{}
}

// Load TODO doc
func Load(path string, frames int) (*Context, error) {
	// TODO: prob should rename this to func to New
	c := Context{}
	// data, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not open file %s: %v", path, err)
	// }
	//
	// img, _, err := image.Decode(bytes.NewReader(data))
	// if err != nil {
	// 	return nil, fmt.Errorf("could not decode file %s: %v", path, err)
	// }
	imgFile, err := os.Open(path)
	if err != nil {
		return &c, fmt.Errorf("could not open file %s: %v", path, err)
	}
	c.Image, _, err = image.Decode(imgFile)
	if err != nil {
		return &c, fmt.Errorf("could not decode file %s: %v", path, err)
	}

	rgba := image.NewRGBA(c.Image.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return &c, fmt.Errorf("unsupported stride")
	}
	c.frames = frames
	c.Width = int(float32(rgba.Rect.Size().X) / float32(frames))
	c.Height = rgba.Rect.Size().Y

	draw.Draw(rgba, rgba.Bounds(), c.Image, image.Point{0, 0}, draw.Src)

	gl.GenTextures(1, &c.texLoc)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.texLoc)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return &c, nil
}

var vertices = []float32{
	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 0.0,
}
