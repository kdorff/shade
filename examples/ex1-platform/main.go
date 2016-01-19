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
// Package app manages the main game loop.

package main

import (
	"log"
	"runtime"

	"github.com/hurricanerix/transylvania/display"
	"github.com/hurricanerix/transylvania/examples/ex1-platform/game"
	"github.com/hurricanerix/transylvania/splash"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	screen, err := display.SetMode("ex1-platform", 640, 480)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	g, err := game.New(screen)
	if err != nil {
		log.Fatalln("failed to create game:", err)
	}

	// Please see transylvania/splash/splash.go for details
	splash.Main(screen)

	g.Main(screen)
}
