# OmegaDoc

OmegaDoc provides one solution to the documentation problems even medium-size
organizations face. We'd like to keep documentation located nearby the things
their documenting, but doing that means actually finding and reading that
documentation requires going to where it's located across potentially many
codebases. OmegaDoc is meant to solve this by bringing together documentation
from anywhere, text files of any type, into a single collected directory.

An "OmegaDoc" is a series of bytes which can be recognized by the OmegaDoc
program and extracted into a separate file. In concept it behaves like a
specialized "here document"; an OmegaDoc is meant to be defined in-band with
code and configuration, so that it's nearby the things it documents.

Specification
-------------

An OmegaDoc is composed of an opening statement, then an output path, then a
body, and then a delimiting identifier, in that order. The opening statement is
a "magic string" followed by a delimiting identifier (e.g. the word EOF or
END), then some amount of whitespace, then an output path. The output path is
followed, starting on the next line, by the text to be quoted, and then closed
by the same delimiting identifier on its own line. The "magic string" which
marks the beginning of an OmegaDoc is the string:

	#!/usr/bin/` + `env omegadoc <<

An example of an OmegaDoc then is like so:

	#!/usr/bin/` + `env omegadoc <<DELIMIDENT exampleoutput/readme.md
	Hello I am a markdown document which will be recorded to
	a file at the relative path exampleoutput/readme.md.
	DELIMIDENT

Additionally, if the file ends before the delimiting identifier is reached,
that is considered to be the end of the omegadoc.
