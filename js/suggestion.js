document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');

    if (searchInput && searchResults) {
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

        searchInput.addEventListener('input', function (e) {
            const query = e.target.value.trim();
            if (query) {
                debouncedFetchSuggestions(query);
            } else {
                searchResults.innerHTML = ''; // Clear suggestions if input is empty
            }
        });
    }

    function showSuggestions(suggestionsData) {
        searchResults.innerHTML = '';
        console.log("received suggestionsData:=", suggestionsData);

        const categories = {
            Artist: document.createElement('div'),
            Member: document.createElement('div'),
            Album: document.createElement('div'),
            Concert: document.createElement('div')
        };

        for (const category in categories) {
            const categoryContainer = categories[category];
            categoryContainer.className = 'col col1';

            const resultContainer = document.createElement('div');
            resultContainer.className = 'result-container';
            resultContainer.style.display = 'grid';
            resultContainer.style.gridTemplateColumns = 'repeat(3, minmax(0, 1fr))';
            resultContainer.style.gap = '1.2rem';

            categoryContainer.appendChild(resultContainer);
            categoryContainer.resultContainer = resultContainer;
            categoryContainer.resultsCount = 0;
        }

        if (!suggestionsData || suggestionsData.length === 0) {
            const noResultsMessage = document.createElement('div');
            noResultsMessage.textContent = 'No results found.';
            searchResults.appendChild(noResultsMessage);
            return;
        } else {
            const resultsHeader = document.querySelector('.filters .small.light.center');
            if (resultsHeader) {
                resultsHeader.textContent = `Showing ${suggestionsData.length} results`;
            }
        }

        suggestionsData.forEach(function (suggestion) {
            let artistName;
            let contentText;
            let content;
            const a = document.createElement('a');
            a.setAttribute('href', `/artist?name=${encodeURIComponent(suggestion.artist.name)}`);
            a.dataset.artistName = suggestion.artist.name;

            if (suggestion.category === 'Concert' && suggestion.matchitem && suggestion.matchitem.dates) {
                content = document.createElement('div');
                content.className = 'content';

                const dateDiv = document.createElement('div');
                dateDiv.className = 'pic date';

                const monthYearDiv = document.createElement('div');
                monthYearDiv.className = 'month-year';

                const monthDiv = document.createElement('div');
                monthDiv.textContent = suggestion.matchitem.dates.Month;

                const yearDiv = document.createElement('div');
                yearDiv.textContent = suggestion.matchitem.dates.Year;

                monthYearDiv.appendChild(monthDiv);
                monthYearDiv.appendChild(yearDiv);

                const dayDiv = document.createElement('div');
                dayDiv.textContent = suggestion.matchitem.dates.Day;

                dateDiv.appendChild(dayDiv);
                dateDiv.appendChild(monthYearDiv);

                contentText = document.createElement('div');
                contentText.className = 'content-text go-down-home';

                artistName = document.createElement('div');
                artistName.className = 'p--bold cut concert';
                artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist';

                const locationName = document.createElement('div');
                locationName.className = 'small light cut concert';
                locationName.textContent = suggestion.matchitem.location || 'Unknown Location';

                contentText.appendChild(artistName);
                contentText.appendChild(locationName);
                content.appendChild(dateDiv);
                content.appendChild(contentText);
                a.appendChild(content);

                if (categories[suggestion.category]) {
                    categories[suggestion.category].resultContainer.appendChild(a);
                    categories[suggestion.category].resultsCount++;
                } else {
                    console.warn(`Category ${suggestion.category} not found in categories object.`);
                }
            } else {
                content = document.createElement('div');
                content.className = 'content';

                const img = document.createElement('img');
                if (suggestion.category === 'Album') {
                    img.className = 'pic album';
                    img.src = suggestion.artist.strAlbumThumb ? suggestion.artist.strAlbumThumb : 'default-album-image-url.jpg';
                    img.alt = 'Album cover of ' + (suggestion.matchitem || 'Unknown Album');
                } else {
                    img.className = 'pic user';
                    img.src = suggestion.artist && suggestion.artist.strArtistThumb ? suggestion.artist.strArtistThumb : 'default-image-url.jpg';
                    img.alt = 'Profile image of ' + (suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist');
                }

                contentText = document.createElement('div');
                contentText.className = 'content-text go-down';

                let boldCut = document.createElement('div');
                boldCut.className = 'p--bold cut';
                boldCut.textContent = suggestion.matchitem || '';



                if (suggestion.category === 'Artist' && suggestion.matchitem && suggestion.matchitem.toLowerCase() === suggestion.artist.name.toLowerCase()) {
                    contentText.appendChild(boldCut);
                } else if (suggestion.category === 'Member' && suggestion.matchitem && suggestion.artist.Members !== "" && suggestion.matchitem.toLowerCase() === suggestion.artist.name.toLowerCase()) {
                            contentText.appendChild(boldCut);
                } else {
                    artistName = document.createElement('div');
                    artistName.className = 'p--normal';
                    artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist';
                    contentText.appendChild(boldCut);
                    contentText.appendChild(artistName);
                }

                content.appendChild(img);
                content.appendChild(contentText);
                a.appendChild(content);

                if (categories[suggestion.category]) {
                    categories[suggestion.category].resultContainer.appendChild(a);
                    categories[suggestion.category].resultsCount++;
                } else {
                    console.warn(`Category ${suggestion.category} not found in categories object.`);
                }
            }
        });

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
});
