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
    updateDoubleSliderBackground(slider1, slider2);

    slider1.addEventListener('input', () => {
        updateDoubleSliderBackground(slider1, slider2);
    });

    slider2.addEventListener('input', () => {
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