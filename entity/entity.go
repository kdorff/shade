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
// Package entity TODO doc

package entity

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/shapes"
)

// Entity TODO doc
type Entity interface {
	Type() string
	Label() string
	Pos2() *mgl32.Vec3 // TODO: change this to Pos (will take a huge refactor)
	Bounds() *shapes.Shape
	Update(dt float32, g *[]Entity)
	Draw()
}
