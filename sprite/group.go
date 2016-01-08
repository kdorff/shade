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

package sprite

import "github.com/hurricanerix/transylvania/interfaces"

// Group TODO doc
type Group struct {
	sprites []interfaces.Sprite
}

// NewGroup TODO doc
func NewGroup() *Group {
	g := Group{}
	return &g
}

// Add TODO doc
func (g *Group) Add(s interfaces.Sprite) {
	g.sprites = append(g.sprites, s)
}

// Update TODO doc
func (g *Group) Update(dt float32) {
	for _, s := range g.sprites {
		s.Update(dt)
	}
}

func (g *Group) Bind(program uint32) error {
	for _, s := range g.sprites {
		s.Bind(program)
	}
	return nil
}

// Draw TODO doc
func (g *Group) Draw() {
	for _, s := range g.sprites {
		s.Draw()
	}
}
