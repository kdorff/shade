Transylvania SDK
================

A simple and easy to use 2.5D game SDK for the Go programming language.

TODOs 
-----

- [x] load sprites
- [x] load animated sprites
- [x] detect collisions
- [ ] load normal maps
- [ ] set ambient lighting
- [ ] add simple light source (radiates from a central point in all directions)
- [ ] add directional light source (radiates from a point in a specific direction)
- [ ] add depth based drawing rather than painters algorithm
- [ ] 03-collisions: create example
- [ ] 04-lighting: Add ability to set lights in SDK
- [ ] 04-lighting: Add ability to load normal maps into sprite
- [ ] 04-lighting: Update shaders to calculate lighting
- [ ] 05-text: Fix text alignment
- [ ] ex1-platform: Fix boundry detection (so player does not run up walls)
- [ ] ex1-platform: Fix input (sometimes jumps do not register and the guy runs faster to the right)
- [ ] ex1-platform: Add basic menu after splash screen
- [ ] ex1-platform: Add more interesting sprites
- [ ] ex1-platform: Add more interesting level
- [ ] ex1-platform: Add enemies
- [ ] ex1-platform: Allow player to shoot
- [ ] ex1-platform: Add door to open when all enemies have been killed (allowing the player to win)


NOTE: This SDK should be considered very experimental as it is still under development.  It is currently being modeled after some aspects of the PyGame SDK, but this will probably change some as it matures.  The project will not have its "experimental" status removed until at least all of the above items have been completed.

While the above should work without needing to work with the OpenGL SDK, the packages of this SDK should be extendable such that more advanced uses are possible.


Attribution
-----------

This project was inspired by the article ["Normal Mapping with Javascript and Canvas"](https://29a.ch/2010/3/24/normal-mapping-with-javascript-and-canvas-tag).

Some aspects of the SDK are inspired by the [PyGame SDK](http://www.pygame.org/).


Helpful Tools
-------------

[Pyxel Edit](http://pyxeledit.com/) - Very nice pixel art editor.

[Sprite DLight](https://www.kickstarter.com/projects/2dee/sprite-dlight-instant-normal-maps-for-2d-graphics) - Instant normal maps for 2D graphics

[Tiled](http://www.mapeditor.org/) - Your free, easy to use and flexible tile map editor.

[sfxr](http://www.drpetter.se/project_sfxr.html)/[cfxr](http://thirdcog.eu/apps/cfxr) - Simple means of getting basic sound effects.


Troubleshooing
--------------

#### Slow saves in VIM

Sometimes, go-imports inserts the the wrong things.  VIM was hanging for me for about 10~20 seconds after saving.  I'm not sure why, but switching from v2.1 to v4.1-core fixed the issue.

```
-       "github.com/go-gl/gl/v2.1/gl"
+       "github.com/go-gl/gl/v4.1-core/gl"
```

#### Error: "could not decode file player.png: image: unknown format"

The following line will not be auto-imported since our code does not call anything in the png, library directly.  In order to ensure that the library loads (so image.decode("somefile.png") works), we must have the following import to load the PNG package.

```
+	_ "image/png" // register PNG decode
```

#### ... has no field or method ...

And it just does not make sense why it is not found, check that go-imports imported the correct thing.
