/* The MIT License (MIT)
 *
 * Copyright (c) 2014 Richard Hawkins
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */
#version 330

uniform mat4 model_matrix;
uniform mat4 view_matrix;
uniform mat4 proj_matrix;
//uniform mat3 normal_matrix;
uniform mat3 tex_matrix;
uniform vec3 light_position;

in vec4 mc_vertex;
in vec3 mc_normal;
in vec3 mc_tangent;
in vec3 TexCoord0;

smooth out vec2 tex_coord;
smooth out vec3 pos;
smooth out vec3 light_dir;
smooth out vec3 eye_dir;


void main()
{
  mat4 mv_matrix = view_matrix * model_matrix;
  vec4 cc_vertex = mv_matrix * mc_vertex;
  gl_Position = proj_matrix * cc_vertex;
  pos = vec3(model_matrix * mc_vertex);

  tex_coord = vec3(tex_matrix * TexCoord0.stp).st;

  mat3 normal_matrix = mat3x3(mv_matrix);
  normal_matrix = inverse(normal_matrix);
  normal_matrix = transpose(normal_matrix);

  mat3 mv3_matrix = mat3(mv_matrix);
  vec3 n = normalize(mv3_matrix * mc_normal);
  vec3 t = normalize(mv3_matrix * mc_tangent);
  vec3 b = normalize(mv3_matrix * cross(n, t));

  light_dir = vec3(view_matrix * vec4(light_position, 0.0)) - vec3(cc_vertex);

  vec3 v;
  v.x = dot(light_dir, t);
  v.y = dot(light_dir, b);
  v.z = dot(light_dir, n);
  light_dir = v;

  eye_dir = vec3(-cc_vertex);
  v.x = dot(eye_dir, t);
  v.y = dot(eye_dir, b);
  v.z = dot(eye_dir, n);
  eye_dir = v;
}
