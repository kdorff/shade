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

import "github.com/hurricanerix/transylvania/shapes"

// Collide TODO doc
func Collide(t Sprite, g *Group, dokill bool) []Sprite {
	var hits []Sprite
	if g == nil {
		return nil
	}
	tb := t.Bounds()
	for _, s := range g.sprites {
		if testBounds(tb, s.Bounds()) {
			hits = append(hits, s)
		}
	}
	return hits
}

func testBounds(a, b shapes.Rect) bool {
	return testHorz(a, b) && testVert(a, b)
}

func testHorz(a, b shapes.Rect) bool {
	arbl := a.Right() >= b.Left() && a.Right() < b.Right()
	albr := a.Left() <= b.Right() && a.Left() > b.Left()
	return arbl || albr
}

func testVert(a, b shapes.Rect) bool {
	atbb := a.Top() >= b.Bottom() && a.Top() < b.Top()
	abbt := a.Bottom() <= b.Top() && a.Bottom() > b.Bottom()
	return atbb || abbt
}
