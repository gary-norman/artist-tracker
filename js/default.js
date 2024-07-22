import { formatDateToUK } from './calendar.js';

function wait(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

document.getElementById('form-homepage').addEventListener('keydown', function(event) {
    // Check if the key pressed is Enter (key code 13)
    if (event.key === 'Enter' || event.code === 'Enter') {
        // Prevent the default action (form submission)
        event.preventDefault();
    }
});

const date = new Date();
let currentDate = formatDateToUK(date);
let currentYear = date.getFullYear();

document.addEventListener('DOMContentLoaded', function () {
    const filterSubmit = document.getElementById("search-submit-filter");
    let openFilters = 0; // Initialize counter outside the function
    let filterNumber = document.getElementById("button-filter-number");

    // Function to toggle the visibility of filter-open containers and the submit button
    function toggleFilterContainers() {
        openFilters = 0; // Reset the counter

        // Select all input elements that match the criteria
        const filterInputs = document.querySelectorAll('input[id^="filter-"]');


        // Loop through each input element
        filterInputs.forEach(input => {
            // Increment the counter if the filter is checked
            if (input.checked) {
                openFilters++;
                let parent = input.closest('.filter');
                if (parent) {
                    const filterOpenElements = parent.querySelectorAll('[class^="filter-open"]');
                    parent.classList.add('open');
                    filterOpenElements.forEach(filter => filter.classList.remove('hide'));
                }
            } else {
                // Handle the case where the filter is unchecked
                let parent = input.closest('.filter');
                if (parent) {
                    const filterOpenElements = parent.querySelectorAll('[class^="filter-open"]');
                    parent.classList.remove('open');
                    filterOpenElements.forEach(filter => filter.classList.add('hide'));
                }
            }
        });

        // Show or hide the filter submit button based on the number of open filters
        if (openFilters > 0) {
            filterSubmit.classList.remove('hide');
            filterNumber.textContent = `Filters (${openFilters})`;
        } else {
            filterSubmit.classList.add('hide');
            filterNumber.textContent = `Filters (0)`;
        }

        console.log("Number of open filters: ", openFilters); // For debugging
    }

    // Select all input elements that match the criteria
    const filterInputs = document.querySelectorAll('input[id^="filter-"]');

    // Find all the end date elements
    const artistCreationEndDateInput = document.getElementById('creation-year-end');
    const concertEndDateInput = document.getElementById('concert-end-date');
    const albumEndDateInput = document.getElementById('album-end-date');

    // Set the end date of filters to today's date, and put inside placeholder
    //const currentDate = new Date().toISOString().split('T')[0]; // Get today's date in YYYY-MM-DD format
    artistCreationEndDateInput.placeholder = currentYear;
    concertEndDateInput.placeholder = currentDate;
    albumEndDateInput.placeholder = currentDate;

    // Add event listeners to each input element
    filterInputs.forEach(input => {
        input.addEventListener('change', toggleFilterContainers);
        // Initialize counter based on current state of inputs
        if (input.checked) {
            openFilters++;
        }
    });

    // Initial check to set the correct visibility state for the submit button
    if (openFilters > 0) {
        filterSubmit.classList.remove('hide');
    } else {
        filterSubmit.classList.add('hide');
    }

    console.log("Initial number of open filters: ", openFilters); // For debugging

    // Initial check to set the correct visibility state
    toggleFilterContainers();
});

// document.addEventListener('DOMContentLoaded', () => {
//     const members = document.querySelectorAll('[class^="member-item"]');
//
//     if (!members.length) {
//         console.log("No members present");
//         return; // Exit if no members are found
//     } else {
//         console.log(members.length, " members present");
//
//     }
//     console.log("Members present__", members.length);
//     console.error("Members present__", members.length);
//
//
//     members.forEach(member => {
//         member.addEventListener('mouseover', () => toggleMemberCard(member, true));
//         member.addEventListener('mouseleave', () => toggleMemberCard(member, false));
//     });
//
//     function toggleMemberCard(member, hover) {
//         if (hover) {
//             console.log("Mouse firing on member");
//             member.classList.remove("cut");
//             member.style.setAttribute('text-wrap','wrap'); // Correct the text-wrap setting
//         } else {
//             member.classList.add("cut");
//             member.style.setAttribute('text-wrap','nowrap');        }
//     }
// });

// document.addEventListener('DOMContentLoaded', () => {
//     const globeIcon = document.querySelector('.globe');
//     const mapIcon = document.querySelector('.kaart');
//     const parent = document.getElementById('mapProject');
//     let switchButtonText = document.querySelector('#mapProject .button-text');
//     let resetButtonText = document.querySelector('#resetMap .button-text');
//     const mapControls = document.getElementById('map-controls');
//     const mapControlsContainer = document.querySelector('.mapControls');
//
//     parent.addEventListener('click', () => {
//         toggleMapView();
//     });
//
//     //initialise button text on load
//     toggleButtonText();
//     toggleControlsView();
//
//     window.addEventListener('resize', () => {
//         toggleButtonText();
//         toggleControlsView();
//     });
//
//
//     async function toggleMapView() {
//         if (globeIcon.classList.contains('hide-icon')){
//             globeIcon.classList.remove('hide-icon');
//             mapIcon.classList.add('hide-icon');
//             toggleButtonText();
//
//         } else {
//             mapIcon.classList.remove('hide-icon');
//             globeIcon.classList.add('hide-icon');
//             toggleButtonText();
//         }
//     }
//
//     function toggleButtonText() {
//         let isGlobe = false;
//
//         if (globeIcon.classList.contains('hide-icon')) {
//             isGlobe = true;
//         }
//         if (window.innerWidth < 380) {
//             if (isGlobe) {
//                 switchButtonText.textContent = '2D view';
//             } else {
//                 switchButtonText.textContent = '3D view';
//             }
//             resetButtonText.textContent = 'Reset';
//         } else if (window.innerWidth < 650) {
//             if (isGlobe) {
//                 switchButtonText.textContent = 'Switch to 2D';
//             } else {
//                 switchButtonText.textContent = 'Switch to 3D';
//             }
//             resetButtonText.textContent = 'Reset map';
//         } else {
//             if (isGlobe) {
//                 switchButtonText.textContent = 'Switch to 2D map';
//             } else {
//                 switchButtonText.textContent = 'Switch to 3D map';
//             }
//             resetButtonText.textContent = 'Reset map position';
//         }
//     }
//
//     function toggleControlsView(){
//         if (window.innerWidth < 550) {
//             mapControls.classList.add('stretch');
//             mapControlsContainer.classList.remove('space');
//         } else {
//             mapControls.classList.remove('stretch');
//             mapControlsContainer.classList.add('space');
//         }
//     }
// });

document.addEventListener('DOMContentLoaded', () => {
    const searchButton = document.getElementById('button-filter-concert-location');
    const searchLocResultsContainer = document.getElementById('loc-search-result')
    const searchLocResults = document.getElementById('filter-checkbox-locations');
    const searchIcon = document.getElementById('search-loc-icon');
    let isSearching = false;

    function debounce(func, wait) {
        let timeout;
        return function(...args) {
            const context = this;
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(context, args), wait);
        };
    }

    function showSection(element) {
        element.classList.remove('hide');
    }

    function hideSection(element) {
        element.classList.add('hide');
    }

    const updateSearchCancelIcon = (button) => {
        const search = "url('../icons/search_x16.svg')";
        const cancel = "url('../icons/close_x16.svg')";

        if (button === "cancel") {
            document.documentElement.style.setProperty("--search-loc-icon", cancel);
        } else {
            document.documentElement.style.setProperty("--search-loc-icon", search);
        }
    };

    searchButton.addEventListener('focus', debounce(() => {

        // Clear suggestions if input is empty
        if (searchButton.value === '') {
            searchLocResults.innerHTML = '';
            searchButton.placeholder = 'Start typing location...';
        }
    }, 300));

    searchButton.addEventListener('input', debounce(function() {
        if (searchButton.value.trim() !== '') {
            isSearching = true;
            updateSearchCancelIcon("cancel");
            showSection(searchLocResultsContainer);
            console.log("showing search loc results on input")
        } else {
            isSearching = true;
            updateSearchCancelIcon("search");
            hideSection(searchLocResultsContainer);
            console.log("hiding search loc results")
        }
    }, 300));

    searchIcon.addEventListener('click', debounce (function(e) {
        e.stopPropagation(); // Prevent the click from propagating to the parent element**

        // Handle the click on the clear icon
        searchButton.value = '';
        searchLocResults.innerHTML = ''; // Clear suggestions if input is empty
        searchButton.placeholder = 'Start typing location...';
        updateSearchCancelIcon("search");
        hideSection(searchLocResultsContainer);
        console.log("hiding from searchIcon click")
    }, 300));




    // document.addEventListener('click', debounce(function(event) {
    //     let clickInsideAnyElement = false;
    //
    //     allSearchElements.forEach(element => {
    //         if (element.contains(event.target)) {
    //             clickInsideAnyElement = true;
    //             // console.log("inside: ", event.target)
    //         } else {
    //
    //         }
    //     });
    //
    //     if (!clickInsideAnyElement) {
    //         showSection(homeElements);
    //         hideSection(searchElements);
    //         console.log("hiding from click outside")
    //         changeLogo(logo, subLogo, "large");
    //         updateSearchCancelIcon("search");
    //     }
    // }, 300));
});


