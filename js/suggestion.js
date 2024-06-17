let latestRequestTimestamp = 0;

async function fetchSuggestions(query) {
    const response = await fetch(`/suggest?query=${encodeURIComponent(query)}`);
    if (response.ok) {
        return await response.json();
    }
    throw new Error('Failed to fetch suggestions');
}


// delay the execution of the fetch request, 
// ensuring that rapid inputs do not result in multiple network requests being sent in quick succession.
function debounce(fn, delay) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => fn.apply(this, args), delay);
    };
}

const debouncedFetchSuggestions = debounce(async function (query) {
    const currentRequestTimestamp = Date.now();
    latestRequestTimestamp = currentRequestTimestamp;

    try {
        const suggestionsData = await fetchSuggestions(query);
        if (latestRequestTimestamp === currentRequestTimestamp) {
            showSuggestions(suggestionsData);
        }
    } catch (error) {
        console.error('Error fetching suggestions:', error);
    }
}, 300);

var searchInput = document.getElementById('search-input');
var searchResults = document.getElementById('search-results');

searchInput.addEventListener('input', function (e) {
    const query = e.target.value.trim();
    if (query) {
        debouncedFetchSuggestions(query);
    } else {
        searchResults.innerHTML = ''; // Clear suggestions if input is empty
    }
});

function showSuggestions(suggestionsData) {
    // Clear previous search results
    searchResults.innerHTML = '';

    // Log all suggestionsData details to console for debugging
    console.log("recieved suggestionsData:=",suggestionsData);

    // Create containers for each category
    const categories = {
        Artist: document.createElement('div'),
        Member: document.createElement('div'),
        Album: document.createElement('div'),
        Concert: document.createElement('div')
    };

    // Initialize category containers with their headers
    for (const category in categories) {
        const categoryContainer = categories[category];
        categoryContainer.className = 'col col1';
        categoryContainer.resultsCount = 0; // Initialize a counter for each category
    }

    // Iterate through each suggestion in suggestionsData
    suggestionsData.forEach(function (suggestion) {
        if (suggestion.category === 'Concert' && suggestion.matchitem && suggestion.matchitem.dates) {
            // Iterate over each date for the concert location
            suggestion.matchitem.dates.forEach(function (date) {
                // Create elements for each suggestion
                var content = document.createElement('div');
                content.className = 'content';

                var dateDiv = document.createElement('div');
                dateDiv.className = 'pic date';
                dateDiv.textContent = date; // Display the date

                var contentText = document.createElement('div');
                contentText.className = 'content-text go-down-home';

                var artistName = document.createElement('div');
                artistName.className = 'p--bold cut concert';
                artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist'; // Handle undefined artist name

                var locationName = document.createElement('div');
                locationName.className = 'small light cut concert';
                locationName.textContent = suggestion.matchitem.location || 'Unknown Location';

                contentText.appendChild(artistName);
                contentText.appendChild(locationName);
                content.appendChild(dateDiv);
                content.appendChild(contentText);

                // Append content to the corresponding category container if it exists
                if (categories[suggestion.category]) {
                    categories[suggestion.category].appendChild(content);
                    categories[suggestion.category].resultsCount++; // Increment the counter
                } else {
                    console.warn(`Category ${suggestion.category} not found in categories object.`);
                }
            });
        } else {
            // Handle other categories
            var content = document.createElement('div');
            content.className = 'content';

            var img = document.createElement('img');
            img.className = 'pic user';
            img.src = suggestion.artist && suggestion.artist.strArtistThumb ? suggestion.artist.strArtistThumb : 'default-image-url.jpg'; // Handle undefined image
            img.alt = 'Profile image of ' + (suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist'); // Handle undefined name

            var contentText = document.createElement('div');
            contentText.className = 'content-text go-down';

            var boldCut = document.createElement('div');
            boldCut.className = 'p--bold cut';
            boldCut.textContent = suggestion.matchitem || ''; // Handle undefined match item

            // Display artist name only once if it matches exactly
            if (suggestion.category === 'Artist' && suggestion.matchitem && suggestion.matchitem.toLowerCase() === suggestion.artist.name.toLowerCase()) {
                contentText.appendChild(boldCut);
                // debug print
               /*  console.log("artist matchitem:", suggestion.matchitem);
                console.log("artist name:", suggestion.artist.name); */
            } else if (suggestion.category === 'Member' && suggestion.matchitem) {
                // Check if suggestion.artist.Members exists and is an array
                if (suggestion.artist.Members !== "") {
                     // debug print
               /*  console.log("artist matchitem:", suggestion.matchitem);
                console.log("artist name:", suggestion.artist.name); */
                    if (suggestion.matchitem.toLowerCase() === suggestion.artist.name.toLowerCase()) {
                        contentText.appendChild(boldCut);
                    } else {
                        var artistName = document.createElement('div');
                        artistName.className = 'p--bold cut'; // Use bold cut class for consistent styling
                        artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist'; // Handle undefined artist name
                        contentText.appendChild(boldCut);
                        contentText.appendChild(artistName);
                    }
                } else {
                    var artistName = document.createElement('div');
                    artistName.className = 'p--normal';
                    artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist'; // Handle undefined artist name
                    contentText.appendChild(boldCut);
                    contentText.appendChild(artistName);
                }
            } else {
                var artistName = document.createElement('div');
                artistName.className = 'p--normal';
                artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist'; // Handle undefined artist name
                contentText.appendChild(boldCut);
                contentText.appendChild(artistName);
            }

            content.appendChild(img);
            content.appendChild(contentText);

            // Append content to the corresponding category container if it exists
            if (categories[suggestion.category]) {
                categories[suggestion.category].appendChild(content);
                categories[suggestion.category].resultsCount++; // Increment the counter
            } else {
                console.warn(`Category ${suggestion.category} not found in categories object.`);
            }
        }
    });

    // Append all category containers with results to the searchResults container
    for (const category in categories) {
        const categoryContainer = categories[category];
        if (categoryContainer.resultsCount > 0) {
            const header = document.createElement('h2');
            header.textContent = category + 's';
            categoryContainer.insertBefore(header, categoryContainer.firstChild);
            searchResults.appendChild(categoryContainer);
        }
    }
}
