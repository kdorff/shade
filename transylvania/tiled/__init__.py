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

"""Implementation of Tiled's Map format.

More information about TMX files can be found here:
https://github.com/bjorn/tiled/wiki/TMX-Map-Format
"""

import os
import xml.etree.ElementTree as ET


class TileSet(object):
    def __init__(self, name=None, width=None, height=None, spacing=None,
                 margin=None, offset=None, properties=None, image=None,
                 terraintypes=None, tile=None, **kwargs):
        """
        name: The name of this tileset.
        tilewidth: The (maximum) width of the tiles in this tileset.
        tileheight: The (maximum) height of the tiles in this tileset.
        spacing: The spacing in pixels between the tiles in this tileset
                 (applies to the tileset image).
        margin: The margin around the tiles in this tileset (applies to the
                tileset image).
        tileoffset: (optional) (since 0.8.0)
        properties: (optional)(since 0.8.0)
        image: (optional)
        terraintypes: (optional)(since 0.9.0)
        tile: (optional)
        """
        self.name = name
        self.width = width
        self.height = height
        self.spacing = spacing
        self.margin = margin
        self.offset = offset
        self.properties = properties
        self.image = image
        self.terraintypes = terraintypes
        self.tile = tile

        for k, v in kwargs.iteritems():
            setattr(self, k, v)
    def __repr__(self):
        s = ''
        s += 'Name: {}\n'.format(self.name)
        s += 'Width: {}\n'.format(self.width)
        s += 'Height: {}\n'.format(self.height)
        s += 'Spacing: {}\n'.format(self.spacing)
        s += 'Margin: {}\n'.format(self.margin)
        s += 'Offset: {}\n'.format(self.offset)
        s += 'Properties: {}\n'.format(self.properties)
        s += 'Image: {}\n'.format(self.image)
        s += 'Terrain Types: {}\n'.format(self.terraintypes)
        s += 'Tile: {}\n'.format(self.tile)

        return s

class TileMap (object):
    def __init__(self, version=None, orientation=None, width=None, height=None,
                 tilewidth=None, tileheight=None, backgroundcolor=None,
                 renderorder=None,**kwargs):
        """
        version: The TMX format version, generally 1.0.
        orientation: Map orientation. Tiled supports "orthogonal", "isometric"
                     and "staggered" (since 0.9) at the moment.
        width: The map width in tiles.
        height: The map height in tiles.
        tilewidth: The width of a tile.
        tileheight: The height of a tile.
        backgroundcolor: The background color of the map. (since 0.9, optional)
        renderorder: The order in which tiles on tile layers are rendered.
                     Valid values are right-down (the default), right-up,
                     left-down and left-up. In all cases, the map is drawn
                     row-by-row. (since 0.10, but only supported for orthogonal
                     maps at the moment)
        kwargs: Unexpected map attributes.
        """

        if renderorder is None:
            renderorder = 'right-down'

        self.version = version
        self.orientation = orientation
        self.width = width
        self.height = height
        self.tilewidth = tilewidth
        self.tileheight = tileheight
        self.backgroundcolor = backgroundcolor
        self.renderorder = renderorder
        self.tilesets = {}

        for k, v in kwargs.iteritems():
            setattr(self, k, v)

    def __repr__(self):
        s = ''
        s += 'Version: {}\n'.format(self.version)
        s += 'Orientation: {}\n'.format(self.orientation)
        s += 'Width: {}\n'.format(self.width)
        s += 'Height: {}\n'.format(self.height)
        s += 'Tile Width: {}\n'.format(self.tilewidth)
        s += 'Tile Height: {}\n'.format(self.tileheight)
        s += 'Background Color: {}\n'.format(self.backgroundcolor)
        s += 'Render Order: {}\n'.format(self.renderorder)
        for i, ts in self.tilesets.iteritems():
            s += '### TileSet({}) ###\n'.format(i)
            s += str(ts)
            s += '### END TileSet ###\n'
        return s

    def add_tileset(self, firstgid, ts):
        if firstgid in self.tilesets:
            return
        self.tilesets[firstgid] = ts


def read_tsx(source):
    tsx = ET.parse(source)
    return tsx.getroot()


def read_tileset(game_dir, ts):
    if ts.get('source'):
        # if set, tileset is stored in an external TSX (Tile Set XML) file
        ts = read_tsx(os.path.join(game_dir, ts.get('source')))
    return TileSet(**dict(ts.attrib))


def read_tile_map(game_dir, filename):
    tree = ET.parse(os.path.join(game_dir, filename))
    root = tree.getroot()
    m = TileMap(**dict(root.attrib))

    for ts in root.iter('tileset'):
        i = ts.get('firstgid')
        m.add_tileset(i, read_tileset(game_dir, ts))

    return m
