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

import ctypes
import numpy
import os

from OpenGL import GL as gl
from OpenGL.GL import shaders as glsl
from PIL import Image as image

import gmath

vert_shader = """#version 330
uniform mat4 model_matrix;
uniform mat4 view_matrix;
uniform mat4 proj_matrix;
uniform mat3 tex_matrix;
in vec4 mc_vertex;
in vec3 TexCoord0;
smooth out vec2 tex_coord;

void main() {
  mat4 mv_matrix = view_matrix * model_matrix;
  vec4 cc_vertex = mv_matrix * mc_vertex;
  gl_Position = proj_matrix * cc_vertex;
  tex_coord = vec3(tex_matrix * TexCoord0.stp).st;
}
"""

frag_shader = """#version 330
uniform sampler2D ColorMap;
smooth in vec2 tex_coord;
out vec4 frag_color;

void main() {
    frag_color = vec4(texture(ColorMap, tex_coord.st).rgba);
}
"""
vertices = [0.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 0.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0]
vertices = numpy.array(vertices, dtype=numpy.float32)
tex_coords = [0.0, 1.0, 1.0,
              1.0, 1.0, 1.0,
              1.0, 0.0, 1.0,
              0.0, 0.0, 1.0,
              0.0, 1.0, 1.0,
              1.0, 0.0, 1.0]
tex_coords = numpy.array(tex_coords, dtype=numpy.float32)
vao = None
shader = None
shader_locs = None


def _get_vao():
    global vao
    global shader_locs
    if vao:
        return vao
    shader_locs = {'mc_vertex': 0, 'mc_normal': 1, 'TexCoord0': 2}
    vao = gl.glGenVertexArrays(1)
    gl.glBindVertexArray(vao)

    vertex_buffer = gl.glGenBuffers(1)
    gl.glBindBuffer(gl.GL_ARRAY_BUFFER, vertex_buffer)
    gl.glVertexAttribPointer(shader_locs['mc_vertex'], 4, gl.GL_FLOAT, False,
                             0, ctypes.c_void_p(0))
    gl.glBufferData(gl.GL_ARRAY_BUFFER, 4 * len(vertices), vertices,
                    gl.GL_STATIC_DRAW)

    tex_coords_buf = gl.glGenBuffers(1)
    gl.glBindBuffer(gl.GL_ARRAY_BUFFER, tex_coords_buf)
    gl.glVertexAttribPointer(shader_locs['TexCoord0'], 3, gl.GL_FLOAT, False,
                             0, ctypes.c_void_p(0))
    gl.glBufferData(gl.GL_ARRAY_BUFFER, 4 * len(tex_coords), tex_coords,
                    gl.GL_STATIC_DRAW)

    gl.glBindBuffer(gl.GL_ARRAY_BUFFER, 0)
    gl.glBindVertexArray(0)
    return vao


def _get_shader():
    global shader
    global shader_locs
    if shader:
        return (shader, shader_locs)

    vert_prog = glsl.compileShader(
        vert_shader, gl.GL_VERTEX_SHADER)
    frag_prog = glsl.compileShader(
        frag_shader, gl.GL_FRAGMENT_SHADER)

    shader = gl.glCreateProgram()
    gl.glAttachShader(shader, vert_prog)
    gl.glAttachShader(shader, frag_prog)

    for attrib in shader_locs:
        gl.glBindAttribLocation(shader, shader_locs[attrib], attrib)

    gl.glLinkProgram(shader)
    if gl.glGetProgramiv(shader, gl.GL_LINK_STATUS) != gl.GL_TRUE:
        raise RuntimeError(gl.glGetProgramInfoLog(shader))

    for name in ['model_matrix', 'view_matrix', 'proj_matrix', 'tex_matrix',
                 'ColorMap']:
        shader_locs[name] = gl.glGetUniformLocation(shader, name)

    return (shader, shader_locs)


