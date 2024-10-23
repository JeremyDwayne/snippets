(() => {
  // ui/static/js/custom.js
  var navLinks = document.querySelectorAll("nav a");
  for (i = 0; i < navLinks.length; i++) {
    link = navLinks[i];
    if (link.getAttribute("href") == window.location.pathname) {
      link.classList.add("live");
      break;
    }
  }
  var link;
  var i;
})();
