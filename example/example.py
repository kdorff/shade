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

    app = ReferenceApp(config=config, display=display)
    # TODO(hurricanerix): This should not be a hard coded path which
    # probably only exists on my computer.
    app.add_object(Sprite('/Users/rhawkins/workspace/transylvania/example/'
                          'resources/sprites/bimon_selmont'))

    app.add_object(Sprite('/Users/rhawkins/workspace/transylvania/example/'
                          'resources/sprites/bimon_selmont', pos_x=100, pos_y=100, layer=-1))

    app.run()


if __name__ == '__main__':
    start_app()
