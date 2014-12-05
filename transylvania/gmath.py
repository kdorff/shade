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

from gameobjects.matrix44 import Matrix44


def get_4x4_transform(scale_x, scale_y, trans_x, trans_y, layer):
    """
    Transform the from local coordinates to world coordinates.

    @return: transformation matrix used to transform from local coords
             to world coords.
    @rtype: Matrix44
    """
    transform = Matrix44()
    transform.set_row(0, [scale_x, 0.0, 0.0, trans_x])
    transform.set_row(1, [0.0, scale_y, 0.0, trans_y])
    transform.set_row(2, [0.0, 0.0, 1.0, layer])
    transform.set_row(3, [0.0, 0.0, 0.0, 1.0])

    return transform.to_opengl()


def get_3x3_transform(scale_x, scale_y, trans_x, trans_y):
    """
    Transform the tex coords to shrink the texutre down to a single frame.

    @return: matrix used to transform texture coords to the current
             frame in the texture to be displayed.
    @rtype: list
    """
    transform = [scale_x, 0.0, trans_x,
                 0.0, scale_y, trans_y,
                 0.0, 0.0, 1.0]
    return transform
