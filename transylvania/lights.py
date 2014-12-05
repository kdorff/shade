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


class PointLight(object):
    """
    Point Light
    """
    def __init__(self, x=0.0, y=0.0, z=0.0, r=1.0, g=1.0, b=1.0, i=40.0):
        self.x = x
        self.y = y
        self.z = z
        self.r = r
        self.g = g
        self.b = b
        self.i = i

    def update(self, timedelta):
        """
        """
        pass

    def get_position(self):
        return [self.x, self.y, self.z]

    def get_color(self):
        return [self.r, self.g, self.b]

    def get_power(self):
        return self.i
