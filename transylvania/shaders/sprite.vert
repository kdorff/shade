/* Copyright 2014 Richard Hawkins
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
