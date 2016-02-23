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
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/shapes"
)

// Entity ...
type Entity interface{}

// Updater ...
type Updater interface {
	Update(dt float32, group []Entity)
}

// Drawer ...
type Drawer interface {
	Pos() mgl32.Vec3
	Draw()
}

// Collider ...
type Collider interface {
	Pos() mgl32.Vec3
	Bounds() shapes.Shape
}

type Collision struct {
	Hit Collider
	Dir mgl32.Vec3
}

// Collide target with all enttities in group, returning all hits.
func Collide(target Collider, group []Collider) (hits []Collision) {
	if target == nil || group == nil {
		return hits
	}

	tBounds := target.Bounds()
	tPos := target.Pos()
	///tType := reflect.TypeOf(target).String()
	tType := tBounds.Type

	var hit bool
	var dir mgl32.Vec3
	var ep mgl32.Vec3
	var eb shapes.Shape
	for _, e := range group {
		hit = false
		ep = e.Pos()
		eb = e.Bounds()

		if tType == eb.Type && tType == "rect" {
			hit, dir = rectTest(tPos, ep, tBounds, eb, true)
			//} else if tType == reflect.TypeOf(eb).String() && tType == "circle" {
		} else if tType == eb.Type && tType == "circle" {
			hit, dir = circleTest(tPos, ep, tBounds, eb, true)
		} else {
			hit, dir = mixedTest(tPos, ep, tBounds, eb, true)
		}

		if hit {
			hits = append(hits, Collision{Hit: e, Dir: dir})
		}
	}
	return hits
}

func rectTest(ap, bp mgl32.Vec3, ab, bb shapes.Shape, ignoreZ bool) (bool, mgl32.Vec3) {
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
		return true, getDir(apt, bpt)
	}
	return false, mgl32.Vec3{}
}

func circleTest(ap, bp mgl32.Vec3, ab, bb shapes.Shape, ignoreZ bool) (bool, mgl32.Vec3) {
	if ap[0] == bp[0] && ap[1] == bp[1] && ab.Data[0] == bb.Data[0] {
		return false, mgl32.Vec3{}
	}

	apt := mgl32.Vec3{ap[0] + ab.Data[0], ap[1] + ab.Data[1]}
	bpt := mgl32.Vec3{bp[0] + bb.Data[0], bp[1] + bb.Data[1]}

	// distance between centers
	d := math.Sqrt(
		math.Pow(float64(apt[0])-float64(bpt[0]), 2) +
			math.Pow(float64(apt[1])-float64(bpt[1]), 2))
	// Sum of radiuses
	if d <= float64(ab.Data[2])+float64(bb.Data[2]) {
		return true, getDir(apt, bpt)
	}
	return false, mgl32.Vec3{}
}

func mixedTest(ap, bp mgl32.Vec3, ab, bb shapes.Shape, ignoreZ bool) (bool, mgl32.Vec3) {
	var rp mgl32.Vec3
	var rs shapes.Shape
	var cp mgl32.Vec3
	var cs shapes.Shape
	mod := float32(1.0)
	if ab.Type == "rect" {
		rp = ap
		rs = ab
		cp = mgl32.Vec3{bp[0] + bb.Data[0] - bb.Data[2], bp[1] + bb.Data[1] - bb.Data[2]}
		cs = bb
	} else {
		rp = bp
		rs = bb
		cp = mgl32.Vec3{ap[0] + ab.Data[0] - ab.Data[2], ap[1] + ab.Data[1] - ab.Data[2]}
		cs = ab
		mod = -1.0
	}
	w := float32(cs.Data[2]) * 2
	result, dir := rectTest(rp, cp, rs, *shapes.NewRect(0, w, 0, w), ignoreZ)
	dir[0] *= mod
	dir[1] *= mod
	dir[2] *= mod
	return result, dir
}

func getDir(aPos, bPos mgl32.Vec3) mgl32.Vec3 {
	bPos[0] -= aPos[0]
	bPos[1] -= aPos[1]
	bPos[2] -= aPos[2]
	m := float32(math.Sqrt(float64((bPos[0] * bPos[0]) + (bPos[1] * bPos[1]))))
	if m == 0 {
		return mgl32.Vec3{}
	}
	uv := mgl32.Vec3{bPos[0] / m, bPos[1] / m, bPos[2] / m}
	return uv
}
