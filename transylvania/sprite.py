## The MIT License (MIT)
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
import ctypes
import math
import numpy

import OpenGL.GL.shaders
from OpenGL.GL import (
    glActiveTexture, glBindBuffer, glBindVertexArray, glBindTexture,
    glBlendFunc, glBufferData, glDrawArrays, glEnable,
    glEnableVertexAttribArray, glGenBuffers, glGenTextures, glGenVertexArrays,
    glGetAttribLocation, glGetUniformLocation, glTexImage2D, glTexParameteri,
    glUniform1i, glUniformMatrix4fv, glUseProgram, glVertexAttribPointer)
from OpenGL.GL import (
    GL_ARRAY_BUFFER, GL_BLEND, GL_CULL_FACE, GL_FALSE, GL_FLOAT,
    GL_FRAGMENT_SHADER, GL_ONE_MINUS_SRC_ALPHA, GL_LINEAR, GL_REPEAT, GL_RGBA,
    GL_SRC_ALPHA, GL_TEXTURE0, GL_TEXTURE_2D, GL_TEXTURE_BASE_LEVEL,
    GL_TEXTURE_MAG_FILTER, GL_TEXTURE_MAX_LEVEL, GL_TEXTURE_MIN_FILTER,
    GL_TEXTURE_WRAP_S, GL_TEXTURE_WRAP_T, GL_TRIANGLES, GL_STATIC_DRAW,
    GL_UNSIGNED_BYTE, GL_VERTEX_SHADER)
from PIL import Image

from gameobjects.matrix44 import Matrix44

current_frame = 0


vertex_shader = """
#version 330

uniform mat4 proj_mat;
uniform mat4 offset;

uniform int x;
uniform int y;

in vec2 TexCoord0;
in vec4 position;
smooth out vec2 TexCoord;

void main()
{
   float offset_x = 1.0 / 3.0;
   float offset_y = 1.0 / 10.0;

   TexCoord = (TexCoord0.st / vec2(3, 10)) +
        vec2(offset_x * (x), offset_y * (y));
   gl_Position = position * offset * proj_mat;
}
"""

fragment_shader = """
#version 330

out vec4 MyFragColor;
in vec2 TexCoord;

uniform sampler2D ColorMap;

void main()
{
   vec4 color = vec4(texture(ColorMap, TexCoord.st).rgba);
   MyFragColor = vec4(color);
}
"""

vertices = [0.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 0.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0]

tex_coords = [0.0, 1.0,
              1.0, 1.0,
              1.0, 0.0,
              0.0, 0.0,
              0.0, 1.0,
              1.0, 0.0]

vertices = numpy.array(vertices, dtype=numpy.float32)
tex_coords = numpy.array(tex_coords, dtype=numpy.float32)


