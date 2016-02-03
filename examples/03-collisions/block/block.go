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

	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Block struct {
	Sprite *sprite.Context
	Rect   *shapes.Rect
}

// New TODO doc
func New(x, y float32, s *sprite.Context, group *sprite.Group) (*Block, error) {
	// TODO should take a group in as a argument
	b := Block{
		Sprite: s,
	}

	rect, err := shapes.NewRect(x, y, float32(b.Sprite.Width), float32(b.Sprite.Height))

	if err != nil {
		return &b, err
	}
	b.Rect = rect

	// TODO: this should probably be added outside of player
	group.Add(&b)
	return &b, nil
}

// Bind TODO doc
func (b *Block) Bind(program uint32) error {
	return b.Sprite.Bind(program)
}

// Update TODO doc
func (b *Block) Update(dt float32, g *sprite.Group) {
	// Blocks don't do anything
}

// Draw TODO doc
func (b *Block) Draw() {
	b.Sprite.Draw(b.Rect.X, b.Rect.Y)
}

// Bounds TODO doc
func (b *Block) Bounds() chan shapes.Rect {
	ch := make(chan shapes.Rect, 1)
	ch <- *(b.Rect)
	close(ch)
	return ch
}
