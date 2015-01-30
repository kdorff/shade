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

from sdl2 import SDL_Event, SDL_PollEvent, SDL_Quit, events
from sdl2 import SDL_MOUSEBUTTONDOWN, SDL_MOUSEMOTION, SDL_QUIT
import time
import sdl2

from OpenGL import GL as gl

# Version information (major, minor, revision[, 'dev']).
version_info = (0, 0, 1)

# Version string 'major.minor.revision'.
version = __version__ = ".".join(map(str, version_info))
#
#
# class Application(object):
#     """
#     Application class for the Transylvania Engine.
#     """
#
#     def __init__(self, config=None, display=None, sprite_manager=None):
#         """
#         Initialize the application.
#
#         @param config: Application configuration.
#         @type config: dict
#         @param display: Display Manager to be used.
#         @type display: DisplayManager
#         """
#         self.config = config
#         self.display = display
#         self.sprite_manager = sprite_manager
#         self.running = True
#         self.objects = []
#         self.lights = []
#
#     def add_object(self, obj):
#         self.objects.append(obj)
#
#     def add_light(self, light):
#         self.lights.append(light)
#
#     def __del__(self):
#         """
#         Cleanup the application.
#         """
#         del self.display
#         SDL_Quit()
#
#     def handle_mousemotion(self, motion):
#         pass
#
#     def handle_mousebuttondown(self, button):
#         pass
#
#     def handle_keydown(self, key):
#         pass
#
#     def run(self):
#         """
#         Start the application event loop.
#         """
#         self.display.init_window()
#
#         event = SDL_Event()
#
#         prev_time = 0
#         current_time = current_time = time.time()
#         while self.running:
#             prev_time = current_time
#             current_time = time.time()
#             timedelta = current_time - prev_time
#
#             while SDL_PollEvent(ctypes.byref(event)) != 0:
#                 if event.type == SDL_QUIT:
#                     return
#                 if event.type == events.SDL_KEYDOWN:
#                     self.handle_keydown(event.key.keysym.sym)
#                 if (event.type == SDL_MOUSEMOTION):
#                     self.handle_mousemotion(event.motion)
#                 if (event.type == SDL_MOUSEBUTTONDOWN):
#                     self.handle_mousebuttondown(event.button.button)
#
#             for light in self.lights:
#                 light.update(timedelta)
#
#             for obj in self.objects:
#                 obj.update(timedelta)
#
#             self.display.start_render()
#             self.objects.sort(key=lambda obj: obj.layer)
#             proj_matrix = self.display.get_proj_matrix()
#             view_matrix = self.display.get_view_matrix(0, 0)
#
#             for obj in self.objects:
#                 obj.draw(proj_matrix, view_matrix, self.lights)
#
#             self.display.stop_render()


def init():
    sdl2.SDL_Init(sdl2.SDL_INIT_EVERYTHING)

    sdl2.SDL_GL_SetAttribute(sdl2.SDL_GL_CONTEXT_MAJOR_VERSION, 3)
    sdl2.SDL_GL_SetAttribute(sdl2.SDL_GL_CONTEXT_MINOR_VERSION, 2)
    sdl2.SDL_GL_SetAttribute(sdl2.SDL_GL_CONTEXT_PROFILE_MASK,
                             sdl2.SDL_GL_CONTEXT_PROFILE_CORE)


