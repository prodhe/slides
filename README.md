# Slides

Create a web browser based slide show from extremely simple markup. Have a new slide show ready in seconds!

![Slides](./slides.gif)

This is pretty much a go clone of *sent* from suckless.org (https://tools.suckless.org/sent/) for the browser instead of X11.

It also serves as a way for me to try out lexing/parsing in Go. The lexer is based on Rob Pike's talk "Lexical Scanning in Go" ([youtube](https://www.youtube.com/watch?v=HxaD_trXwRE)) and a lot of the methods are almost direct copies of the template/parse package.

## Install

```
go get -u github.com/prodhe/slides
```

## Syntax

- Paragraphs (multiple newlines: `\n\n`) separates slides.

- A dot (`.`) on a single line will create a linebreak within a paragraph.

- Lines beginning with `#` is a comment and will not show up in the final slide.

- Include images using `@`.

- Escape special characters at the beginning of a line using `\`. This can also be used to create an empty slide.

- Blanks (spaces and tabs) at the beginning of a line will be preserved. Tabs equals four spaces.

### Example

```
slides

Every paragraph is a new slide

- Blanks at the beginning of a
  line will be preserved
- Tab equals four spaces
	1 2 3 4

Special characters at the beginning of a line:
.
\. - print empty line within a paragraph
\# - source file comments
\@ - include images from local folder
\\ - escape the above or empty slide

# This produces an empty slide.
\

Include images using @name.png

@r2d2_1.png

<b>HTML tags are escaped</b>

UTF-8 support (emojis!)

😀😁😂😃😄😅😆😇😈😉😊😋😌😍😎😏
😐😑😒😓😔😕😖😗😘😙😚😛😜😝😞😟
😠😡😢😣😥😦😧😨😩😪😫😭😮😯😰😱
😲😳😴😵😶😷😸😹😺😻😼😽😾😿🙀☠

Simple and effective

Shamelessly based on 'sent'
tools.suckless.org/sent

Questions?

github.com/prodhe/slides
```
