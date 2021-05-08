Sections to do:

- Grab & Push
  - BRLOGENSHFEGLE
    - Need to transport platform while NOT transforming the floor into shiny yet deadly coins.
    - I.e. time to hold button is limited.

- Grab & Stand
  - Music: mythica.mp3
  - The Butterfly Effect (from Antichamber)
    - Lots of doors that are LogicalOr of three inputs (LeftSw, DoorBlocker, RightSw)
    - At the end you get stuck UNLESS you have a platforms left
    - Solution is to go to end, get the one platform that is there, walk back, block doors and collect all platforms you can find.

- Push & Stand
  - Music: 2012_november_fakeAwake04 back to A minor.wav
  - A platform race! Steer left/right, overtake other "drivers"...
    - need "evil" platform object: solid, moves, thereby kicking us off our platform

- End of each: a fourth ability (surge protector) that lets you go through electric shock things.
  - Without the ability, going into one throws you back to the side you came from. Also, just to make sure, the beams are solid then.
  - One such place at start of game, right behind the thing one can open with ANY of three abilities.
  - Also, one after each part.
  - Maybe hide another place that has electric beams?
  - After that, dump player right in top of the third ability (the missing one).
  - Means there are 3 copies of the obtain-fourth-ability room. Yeah, could also warpzone-hack around that, but why.

- Finale
  - Music: Juhani Junkala [Retro Game Music Pack] Level 3.wav
  - There has to be a "filigran teleport pattern" room. Impl wise, when standing right in the center, the screen is guaranteed to not contain anything else, so we "teleport" by enabling warpzones offscreen that connect the room to another room.
  - Maybe add a small section with gravity flipping instead of jumping? Will have to change some isAbove checks but shouldn't be too hard.
    - Should gravity flipping also flip the direction platforms go? Probably yes.
    - Or should platforms have their own gravity direction, and only flip if player flips while carrying them?

- Credits
  - Music: need to find some
