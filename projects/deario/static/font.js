(function(){
  window.applyUserFont = function(familySpec) {
    if (!familySpec) return;
    familySpec = familySpec.replace(/\s+/g, '+');
    localStorage.setItem('uiFontFamily', familySpec);
    var href = 'https://fonts.googleapis.com/css2?family=' + familySpec + '&display=swap';
    var link = document.getElementById('user-font');
    if (link) {
      link.href = href;
    } else {
      link = document.createElement('link');
      link.id = 'user-font';
      link.rel = 'stylesheet';
      link.href = href;
      document.head.appendChild(link);
    }
    var familyName = familySpec.split(':')[0].replace(/\+/g, ' ');
    document.documentElement.style.setProperty('--font-family', familyName + ', system-ui, sans-serif');
  };

  window.addEventListener('DOMContentLoaded', function() {
    var select = document.getElementById('font-select');
    if (!select) return;
    var fam = localStorage.getItem('uiFontFamily') || 'Gamja+Flower';
    select.value = fam.replace(/\s+/g, '+');
  });
})();
