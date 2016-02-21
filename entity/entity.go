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

// Entity is the interface for a basic game object.
type Entity interface {
	// Type returns the name of the struct implenting the interface.  This can be used to cast a pointer to the interface to a pointer of that struct.
	Type() string
	// Label returns an identifier useful to the program.
	Label() string
	Update(dt float32, g []Collider)
	Draw()
}
