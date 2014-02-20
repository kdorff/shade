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
from sdl2 import SDL_Event, SDL_PollEvent, SDL_Quit
from sdl2 import SDL_QUIT

# Version information (major, minor, revision[, 'dev']).
version_info = (0, 0, 1)

# Version string 'major.minor.revision'.
version = __version__ = ".".join(map(str, version_info))


class Application(object):
    """
    Application class for the Transylvania Engine.
    """

    def __init__(self, config=None, display=None):
        """
        Initialize the application.

        @param config: Application configuration.
        @type config: dict
        @param display: Display Manager to be used.
        @type display: DisplayManager
        """
        self.config = config
        self.display = display
        self.running = True
        self.objects = []

    def add_object(self, obj):
        self.objects.append(obj)

    def __del__(self):
        """
        Cleanup the application.
        """
        del self.display
        SDL_Quit()

    def run(self):
        """
        Start the application event loop.
        """
        event = SDL_Event()

        while self.running:
            while SDL_PollEvent(ctypes.byref(event)) != 0:
                if event.type == SDL_QUIT:
                    return

            self.display.start_render()
            self.objects.sort(key=lambda obj: obj.layer)
            proj_mat = self.display.get_proj_matrix()

            for obj in self.objects:
                obj.draw(proj_mat)

            self.display.stop_render()