document.addEventListener('DOMContentLoaded', () => {
    const searchButton = document.getElementById("search-input");
    const home = document.querySelectorAll('[id^="home"]');
    const filters = document.getElementById("search-filters");
    const resultsHeader = document.querySelector('.filters .small.light.center');
    const searchResults = document.getElementById("search-results");
    const searchIcon = document.getElementById("search-icon");
    const logo = document.querySelector('.logo');
    const subLogo = document.querySelector('.sub-logo');
    let isSearching = false;
    const homeElements = [...home, subLogo];
    const searchElements = [searchResults, filters];
    const allSearchElements = document.querySelectorAll('[id^="search"]');

    function debounce(func, wait) {
        let timeout;
        return function(...args) {
            const context = this;
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(context, args), wait);
        };
    }

    function showSections(elements) {
        elements.forEach(element => {
            // if (element.classList.contains('hide')) {
                element.classList.remove('hide');
            // }
        });
    }

    function hideSections(elements) {
        elements.forEach(element => {
            // if (!element.classList.contains('hide')) {
                element.classList.add('hide');
            // }
        });
    }

    function changeLogo(logo, subLogo, size) {
        if (size === "large") {
            logo.style.fontSize = '8rem';
            logo.style.lineHeight = '1.2';
            subLogo.style.lineHeight = 'normal';
        } else if (size === "small") {
            logo.style.fontSize = '4rem';
            logo.style.lineHeight = '1.5';
            subLogo.style.lineHeight = '1.5';
        }
    }

    const updateSearchCancelIcon = (button) => {
        const search24 = "url('../icons/search_x24.svg')";
        const cancel24 = "url('../icons/close_x24.svg')";

        if (button === "cancel") {
            document.documentElement.style.setProperty("--search-main-icon", cancel24);
        } else {
            document.documentElement.style.setProperty("--search-main-icon", search24);
        }
    };

    searchButton.addEventListener('focus', debounce(() => {
            // showSections(searchElements);
            filters.classList.remove('hide');
            hideSections(homeElements);
            searchButton.placeholder = 'Start typing...';

            // Clear suggestions if input is empty
            if (searchButton.value === '') {
                searchResults.innerHTML = '';
            } else {
                // searchResults.innerHTML = searchButton.value;   //make suggestions populate as for already entered text
            }
            changeLogo(logo, subLogo, "small");
            updateSearchCancelIcon("cancel");
    }, 300));

    searchButton.addEventListener('input', debounce(function() {
        if (searchButton.value.trim() !== '') {
            isSearching = true;
            updateSearchCancelIcon("cancel");
            showSections(searchElements);
            hideSections(homeElements);
            console.log("hiding home elements; showing search elements")
            changeLogo(logo, subLogo, "small");
        } else {
            isSearching = true;
            updateSearchCancelIcon("search");
            showSections(homeElements);
            hideSections(searchElements);
            console.log("hiding from click 2")
            changeLogo(logo, subLogo, "large");
        }
    }, 1));

    searchIcon.addEventListener('click', debounce (function(e) {
        e.stopPropagation(); // Prevent the click from propagating to the parent element**

        // Handle the click on the clear icon
        searchButton.value = '';
        resultsHeader.textContent = `Showing 0 results`; // reset numbers of search results
        searchResults.innerHTML = ''; // Clear suggestions if input is empty
        searchButton.placeholder = 'Search an artist, member, album or concert';
        updateSearchCancelIcon("search");
        showSections(homeElements);
        hideSections(searchElements);
        console.log("hiding from searchIcon click")
        changeLogo(logo, subLogo, "large");
    }, 300));




    // document.addEventListener('click', debounce(function(event) {
    //     let clickInsideAnyElement = false;
    //
    //     allSearchElements.forEach(element => {
    //         if (element.contains(event.target)) {
    //             clickInsideAnyElement = true;
    //             // console.log("inside: ", event.target)
    //         } else {
    //
    //         }
    //     });
    //
    //     if (!clickInsideAnyElement) {
    //         showSections(homeElements);
    //         hideSections(searchElements);
    //         console.log("hiding from click outside")
    //         changeLogo(logo, subLogo, "large");
    //         updateSearchCancelIcon("search");
    //     }
    // }, 300));
});

