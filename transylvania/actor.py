# Copyright 2014 Richard Hawkins
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


class Actor(object):

    def __init__(self, sprite, x=0, y=0, width=None, height=None, layer=0):
        self.sprite = sprite
        self.x = x
        self.y = y
        self.width = width
        self.height = height
        self.layer = layer
        self.current_animation = 'default'
        self.current_time = 0.0
        self.animation_speed = 0.15
        self.current_frame = 0
        self.running = True

    def _get_frames(self):
        return self.sprite.data['animations'][self.current_animation]['frames']

    def _get_frame(self):
        frames = self._get_frames()
        return frames[int(self.current_frame)]

    def get_animations(self):
        return [x for x in self.sprite.data['animations'].iterkeys()]

    def set_animation(self, name):
        self.current_animation = name
        self.current_frame = 0
        self.current_time = 0.0
        self.running = True

    def update(self, timedelta):
        if not self.running:
            return
        self.current_time = self.current_time + timedelta
        frames = self._get_frames()
        self.current_frame = (
            int(self.current_time / self.animation_speed) % len(frames))
        if self.current_frame > len(frames) - 1:
            self.current_frame = 0
        if frames[self.current_frame] is None:
            self.running = False
            self.current_frame = len(frames) - 2

    def draw(self, proj_matrix, view_matrix, lights=None):
        frame_x, frame_y = self._get_frame()
        self.sprite.draw(proj_matrix, view_matrix,
                         self.x, self.y, self.layer,
                         frame_x=frame_x, frame_y=frame_y,
                         lights=lights)
