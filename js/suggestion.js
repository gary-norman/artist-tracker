document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');
    const populateResults = document.getElementById('populate-results');
    
    const locSearchInput = document.getElementById('button-filter-concert-location');
    const locSearchResults = document.getElementById('loc-search-result')
    const locationsContainer = document.getElementById('filter-checkbox-locations')
    
    if (locSearchInput && locSearchResults) {
        let latestRequestTimestamp = 0;

        async function fetchLocationSuggestions(query) {
            const response = await fetch(`/locationSuggest?query=${encodeURIComponent(query)}`);
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

        const debouncedFetchLocationSuggestions = debounce(async function (query) {
            const currentRequestTimestamp = Date.now();
            latestRequestTimestamp = currentRequestTimestamp;

            try {
                const suggestionsData = await fetchLocationSuggestions(query);
                if (latestRequestTimestamp === currentRequestTimestamp) {
                    showLocationSuggestions(suggestionsData);
                }
            } catch (error) {
                console.error('Error fetching suggestions:', error);
            }
        }, 300);

        locSearchInput.addEventListener('input', function (e) {
            const query = e.target.value.trim();
            if (query) {
                debouncedFetchLocationSuggestions(query);
            } else {
                locationsContainer.innerHTML = ''; 
            }
        });
    }

    function showLocationSuggestions(locSuggestionsData) {
        
        // console.log("received slocSuggestionsData:=", locSuggestionsData);
            
        locationsContainer.innerHTML = ''; 

        if (!locSuggestionsData || locSuggestionsData.length === 0) {
            locationsContainer.innerHTML = '<p style="text-align:center">No results found.</p>';
            return;
        }

        locSuggestionsData.forEach(function (location) {
            const checkboxContainer = document.createElement('div');
            checkboxContainer.className = 'checkbox go-across-sm';
            
            const input = document.createElement('input');
            input.className = 'checkbox checkbox-loc';
            input.id = `loc-${location.replace(/[\s, ]+/g, '-').toLowerCase()}`;
            input.name = 'loc';
            input.type = 'checkbox';
            input.value = location;

            const label = document.createElement('label');
            label.className = 'checkbox small';
            label.htmlFor = input.id;
            label.textContent = location;

            checkboxContainer.appendChild(input);
            checkboxContainer.appendChild(label);
            locationsContainer.appendChild(checkboxContainer);
        });
        locSearchResults.appendChild(locationsContainer);
    }
    
    if (searchInput && populateResults) {
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
                searchResults.classList.remove('hide');
                debouncedFetchSuggestions(query);
            } else {
                populateResults.innerHTML = ''; // Clear suggestions if input is empty
                searchResults.classList.add('hide');
            }
        });
    }

    function showSuggestions(suggestionsData) {
        populateResults.innerHTML = '';
        // console.log("received suggestionsData:=", suggestionsData);

        const resultsHeader = document.querySelector('.filters .small.light.center');
        if (!suggestionsData || suggestionsData.length === 0) {
            resultsHeader.textContent = `Showing 0 results`;
            populateResults.innerHTML = 'No results found.';
            searchResults.appendChild(populateResults);
            searchResults.classList.remove('hide'); 
            return;
        } else if (resultsHeader) {
            resultsHeader.textContent = `Showing ${suggestionsData.length} results`;
        }
        
        const categories = {
            Artist: document.createElement('div'),
            Member: document.createElement('div'),
            Album: document.createElement('div'),
            Concert: document.createElement('div')
        };

        for (const category in categories) {
            const categoryContainer = categories[category];
            categoryContainer.className = 'col col1';

            const outerContainer = document.createElement('div');
            outerContainer.id = 'search-results-v2';
            outerContainer.className = 'container';

            const scrollContainer = document.createElement('div');
            scrollContainer.className = 'container scroll';
            scrollContainer.id = 'populate-results'
            outerContainer.appendChild(scrollContainer);

            const resultContainerTwo = document.createElement('div');
            resultContainerTwo.className = 'col col2';

            const resultContainerThree = document.createElement('div');
            resultContainerThree.className = 'col col3';

            if (category === 'Artist' || category === 'Member' || category === 'Concert') {
                categoryContainer.appendChild(resultContainerThree);
                categoryContainer.resultContainer = resultContainerThree;
            } else {
                categoryContainer.appendChild(resultContainerTwo);
                categoryContainer.resultContainer = resultContainerTwo;
            }
            categoryContainer.resultsCount = 0;
            searchResults.classList.add('hide');
        }

        suggestionsData.forEach(function (suggestion) {
            let artistName;
            let contentText;
            let content;
            const a = document.createElement('a');
            
            // Add the concert part for the Concert category
            if (suggestion.category === 'Concert') {
                a.setAttribute('href', `/artist?name=${encodeURIComponent(suggestion.artist.name)}#artist-concerts`);
            } else {
                a.setAttribute('href', `/artist?name=${encodeURIComponent(suggestion.artist.name)}`);
            }
            
            a.dataset.artistName = suggestion.artist.name;

            if (suggestion.category === 'Concert' && suggestion.matchitem && suggestion.matchitem.dates) {
                content = document.createElement('div');
                content.className = 'content';

                const dateDiv = document.createElement('div');
                dateDiv.className = 'pic date';

                const dayMonthDiv = document.createElement('div');


                dayMonthDiv.textContent = suggestion.matchitem.dates.Day + ' ' + suggestion.matchitem.dates.Month;

                const yearDiv = document.createElement('div');
                yearDiv.textContent = suggestion.matchitem.dates.Year;

                dateDiv.appendChild(dayMonthDiv);
                dateDiv.appendChild(yearDiv);

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
                    img.src = suggestion.matchitem.imgLink ?  suggestion.matchitem.imgLink : '/icons/blank_cd_icon.png';
                    img.alt = 'Album cover of ' + (suggestion.matchitem.AlbumName || 'Unknown Album');
                } else if (suggestion.category ==='Member') {
                    img.className = 'pic user';
                    for (const [memberName,memberPic] of Object.entries(suggestion.artist.memberPics)){
                        if ( memberName === suggestion.matchitem){
                            img.src = suggestion.artist && memberPic ? memberPic : '/icons/artist_placeholder.svg';
                            img.alt = 'Profile image of ' + (suggestion.artist && memberName ? memberName : 'Unknown Artist');
                        }
                    }
                }else {
                    img.className = 'pic user';
                    img.src = suggestion.artist && suggestion.artist.strArtistThumb ? suggestion.artist.strArtistThumb : 'default-image-url.jpg';
                    img.alt = 'Profile image of ' + (suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist');
                }

                contentText = document.createElement('div');
                contentText.className = 'content-text go-down';

                //TODO suggestion populates search term inside bold-cut
                let boldCut = document.createElement('div');
                boldCut.className = 'p--bold cut';
                if (suggestion.category === 'Album') {
                    boldCut.textContent = suggestion.matchitem.AlbumName || '';
                } else{
                    boldCut.textContent = suggestion.matchitem || '';
                }

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
                populateResults.appendChild(categoryContainer);
                searchResults.classList.remove('hide');
                searchResults.appendChild(populateResults)
            }
        }
    }
});
