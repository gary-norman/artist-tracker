import {toggleIconFiltersButton, toggleVisFilters, toggleVisSubmit } from './default.js';
import { setupGridMouseMoveListener } from './mouseMove.js';

const pillContainer = document.querySelector('.pills')

document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');
    const populateResults = document.getElementById('populate-results');
    const filterInputs = document.querySelectorAll('input[id^="filter-"]');
    const filters = document.querySelectorAll('input[class^="filter-"]');

    let openFilters = 0;

    const locSearchInput = document.getElementById('button-filter-concert-location');
    const locSearchResults = document.getElementById('loc-search-result');
    const locationsContainer = document.getElementById('filter-checkbox-locations');
    const form = document.getElementById('form-homepage');

    locSearchResults.classList.add("hide");
    function debounce(fn, delay) {
        let timer;
        return function (...args) {
            clearTimeout(timer);
            timer = setTimeout(() => fn.apply(this, args), delay);
        };
    }
    
    if (locSearchInput && locSearchResults) {
        let latestRequestTimestamp = 0;

        async function fetchLocationSuggestions(query) {
            const response = await fetch(`/locationSuggest?query=${encodeURIComponent(query)}`);
            if (response.ok) {
                return await response.json();
            }
            throw new Error('Failed to fetch suggestions');
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
                locSearchResults.classList.remove("hide");
            } else {
                locationsContainer.innerHTML = '';
                locSearchResults.classList.add("hide");
            }
        });
    }

    function showLocationSuggestions(locSuggestionsData) {
        
        // console.log("received locSuggestionsData:=", locSuggestionsData);
            
        locationsContainer.innerHTML = ''; 

        if (!locSuggestionsData || locSuggestionsData.length === 0) {
            locationsContainer.innerHTML = '<p style="text-align:center">No results found.</p>';
            locSearchResults.appendChild(locationsContainer);
            return;
        }

        locSuggestionsData.forEach(function (location) {
            const checkboxContainer = document.createElement('label');
            checkboxContainer.className = 'checkbox go-across-sm';

            const input = document.createElement('input');
            input.className = 'checkbox checkbox-loc';
            input.htmlFor = input.id;
            input.id = `loc-${location.replace(/[\s, ]+/g, '-').toLowerCase()}`;
            input.type = 'checkbox';
            input.value = location;
            // check if already selected the location,if selected, then mark as checked
            if (isPillExist(input.id)){
                input.checked = true;
            }

            const label = document.createElement('label');
            label.className = 'checkbox small';
            label.htmlFor = input.id;
            label.textContent = location;

            checkboxContainer.appendChild(input);
            checkboxContainer.appendChild(label);
            locationsContainer.appendChild(checkboxContainer);
            
        });
        locSearchResults.appendChild(locationsContainer);

        // Assuming locationsContainer is defined
        const locCheckBoxes = locationsContainer.querySelectorAll('[id^="loc-"]');
        
        const pillContainer = document.querySelector('.pills'); 

        function createPill(location, id) {
            const pill = document.createElement('div');
            const pillText = document.createElement('p');
            const removePill = document.createElement('div');
        
            pill.className = 'pill';
            pill.id = 'pill_' + id;
            pillText.className = 'small';
            pillText.textContent = location;
            removePill.className = 'removePill';
            removePill.id = 'removePill_' + id;
        
            pill.appendChild(pillText);
            pill.appendChild(removePill);
        
            
            const checkboxContainer = document.createElement('label');
            checkboxContainer.className = 'checkbox go-across-sm';
            // Create hidden input to submit the selected location
            const hiddenInput = document.createElement('input');
            hiddenInput.type = 'hidden';
            hiddenInput.name = 'loc';
            hiddenInput.value = location;
            hiddenInput.id = 'hidden_' + id;
            
            // put those hidden location input directly into form 
            form.appendChild(hiddenInput);
        
            // Add click listener to the removePill button
            removePill.addEventListener('click', function() {
                pill.remove(); // Remove the pill from the container
                document.getElementById('hidden_' + id).remove(); // Remove hidden input
                const checkbox = document.getElementById(id);
                if (checkbox) checkbox.checked = false; // Uncheck the corresponding checkbox if it exists
            });
        
            pillContainer.appendChild(pill);
        }
        
        function isPillExist(id){
            const pillId = 'pill_' + id;
            const pillContainer = document.getElementById(pillId);
            if (pillContainer){
                return true
            } else{
                return false
            }
        }

        function handleCheckboxInput() {
            locCheckBoxes.forEach(input => {
                input.addEventListener('input', function() {
                    const id = input.id;
                    const location = input.value;
        
                    if ((input.checked) && !(isPillExist(input.id))) {
                        console.log("create pill for locataion:",location)
                        createPill(location, id);
                    } else {
                        const pill = document.getElementById('pill_' + id);
                        if (pill) pill.remove();
                        const hiddenInput = document.getElementById('hidden_' + id);
                        if (hiddenInput) hiddenInput.remove();
                    }
                        // after once click, claer the input  value
                    input.addEventListener('change', function () {
                        locSearchInput.value = ''; 
                    });   
                });
            });
        } 
        
// Call the function to set up initial event listeners
        handleCheckboxInput();
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
        // debug print 
        console.log("received suggestionsData:=", suggestionsData);
        
        const resultsHeader = document.querySelector('.filters .small.light.center');
        if (!suggestionsData || suggestionsData.length === 0) {
            resultsHeader.textContent = ``;
            populateResults.innerHTML = 'No results found.<br>Search for something else 🎹';
            searchResults.appendChild(populateResults);
            searchResults.classList.remove('hide');
            return;
        } else if (resultsHeader) {
            filterInputs.forEach(input => {
                if (input.checked){
                    openFilters++
                }
            });
            if (searchInput.value.trim() === '' && openFilters > 0) {
                resultsHeader.textContent = `Showing ${suggestionsData.length} results`;
            } else {
                resultsHeader.textContent = ``;
            }

        }
        
        const categories = {
            Artist: document.createElement('div'),
            Member: document.createElement('div'),
            Album: document.createElement('div'),
            Concert: document.createElement('div')
        };

        for (const category in categories) {
            const categoryContainer = categories[category];
            categoryContainer.className = 'grid col col1';

            const outerContainer = document.createElement('div');
            outerContainer.id = 'search-results-v2';
            outerContainer.className = 'container';

            const scrollContainer = document.createElement('div');
            scrollContainer.className = 'container scroll';
            scrollContainer.id = 'populate-results'
            outerContainer.appendChild(scrollContainer);

            const resultContainerTwo = document.createElement('div');
            resultContainerTwo.className = 'grid col col2';

            const resultContainerThree = document.createElement('div');
            resultContainerThree.className = 'grid col col3';

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
            let content = document.createElement('a');
            
            // Add the concert part for the Concert category
            if (suggestion.category === 'Concert') {
                content.setAttribute('href', `/artist?name=${encodeURIComponent(suggestion.artist.name)}#artist-concerts`);
            } else {
                content.setAttribute('href', `/artist?name=${encodeURIComponent(suggestion.artist.name)}`);
            }
            
            content.dataset.artistName = suggestion.artist.name;

            if (suggestion.category === 'Concert' && suggestion.matchitem && suggestion.matchitem.dates) {
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
                setupGridMouseMoveListener();

                if (categories[suggestion.category]) {
                    categories[suggestion.category].resultContainer.appendChild(content);
                    categories[suggestion.category].resultsCount++;
                } else {
                    console.warn(`Category ${suggestion.category} not found in categories object.`);
                }
            } else {
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
                } else {
                    img.className = 'pic user';
                    img.src = suggestion.artist && suggestion.artist.strArtistThumb ? suggestion.artist.strArtistThumb : 'default-image-url.jpg';
                    img.alt = 'Profile image of ' + (suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist');
                }

                contentText = document.createElement('div');
                contentText.className = 'content-text go-down';

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
                    artistName.className = 'small light cut';
                    artistName.textContent = suggestion.artist && suggestion.artist.name ? suggestion.artist.name : 'Unknown Artist';
                    contentText.appendChild(boldCut);
                    contentText.appendChild(artistName);
                }

                content.appendChild(img);
                content.appendChild(contentText);
                setupGridMouseMoveListener();

                if (categories[suggestion.category]) {
                    categories[suggestion.category].resultContainer.appendChild(content);
                    categories[suggestion.category].resultsCount++;
                } else {
                    console.warn(`Category ${suggestion.category} not found in categories object.`);
                }
            }
        });



        for (const category in categories) {
            const categoryContainer = categories[category];

            if (categoryContainer.resultsCount > 0) {
            // if (categoryContainer.resultsCount > 0) {
                const header = document.createElement('h2');
                header.textContent = category + 's';
                categoryContainer.insertBefore(header, categoryContainer.firstChild);
                populateResults.appendChild(categoryContainer);
                if (searchInput.value.trim() === ''){
                    filterInputs.forEach(input => {
                        if (input.checked){
                            openFilters++
                        }
                    });
                    if (openFilters > 0) {
                        searchResults.classList.remove('hide');
                    }
                } else {
                    searchResults.classList.remove('hide');
                }
                searchResults.appendChild(populateResults)
            }
        }
    }
    
    // event lestine for submit form
    form.addEventListener('submit', function (e) {
        e.preventDefault(); // Prevent the default form submission

        // Serialize form data
        const formData = new FormData(form);

        // Optional: You can add additional processing or validation here before sending the request

        // Fetch search results
        async function fetchSearchResult() {
            try {
                const response = await fetch(`/search/?${new URLSearchParams(formData).toString()}`);
                if (response.ok) {
                    const searchData = await response.json();
                    //debug print
                    console.log("success!!! all search data:", searchData);
                    
                    // Handle the received searchData, e.g., update the UI
                    showSuggestions(searchData);
                    toggleVisFilters();
                    toggleIconFiltersButton();
                    toggleVisSubmit();
                } else {
                    throw new Error('Failed to fetch search results');
                }
            } catch (error) {
                console.error('Error fetching search results:', error);
            }
        }

        fetchSearchResult();
    });
    
});