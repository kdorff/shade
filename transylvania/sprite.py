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
import ctypes
import numpy

import OpenGL.GL.shaders
from OpenGL.GL import (
    glBindBuffer, glBindVertexArray, glBufferData, glDrawArrays,
    glEnableVertexAttribArray, glGenBuffers, glGenVertexArrays,
    glGetAttribLocation, glUseProgram, glVertexAttribPointer)
from OpenGL.GL import (
    GL_ARRAY_BUFFER, GL_FLOAT, GL_FRAGMENT_SHADER, GL_TRIANGLES,
    GL_STATIC_DRAW, GL_VERTEX_SHADER)

vertex_shader = """
#version 330

in vec4 position;
void main()
{
   gl_Position = position;
}
"""

fragment_shader = """
#version 330

out vec4 MyFragColor;

void main()
{
   MyFragColor = vec4(1.0f, 0.0f, 0.0f, 1.0f);
}
"""

vertices = [0.6,  0.6, 0.0, 1.0,
            -0.6,  0.6, 0.0, 1.0,
            0.0, -0.6, 0.0, 1.0]

vertices = numpy.array(vertices, dtype=numpy.float32)


class Sprite(object):
    """
    """

    def __init__(self):
        """
        """
        # Create a new VAO (Vertex Array Object) and bind it
        self.vertex_array_object = glGenVertexArrays(1)
        glBindVertexArray(self.vertex_array_object)

        # Generate buffers to hold our vertices
        vertex_buffer = glGenBuffers(1)
        glBindBuffer(GL_ARRAY_BUFFER, vertex_buffer)

        self.shader = OpenGL.GL.shaders.compileProgram(
            OpenGL.GL.shaders.compileShader(
                vertex_shader, GL_VERTEX_SHADER),
            OpenGL.GL.shaders.compileShader(
                fragment_shader, GL_FRAGMENT_SHADER))

        # Get the position of the 'position' in
        # parameter of our shader and bind it.
        position = glGetAttribLocation(self.shader, 'position')
        glEnableVertexAttribArray(position)

        # Describe the position data layout in the buffer
        glVertexAttribPointer(
            position, 4, GL_FLOAT, False, 0, ctypes.c_void_p(0))

        # Send the data over to the buffer
        glBufferData(GL_ARRAY_BUFFER, 48, vertices, GL_STATIC_DRAW)

        # Unbind the VAO first (Important)
        glBindVertexArray(0)

        # Unbind other stuff
        glBindBuffer(GL_ARRAY_BUFFER, 0)

    def draw(self):
        """
        """
        glUseProgram(self.shader)

        glBindVertexArray(self.vertex_array_object)
        glDrawArrays(GL_TRIANGLES, 0, 3)
        glBindVertexArray(0)

        glUseProgram(0)
