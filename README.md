# midiseq

Stream editor for midicat.

## commands

### record

midicat in -i=1 | midiseq timestamp > recording.txt

### play

midiseq play recording.txt | midicat out -i=1

## todo

- 'play' should send note off on ctrl+c.