def load(filename):
    img = image.open(filename)
    width, height = img.size
    return img.convert("RGBA").tostring("raw", "RGBA")


def build(filename):
    if not os.path.isfile(filename):
        raise Exception('file {0} does not exist.'.format(filename))
    img = image.open(filename)
    width, height = img.size
    data = img.convert("RGBA").tostring("raw", "RGBA")
    return Sprite(width, height, data)


class Sprite(object):
    """Handles texturing images on a polygon."""

    def __init__(self, width, height, data):
        """Initialize the OpenGL things needed to render the polygon.

        @param data:
        @type data:
        """
        self.width = width
        self.height = height
        self.data = data
        self.texture_id = None

    def get_size(self):
        return self.width, self.height

    def _bind_textures(self):
        if self.texture_id is not None:
            gl.glActiveTexture(gl.GL_TEXTURE0)
            gl.glBindTexture(gl.GL_TEXTURE_2D, self.texture_id)
            return
        self.texture_id = gl.glGenTextures(1)
        gl.glActiveTexture(gl.GL_TEXTURE0)
        gl.glBindTexture(gl.GL_TEXTURE_2D, self.texture_id)
        gl.glTexParameteri(gl.GL_TEXTURE_2D, gl.GL_TEXTURE_BASE_LEVEL, 0)
        gl.glTexParameteri(gl.GL_TEXTURE_2D, gl.GL_TEXTURE_MAX_LEVEL, 0)
        gl.glTexParameteri(gl.GL_TEXTURE_2D, gl.GL_TEXTURE_WRAP_S,
                           gl.GL_REPEAT)
        gl.glTexParameteri(gl.GL_TEXTURE_2D, gl.GL_TEXTURE_WRAP_T,
                           gl.GL_REPEAT)
        gl.glTexParameteri(gl.GL_TEXTURE_2D, gl.GL_TEXTURE_MAG_FILTER,
                           gl.GL_LINEAR)
        gl.glTexParameteri(gl.GL_TEXTURE_2D, gl.GL_TEXTURE_MIN_FILTER,
                           gl.GL_LINEAR)
        gl.glTexImage2D(gl.GL_TEXTURE_2D, 0, gl.GL_RGBA, self.width,
                        self.height, 0, gl.GL_RGBA, gl.GL_UNSIGNED_BYTE,
                        self.data)

    def draw(self, proj_matrix, view_matrix, x, y, z=0):
        """Draw the sprite.

        @param proj_mat: projection matrix to be passed to the shader.
        @type proj_mat: 4x4 matrix
        """
        x = int(x)
        y = int(y)
        vao = _get_vao()
        (shader, shader_locs) = _get_shader()
        self._bind_textures()

        gl.glBindVertexArray(vao)
        gl.glUseProgram(shader)

        gl.glEnableVertexAttribArray(shader_locs['mc_vertex'])
        gl.glEnableVertexAttribArray(shader_locs['TexCoord0'])

        model_matrix = gmath.get_4x4_transform(scale_x=self.width, scale_y=self.height,
                                               trans_x=x, trans_y=y, trans_z=z)
        gl.glUniformMatrix4fv(shader_locs['model_matrix'], 1, gl.GL_TRUE,
                              model_matrix)

        gl.glUniformMatrix4fv(shader_locs['view_matrix'], 1, gl.GL_TRUE,
                              view_matrix)

        gl.glUniformMatrix4fv(shader_locs['proj_matrix'], 1, gl.GL_TRUE,
                              proj_matrix)

        tex_matrix = gmath.get_3x3_transform()
        gl.glUniformMatrix3fv(shader_locs['tex_matrix'], 1, gl.GL_TRUE,
                              tex_matrix)

        gl.glUniform1i(shader_locs['ColorMap'], 0)

        gl.glDrawArrays(gl.GL_TRIANGLES, 0, int(len(vertices) / 4.0))

        gl.glBindVertexArray(0)
        gl.glUseProgram(0)
