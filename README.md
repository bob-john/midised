# midised

Stream editor for midicat.

## commands

### record

midicat in -i=1 | midised timestamp > recording.txt

### play

midised play recording.txt | midicat out -i=1

## todo

- 'play' should send note off on ctrl+c.
