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
	"fmt"
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
	Image *sprite.Context
	Rect  *shapes.Rect
}

// New TODO doc
func New(group *sprite.Group) (*Block, error) {
	// TODO should take a group in as a argument
	b := Block{}

	block, err := sprite.Load("block.png", 1, 1)
	if err != nil {
		return &b, fmt.Errorf("could not load block: %v", err)
	}
	b.Image = block

	rect, err := shapes.NewRect(0.0, 0.0, float32(b.Image.Width), float32(b.Image.Height))
	if err != nil {
		return &b, fmt.Errorf("could create rect: %v", err)
	}
	b.Rect = rect

	// TODO: this should probably be added outside of player
	group.Add(&b)
	return &b, nil
}

// Bind TODO doc
func (b *Block) Bind(program uint32) error {
	return b.Image.Bind(program)
}

// Update TODO doc
func (b *Block) Update(dt float32, g *sprite.Group) {
	// TODO: Myabe handeling events should be done here, and not in a seperate "HandleEvents" func?
}

// Draw TODO doc
func (b *Block) Draw() {
	b.Image.Draw(b.Rect.X, b.Rect.Y)
}

// Bounds TODO doc
func (b *Block) Bounds() chan shapes.Rect {
	ch := make(chan shapes.Rect, 1)
	ch <- *(b.Rect)
	close(ch)
	return ch
}
