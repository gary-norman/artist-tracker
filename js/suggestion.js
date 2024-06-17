let latestRequestTimestamp = 0;

async function fetchSuggestions(query) {
    const response = await fetch(`/suggest?query=${encodeURIComponent(query)}`);
    if (response.ok) {
        return await response.json();
    }
    throw new Error('Failed to fetch suggestions');
}

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
    console.log("received suggestionsData:=", suggestionsData);

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

        // Create a result container for grid display
        const resultContainer = document.createElement('div');
        resultContainer.className = 'result-container';
        resultContainer.style.display = 'grid';
        resultContainer.style.gridTemplateColumns = 'repeat(3, minmax(0, 1fr))'; // Single column layout
        resultContainer.style.gap = '1.2rem'; // Gap between grid items

        categoryContainer.appendChild(resultContainer);
        categoryContainer.resultContainer = resultContainer; // Store reference for later use
        categoryContainer.resultsCount = 0; // Initialize a counter for each category
    }

    // Check if suggestionsData is null or undefined or empty
    if (!suggestionsData || suggestionsData.length === 0) {
        const noResultsMessage = document.createElement('div');
        noResultsMessage.textContent = 'No results found.';
        searchResults.appendChild(noResultsMessage);
        return; // Exit function early
    } else{
        // Update the count of search results in the header
    const resultsHeader = document.querySelector('.filters .small.light.center');
    if (resultsHeader) {
        resultsHeader.textContent = `Showing ${suggestionsData.length} results`;
    }
    }

    suggestionsData.forEach(function (suggestion) {
        if (suggestion.category === 'Concert' && suggestion.matchitem && suggestion.matchitem.dates) {
            // Create elements for each suggestion
            var content = document.createElement('div');
            content.className = 'content';

            var dateDiv = document.createElement('div');
            dateDiv.className = 'pic date';

            // Create divs for month and year
            var monthYearDiv = document.createElement('div');
            monthYearDiv.className = 'month-year';

            var monthDiv = document.createElement('div');
            monthDiv.textContent = suggestion.matchitem.dates.Month;

            var yearDiv = document.createElement('div');
            yearDiv.textContent = suggestion.matchitem.dates.Year;

            monthYearDiv.appendChild(monthDiv);
            monthYearDiv.appendChild(yearDiv);

            var dayDiv = document.createElement('div');
            dayDiv.textContent = suggestion.matchitem.dates.Day;

            dateDiv.appendChild(dayDiv);
            dateDiv.appendChild(monthYearDiv);

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
                categories[suggestion.category].resultContainer.appendChild(content);
                categories[suggestion.category].resultsCount++; // Increment the counter
            } else {
                console.warn(`Category ${suggestion.category} not found in categories object.`);
            }
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
            } else if (suggestion.category === 'Member' && suggestion.matchitem) {
                if (suggestion.artist.Members !== "") {
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
                categories[suggestion.category].resultContainer.appendChild(content);
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
