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
// Package display TODO doc

package display

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/events"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Context TODO doc
type Context struct {
	Window     *glfw.Window
	Width      float32
	Height     float32
	Program    uint32
	ProjMatrix mgl32.Mat4
	ViewMatrix mgl32.Mat4
}

// SetMode TODO doc
func SetMode(title string, width, height int) (*Context, error) {
	c := Context{
		Width:  float32(width),
		Height: float32(height),
	}
	if err := glfw.Init(); err != nil {
		return &c, fmt.Errorf("failed to initialize glfw: %v", err)
	}
	// TODO: move this to a terminate function
	//defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(640, 480, title, nil, nil)
	if err != nil {
		return &c, fmt.Errorf("failed to create window: %v", err)
	}
	c.Window = window

	c.Window.MakeContextCurrent()
	c.Window.SetKeyCallback(events.KeyCallback)
	c.Window.SetCursorPosCallback(events.CursorPositionCallback)

	if err := gl.Init(); err != nil {
		return &c, fmt.Errorf("failed to init glow: %v", err)
	}

	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	c.Program, err = newProgram(vertexShader, fragmentShader)
	if err != nil {
		return &c, fmt.Errorf("error loading program: %v", err)
	}

	gl.UseProgram(c.Program)

	var left, right, top, bottom, near, far float32
	right = float32(width)
	top = float32(height)
	near = 0.1
	far = 100.0
	c.ProjMatrix = mgl32.Ortho(left, right, bottom, top, near, far)
	projUniform := gl.GetUniformLocation(c.Program, gl.Str("ProjMatrix\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &c.ProjMatrix[0])

	var eye, center, up mgl32.Vec3
	eye = mgl32.Vec3{0.0, 0.0, 7.0}
	center = mgl32.Vec3{0.0, 0.0, -1.0}
	up = mgl32.Vec3{0.0, 1.0, 0.0}
	c.ViewMatrix = mgl32.LookAtV(eye, center, up)
	viewUniform := gl.GetUniformLocation(c.Program, gl.Str("ViewMatrix\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &c.ViewMatrix[0])

	gl.BindFragDataLocation(c.Program, 0, gl.Str("outputColor\x00"))

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	// TODO: Figure out why "layering" using z-buffer does not work.
	// gl.DepthFunc(gl.LESS)
	// gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.NEVER)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//window.SwapBuffers()
	return &c, nil
}

// Fill TODO doc
func (c *Context) Fill(r, g, b float32) {
	gl.ClearColor(r, g, b, 1.0)
}

// Flip TODO doc
func (c *Context) Flip() {
	c.Window.SwapBuffers()
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("can not create vert shader: %s", err)
	}

	fragShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("can not create frag shader: %s", err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 330

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;
uniform mat3 TexMatrix;

in vec3 MCVertex;
in vec2 TexCoord0;

out vec2 TexCoord;

void main() {
  TexCoord = vec3(TexMatrix * vec3(TexCoord0, 1.0)).st;
  gl_Position = ProjMatrix * ViewMatrix * ModelMatrix * vec4(MCVertex, 1);
}
` + "\x00"

var fragmentShader = `
#version 330

uniform int AddColor;
uniform vec4 AColor;
uniform int SubColor;
uniform vec4 SColor;
uniform sampler2D ColorMap;
uniform vec4 AmbientColor;

in vec2 TexCoord;

out vec4 outputColor;

void main() {
	vec4 diffuse = texture(ColorMap, TexCoord);
    vec4 ambient = AmbientColor * diffuse;

	if (AddColor == 1) {
		diffuse = clamp(diffuse + AColor, 0.0, 1.0);
	}
	if (SubColor == 1) {
		diffuse = clamp(diffuse - SColor, 1.0, 1.0);
	}

	outputColor = ambient; 
}
` + "\x00"