function updateSliderBackground(slider) {
    const value = slider.value;
    const min = slider.min;
    const max = slider.max;
    const percentage = ((value - min) / (max - min)) * 100;

    slider.style.background = `linear-gradient(to right, var(--green-0) ${percentage}%, var(--white-4) ${percentage}%)`;
}



document.addEventListener('DOMContentLoaded', () => {
    const filterButton = document.getElementById('button-filter');
    const filters = document.querySelectorAll('.filter:not(:first-child)');
    const filterInputs = document.querySelectorAll('input[id^="filter-"]');
    let openFilters = 0;
    const filterSubmit = document.getElementById("search-submit-filter");
    let filterNumber = document.getElementById("button-filter-number");

    const show24 = "url('../icons/show_x24.svg')";
    const hide24 = "url('../icons/hide_x24.svg')";
    const show16 = "url('../icons/show_x16.svg')";
    const hide16 = "url('../icons/hide_x16.svg')";

    console.log("Initializing show-hide-icon to hide24");
    document.documentElement.style.setProperty("--show-hide-icon", show24);

    const updateSearchSubmitButtonVisibility = () => {
        const filterSubmit = document.getElementById("search-submit-filter");


        // Check if any filters are visible (not hidden)
        const anyFilterVisible = Array.from(filters).some(filter => !filter.classList.contains('hide'));

        if (anyFilterVisible) {
            filterSubmit.classList.remove('hide');
        } else {
            filterSubmit.classList.add('hide');
        }
    }

    const updateShowHideIcon = () => {
        const isHidden = filters[0].classList.contains('hide');
        if (window.innerWidth < 500) {
            console.log(`Setting show-hide-icon to ${isHidden ? 'show16' : 'hide16'}`);
            document.documentElement.style.setProperty("--show-hide-icon", isHidden ? show16 : hide16);
        } else {
            console.log(`Setting show-hide-icon to ${isHidden ? 'show24' : 'hide24'}`);
            document.documentElement.style.setProperty("--show-hide-icon", isHidden ? show24 : hide24);
        }
    };

    filterButton.addEventListener('click', () => {
        filters.forEach(filter => {
            filter.classList.toggle('hide');
        });
        updateShowHideIcon();
        updateSearchSubmitButtonVisibility();
    });

    window.addEventListener('resize', updateShowHideIcon);

    // Initial update based on the current window size
    updateShowHideIcon();
});



