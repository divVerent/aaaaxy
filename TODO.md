# TODO

## Level sections

Graph:

* 1-1 (not named)
  * Stop and Stair (infinite staircase)
    * Endless Eight

* (some hub room)
  * Short Circuited (two straight paths, one much shorter than the other, to same destination)
    * Switched Around (two straight paths, one above the other, reverse relative order)
      * M. C. Waterfall (S formed path upwards, then down)
        * Silver City ("port" of Nexuiz map, basically a square connected by jump pads)
          * Nine boxes room from Antichamber (boxes can only be looked into, have different content from different sides)
            * Hogwarts Express (track 9 3/4 puzzle)
              * -> Part 2: world of magic. Player learns to move balloons (i.e. carry). ADDS A BUTTON (pick up is auto, release with button).
                * Various sections to train ability. Switch/door puzzles!
                * Entrance to part 2+3 - requiring both abilities (idea: impossible platforming puzzle that works once a balloon is moved). Offers way to gain ability 3 first.
                * A "dev room" that requires all three abilities.
  * The Strip
    * The Torus | The Moebius Strip
      * The Klein Bottle
        * The Projective Plane
          * The Sphere
            * The House With Two Rooms
              * -> Part 3: Topology 101 passed. Player can now stand on balloons w/o falling down.
                * Various sections to train ability. Platforming!
                * Part that starts requiring another ability connects to previous area to provide it.
                * Entrance to part 3+4 - requiring both abilities (idea: misplaced balloons too, but can push them away). Offers way to gain ability 4 first.
                * A "dev room" that requires all three abilities.
  * Hilbert's Hotel
    * Choices, choices! (Endless sequence of two paths; have to turn back to proceed; actually more Antichamber reference)
      * Shepard Tone (player has to play melody on an infinite piano by jumping on the right keys) melody idea: E C E D to open (Loom), alterates: C D F D A A G never gonna give you up
        * Turtles all the way down (have to enter mouth of tutle, ends up in front of another turtle, but is actually a copied room).
          * Don't Look Up (if you do, you fall down ad infinitum, must escape)
            * More infinity stuff?
              * -> Part 4: Infinity understood. Player can now push away balloons from distance (same button).
                * Various sections to train ability. Sokoban!
                * Part that starts requiring another ability connects to previous area to provide it.
                  * Entrance to part 4+2 - requiring both abilities (idea: need to both move and push away). Offers way to gain ability 2 first.
                  * A "dev room" that requires all three abilities.
  * The End
    * -> Part 5: Series of puzzles that require all three abilities.

Abilities:
- Carry item
  - Press button while touching to grab, release to release.
  - No hold-to-grab
  - Animation: idle when carried, up when released.
- Stand on item
  - Before this, items are nonsolid
  - After this, items are solid IFF player is above (i.e. semisolids)
  - Animation: up (blink).
  - While standing, it moves up faster!
- Push item away
  - Hold button to push away
  - Animation: left/right (blink), up when released.

Layout:

1
-> 2
   -> new ability
   -> 2b
      -> middle of 3
      -> 2+3
         -> dev room that requires 4
         -> middle of 4
         -> end
-> 3
   -> new ability
   -> 3b
      -> middle of 4
      -> 3+4
         -> dev room that requires 2
         -> middle of 2
         -> end
-> 4
   -> new ability
   -> 4b
      -> middle of 2
      -> 4+2
         -> dev room that requires 3
         -> middle of 3
         -> end


* At end of game, print the time the game took from the frames counter in the player. Also, halt the frames counter then.

ADD:
- Entity that gives player an ability.
- Sound file for new ability.
- Icons for abilities (blinking or even Z-rotating sprites).
- Puzzles!
