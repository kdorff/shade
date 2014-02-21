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
import json
import math
import numpy

import OpenGL.GL.shaders
from OpenGL.GL import (
    glActiveTexture, glBindBuffer, glBindVertexArray, glBindTexture,
    glBlendFunc, glBufferData, glDrawArrays, glEnable,
    glEnableVertexAttribArray, glGenBuffers, glGenTextures, glGenVertexArrays,
    glGetAttribLocation, glGetUniformLocation, glTexImage2D, glTexParameteri,
    glUniform1i, glUniformMatrix3fv, glUniformMatrix4fv, glUseProgram,
    glVertexAttribPointer)
from OpenGL.GL import (
    GL_ARRAY_BUFFER, GL_BLEND, GL_CULL_FACE, GL_FALSE, GL_FLOAT,
    GL_FRAGMENT_SHADER, GL_ONE_MINUS_SRC_ALPHA, GL_LINEAR, GL_REPEAT, GL_RGBA,
    GL_SRC_ALPHA, GL_TEXTURE0, GL_TEXTURE_2D, GL_TEXTURE_BASE_LEVEL,
    GL_TEXTURE_MAG_FILTER, GL_TEXTURE_MAX_LEVEL, GL_TEXTURE_MIN_FILTER,
    GL_TEXTURE_WRAP_S, GL_TEXTURE_WRAP_T, GL_TRIANGLES, GL_STATIC_DRAW,
    GL_UNSIGNED_BYTE, GL_VERTEX_SHADER)
from PIL import Image

from gameobjects.matrix44 import Matrix44


vertex_shader = """
#version 330

uniform int debug;
in vec3 TexCoord0;
uniform mat3 tex_trans_mat;
uniform int use_alt_tex;
uniform mat3 alt_tex_trans_mat;
smooth out vec4 debug_position;
smooth out vec2 TexCoord;
smooth out vec2 AltTexCoord;

in vec4 position;
uniform mat4 proj_mat;
uniform mat4 offset;

void main()
{
  TexCoord = vec3(TexCoord0 * tex_trans_mat).st;
  if (use_alt_tex == 1) {
    AltTexCoord = vec3(TexCoord0 * alt_tex_trans_mat).st;
  }

  if (debug == 1) {
    debug_position = position;
  }
  gl_Position = position * offset * proj_mat;
}
"""

fragment_shader = """
#version 330

uniform int debug;
uniform int use_alt_tex;
uniform int layer;
uniform sampler2D ColorMap;
in vec4 debug_position;
in vec2 TexCoord;
in vec2 AltTexCoord;
out vec4 MyFragColor;

void main()
{
  vec4 color = texture(ColorMap, TexCoord.st).rgba;
  if (use_alt_tex == 1) {
    vec4 alt_color = texture(ColorMap, AltTexCoord.st).rgba;
    if (alt_color.a == 1) {
    color = alt_color;
    } else {
    color = color + alt_color;
    }
  }
  if (debug == 1) {
    if (debug_position.x >= 0.995 || debug_position.x <= 0.005 ||
        debug_position.y >= 0.995 || debug_position.y <= 0.005) {
      color = vec4(1.0, 1.0, 1.0, 1.0);
      if (layer < 0) {
        color = vec4(0.0, 0.0, 0.0, 1.0);
      } else if (layer == 0) {
        color = vec4(0.0, 0.0, 1.0, 1.0);
      } else if (layer == 1) {
        color = vec4(0.0, 1.0, 0.0, 1.0);
      } else if (layer == 2) {
        color = vec4(0.0, 1.0, 1.0, 1.0);
      } else if (layer == 3) {
        color = vec4(1.0, 0.0, 0.0, 1.0);
      } else if (layer == 4) {
        color = vec4(1.0, 0.0, 1.0, 1.0);
      } else if (layer == 5) {
        color = vec4(1.0, 1.0, 0.0, 1.0);
      } else if (layer == 6) {
        color = vec4(1.0, 1.0, 1.0, 1.0);
      }
    }
  }
  MyFragColor = color;
}
"""

vertices = [0.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 0.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0]

tex_coords = [0.0, 1.0, 1.0,
              1.0, 1.0, 1.0,
              1.0, 0.0, 1.0,
              0.0, 0.0, 1.0,
              0.0, 1.0, 1.0,
              1.0, 0.0, 1.0]

vertices = numpy.array(vertices, dtype=numpy.float32)
tex_coords = numpy.array(tex_coords, dtype=numpy.float32)


class SpriteBuilder(object):
    """
    Handles reading in sprite resources and creating a sprite from them.
    """

    @staticmethod
    def build(path, pos_x=0, pos_y=0, layer=0):
        data_path = '{0}/data.json'.format(path)
        # TODO: Test if file exists first
        data = json.loads(open(data_path).read())
        # TODO: create shader and pass it to the sprite.
        # TODO: Read in color data and pass it to the sprite.
        return Sprite(path, pos_x=pos_x, pos_y=pos_y, layer=layer, data=data)


