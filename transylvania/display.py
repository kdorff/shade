# The MIT License (MIT)
#
# Copyright (c) 2014 Richard Hawkins
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

from sdl2 import (
    SDL_CreateWindow, SDL_DestroyWindow, SDL_GL_CreateContext,
    SDL_GL_SetAttribute, SDL_GL_SetSwapInterval, SDL_GL_SwapWindow, SDL_Init)
from sdl2 import (
    SDL_GL_CONTEXT_MAJOR_VERSION, SDL_GL_CONTEXT_MINOR_VERSION,
    SDL_GL_CONTEXT_PROFILE_CORE, SDL_GL_CONTEXT_PROFILE_MASK,
    SDL_GL_DEPTH_SIZE, SDL_GL_DOUBLEBUFFER, SDL_INIT_EVERYTHING,
    SDL_WINDOW_OPENGL, SDL_WINDOW_SHOWN, SDL_WINDOWPOS_CENTERED)

from OpenGL.GL import glClear, glClearColor, glEnable, glViewport
from OpenGL.GL import (
    GL_COLOR_BUFFER_BIT, GL_CULL_FACE, GL_DEPTH_BUFFER_BIT, GL_DEPTH_TEST)

from transylvania.sprite import Sprite


class DisplayManager(object):
    """
    Manages screen related things.
    """

    def __init__(self, width=0, height=0):
        """
        Initialize the display manager.
        """
        self.width = width
        self.height = height
        self.window = None
        self.glcontext = None

        SDL_Init(SDL_INIT_EVERYTHING)

        SDL_GL_SetAttribute(SDL_GL_CONTEXT_MAJOR_VERSION, 3)
        SDL_GL_SetAttribute(SDL_GL_CONTEXT_MINOR_VERSION, 2)
        SDL_GL_SetAttribute(
            SDL_GL_CONTEXT_PROFILE_MASK, SDL_GL_CONTEXT_PROFILE_CORE)

        SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER, 1)
        SDL_GL_SetAttribute(SDL_GL_DEPTH_SIZE, 24)

        self._create_window()
        self.sprite = Sprite()

    def __del__(self):
        """
        Cleanup the display manager.
        """
        SDL_DestroyWindow(self.window)

    def _create_window(self):
        """
        Handles creating the SDL window and creating a GL context for it.
        """
        self.window = SDL_CreateWindow(
            "Transylvania Engine", SDL_WINDOWPOS_CENTERED,
            SDL_WINDOWPOS_CENTERED, self.width, self.height,
            SDL_WINDOW_OPENGL | SDL_WINDOW_SHOWN)
        if not self.window:
            raise Exception('Could not create window')

        self.glcontext = SDL_GL_CreateContext(self.window)
        SDL_GL_SetSwapInterval(1)
        glEnable(GL_DEPTH_TEST)
        glEnable(GL_CULL_FACE)
        glClearColor(0.2, 0.2, 0.2, 1.0)

    def resize(self, width, height):
        """
        Resize the display.

        @param width: pixel count horizontally
        @type width: int
        @param height: pixel count vertically
        @type height: int
        """
        self.width = width
        self.height = height

        SDL_DestroyWindow(self.window)
        self._create_window()

        glViewport(0, 0, self.width, self.height)

    def render(self):
        """
        Renders the scene.
        """
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)

        self.sprite.draw()

        SDL_GL_SwapWindow(self.window)
