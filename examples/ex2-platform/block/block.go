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
// Package block TODO doc

package block

import (
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Block struct {
	Pos    mgl32.Vec3
	Sprite *sprite.Context
	Shape  *shapes.Shape
}

// New TODO doc
func New(x, y float32, s *sprite.Context, group *[]entity.Entity) (*Block, error) {
	// TODO should take a group in as a argument
	b := Block{
		Pos:    mgl32.Vec3{x, y, 1.0},
		Sprite: s,
		Shape:  shapes.NewRect(0, float32(s.Width), 0, float32(s.Height)),
	}

	// TODO: this should probably be added outside of player
	if group != nil {
		*group = append(*group, &b)
	}
	return &b, nil
}

func (b Block) Type() string {
	return "block"
}

func (b Block) Label() string {
	return ""
}

func (b Block) Bounds() *shapes.Shape {
	return b.Shape
}

func (b Block) Pos2() *mgl32.Vec3 {
	return &b.Pos
}

// Bind TODO doc
func (b *Block) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

// Update TODO doc
func (b *Block) Update(dt float32, g []entity.Entity) {
	// TODO: Myabe handeling events should be done here, and not in a seperate "HandleEvents" func?
}

// Draw TODO doc
func (b Block) Draw() {
	//e *sprite.Effects) {
	//b.Sprite.Draw(b.Pos, e)
	b.Sprite.Draw(b.Pos, nil)
}