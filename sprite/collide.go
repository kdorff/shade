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
	al := a.X
	ar := a.X + a.Width
	bl := b.X
	br := b.X + b.Width

	albl := al >= bl && al < br
	arbl := ar >= bl && ar < br

	albr := al <= br && al > bl
	arbr := ar <= br && ar > bl

	return albl || arbl || albr || arbr
}

func testVert(a, b shapes.Rect) bool {
	ab := a.Y
	at := a.Y + a.Height
	bb := b.Y
	bt := b.Y + b.Height

	abbb := ab >= bb && ab < bt
	atbb := at >= bb && at < bt

	abbt := ab <= bt && ab > bb
	atbt := at <= bt && at > bb

	return abbb || atbb || abbt || atbt
}
