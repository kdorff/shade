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

from transylvania import Application
from transylvania.display import DisplayManager
from transylvania.sprite import SpriteBuilder


class ReferenceApp(Application):
    """
    Example Game using the Transylvania Engine.
    """


def start_app():
    config = {}
    display = DisplayManager(width=800, height=600)

    resource_dir = path.realpath(__file__)
    resource_dir = resource_dir.split('/')
    resource_dir.pop()
    resource_dir.append('resources')
    resource_dir.append('')
    resource_dir = '/'.join(resource_dir)

    app = ReferenceApp(config=config, display=display)

    sprite_path = '{0}/sprites/bimon_selmont'.format(resource_dir)

    sprite1 = SpriteBuilder.build(sprite_path, pos_x=0, pos_y=0, layer=-1)
    sprite2 = SpriteBuilder.build(sprite_path, pos_x=25, pos_y=25, layer=0)
    sprite3 = SpriteBuilder.build(sprite_path, pos_x=50, pos_y=50, layer=1)
    sprite4 = SpriteBuilder.build(sprite_path, pos_x=250, pos_y=250, layer=2)

    sprite2.set_animation('walking')
    sprite3.set_animation('die') # black
    sprite4.set_animation('attack') # balck

    app.add_object(sprite1) # black
    app.add_object(sprite2) # blue
    app.add_object(sprite3) # green
    app.add_object(sprite4) # yellow

    app.run()


if __name__ == '__main__':
    start_app()