// Rin commentout
// Initialize the background on page load
/* document.addEventListener('DOMContentLoaded', () => {
    const slider = document.getElementById('album-date-range');
    updateSliderBackground(slider);

    slider.addEventListener('input', () => {
        updateSliderBackground(slider);
    });
}); 

// Initialize the background on page load
 document.addEventListener('DOMContentLoaded', () => {
    const slider = document.getElementById('artist-date-range');
    updateSliderBackground(slider);

    slider.addEventListener('input', () => {
        updateSliderBackground(slider);
    });
}); */

function updateDoubleSliderBackground(slider1, slider2) {
    const value1 = slider1.value;
    const value2 = slider2.value;
    const min = slider1.min;
    const max = slider2.max;
    const percentageLeft = ((value1 - min) / (max - min)) * 100;
    const percentageRight = ((value2 - min) / (max -min)) * 100;

    slider1.style.background = `linear-gradient(
    to right, 
    var(--white-4) ${percentageLeft}%, 
    var(--green-0) ${percentageLeft}%, 
    var(--green-0) ${percentageRight}%, 
    var(--white-4) ${percentageRight}%)`;

    slider2.style.background = `linear-gradient(
    to right, 
    var(--white-4) ${percentageLeft}%, 
    var(--green-0) ${percentageLeft}%, 
    var(--green-0) ${percentageRight}%, 
    var(--white-4) ${percentageRight}%)`;

    const filterBar = `linear-gradient(to right, var(--white-4) ${percentageLeft}%, var(--green-0) ${percentageLeft}%, var(--green-0) ${percentageRight}%, var(--white-4) ${percentageRight}%)`;
    document.documentElement.style.setProperty("--filter-bar", filterBar);



}

