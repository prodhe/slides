package main

const (
	HTML_TMPL = `
<!doctype html>
<html>
<head>
<title>Slides</title>
<style type="text/css">
{{style}}
</style>
</head>
<body>
<div id="slides" class="font-mono">
{{data}}
</div>
<script>
{{javascript}}
</script>
</body>
</html>
`
	STYLESHEET = `
* { border: 0; margin: 0; padding: 0; box-sizing: border-box; }
body {
	background-color: #ffffea;
	font-size: 24pt;
	line-height: 1.7;
}
#slides {
	display: flex;
	flex-flow: row nowrap;
}
section {
	display: none;
	flex: none;
	width: 100vw;
	height: 100vh;
	align-items: center;
	justify-content: center;
	flex-flow: column wrap;
	text-align: left;
	padding: 1em;
	user-select: none;
	cursor: default;
}
section > div {
	display: block;
	position: relative;
}
.current {
	display: flex;
}
img {
	display: block;
	position: relative;
	max-width: 100vw;
	max-height: 100vh;
}
.font-mono {
	font-family: monospace;
}
.font-sans {
	font-family: sans-serif;
}
.font-serif {
	font-family: serif;
}
`
	JAVASCRIPT = `
var slideIndex = 0;

initSlides();
showSlides(slideIndex);

function initSlides() {
	var i;
	var slides = document.getElementsByTagName("section");
	for (i = 0; i < slides.length; i++) {
			slides[i].onclick = function (){
				plusSlides(1);
			};
	}
}

function plusSlides(n) {
	showSlides(slideIndex += n);
}

function currentSlide(n) {
	showSlides(slideIndex = n);
}

function showSlides(n) {
	var i;
	var slides = document.getElementsByTagName("section");

	if (n > slides.length-1) { slideIndex = slides.length-1; }
	if (n < 0) { slideIndex = 0; }

	for (i = 0; i < slides.length; i++) {
		slides[i].classList.remove("current");
	}
	slides[slideIndex].classList.add("current");
}

function cycleFont() {
	cl = document.getElementById("slides").classList;
	if (cl.contains("font-mono")) {
		cl.remove("font-mono");
		cl.add("font-sans");
	} else if (cl.contains("font-sans")) {
		cl.remove("font-sans")
		cl.add("font-serif");
	} else if (cl.contains("font-serif")) {
		cl.remove("font-serif")
		cl.add("font-mono");
	}
}

window.onkeydown = function (e) {
	var e=window.event || e;
	// space, left, up, right, down: 32 37 38 39 40
	// f: 70
	switch (e.keyCode) {
		case 37:
		case 38:
			plusSlides(-1);
			break;
		case 32:
		case 39:
		case 40:
			plusSlides(1);
			break;
		case 70:
			cycleFont();
			break;
	}	
};
`
)