class Sprite(object):
    """
    Handles texturing images on a polygon.
    """

    def __init__(self, path, pos_x=0, pos_y=0, layer=0, animations=None):
        """
        Initialize the OpenGL things needed to render the polygon.
        """
        # TODO(hurricanerix): position stuff should probably be moved outside
        # of the sprite class.
        self.path = path
        self.pos_x = pos_x
        self.pos_y = pos_y
        self.layer = layer
        self.width = 308
        self.height = 132
        self.current_animation = 'default'

        if not animations:
            animations = {'default': [(0, 0)]}
        self.animations = animations

        glEnable(GL_CULL_FACE)
        glEnable(GL_BLEND)
        glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)

        # Create a new VAO (Vertex Array Object) and bind it
        self.vertex_array_object = glGenVertexArrays(1)
        glBindVertexArray(self.vertex_array_object)

        self.shader = OpenGL.GL.shaders.compileProgram(
            OpenGL.GL.shaders.compileShader(
                vertex_shader, GL_VERTEX_SHADER),
            OpenGL.GL.shaders.compileShader(
                fragment_shader, GL_FRAGMENT_SHADER))

        # Generate buffers to hold our vertices
        vertex_buffer = glGenBuffers(1)
        glBindBuffer(GL_ARRAY_BUFFER, vertex_buffer)
        # Get the position of the 'position' in
        # parameter of our shader and bind it.
        position = glGetAttribLocation(self.shader, 'position')
        glEnableVertexAttribArray(position)
        # Describe the position data layout in the buffer
        glVertexAttribPointer(
            position, 4, GL_FLOAT, False, 0, ctypes.c_void_p(0))
        # Send the data over to the buffer
        glBufferData(
            GL_ARRAY_BUFFER, 4 * len(vertices), vertices, GL_STATIC_DRAW)

        tex_coords_buf = glGenBuffers(1)
        glBindBuffer(GL_ARRAY_BUFFER, tex_coords_buf)
        glBufferData(
            GL_ARRAY_BUFFER, 4 * len(tex_coords), tex_coords, GL_STATIC_DRAW)
        tex_coords_loc = glGetAttribLocation(self.shader, "TexCoord0")
        glEnableVertexAttribArray(tex_coords_loc)
        glVertexAttribPointer(
            tex_coords_loc, 2, GL_FLOAT, False, 0, ctypes.c_void_p(0))

        self.load_2d_texture()

        #glBindVertexArray(0)
        #glBindBuffer(GL_ARRAY_BUFFER, 0)

        # Unbind the VAO first (Important)
        glBindVertexArray(0)
        # Unbind other stuff
        glBindBuffer(GL_ARRAY_BUFFER, 0)

    def set_animation(self, name):
        self.current_animation = name

    def _get_transform(self):
        """
        Transform the from local coordinates to world coordinates.

        @return: transformation matrix used to transform from local coords
                 to world coords.
        @rtype: Matrix44
        """
        transform = Matrix44()

        transform.set_row(0, [self.width, 0.0, 0.0, self.pos_x])
        transform.set_row(1, [0.0, self.height, 0.0, self.pos_y])
        transform.set_row(2, [0.0, 0.0, 1.0, self.layer])
        transform.set_row(3, [0.0, 0.0, 0.0, 1.0])

        return transform

    def load_2d_texture(self):
        tex_data = Image.open('{0}/color.png'.format(self.path))
        t_id = glGenTextures(1)
        t_width, t_height = tex_data.size
        t_data = tex_data.convert("RGBA").tostring("raw", "RGBA")

        glActiveTexture(GL_TEXTURE0)
        glBindTexture(GL_TEXTURE_2D, t_id)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_BASE_LEVEL, 0)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAX_LEVEL, 0)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_REPEAT)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_REPEAT)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR)
        #glPixelStorei(GL_UNPACK_ALIGNMENT, 1)
        #glPixelStorei(GL_UNPACK_ROW_LENGTH, 0)
        #glPixelStorei(GL_UNPACK_SKIP_PIXELS, 0)
        #glPixelStorei(GL_UNPACK_SKIP_ROWS, 0)
        glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, t_width, t_height, 0, GL_RGBA,
                     GL_UNSIGNED_BYTE, t_data)

    def draw(self, proj_mat):
        """
        Draw the sprite.

        @param proj_mat: projection matrix to be passed to the shader.
        @type proj_mat: 4x4 matrix
        """
        global current_frame

        # TODO(hurricanerix): use a timer, but for now, slow things down some.
        frames = self.animations[self.current_animation]
        current_frame = current_frame + 1
        if current_frame == len(frames) * 30:
            current_frame = 0
        tmp = int(math.floor(current_frame / 30) % len(frames))

        (x, y) = frames[tmp]

        glUseProgram(self.shader)

        loc_sprite_x = glGetUniformLocation(self.shader, 'x')
        glUniform1i(loc_sprite_x, x)
        loc_sprite_y = glGetUniformLocation(self.shader, 'y')
        glUniform1i(loc_sprite_y, y)

        loc_proj_mat = glGetUniformLocation(self.shader, 'proj_mat')
        glUniformMatrix4fv(loc_proj_mat, 1, GL_FALSE, proj_mat.to_opengl())

        offset = self._get_transform()
        loc_offset = glGetUniformLocation(self.shader, 'offset')
        glUniformMatrix4fv(loc_offset, 1, GL_FALSE, offset.to_opengl())

        glBindVertexArray(self.vertex_array_object)
        glDrawArrays(GL_TRIANGLES, 0, int(len(vertices) / 4.0))
        glBindVertexArray(0)

        glUseProgram(0)
