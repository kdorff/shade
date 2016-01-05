Transylvania SDK
================

A simple and easy to use 2.5D game SDK for the Go programming language.

NOTE: This SDK should be considered very unstable as it is still under development.  It is currently being modeled after some aspects of the PyGame SDK, but this will probably change some as it matures.  Additionally, it will eventually have built in support for easily handeling sprites w/ normal maps and simple lighting (allowing to be extended for more advanced lighting).

Attribution
-----------

This project was inspired by the article ("Normal Mapping with Javascript and Canvas")[https://29a.ch/2010/3/24/normal-mapping-with-javascript-and-canvas-tag].

Some aspects of the SDK are inspired by the (PyGame SDK)[http://www.pygame.org/]


Helpful Tools
-------------

[Pyxel Edit](http://pyxeledit.com/) - Very nice pixel art editor.

[Sprite DLight](https://www.kickstarter.com/projects/2dee/sprite-dlight-instant-normal-maps-for-2d-graphics) - Instant normal maps for 2D graphics

[Tiled](http://www.mapeditor.org/) - Your free, easy to use and flexible tile map editor.

[sfxr](http://www.drpetter.se/project_sfxr.html)/[cfxr](http://thirdcog.eu/apps/cfxr) - Simple means of getting basic sound effects.


Troubleshooing
--------------

# Slow saves in VIM

Sometimes, go-imports inserts the the wrong things.  VIM was hanging for me for about 10~20 seconds after saving.  I'm not sure why, but switching from v2.1 to v4.1-core fixed the issue.

```
-       "github.com/go-gl/gl/v2.1/gl"
+       "github.com/go-gl/gl/v4.1-core/gl"
```

# could not decode file player.png: image: unknown format

The following line will not be auto-imported since our code does not call anything in the png, library directly.  In order to ensure that the library loads (so image.decode("somefile.png") works), we must have the following import to load the PNG package.

```
+	_ "image/png" // register PNG decode
```
