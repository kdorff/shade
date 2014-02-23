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
import numpy
import os.path

import OpenGL.GL.shaders
from OpenGL.GL import (
    glActiveTexture, glAttachShader, glBindAttribLocation, glBindBuffer,
    glBindVertexArray,
    glBindTexture, glBufferData, glCreateProgram, glDrawArrays,
    glEnableVertexAttribArray, glGenBuffers, glGenTextures, glGenVertexArrays,
    glGetProgramInfoLog, glGetProgramiv,
    glGetUniformLocation, glLinkProgram,
    glTexImage2D,
    glTexParameteri,
    glUniform1i, glUniform3fv, glUniformMatrix3fv, glUniformMatrix4fv,
    glUseProgram,
    glVertexAttribPointer)
from OpenGL.GL import (
    GL_ARRAY_BUFFER, GL_FALSE, GL_FLOAT,
    GL_FRAGMENT_SHADER, GL_LINEAR, GL_LINK_STATUS, GL_REPEAT, GL_RGBA,
    GL_TEXTURE0, GL_TEXTURE1, GL_TEXTURE2, GL_TEXTURE_2D,
    GL_TEXTURE_BASE_LEVEL,
    GL_TEXTURE_MAG_FILTER, GL_TEXTURE_MAX_LEVEL, GL_TEXTURE_MIN_FILTER,
    GL_TEXTURE_WRAP_S, GL_TEXTURE_WRAP_T, GL_TRIANGLES, GL_TRUE,
    GL_STATIC_DRAW,
    GL_UNSIGNED_BYTE, GL_VERTEX_SHADER)
from PIL import Image

from transylvania.gmath import get_3x3_transform, get_4x4_transform


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


vertex_shader = """
#version 330

in vec4 mc_vertex;
in vec3 TexCoord0;

uniform mat4 view_matrix;
uniform mat4 proj_matrix;
uniform mat3 tex_matrix;
uniform vec3 light_position;

smooth out vec3 light_dir;
smooth out vec3 eye_dir;
smooth out vec2 tex_coord;

void main()
{
  gl_Position = mc_vertex * view_matrix * proj_matrix;
  tex_coord = vec3(TexCoord0 * tex_matrix).st;

  vec4 pos = view_matrix * mc_vertex;
  eye_dir = pos.xyz;

  vec3 n = vec3(0.0, 0.0, 1.0);
  vec3 t = vec3(0.0, 1.0, 0.0);
  vec3 b = cross(n, t);

  vec3 v;
  v.x = dot(light_position, t);
  v.y = dot(light_position, b);
  v.z = dot(light_position, n);
  light_dir = normalize(v);

  v.x = dot(eye_dir, t);
  v.y = dot(eye_dir, b);
  v.z = dot(eye_dir, n);
  eye_dir = normalize(v);
}
"""

fragment_shader = """
#version 330

uniform sampler2D ColorMap;
uniform sampler2D BumpMap;
uniform sampler2D NormalMap;

in vec3 light_dir;
in vec3 eye_dir;
in vec2 tex_coord;
out vec4 frag_color;

void main()
{
  vec4 color = texture(ColorMap, tex_coord.st).rgba;
  vec3 normDelta = texture(NormalMap, tex_coord.st).rgb;

  vec3 litColor = color.rgb * max(dot(normDelta, light_dir), 0.0);
  vec3 reflect_dir = reflect(light_dir, normDelta);

  float spec = max(dot(eye_dir, reflect_dir), 0.0);
  spec = pow(spec, 6.0);
  litColor = min(litColor + spec, vec3(1.0));
  frag_color = vec4(litColor, color.a);
}
"""

vao = None
shader = None
shader_locs = None


def get_vao():
    global vao
    global shader_locs
    if vao:
        return vao
    shader_locs = {'mc_vertex': 0, 'TexCoord0': 1}
    vao = glGenVertexArrays(1)
    glBindVertexArray(vao)

    vertex_buffer = glGenBuffers(1)
    glBindBuffer(GL_ARRAY_BUFFER, vertex_buffer)
    glVertexAttribPointer(shader_locs['mc_vertex'], 4, GL_FLOAT, False, 0,
                          ctypes.c_void_p(0))
    glBufferData(GL_ARRAY_BUFFER, 4 * len(vertices), vertices, GL_STATIC_DRAW)

    tex_coords_buf = glGenBuffers(1)
    glBindBuffer(GL_ARRAY_BUFFER, tex_coords_buf)
    glVertexAttribPointer(shader_locs['TexCoord0'], 3, GL_FLOAT, False, 0,
                          ctypes.c_void_p(0))
    glBufferData(GL_ARRAY_BUFFER, 4 * len(tex_coords), tex_coords,
                 GL_STATIC_DRAW)

    glBindBuffer(GL_ARRAY_BUFFER, 0)
    glBindVertexArray(0)
    return vao


def get_shader():
    global shader
    global shader_locs
    if shader:
        return (shader, shader_locs)

    vert_prog = OpenGL.GL.shaders.compileShader(
        vertex_shader, GL_VERTEX_SHADER)
    frag_prog = OpenGL.GL.shaders.compileShader(
        fragment_shader, GL_FRAGMENT_SHADER)

    shader = glCreateProgram()
    glAttachShader(shader, vert_prog)
    glAttachShader(shader, frag_prog)

    for attrib in shader_locs:
        glBindAttribLocation(shader, shader_locs[attrib], attrib)

    glLinkProgram(shader)
    if glGetProgramiv(shader, GL_LINK_STATUS) != GL_TRUE:
        raise RuntimeError(glGetProgramInfoLog(shader))

    for name in ['view_matrix', 'proj_matrix', 'tex_matrix', 'ColorMap',
                 'BumpMap', 'NormalMap', 'light_position']:
        shader_locs[name] = glGetUniformLocation(shader, name)

    return (shader, shader_locs)


