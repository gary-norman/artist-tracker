document.addEventListener('DOMContentLoaded', function () {
    // Function to toggle the visibility of filter-open containers
    function toggleFilterContainers() {
        // Select all input elements that match the criteria
        const filterInputs = document.querySelectorAll('input[id^="filter-"]');

        // Loop through each input element
        filterInputs.forEach(input => {
            // Traverse up the DOM to find the parent with the class 'filter'
            let parent = input.closest('.filter');
            if (parent) {
                // Find all elements within the parent that have a class starting with 'filter-open'
                const filterOpenElements = parent.querySelectorAll('[class^="filter-open"]');


                // Toggle the 'hide' class on each element based on checkbox state
                filterOpenElements.forEach(element => {
                    if (input.checked) {

                        parent.classList.add('open');
                        element.classList.remove('hide');
                    } else {
                        parent.classList.remove('open');
                        element.classList.add('hide');
                    }
                });
            }
        });
    }

    // Select all input elements that match the criteria
    const filterInputs = document.querySelectorAll('input[id^="filter-"]');

    // Add event listeners to each input element
    filterInputs.forEach(input => {
        input.addEventListener('change', toggleFilterContainers);
    });

    // Initial check to set the correct visibility state
    toggleFilterContainers();
});

document.addEventListener('DOMContentLoaded', () => {
    const searchButton = document.getElementById("search-input");
    const home = document.querySelectorAll('[id^="home"]');
    const recent = document.getElementById("search-recent-text");
    const filters = document.getElementById("search-filters");
    const searchResults = document.getElementById("search-results");
    const searchIcon = document.getElementById("search-icon");
    const logo = document.querySelector('.logo');
    const subLogo = document.querySelector('.sub-logo');
    let isSearching = false;
    const homeElements = [...home, subLogo];
    const searchElements = [searchResults, recent, filters];
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
            if (element.classList.contains('hide')) {
                element.classList.remove('hide');
            }
        });
    }

    function hideSections(elements) {
        elements.forEach(element => {
            if (!element.classList.contains('hide')) {
                element.classList.add('hide');
            }
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
            document.documentElement.style.setProperty("--search-cancel-icon", cancel24);
        } else {
            document.documentElement.style.setProperty("--search-cancel-icon", search24);
        }
    };

    searchButton.addEventListener('focus', debounce(() => {
            showSections(searchElements);
            hideSections(homeElements);
            changeLogo(logo, subLogo, "small");
            updateSearchCancelIcon("cancel");
    }, 200));

    searchButton.addEventListener('input', debounce(function() {
        if (searchButton.value.trim() !== '') {
            isSearching = true;
            updateSearchCancelIcon("cancel");
            showSections(searchElements);
            hideSections(homeElements);
            console.log("hiding from click 2")
            changeLogo(logo, subLogo, "small");
        } else {
            isSearching = false;
            updateSearchCancelIcon("search");
            showSections(homeElements);
            hideSections(searchElements);
            console.log("hiding from click 2")
            changeLogo(logo, subLogo, "large");
        }
    }, 100));

    searchButton.addEventListener('click', debounce(function() {
        if (isSearching) {
            updateSearchCancelIcon("search");
            showSections(homeElements);
            hideSections(searchElements);
            console.log("hiding from click 1")
            changeLogo(logo, subLogo, "large");
        } else {
            updateSearchCancelIcon("cancel");
            showSections(searchElements);
            hideSections(homeElements);
            console.log("hiding from click 2")
            changeLogo(logo, subLogo, "small");
        }
    }, 200));

    searchIcon.addEventListener('click', debounce (function(e) {
        e.stopPropagation(); // Prevent the click from propagating to the parent element**

        // Handle the click on the clear icon
        searchButton.value = '';
        searchResults.innerHTML = ''; // Clear suggestions if input is empty
        updateSearchCancelIcon("search");
        showSections(homeElements);
        hideSections(searchElements);
        console.log("hiding from searchIcon click")
        changeLogo(logo, subLogo, "large");
    }, 200));


    document.addEventListener('click', debounce(function(event) {
        let clickInsideAnyElement = false;

        allSearchElements.forEach(element => {
            if (element.contains(event.target)) {
                clickInsideAnyElement = true;
                // console.log("inside: ", event.target)
            } else {

            }
        });

        if (!clickInsideAnyElement) {
            showSections(homeElements);
            hideSections(searchElements);
            console.log("hiding from click outside")
            changeLogo(logo, subLogo, "large");
            updateSearchCancelIcon("search");
        }
    }, 400));
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
    const filterSubmit = document.getElementById("search-submit-filter");
    const formContainer = $(filterSubmit).parent();

    const show24 = "url('../icons/show_x24.svg')";
    const hide24 = "url('../icons/hide_x24.svg')";
    const show16 = "url('../icons/show_x16.svg')";
    const hide16 = "url('../icons/hide_x16.svg')";

    console.log("Initializing show-hide-icon to hide24");
    document.documentElement.style.setProperty("--show-hide-icon", show24);

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
        filterSubmit.classList.toggle('hide');
        // formContainer.classList.toggle('open')
        filters.forEach(filter => {
            filter.classList.toggle('hide');
        });
        updateShowHideIcon();
    });

    window.addEventListener('resize', updateShowHideIcon);

    // Initial update based on the current window size
    updateShowHideIcon();
});

// Initialize the background on page load
document.addEventListener('DOMContentLoaded', () => {
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
});

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




        if (window.innerWidth < 400) {
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