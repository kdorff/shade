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
from os import path
from sdl2 import keycode

from transylvania import Application
from transylvania.display import DisplayManager
from transylvania.sprite import SpriteManager


class ReferenceApp(Application):
    """
    Example Game using the Transylvania Engine.
    """

    def handle_input(self, key):
        if key == keycode.SDLK_LEFT:
            self.sprite1.pos_x = self.sprite1.pos_x - 15
            return
        if key == keycode.SDLK_RIGHT:
            self.sprite1.pos_x = self.sprite1.pos_x + 15
            return

    def __init__(self, config=None, display=None, sprite_manager=None):
        super(ReferenceApp, self).__init__(config=config, display=display,
                                           sprite_manager=sprite_manager)
        # do some custom app stuff here.


def start_app():
    resource_dir = path.realpath(__file__)
    resource_dir = resource_dir.split('/')
    resource_dir.pop()
    resource_dir.append('resources')
    resource_dir = '/'.join(resource_dir)

    config = {}
    display = DisplayManager(width=800, height=600)
    sprite_manager = SpriteManager('{0}/sprites'.format(resource_dir))
    app = ReferenceApp(config=config, display=display,
                       sprite_manager=sprite_manager)
    app.run()


if __name__ == '__main__':
    start_app()
