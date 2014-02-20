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
from transylvania.sprite import Sprite


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

    animations = {'default': [(1, 0)],
                  'walking': [(0, 0), (1, 0), (2, 0)],
                  'die': [(1, 0), (1, 2), (2, 2)]}

    sprite1 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     pos_x=25, pos_y=25, layer=1, animations=animations)

    sprite2 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     animations=animations)
    sprite2.set_animation('walking')

    sprite3 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     pos_x=50, pos_y=50, layer=2, animations=animations)
    sprite3.set_animation('die')

    app.add_object(sprite1)
    app.add_object(sprite2)
    app.add_object(sprite3)




    app.run()


if __name__ == '__main__':
    start_app()
