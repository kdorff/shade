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

uniform sampler2D ColorMap;
uniform sampler2D NormalMap;
uniform vec3 light_position;
uniform vec3 light_color;
uniform float light_power;

smooth in vec3 pos;
smooth in vec3 light_dir;
smooth in vec3 eye_dir;
smooth in vec2 tex_coord;

out vec4 frag_color;


void main()
{
    float alpha = texture(ColorMap, tex_coord.st).a;
    vec3 diffuse = texture(ColorMap, tex_coord.st).rgb;
    vec3 ambient = vec3(0.2, 0.2, 0.2) * diffuse;
    vec3 specular = diffuse / 8;

    vec3 normal = texture(NormalMap, tex_coord.st).rgb * 2 - 1;
    float distance = length(light_position - pos);

    vec3 n = normalize(normal);
    vec3 l = normalize(light_dir);

    float cos_theta = clamp(dot(n, l), 0.0, 1.0);

    vec3 e = normalize(eye_dir);
    vec3 r = reflect(-l, n);

    float cos_alpha = clamp(dot(e, r), 0.0, 1.0);

    frag_color = vec4(
        ambient +
        diffuse * light_color * light_power * cos_theta /
            (distance * distance) +
        specular * light_color * light_power * pow(cos_alpha, 5) /
            (distance * distance), alpha);
}