class SpriteBuilder(object):
    """
    Handles reading in sprite resources and creating a sprite from them.
    """

    @staticmethod
    def build(path):
        data_path = '{0}/data.json'.format(path)
        if not os.path.isfile(data_path):
            raise Exception('file {0} does not exist.'.format(data_path))
        data = json.loads(open(data_path).read())

        tex_data = {'color': None, 'normal': None, 'bump': None}
        color_path = '{0}/color.png'.format(path)
        if not os.path.isfile(color_path):
            raise Exception('file {0} does not exist.'.format(color_path))
        img = Image.open(color_path)
        width, height = img.size
        tex_data['color'] = img.convert("RGBA").tostring("raw", "RGBA")

        bump_path = '{0}/bump.png'.format(path)
        if os.path.isfile(bump_path):
            img = Image.open(bump_path)
            tex_data['bump'] = img.convert("RGBA").tostring("raw", "RGBA")

        normal_path = '{0}/normal.png'.format(path)
        if os.path.isfile(normal_path):
            img = Image.open(normal_path)
            tex_data['normal'] = img.convert("RGBA").tostring("raw", "RGBA")

        return Sprite(width, height, data, tex_data)


class SpriteManager(object):
    """
    Manages sprites such that they can be used many times while only being
    loaded once.
    """

    def __init__(self, sprite_dir):
        """
        @param sprite_dir:
        @type sprite_dir: string
        """
        self.sprite_dir = sprite_dir
        self.sprites = {}

    def load(self, *names):
        """
        @param names:
        @type names:
        """
        for name in names:
            path = '{0}/{1}'.format(self.sprite_dir, name)
            self.sprites[name] = SpriteBuilder.build(path)

    def get_sprite(self, path):
        """
        @param path:
        @type path:
        @return:
        @rtype:
        """
        return self.sprites[path]


class Sprite(object):
    """
    Handles texturing images on a polygon.
    """

    def __init__(self, width, height, data, tex_data):
        """
        Initialize the OpenGL things needed to render the polygon.
        @param data:
        @type data:
        @param tex_data:
        @type tex_data:
        """
        self.width = width
        self.height = height
        self.data = data
        self.tex_data = tex_data
        self.texture_ids = None

        self.light_position = [0.0, 0.0, 1.0]
        self.light_dir = 0.01

    def _bind_textures(self, data):
        samplers = {
            'color': GL_TEXTURE0,
            'bump': GL_TEXTURE1,
            'normal': GL_TEXTURE2}

        for map_type in data:
            if data[map_type] is None:
                continue
            if self.texture_ids is None:
                self.texture_ids = {'color': None, 'bump': None,
                                    'normal': None}
            if self.texture_ids[map_type] is not None:
                glBindTexture(GL_TEXTURE_2D, self.texture_ids[map_type])
                continue
            self.texture_ids[map_type] = glGenTextures(1)
            glActiveTexture(samplers[map_type])
            glBindTexture(GL_TEXTURE_2D, self.texture_ids[map_type])
            glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_BASE_LEVEL, 0)
            glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAX_LEVEL, 0)
            glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_REPEAT)
            glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_REPEAT)
            glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR)
            glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR)
            glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, self.width, self.height, 0,
                         GL_RGBA, GL_UNSIGNED_BYTE, data[map_type])

    def draw(self, proj_matrix, x, y, layer=0, frame_x=0, frame_y=0):
        """
        Draw the sprite.

        @param proj_mat: projection matrix to be passed to the shader.
        @type proj_mat: 4x4 matrix
        """

        vao = get_vao()
        (shader, shader_locs) = get_shader()
        self._bind_textures(self.tex_data)

        glBindVertexArray(vao)
        glUseProgram(shader)

        glEnableVertexAttribArray(shader_locs['mc_vertex'])
        glEnableVertexAttribArray(shader_locs['TexCoord0'])

        glUniformMatrix4fv(shader_locs['proj_matrix'], 1, GL_FALSE,
                           proj_matrix)

        view_matrix = get_4x4_transform(
            scale_x=self.data['frame']['size']['width'],
            scale_y=self.data['frame']['size']['height'],
            trans_x=0, trans_y=0, layer=0)
        glUniformMatrix4fv(shader_locs['view_matrix'], 1, GL_FALSE,
                           view_matrix)

        scale_x = 1.0/self.data['frame']['count']['x']
        scale_y = 1.0/self.data['frame']['count']['y']
        tex_matrix = get_3x3_transform(scale_x, scale_y, frame_x, frame_y)
        glUniformMatrix3fv(shader_locs['tex_matrix'], 1, GL_FALSE, tex_matrix)

        if self.light_position[0] > 1.0:
            self.light_position[0] = 1.0
            self.light_dir = -0.01
        if self.light_position[0] < 0.0:
            self.light_position[0] = 0.0
            self.light_dir = 0.01
        self.light_position[0] = self.light_position[0] + self.light_dir

        glUniform3fv(shader_locs['light_position'], 1, self.light_position)

        if self.tex_data['color']:
            glUniform1i(shader_locs['ColorMap'], 0)

        if self.tex_data['bump']:
            glUniform1i(shader_locs['BumpMap'], 1)

        if self.tex_data['normal']:
            glUniform1i(shader_locs['NormalMap'], 2)

        glDrawArrays(GL_TRIANGLES, 0, int(len(vertices) / 4.0))

        glBindVertexArray(0)
        glUseProgram(0)
