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

#from math import cos, sin, tan, sqrt, pi
from sdl2 import (
    SDL_CreateWindow, SDL_DestroyWindow, SDL_GL_CreateContext,
    SDL_GL_SetAttribute, SDL_GL_SetSwapInterval, SDL_GL_SwapWindow, SDL_Init)
from sdl2 import (
    SDL_GL_CONTEXT_MAJOR_VERSION, SDL_GL_CONTEXT_MINOR_VERSION,
    SDL_GL_CONTEXT_PROFILE_CORE, SDL_GL_CONTEXT_PROFILE_MASK,
    SDL_GL_DEPTH_SIZE, SDL_GL_DOUBLEBUFFER, SDL_INIT_EVERYTHING,
    SDL_WINDOW_OPENGL, SDL_WINDOW_SHOWN, SDL_WINDOWPOS_CENTERED)

from OpenGL.GL import glBlendFunc, glClear, glClearColor, glEnable, glViewport
from OpenGL.GL import (
    GL_BLEND, GL_COLOR_BUFFER_BIT, GL_CULL_FACE, GL_DEPTH_BUFFER_BIT,
    GL_DEPTH_TEST, GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)

from gameobjects.matrix44 import Matrix44

from transylvania.gmath import get_4x4_transform


class DisplayManager(object):
    """
    Manages screen related things.
    """

    def __init__(self, width=0, height=0, left=0, bottom=0):
        """
        Initialize the display manager.

        @param width: Width in pixels of the display.
        @type width: int
        @param height: Height in pixels of the display.
        @type height: int
        """
        self.width = width
        self.height = height
        self.proj_matrix = None
        self.view_matrix = None
        self.window = None
        self.glcontext = None

    def __del__(self):
        """
        Cleanup the display manager.
        """
        SDL_DestroyWindow(self.window)

    def _get_view_matrix(self, x, y):
        scale_x = 1.0
        scale_y = 1.0
        trans_x = x
        trans_y = y
        layer = 1.0
        return get_4x4_transform(scale_x, scale_y, trans_x, trans_y, layer)

    def _get_projection_matrix(self, left, right, bottom, top):
        """
        Create a  orthographic projection matrix.

        U{Modern glOrtho2d<http://stackoverflow.com/questions/21323743/
        modern-equivalent-of-gluortho2d>}

        U{Orthographic Projection<http://en.wikipedia.org/wiki/
        Orthographic_projection_(geometry)>}

        @param left: position of the left side of the display
        @type left: int
        @param right: position of the right side of the display
        @type right: int
        @param bottom: position of the bottom side of the display
        @type bottom: int
        @param top: position of the top side of the display
        @type top: int

        @return: orthographic projection matrix
        @rtype: Matrix44
        """
        zNear = -25.0
        zFar = 25.0
        inv_z = 1.0 / (zFar - zNear)
        inv_y = 1.0 / (top - bottom)
        inv_x = 1.0 / (right - left)

        mat = Matrix44()
        mat.set_row(0, [(2.0 * inv_x), 0.0, 0.0, (-(right + left) * inv_x)])
        mat.set_row(1, [0.0, (2.0 * inv_y), 0.0, (-(top + bottom) * inv_y)])
        mat.set_row(2, [0.0, 0.0, (-2.0 * inv_z), (-(zFar + zNear) * inv_z)])
        mat.set_row(3, [0.0, 0.0, 0.0, 1.0])

        return mat.to_opengl()

    def init_window(self):
        """
        Handles creating the SDL window and creating a GL context for it.
        """
        SDL_Init(SDL_INIT_EVERYTHING)

        SDL_GL_SetAttribute(SDL_GL_CONTEXT_MAJOR_VERSION, 3)
        SDL_GL_SetAttribute(SDL_GL_CONTEXT_MINOR_VERSION, 2)
        SDL_GL_SetAttribute(
            SDL_GL_CONTEXT_PROFILE_MASK, SDL_GL_CONTEXT_PROFILE_CORE)

        SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER, 1)
        SDL_GL_SetAttribute(SDL_GL_DEPTH_SIZE, 24)

        self.window = SDL_CreateWindow(
            "Transylvania Engine", SDL_WINDOWPOS_CENTERED,
            SDL_WINDOWPOS_CENTERED, self.width, self.height,
            SDL_WINDOW_OPENGL | SDL_WINDOW_SHOWN)
        if not self.window:
            raise Exception('Could not create window')

        self.proj_matrix = self._get_projection_matrix(0.0, self.width,
                                                       0.0, self.height)

        self.glcontext = SDL_GL_CreateContext(self.window)
        SDL_GL_SetSwapInterval(1)
        glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)
        glEnable(GL_DEPTH_TEST)
        glEnable(GL_CULL_FACE)
        glEnable(GL_BLEND)
        glClearColor(0.3, 0.3, 0.3, 1.0)

    def set_clear_color(self, r=0.0, g=0.0, b=0.0):
        glClearColor(r, g, b, 1.0)

    def resize(self, width, height):
        """
        Resize the display.

        @param width: Width in pixels to change the display to.
        @type width: int
        @param height: Height in pixels to change the display to.
        @type height: int
        """
        self.width = width
        self.height = height

        SDL_DestroyWindow(self.window)
        self.init_window()

        glViewport(0, 0, self.width, self.height)

    def get_proj_matrix(self):
        return self.proj_matrix

    def get_view_matrix(self, x, y):
        self.view_matrix = self._get_view_matrix(x, y)
        return self.view_matrix

    def start_render(self):
        """
        Setup for scene render.
        """
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)

    def stop_render(self):
        """
        Cleanup from scene render.
        """
        SDL_GL_SwapWindow(self.window)
