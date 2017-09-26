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
<div id="slides">
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
body { background-color: #ffffea; }
#slides {
	display: flex;
	flex-flow: row nowrap;
}
p {
	flex: none;
	width: 100vw;
	height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 24pt;
	font-family: monospace;
	line-height: 1.7;
	text-align: left;
	padding: 1em;
}

.slide {
	display: none;
}
`
	JAVASCRIPT = `
var slideIndex = 0;
showSlides(slideIndex);

function plusSlides(n) {
  showSlides(slideIndex += n);
}

function currentSlide(n) {
  showSlides(slideIndex = n);
}

function showSlides(n) {
  var i;
  var slides = document.getElementsByClassName("slide");
//  var dots = document.getElementsByClassName("dot");
  if (n > slides.length-1) {slideIndex = 0} 
  if (n < 0) {slideIndex = slides.length-1}
  for (i = 0; i < slides.length; i++) {
      slides[i].style.display = "none"; 
      slides[i].onclick = function (){
        plusSlides(1);
      };
  }
//  for (i = 0; i < dots.length; i++) {
//      dots[i].className = dots[i].className.replace(" active", "");
//  }
  console.log(slideIndex);
  slides[slideIndex].style.display = "flex"; 
//  dots[slideIndex].className += " active";
}
`
)