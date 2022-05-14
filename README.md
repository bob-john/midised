# midised

Stream editor for midicat.

# record

midicat in -i=1 | midised timestamp > recording.txt

# play

midised play recording.txt | midicat out -i=1
