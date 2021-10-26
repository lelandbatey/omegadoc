package domain

type OmegaAttribute struct {
	Key   string
	Value string
}

type OmegaDoc struct {
	// SourceFilePath is the path on-disk of the file which originally defined this OmegaDoc
	SourceFilePath string
	// DestFilePath is the path to where the OmegDoc should be stored in the
	// output FileTree. Even if this path is absolute (starts with a '/') it
	// will be treated as relative to the output directory when output.
	DestFilePath string
	Attributes   []OmegaAttribute
	// The contents of the OmegaDoc, found between the opening statement (which
	// defines the delimiting identifier) and the delimiting identifier.
	Contents string
	// The line within SourceFilePath on which the OmegaDoc starts
	StartLineNumber int
	// HTTPURL contains a single full HTTP URL where you can read the source of
	// this OmegaDoc in your web-browser. This URL is not present in the
	// original document and if present will have been derived from the git
	// repository which the document belongs within. If the file at
	// "SourceFilePath" is not within a git repository or there is not
	// configuration sufficient to derive the HTTP url for that file inside
	// that repository, then this HTTPUrl will be blank.
	HTTPUrl string
}

/*
#!/usr/bin/env omegadoc <<ENDDOC omegadoc/index.md
# Narrative Purpose of OmegaDoc

I want this tool to exist because I find that there's inherent tension in our
usual goals for documentation, and the existing solutions to that tension are
too limited to solve the problems myself and many other developers have begun
facing.

The tension I've found comes from two competing goals:

1. "I want all my documentation in one place so that I may browse it all together"
2. "I want all my documentation stored next to the implementation which relates
   to that documentation, so that keeping the documentation up to date is easy and
   the documentation is discoverable when I am browsing the implementation."

Within every organization I've ever been a part of, this leads to essentially
two piles of documentation. There's the pile of documentation which is stored
within the source code/implementation of the product itself. Then there's the
pile of documentation which is manually written and updated, usually kept in a
private wiki for the team.

I have seen other solutions to this problem, specifically in the form of
documentation which is generated _from_ the source code of the implementation.
Examples of this include Godoc in Golang, Sphinx in Python, Doxygen for C++ and
Javadoc for Java. These tools go through the source code and extract
specially-tagged pieces of documentation, then they transform that
documentation into a highly readable (usually HTML) form. The limitation of
each of these documentation generation pieces of software is that they are very
focused on their own language/ecosystem niche. It seems that most software
shops are now juggling at least two languages (some backend + JS frontend),
potentially many more. All these documentation generation tools are great, but
they can't bring together all the documentation I'm writing in 4+ languages in
30+ repositories.

I believe the solution to this, and the purpose for OmegaDoc, is buried in that
description of how documentation generation software works in general. Here's
the specific line:

> [Documentation generation] tools go through the source code and **extract**
> ... documentation, then they **transform** that documentation into a highly
> readable (usually HTML) form.

All these tools are first extracting (finding, parsing, etc) documentation and
then they're transforming (formatting, highlighting, rendering, etc) those docs
into a form for us humans to read. What's so interesting is that there are lots
of tools for doing the transforming part of that process in a more general way,
but there's almost nothing for doing the extraction and relocation of
documentation in a general way. By this I mean, Doxygen/Sphinx/GoDoc could be
implemented as specialty frontend for a static site-generator like
Jekyl/Hugo/Pelican, etc. I know, I know, there are very fancy features of some
of the doc generators which are very nice and aren't exactly document-oriented,
things like creating runnable examples in your browser, but in general I feel
like most of our documentation could be some markdown right next to, or even
better, _inside of_ the source code files themselves.

This idea asking "what if we could just have markdown files right next
to/inside of our source code, and no matter where they are they could all be
collected together?" is what prompted me to write OmegaDoc. OmegaDoc allows you
to write documentation/files _inside_ your source code (or next to it if you
like) then collect all that documentation together. Once collected together,
you can do whatever you want with it; browse it in a text editor if you like,
render it all using a static site generator, whatever. The core idea is one of
"collecting together" documentation from everywhere.

To do this, I've taken the classic idea of a [here document](https://en.wikipedia.org/wiki/Here_document)
and made it even more gruesome by saying "you can put a special here doc in any
stream of bytes (a file) and OmegaDoc will see it and extract that here doc
into another file." This is somewhat nasty in concept, but also very clean and
nice in practice. It means that I can use
[ghorg](https://github.com/gabrie30/ghorg) to clone down the hundreds of
repositories across the company I work for, then extract all the documentation
from every repo into one folder using the following command:

```
omegadoc --input-search-path /home/leland/root_folder_with_all_repos --output-path /home/leland/documentation/
```

And what I get out the other side is a perfect tree of all the documentation
from every repo all together. From there I can render it all into HTML (since
most all the documentation files are HTML) or browse them in my text editor. It
means when I want to read documentation, that documentation is collected
together cohesively, and when I want to write documentation, I can write it
directly next to (or inside of) the component/service described by that
documentation.

Of course, to do this I have to write my documentation in the form of
OmegaDocs, but all that means is adding a few lines to existing docs.

Also, I recommend using Pandoc and the following command to render everything
to HTML for your browsing pleasure:

	find ./<OUTPUT FOLDER HERE> -iname "*.md" -type f -exec sh -c 'pandoc "${0}" -o "${0%.md}.html"' {} \;

ENDDOC
*/
