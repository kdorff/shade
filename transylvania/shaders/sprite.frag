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

uniform sampler2D ColorMap;
uniform sampler2D NormalMap;
uniform vec3 light_position;

smooth in vec3 pos;
smooth in vec3 light_dir;
smooth in vec3 eye_dir;
smooth in vec2 tex_coord;

out vec4 frag_color;


void main()
{
    vec3 light_color = vec3(0.0, 1.0, 0.0);
    float light_power = 10000.0;

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
