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

    animations = {'default': [(0, 0)],
                  'walking': [(1, 0), (2, 0)],
                  'duck': [(1, 2), None],
                  'climb_down': [(0, 1), (1, 0)],
                  'climb_up': [(0, 1), (1, 1)],
                  'falling_back': [(0, 2)],
                  'die': [(1, 0), (1, 2), (2, 2), None],
                  'whip': [(0, 3), (1, 3), (2, 3)],
                  'whip_climb_up': [(0, 4), (1, 4), (2, 4)],
                  'whip_climb_down': [(0, 5), (1, 5), (2, 5)],
                  'whip_duck': [(0, 6), (1, 6), (2, 6)],
                  'look_up': [(3, 9), None]}

    sprite1 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     pos_x=25, pos_y=25, layer=1, animations=animations)

    sprite2 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     animations=animations)
    sprite2.set_animation('walking')

    sprite3 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     pos_x=50, pos_y=50, layer=2, animations=animations)
    sprite3.set_animation('die')

    sprite4 = Sprite('{0}/sprites/bimon_selmont'.format(resource_dir),
                     pos_x=250, pos_y=250, layer=2, animations=animations)
    sprite4.set_animation('climb_up')

    app.add_object(sprite1)
    app.add_object(sprite2)
    app.add_object(sprite3)
    app.add_object(sprite4)

    app.run()


if __name__ == '__main__':
    start_app()
