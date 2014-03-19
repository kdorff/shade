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

import ctypes
from sdl2 import SDL_Event, SDL_PollEvent, SDL_Quit, events
from sdl2 import SDL_MOUSEBUTTONDOWN, SDL_MOUSEMOTION, SDL_QUIT
import time

# Version information (major, minor, revision[, 'dev']).
version_info = (0, 0, 1)

# Version string 'major.minor.revision'.
version = __version__ = ".".join(map(str, version_info))


class Application(object):
    """
    Application class for the Transylvania Engine.
    """

    def __init__(self, config=None, display=None, sprite_manager=None):
        """
        Initialize the application.

        @param config: Application configuration.
        @type config: dict
        @param display: Display Manager to be used.
        @type display: DisplayManager
        """
        self.config = config
        self.display = display
        self.sprite_manager = sprite_manager
        self.running = True
        self.objects = []
        self.lights = []

    def add_object(self, obj):
        self.objects.append(obj)

    def add_light(self, light):
        self.lights.append(light)

    def __del__(self):
        """
        Cleanup the application.
        """
        del self.display
        SDL_Quit()

    def handle_mousemotion(self, motion):
        pass

    def handle_mousebuttondown(self, button):
        pass

    def handle_keydown(self, key):
        pass

    def run(self):
        """
        Start the application event loop.
        """
        self.display.init_window()

        event = SDL_Event()

        prev_time = 0
        current_time = current_time = time.time()
        while self.running:
            prev_time = current_time
            current_time = time.time()
            timedelta = current_time - prev_time

            while SDL_PollEvent(ctypes.byref(event)) != 0:
                if event.type == SDL_QUIT:
                    return
                if event.type == events.SDL_KEYDOWN:
                    self.handle_keydown(event.key.keysym.sym)
                if (event.type == SDL_MOUSEMOTION):
                    self.handle_mousemotion(event.motion)
                if (event.type == SDL_MOUSEBUTTONDOWN):
                    self.handle_mousebuttondown(event.button.button)

            for light in self.lights:
                light.update(timedelta)

            for obj in self.objects:
                obj.update(timedelta)

            self.display.start_render()
            self.objects.sort(key=lambda obj: obj.layer)
            proj_matrix = self.display.get_proj_matrix()
            view_matrix = self.display.get_view_matrix(0, 0)

            for obj in self.objects:
                obj.draw(proj_matrix, view_matrix, self.lights)

            self.display.stop_render()
