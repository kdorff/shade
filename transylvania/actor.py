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


class Actor(object):

    def __init__(self, sprite, x=0, y=0, width=None, height=None, layer=0):
        self.sprite = sprite
        self.x = x
        self.y = y
        self.width = width
        self.height = height
        self.layer = layer
        self.current_animation = 'default'
        self.animation_speed = 0.4
        self.current_frame = 0

    def _get_frames(self):
        return self.sprite.data['animations'][self.current_animation]['frames']

    def _get_frame(self):
        frames = self._get_frames()
        return frames[int(self.current_frame)]

    def set_animation(self, name):
        self.current_animation = name
        self.current_frame = 0

    def update(self, timedelta):
        frames = self._get_frames()
        self.current_frame = self.current_frame + 1
        if self.current_frame > len(frames) - 1:
            self.current_frame = 0

    def draw(self, proj_matrix, view_matrix, light_position=None):
        frame_x, frame_y = self._get_frame()
        self.sprite.draw(proj_matrix, view_matrix,
                         self.x, self.y, self.layer,
                         frame_x=frame_x, frame_y=frame_y,
                         light_position=light_position)
