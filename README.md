# OmegaDoc

<!--
All the OmegaDocs defined in this file should be ignored since they're mostly
nonesense. Thus this ignore directive at the top here.
#!/usr/bin/env omegadoc ignore-this-file
-->
OmegaDoc provides one solution to the documentation problems even medium-size
organizations face. We'd like to keep documentation located nearby the things
they're documenting, but doing that means actually finding and reading that
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
END), then some amount of whitespace (without newlines), then attributes (if any), then
an output path.  The output path is followed, starting on the next line, by the
text to be quoted, and then closed by the same delimiting identifier on its own
line. The "magic string" which marks the beginning of an OmegaDoc is the
string:

	#!/usr/bin/env omegadoc <<

Attributes are strings of key-value pairs, separated by a colon and delimited
by spaces. As an example, the string `zip:pow` is an attribute, with `zip`
being the 'key' and `pow` being the 'value'. Specifying two attributes might be
done with a string like `foo:bar fizz:buzz`. An OmegaDoc may have zero
attributes. There is no limit to the number of attributes an OmegaDoc may have.

An example of an OmegaDoc with zero attributes then is like so:

	#!/usr/bin/env omegadoc <<DELIMIDENT exampleoutput/readme.md
	Hello I am a markdown document which will be recorded to
	a file at the relative path exampleoutput/readme.md.
	DELIMIDENT

Additionally, if the file ends before the delimiting identifier is reached,
that is considered to be the end of the omegadoc.

If a file contains an "ignore" directive in its bytes before an OmegaDoc
opening statement, then that file will be considered to have NO OmegaDocs in
it, even if it otherwise contains one or more valid OmegaDocs. If a file
contains an "ignore" directive after one or more valid OmegaDoc directives,
then that "ignore" directive will itself be ignored.

An "ignore" directive is the following string, present anywhere in the bytes of
a file:

    #!/usr/bin/env omegadoc ignore-this-file

Pieces
------
```
			  Delimiting       Section
        Magic string      identifier      attribute            Output path
┌─────────────┴──────────┐┌───┴────┐ ┌────────┴────────┐ ┌──────────┴──────────┐
#!/usr/bin/env omegadoc <<DELIMIDENT section:foo-section exampleoutput/readme.md
Hello I am a markdown document which will be recorded to
a file at the relative path exampleoutput/readme.md. Additionally, if multiple
OmegaDocs were to have the same output path but different sections, then a
single file will be written to the output path, but its contents will be
composed of the contents of each OmegaDoc with that output path concatenated
together in the order of their "section" labels sorted by lexicographically.
DELIMIDENT
```
