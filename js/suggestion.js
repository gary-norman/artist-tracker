var searchInput = document.getElementById('search-input');
var suggestions = document.getElementById('suggestions');

searchInput.addEventListener('input', function() {
  console.log(searchInput.value)
  var searchQuery = searchInput.value.toLowerCase();
  var url = '/suggest?query=' + encodeURIComponent(searchQuery);

  var xhr = new XMLHttpRequest();
  xhr.open('GET', url, true);
  xhr.onload = function() {
    if (xhr.status === 200) {
      console.log(xhr.responseText);
      var suggestionsData = JSON.parse(xhr.responseText);
      showSuggestions(suggestionsData);
    }
  };
  
  xhr.send();
});

function showSuggestions(suggestionsData) {
  suggestions.innerHTML = '';
  if (suggestionsData.length === 0) {
    suggestions.style.display = 'none';
    return;
  }

  for (var i = 0; i < suggestionsData.length; i++) {
    var suggestion = document.createElement('div');
    suggestion.className = 'suggestion';
    suggestion.textContent = suggestionsData[i];
    suggestion.addEventListener('click', function() {
      searchInput.value = this.textContent;
      suggestions.style.display = 'none';
    });
    suggestions.appendChild(suggestion);
  }

  suggestions.style.display = 'block';
}

document.addEventListener('click', function(e) {
  if (!suggestions.contains(e.target) && e.target !== searchInput) {
    suggestions.style.display = 'none';
  }
});
