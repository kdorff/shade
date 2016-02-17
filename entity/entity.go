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
// Package entity provies interfaces for game objects.

package entity

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/shapes"
)

// Entity is the interface for a basic game object.
type Entity interface {
	// Type returns the name of the struct implenting the interface.  This can be used to cast a pointer to the interface to a pointer of that struct.
	Type() string
	// Label returns an identifier useful to the program.
	Label() string
	// Pos2 returns the position of the entity. (TODO: rename to Pos after refactor)
	Pos2() *mgl32.Vec3
	// Boundry of the object for the purpose of collision detection, if nil, the entity is not intended to be considered when detecting collisions.
	Bounds() *shapes.Shape
	// Update the state of the entity.
	Update(dt float32, g []Entity)
	// Draw the entity.
	Draw()
}
