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

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/entity"
	"github.com/hurricanerix/shade/shapes"
)

// Collide TODO doc
func Collide(t entity.Entity, g *[]entity.Entity, dokill bool) *[]entity.Entity {
	// TODO: maybe move this to entity package?
	var hits []entity.Entity

	if t == nil || g == nil {
		return &hits
	}

	tb := t.Bounds()
	if tb == nil {
		return &hits
	}
	tp := t.Pos2()

	var hit bool
	var ep *mgl32.Vec3
	var eb *shapes.Shape
	for _, e := range *g {
		hit = false
		ep = e.Pos2()
		eb = e.Bounds()

		if tb.Type == eb.Type && tb.Type == "rect" {
			hit = rectTest(tp, ep, tb, eb, true)
		} else if tb.Type == eb.Type && tb.Type == "circle" {
			hit = circleTest(tp, ep, tb, eb, true)
		} else {
			hit = mixedTest(tp, ep, tb, eb, true)
		}

		if hit {
			hits = append(hits, e)
		}
	}
	return &hits
}

func rectTest(ap, bp *mgl32.Vec3, ab, bb *shapes.Shape, ignoreZ bool) bool {

	apt := mgl32.Vec3{
		ap[0] + ab.Data[0],
		ap[1] + ab.Data[2],
	}
	aw := float32(math.Abs(float64(ab.Data[0]))) + float32(math.Abs(float64(ab.Data[1])))
	ah := float32(math.Abs(float64(ab.Data[2]))) + float32(math.Abs(float64(ab.Data[3])))

	bpt := mgl32.Vec3{
		bp[0] + bb.Data[0],
		bp[1] + bb.Data[2],
	}
	bw := float32(math.Abs(float64(bb.Data[0]))) + float32(math.Abs(float64(bb.Data[1])))
	bh := float32(math.Abs(float64(bb.Data[2]))) + float32(math.Abs(float64(bb.Data[3])))

	if apt[0] < bpt[0]+bw &&
		apt[0]+aw > bpt[0] &&
		apt[1] < bpt[1]+bh &&
		ah+apt[1] > bpt[1] {
		return true
	}
	return false
}

func circleTest(ap, bp *mgl32.Vec3, ab, bb *shapes.Shape, ignoreZ bool) bool {
	if ap[0] == bp[0] && ap[1] == bp[1] && ab.Data[0] == bb.Data[0] {
		return false
	}
	// distance between centers
	d := math.Sqrt(
		math.Pow(float64(ap[0])-float64(bp[0]), 2) +
			math.Pow(float64(ap[1])-float64(bp[1]), 2))
	// Sum of radiuses
	if d <= float64(ab.Data[2])+float64(bb.Data[2]) {
		return true
	}
	return false
}

func mixedTest(ap, bp *mgl32.Vec3, ab, bb *shapes.Shape, ignoreZ bool) bool {
	var rp *mgl32.Vec3
	var rs *shapes.Shape
	var cp *mgl32.Vec3
	var cs *shapes.Shape

	if ab.Type == "rect" {
		rp = ap
		rs = ab
		cp = bp
		cs = bb
	} else {
		rp = bp
		rs = bb
		cp = ap
		cs = ab
	}
	w := float32(cs.Data[2])
	return rectTest(rp, cp, rs, shapes.NewRect(w*-1, w*-1, w, w), ignoreZ)
}
