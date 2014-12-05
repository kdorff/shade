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


def get_4x4_transform(scale_x=1.0, scale_y=1.0, trans_x=1.0, trans_y=1.0,
                      trans_z=1.0):
    """Returns a 4x4 transform matrix.

    @return: transformation matrix
    @rtype: list
    """
    transform = [
        scale_x, 0.0, 0.0, trans_x,
        0.0, scale_y, 0.0, trans_y,
        0.0, 0.0, 1.0, trans_z,
        0.0, 0.0, 0.0, 1.0]

    return transform


def get_3x3_transform(scale_x=1.0, scale_y=1.0, trans_x=1.0, trans_y=1.0):
    """Returns a 3x3 transform.

    @return: transformation matrix
    @rtype: list
    """
    transform = [scale_x, 0.0, trans_x,
                 0.0, scale_y, trans_y,
                 0.0, 0.0, 1.0]
    return transform


def get_projection_matrix(left, right, bottom, top):
    """Create a  orthographic projection matrix.

    U{Modern glOrtho2d<http://stackoverflow.com/questions/21323743/
    modern-equivalent-of-gluortho2d>}

    U{Orthographic Projection<http://en.wikipedia.org/wiki/
    Orthographic_projection_(geometry)>}

    @param left: position of the left side of the display
    @type left: int
    @param right: position of the right side of the display
    @type right: int
    @param bottom: position of the bottom side of the display
    @type bottom: int
    @param top: position of the top side of the display
    @type top: int

    @return: orthographic projection matrix
    @rtype: list
    """
    zNear = -25.0
    zFar = 25.0
    inv_z = 1.0 / (zFar - zNear)
    inv_y = 1.0 / (top - bottom)
    inv_x = 1.0 / (right - left)

    mat = [
        (2.0 * inv_x), 0.0, 0.0, (-(right + left) * inv_x),
        0.0, (2.0 * inv_y), 0.0, (-(top + bottom) * inv_y),
        0.0, 0.0, (-2.0 * inv_z), (-(zFar + zNear) * inv_z),
        0.0, 0.0, 0.0, 1.0]

    return mat