# def set_display_mode(width, height):
#     window = sdl2.SDL_CreateWindow(
#         "PyCon14 Tutorial (SDL2 Port)", sdl2.SDL_WINDOWPOS_CENTERED,
#         sdl2.SDL_WINDOWPOS_CENTERED, width, height,
#         sdl2.SDL_WINDOW_OPENGL | sdl2.SDL_WINDOW_SHOWN)
#     if not window:
#         raise Exception('Could not create window')
#
#     sdl2.SDL_GL_CreateContext(window)
#     sdl2.SDL_GL_SetSwapInterval(1)
#     gl.glBlendFunc(gl.GL_SRC_ALPHA, gl.GL_ONE_MINUS_SRC_ALPHA)
#     gl.glEnable(gl.GL_DEPTH_TEST)
#     gl.glEnable(gl.GL_CULL_FACE)
#     gl.glEnable(gl.GL_BLEND)
#     gl.glClearColor(0.3, 0.3, 0.3, 1.0)
#     return window
#
#
# def get_keyboard_state():
#     """ Returns a list with the current SDL keyboard state,
#     which is updated on SDL_PumpEvents. """
#     numkeys = ctypes.c_int()
#     keystate = sdl2.keyboard.SDL_GetKeyboardState(ctypes.byref(numkeys))
#     ptr_t = ctypes.POINTER(ctypes.c_uint8 * numkeys.value)
#     return ctypes.cast(keystate, ptr_t)[0]
#
#
# class Clock(object):
#     def __init__(self):
#         self.call_count = 0
#         self.last_time = 0
#         self.last_mod = 0
#         self.dt = 0
#
#     def tick(self, count):
#         # TODO: Fix tick, it makes things run like crap
#         current_time = sdl2.timer.SDL_GetTicks()
#         dt = current_time - self.last_time
# #        current_mod = current_time % 1000
# #        # print 'LAST: {}, CUR: {}'.format(self.last_mod, current_mod)
# #        if current_mod < self.last_mod:
# #            # print 'reset'
# #            self.call_count = 0
# #        self.call_count += 1
# #        if self.call_count > count:
# #            # print 'sleep: {}'.format(1000 - current_mod)
# #            sdl2.timer.SDL_Delay(50)
#
# #        self.last_mod = current_mod
#         self.last_time = current_time
#         return dt
#
#
# def spritecollide(sprite, group, dokill, collided=None):
#     collisions = []
#
#     def default_collided(s1, s2):
#         clip_left = (s1.rect.left <= s2.rect.right and
#                      s1.rect.left >= s2.rect.left)
#         clip_right = (s1.rect.right >= s2.rect.left and
#                       s1.rect.right <= s2.rect.right)
#         clip_top = (s1.rect.top >= s2.rect.bottom and
#                    s1.rect.top <= s2.rect.top)
#         clip_bottom = (s1.rect.bottom <= s2.rect.top and
#                      s1.rect.bottom >= s2.rect.bottom)
#
#         if clip_left and clip_top:
#             return True
#         if clip_right and clip_top:
#             return True
#         if clip_left and clip_bottom:
#             return True
#         if clip_right and clip_bottom:
#             return True
#         if clip_left and clip_right and clip_bottom:
#             return True
#         if clip_left and clip_right and clip_top:
#             return True
#         if clip_left and clip_right and clip_top and clip_bottom:
#             return True
#         return False
#
#     if not collided:
#         collided = default_collided
#
#     for s in group:
#         if collided(sprite, s):
#             if dokill:
#                 return [s]
#             collisions.append(s)
#     return collisions
#
#
# class Rect(object):
#     def __init__(self, pos=None, dims=None):
#         if not pos:
#             pos = (0, 0)
#         if not dims:
#             dims = (0, 0)
#         self.x = pos[0]
#         self.y = pos[1]
#         self.w = dims[0]
#         self.h = dims[1]
#
#     @property
#     def left(self):
#         return self.x
#
#     @property
#     def right(self):
#         return self.x + self.w
#
#     @property
#     def top(self):
#         return self.y + self.h
#
#     @property
#     def bottom(self):
#         return self.y
#
#     @left.setter
#     def left(self, value):
#         self.x = value
#
#     @bottom.setter
#     def bottom(self, value):
#         self.y = value
#
#     @right.setter
#     def right(self, value):
#         self.x = value - self.w
#
#     @top.setter
#     def top(self, value):
#         self.y = value - self.h
#
#     def copy(self):
#         return Rect((self.x, self.y), (self.w, self.h))
#
#     def __str__(self):
#         return 'Rect(({0},{1}),({2},{3}))'.format(self.left, self.right, self.top, self.bottom)
#
#
# class SpriteGroup(object):
#     def __init__(self):
#         self.current = 0
#         self.high = 0
#         self.items = []
#
#     def add(self, g):
#         self.high += 1
#         self.items.append(g)
#
#     def update(self, *args, **kwargs):
#         [s.update(*args, **kwargs) for s in self]
#
#     def draw(self, *args, **kwargs):
#         [s.draw(*args, **kwargs) for s in self]
#
#     def __len__(self):
#         return self.high
#
#     def __iter__(self):
#         return self
#
#     def next(self):
#         if self.current >= self.high:
#             self.current = 0
#             raise StopIteration
#         self.current += 1
#         return self.items[self.current - 1]
