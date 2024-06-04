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
                        element.classList.remove('hide');
                    } else {
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


function updateSliderBackground(slider) {
    const value = slider.value;
    const min = slider.min;
    const max = slider.max;
    const percentage = ((value - min) / (max - min)) * 100;

    slider.style.background = `linear-gradient(to right, var(--green-0) ${percentage}%, var(--white-4) ${percentage}%)`;
}

// Initialize the background on page load
document.addEventListener('DOMContentLoaded', () => {
    const slider = document.getElementById('album-date-range');
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

}

// Initialize the background on page load
document.addEventListener('DOMContentLoaded', () => {
    const slider1 = document.getElementById('members-min-range');
    const slider2 = document.getElementById('members-max-range');
    const labelMin = document.querySelector('.range_min');
    const labelMax = document.querySelector('.range_max');
    updateDoubleSliderBackground(slider1, slider2);

    slider1.addEventListener('input', () => {
        // Update the position of the labels
        updateLabelPosition(slider1, labelMin);

        // Update background colour
        updateDoubleSliderBackground(slider1, slider2);
    });

    slider2.addEventListener('input', () => {
        // Update the position of the labels
        updateLabelPosition(slider2, labelMax);

        // Update background colour
        updateDoubleSliderBackground(slider1, slider2);
    });
});

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
        if(minVal > 9 ) {
            $(range_min).html(addSeparator(minVal) + '+ Members');
        } else if(minVal > 1 ) {
            $(range_min).html(addSeparator(minVal) + ' Members');
        } else {
            $(range_min).html(addSeparator(minVal) + ' Member');
        }


        if(origin === 'max' && maxVal < minVal){
            $(maxBtn).val(minVal);
        }
        maxVal = parseInt($(maxBtn).val());
        if(maxVal > 9 ) {
            $(range_max).html(addSeparator(maxVal) + '+ Members');
        } else if(maxVal > 1 ) {
            $(range_max).html(addSeparator(maxVal) + ' Members');
        }else {
            $(range_max).html(addSeparator(maxVal) + ' Member');
        }
    }

    $('input[type="range"]').on( 'input', rangeInputChangeEventHandler);
})();


var counter = 0
function updateLabelPosition(slider, label) {
    const sliderWidth = slider.offsetWidth;
    const sliderMin = parseInt(slider.min);
    const sliderMax = parseInt(slider.max);
    const sliderValue = parseInt(slider.value);

    const position = ((sliderValue - sliderMin) / (sliderMax - sliderMin)) * 100;
    counter += 6

    label.style.left = `calc(${position}%)`; // Adjust '20px' to center the label correctly
}