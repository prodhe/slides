# Slides

Creates a web browser based slide show from extremely simple markup. Create slides in seconds!

This is pretty much a clone of *sent* from suckless.org (https://tools.suckless.org/sent/) for the browser instead of X11 and a way to try out lexing/parsing in Go. The lexer is based on Rob Pike's talk "Lexical Scanning in Go" ([youtube](https://www.youtube.com/watch?v=HxaD_trXwRE)) and a lot of the routines are actually direct copies of the template/parse package.

## Syntax

Paragraphs (multiple newlines) separates slides.

Lines beginning with `#` is a comment.

See example.txt.

### On-going

Include images using `@`.

Include external text file using `.`.