class Sprite(object):
    """
    Handles texturing images on a polygon.
    """

    def __init__(self, path, pos_x=0, pos_y=0, layer=0, data=None):
        """
        Initialize the OpenGL things needed to render the polygon.
        """
        # TODO(hurricanerix): position stuff should probably be moved outside
        # of the sprite class.
        self.path = path
        self.pos_x = pos_x
        self.pos_y = pos_y
        self.layer = layer
        self.data = data
        self.debug = True

        self.current_frame = 0
        self.current_animation = 'default'
        self.animate = True

        if 'animations' not in data:
            self.data['animations'] = {'default': [[0, 0]]}

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
            tex_coords_loc, 3, GL_FLOAT, False, 0, ctypes.c_void_p(0))

        #glBindVertexArray(0)
        #glBindBuffer(GL_ARRAY_BUFFER, 0)

        # Unbind the VAO first (Important)
        glBindVertexArray(0)
        # Unbind other stuff
        glBindBuffer(GL_ARRAY_BUFFER, 0)

    def set_animation(self, name):
        self.animate = True
        self.current_animation = name

    def _get_transform(self):
        """
        Transform the from local coordinates to world coordinates.

        @return: transformation matrix used to transform from local coords
                 to world coords.
        @rtype: Matrix44
        """
        transform = Matrix44()
        width = self.data['frame']['size']['width']
        height = self.data['frame']['size']['height']

        transform.set_row(0, [width, 0.0, 0.0, self.pos_x])
        transform.set_row(1, [0.0, height, 0.0, self.pos_y])
        transform.set_row(2, [0.0, 0.0, 1.0, self.layer])
        transform.set_row(3, [0.0, 0.0, 0.0, 1.0])

        return transform

    def _get_tex_trans_matrix(self, frame_x, frame_y):
        """
        Transform the tex coords to shrink the texutre down to a single frame.

        @return: matrix used to transform texture coords to the current
                 frame in the texture to be displayed.
        @rtype: list
        """
        scale_x = 1.0/self.data['frame']['count']['x']
        scale_y = 1.0/self.data['frame']['count']['y']
        trans_x = frame_x * scale_x
        trans_y = frame_y * scale_y
        transform = [scale_x, 0.0, trans_x,
                     0.0, scale_y, trans_y,
                     0.0, 0.0, 1.0]
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
        # TODO(hurricanerix): use a timer, but for now, slow things down some.
        self.load_2d_texture()

        frames = self.data['animations'][self.current_animation]
        if self.animate:
            self.current_frame = self.current_frame + 1
            if self.current_frame == len(frames) * 10:
                self.current_frame = 0
        tmp = int(math.floor(self.current_frame / 10) % len(frames))

        if frames[tmp]:
            (x, y) = frames[tmp]
        else:
            self.animate = False
            (x, y) = frames[tmp - 1]

        tex_trans_mat = self._get_tex_trans_matrix(x, y)

        glUseProgram(self.shader)

        loc_tex_trans_mat = glGetUniformLocation(self.shader, 'tex_trans_mat')
        glUniformMatrix3fv(loc_tex_trans_mat, 1, GL_FALSE, tex_trans_mat)

        if use_alt:
            alt_frames = self.data['animations'][use_alt]['frames']
            (alt_x, alt_y) = alt_frames[tmp - 1]
            alt_tex_trans_mat = self._get_tex_trans_matrix(alt_x, alt_y)
            loc_alt_tex_trans_mat = glGetUniformLocation(
                self.shader, 'alt_tex_trans_mat')
            glUniformMatrix3fv(
                loc_alt_tex_trans_mat, 1, GL_FALSE, alt_tex_trans_mat)

        loc_debug = glGetUniformLocation(self.shader, 'debug')
        glUniform1i(loc_debug, self.debug)

        loc_layer = glGetUniformLocation(self.shader, 'layer')
        glUniform1i(loc_layer, self.layer)

        loc_use_alt_tex = glGetUniformLocation(self.shader, 'use_alt_tex')
        glUniform1i(loc_use_alt_tex, use_alt != None)

        loc_proj_mat = glGetUniformLocation(self.shader, 'proj_mat')
        glUniformMatrix4fv(loc_proj_mat, 1, GL_FALSE, proj_mat.to_opengl())

        offset = self._get_transform()
        loc_offset = glGetUniformLocation(self.shader, 'offset')
        glUniformMatrix4fv(loc_offset, 1, GL_FALSE, offset.to_opengl())

        glBindVertexArray(self.vertex_array_object)
        glDrawArrays(GL_TRIANGLES, 0, int(len(vertices) / 4.0))
        glBindVertexArray(0)

        glUseProgram(0)