(function() {
    function addSeparator(nStr) {
        nStr += '';
        var x = nStr.split('.');
        var x1 = x[0];
        var x2 = x.length > 1 ? '.' + x[1] : '';
        var rgx = /(\d+)(\d{3})/;
        while (rgx.test(x1)) {
            x1 = x1.replace(rgx, '$1' + '.' + '$2');
        }
        return x1 + x2;
    }

    function updateRangeLabel(range_min, range_max, minVal, maxVal) {
        if (window.innerWidth < 800) {
            if (minVal === 10) {
                $(range_min).html(addSeparator(minVal) + '+');
            } else {
                $(range_min).html(addSeparator(minVal));
            }

            if (maxVal === 10) {
                $(range_max).html(addSeparator(maxVal) + '+');
            } else {
                $(range_max).html(addSeparator(maxVal));
            }

        } else {
            var minText = minVal > 9 ? addSeparator(minVal) + '+ Members' : minVal > 1 ? addSeparator(minVal) + ' Members' : addSeparator(minVal) + ' Member';
            var maxText = maxVal > 9 ? addSeparator(maxVal) + '+ Members' : maxVal > 1 ? addSeparator(maxVal) + ' Members' : addSeparator(maxVal) + ' Member';
            $(range_min).html(minText);
            $(range_max).html(maxText);
        }
    }

    function rangeInputChangeEventHandler(e){
        var rangeGroup = $(this).attr('name'),
            minBtn = $(this).parent().children('.min'),
            maxBtn = $(this).parent().children('.max'),
            range_min = $(this).parent().children('.range_min'),
            range_max = $(this).parent().children('.range_max'),
            minVal = parseInt($(minBtn).val()),
            maxVal = parseInt($(maxBtn).val()),
            origin = $(this).context.className;

        if(origin === 'min' && minVal > maxVal){
            $(minBtn).val(maxVal);
        }
        minVal = parseInt($(minBtn).val());

        if(origin === 'max' && maxVal < minVal){
            $(maxBtn).val(minVal);
        }
        maxVal = parseInt($(maxBtn).val());

        updateRangeLabel(range_min, range_max, minVal, maxVal);
    }

    $('input[type="range"]').on( 'input', rangeInputChangeEventHandler);

    function updateLabelPosition(slider, label) {
        const sliderWidth = slider.offsetWidth;
        const sliderMin = parseInt(slider.min);
        const sliderMax = parseInt(slider.max);
        const sliderValue = parseInt(slider.value);

        const slider1 = document.getElementById('members-min-range');
        const slider2 = document.getElementById('members-max-range');
        const label1 = document.getElementById('label-members-min');
        const label2 = document.getElementById('label-members-max');




        if (window.innerWidth < 300) {
            const position = ((sliderValue - sliderMin) / (sliderMax - sliderMin)) * 93;
            label.style.left = `calc(${position}% - 0.4rem)`;
        } else if (window.innerWidth < 500) {
            const position = ((sliderValue - sliderMin) / (sliderMax - sliderMin)) * 94.5;
            label.style.left = `calc(${position}% - 0.4rem)`;
        } else if (window.innerWidth < 600) {
            const position = ((sliderValue - sliderMin) / (sliderMax - sliderMin)) * 95;
            label.style.left = `calc(${position}% - 0.4rem)`;
        } else if (window.innerWidth < 800){
            const position = ((sliderValue - sliderMin) / (sliderMax - sliderMin)) * 96;
            label.style.left = `calc(${position}% - 0.4rem)`;
        } else {
            const position = ((sliderValue - sliderMin) / (sliderMax - sliderMin)) * 97.5;
            label.style.left = `calc(${position}% - 3.2rem)`;
        }

        //set one of the labels to invisible when they overlap
        if (slider1.value === slider2.value) {
            label1.style.visibility = 'hidden';
        } else {
            label1.style.visibility = 'visible';
        }
    }

    window.addEventListener('resize', function() {
        const slider1 = document.getElementById('members-min-range');
        const slider2 = document.getElementById('members-max-range');
        const labelMin = document.querySelector('.range_min');
        const labelMax = document.querySelector('.range_max');

        updateDoubleSliderBackground(slider1, slider2);
        updateLabelPosition(slider1, labelMin);
        updateLabelPosition(slider2, labelMax);
        updateRangeLabel(labelMin, labelMax, parseInt(slider1.value), parseInt(slider2.value));
    });

    // Initialize the background on page load
    document.addEventListener('DOMContentLoaded', () => {
        const slider1 = document.getElementById('members-min-range');
        const slider2 = document.getElementById('members-max-range');
        const labelMin = document.querySelector('.range_min');
        const labelMax = document.querySelector('.range_max');

        // Update the background and position for each element on load
        updateDoubleSliderBackground(slider1, slider2);
        updateLabelPosition(slider1, labelMin);
        updateLabelPosition(slider2, labelMax);
        updateRangeLabel(labelMin, labelMax, parseInt(slider1.value), parseInt(slider2.value));

        // Event listener for slider 1 input
        slider1.addEventListener('input', () => {
            // Update the position of the labels
            updateLabelPosition(slider1, labelMin);

            // Update background color and labels
            updateDoubleSliderBackground(slider1, slider2);
            updateRangeLabel(labelMin, labelMax, parseInt(slider1.value), parseInt(slider2.value));
        });

        // Event listener for slider 2 input
        slider2.addEventListener('input', () => {
            // Update the position of the labels
            updateLabelPosition(slider2, labelMax);

            // Update background color and labels
            updateDoubleSliderBackground(slider1, slider2);
            updateRangeLabel(labelMin, labelMax, parseInt(slider1.value), parseInt(slider2.value));
        });
    });
})();