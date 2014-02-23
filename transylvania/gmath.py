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
