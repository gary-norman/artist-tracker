document.addEventListener('DOMContentLoaded', function () {
    const artistFilter = document.getElementById('artist-creation-date-filter');
    const albumFilter = document.getElementById('album-creation-date-filter');
    const membersFilter = document.getElementById('members-filter');
    const concertsFilter = document.getElementById('concert-location-filter');
    const submitFilter = document.getElementById('search-submit-filter');
    const resultContainer = document.getElementById('search-results');
  
    let isStartCalendarOpen = false;
    let isEndCalendarOpen = false;
    let isAlbumStartCalendarOpen = false;
    let isAlbumEndCalendarOpen = false;
    
    let hideElements;
    let showElements;
    
     // artist
    const startDateInput = document.getElementById('artist-start-date');
    const endDateInput = document.getElementById('artist-end-date');
    const startDateContainer = document.getElementById('startDateContainer');
    const endDateContainer = document.getElementById('endDateContainer');
    
     // album
    const albumStartDateInput = document.getElementById('album-start-date');
    const albumEndDateInput = document.getElementById('album-end-date');
    const albumStartDateContainer = document.getElementById('albumStartDateContainer');
    const albumEndDateContainer = document.getElementById('albumEndDateContainer');
    
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    const days = ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su'];

    let selectedDate = new Date();
    let currentMonth = selectedDate.getMonth();
    let currentYear = selectedDate.getFullYear();
    let startYear = Math.floor(currentYear / 28) * 28; 

    const minYear = 1900; // Minimum year in the year selection
    const maxYear = new Date().getFullYear(); // Maximum year is the current year
    
    // Create calendar elements
    createCalendarElements(startDateContainer, 'Start');
    createCalendarElements(endDateContainer, 'End');
    createCalendarElements(albumStartDateContainer, 'AlbumStart');
    createCalendarElements(albumEndDateContainer, 'AlbumEnd');
    
     // artsit
    const dayCalendarStart = document.getElementById('dayCalendarStart');
    const monthCalendarStart = document.getElementById('monthCalendarStart');
    const yearCalendarStart = document.getElementById('yearCalendarStart');
    
    const dayCalendarEnd = document.getElementById('dayCalendarEnd');
    const monthCalendarEnd = document.getElementById('monthCalendarEnd');
    const yearCalendarEnd = document.getElementById('yearCalendarEnd');
    
    // album
    const dayCalendarAlbumStart = document.getElementById('dayCalendarAlbumStart');
    const monthCalendarAlbumStart = document.getElementById('monthCalendarAlbumStart');
    const yearCalendarAlbumStart = document.getElementById('yearCalendarAlbumStart');

    const dayCalendarAlbumEnd = document.getElementById('dayCalendarAlbumEnd');
    const monthCalendarAlbumEnd = document.getElementById('monthCalendarAlbumEnd');
    const yearCalendarAlbumEnd = document.getElementById('yearCalendarAlbumEnd');
    
    
    // Initialize all calendars,artsit and album
    initializeCalendar('start', 'day');
    initializeCalendar('end', 'day');
    initializeCalendar('start', 'month');
    initializeCalendar('end', 'month');
    initializeCalendar('start', 'year');
    initializeCalendar('end', 'year');
    
    initializeCalendar('albumStart', 'month');
    initializeCalendar('albumEnd', 'month');
    initializeCalendar('albumStart', 'day');
    initializeCalendar('albumEnd', 'day');
    initializeCalendar('albumStart', 'year');
    initializeCalendar('albumEnd', 'year');
    
    // Toggle calendar visibility for start date
    startDateInput.addEventListener('click', function () {
        
        dayCalendarStart.classList.toggle('hidden', false);
        
        // toggle other calendars all hidden 
        hideElements = [
            monthCalendarStart, yearCalendarStart,
            dayCalendarEnd,monthCalendarEnd, yearCalendarEnd,
        ];
        toggleElementVisibility(hideElements,false);
        
        // Hide other filters
        hideElements = [
            albumFilter, membersFilter, concertsFilter, submitFilter, resultContainer
        ];
        toggleElementVisibility(hideElements,false);
        isStartCalendarOpen = true;
        isEndCalendarOpen = false
        
        // check first if user select a date
        const endDateValue = endDateInput.value ? parseUKDate(endDateInput.value) : new Date();
        currentYear = endDateValue.getFullYear();
        currentMonth = endDateValue.getMonth();
        const currentDay = endDateValue.getDate();
        
        // Check if the start date is the first day of the month
        if (currentDay === 1) {
            currentMonth--;
            if (currentMonth < 0) { 
                currentMonth = 11; // Dec
                currentYear--;
            }
        }
        
        // debug print
        console.log("year now is:",currentYear)
        console.log("month now after formatting is:",currentMonth)
        
        renderDayCalendar(currentMonth, currentYear, "start");
    });

    // Toggle calendar visibility for end date
    endDateInput.addEventListener('click', function () {
       
        
        dayCalendarEnd.classList.toggle('hidden', false);
        // toggle other calendar all hidden
        hideElements = [
            monthCalendarEnd, yearCalendarEnd,
            dayCalendarStart,monthCalendarStart, yearCalendarStart,
        ];
        toggleElementVisibility(hideElements,false);
                
        // Hide other filters
        hideElements = [
            albumFilter, membersFilter, concertsFilter, submitFilter, resultContainer
        ];
        toggleElementVisibility(hideElements,false);
   
        isEndCalendarOpen = true;
        isStartCalendarOpen = false;
        
        // check first if user select a date, if selected then convert the UKformat back to normal
        const startDateValue = startDateInput.value ? parseUKDate(startDateInput.value) : new Date();
        currentYear = startDateValue.getFullYear();
        currentMonth = startDateValue.getMonth();
        const currentDay = startDateValue.getDate();
        
        console.log("current Day:",currentDay)
        
        // Check if the selected start date is the last day of the month
        const lastDayOfMonth = new Date(currentYear, currentMonth + 1, 0).getDate();
        
        console.log("last Day of month:",lastDayOfMonth)
        if (currentDay === lastDayOfMonth) {
            currentMonth++;
            if (currentMonth > 11) {
                currentMonth = 0; // Jan
                currentYear++;
            }
        }
        
        // debug print
        console.log("year now is:",currentYear)
        console.log("month now after formatting is:",currentMonth)
        
        renderDayCalendar(currentMonth, currentYear, "end");
    }); 
    
        // Toggle calendar visibility for start date
    albumStartDateInput.addEventListener('click', function () {
        
        dayCalendarAlbumStart.classList.toggle('hidden', false);
        // toggle other calendar all hidden , make sure if one calendar was opened, then hidden it
        hideElements = [
            monthCalendarAlbumStart, yearCalendarAlbumStart,
            dayCalendarAlbumEnd, monthCalendarAlbumEnd, yearCalendarAlbumEnd
        ];
        toggleElementVisibility(hideElements,false);
        // Hide other filters
        hideElements = [
            artistFilter,membersFilter, concertsFilter, submitFilter, resultContainer
        ];
        toggleElementVisibility(hideElements,false);
  
        isAlbumStartCalendarOpen = true;
        isAlbumEndCalendarOpen = false
        
        // check first if user select a date, if selected then convert the UKformat back to normal
        const albumEndDateValue = albumEndDateInput.value ? parseUKDate(albumEndDateInput.value) : new Date();
        currentYear = albumEndDateValue.getFullYear();
        currentMonth = albumEndDateValue.getMonth();
        let currentDay = albumEndDateValue.getDate();
        
        // Check if the start date is the first day of the month
        if (currentDay === 1) {
            currentMonth--;
            if (currentMonth < 0) { 
                currentMonth = 11; // Dec
                currentYear--;
            }
        }

        // debug print
/*         console.log("year now is:",currentYear)
        console.log("month now after formatting is:",currentMonth) */
         
        renderDayCalendar(currentMonth, currentYear, "albumStart");
    });

    // Toggle calendar visibility for end date
    albumEndDateInput.addEventListener('click', function () {
        dayCalendarAlbumEnd.classList.toggle('hidden', false);
        // toggle other calendar all hidden
        hideElements = [
            dayCalendarAlbumStart, monthCalendarAlbumStart, yearCalendarAlbumStart,
            monthCalendarAlbumEnd, yearCalendarAlbumEnd
        ];
        toggleElementVisibility(hideElements,false);
                
        // Hide other filters
        hideElements = [
            artistFilter,membersFilter, concertsFilter, submitFilter, resultContainer
        ];
        toggleElementVisibility(hideElements,false);
        
        isAlbumEndCalendarOpen = true;
        isAlbumStartCalendarOpen = false;
        
        // check first if user select a date
        const albumStartDateValue = albumStartDateInput.value ? parseUKDate(albumStartDateInput.value) : new Date();
        currentYear = albumStartDateValue.getFullYear();
        currentMonth = albumStartDateValue.getMonth();
        const currentDay =  albumStartDateValue.getDate();
        
        // Check if the selected start date is the last day of the month
        const lastDayOfMonth = new Date(currentYear, currentMonth + 1, 0).getDate();

        console.log("last Day of month:",lastDayOfMonth)
        if (currentDay === lastDayOfMonth) {
            currentMonth++;
            if (currentMonth > 11) { // Handle year increment if month exceeds December
                currentMonth = 0;
                currentYear++;
            }
        }
        // debug print
    /*     console.log("year now is:",currentYear)
        console.log("month now after formatting is:",currentMonth) */
        
        renderDayCalendar(currentMonth, currentYear, "albumEnd");
    }); 
    
    // Function to create calendar elements,* Day Calendar,* Month Calendar , *Year Calendar  
    function createCalendarElements(container, type) {
        container.innerHTML = `
            <div class="date-picker-calendar hidden" id="dayCalendar${type.charAt(0).toUpperCase() + type.slice(1)}">
                <div class="calendar-header">
                    <button class="cal-btn back-year">&lt;</button>
                    <span class="calendar-year"></span>
                    <button class="cal-btn front-year">&gt;</button>
                    <button class="cal-btn back-month">&lt;</button>
                    <span class="calendar-month"></span>
                    <button class="cal-btn front-month">&gt;</button>
                </div>
                <div class="cal-wrapper">
                    <div class="cal-days p--bold"></div>
                    <div class="calendar-main"></div>
                </div>
            </div>
            <div class="date-picker-calendar hidden" id="monthCalendar${type.charAt(0).toUpperCase() + type.slice(1)}">
                <div class="calendar-header">
                    <button class="cal-btn back">&lt;</button>
                    <span class="calendar-year"></span>
                    <button class="cal-btn front">&gt;</button>
                </div>
                <div class="cal-wrapper">
                    <div class="cal-months"></div>
                </div>
            </div>
            <div class="date-picker-calendar hidden" id="yearCalendar${type.charAt(0).toUpperCase() + type.slice(1)}">
                <div class="calendar-header">
                    <button class="cal-btn back">&lt;</button>
                    <button class="cal-btn front">&gt;</button>
                </div>
                <div class="cal-wrapper">
                    <div class="cal-years"></div>
                </div>
            </div>
        `;
    }
    
    // Function to initialize calendar components
    function initializeCalendar(type, viewType) {
        const containerId = `${viewType}Calendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);

        // Ensure calendar is hidden initially
        container.classList.toggle('hidden', true);

        // Example of accessing specific elements within the container
        const backYearButton = container.querySelector('.cal-btn.back-year');
        const forwardYearButton = container.querySelector('.cal-btn.front-year');
        const backMonthButton = container.querySelector('.cal-btn.back-month');
        const forwardMonthButton = container.querySelector('.cal-btn.front-month');
        const backButton = container.querySelector('.cal-btn.back');
        const forwardButton = container.querySelector('.cal-btn.front');

        // Function to render the calendar based on the view type
        function renderCalendar() {
            switch (viewType) {
                case 'day':
                    renderDayCalendar(currentMonth, currentYear, type);
                    break;
                case 'month':
                    renderMonthCalendar(currentYear, type);
                    break;
                case 'year':
                    renderYearCalendar(currentYear, type);
                    break;
                default:
                    console.error(`Invalid viewType: ${viewType}`);
            }
        }

        // Attach event listeners for navigation buttons based on viewType
        if (viewType === 'day') {
            backYearButton.addEventListener('click', function() {
                currentYear--;
                if (currentYear < minYear) {
                    currentYear = minYear;
                }
                renderDayCalendar(currentMonth, currentYear, type);
            });
            forwardYearButton.addEventListener('click', function() {
                currentYear++;
                if (currentYear > maxYear) {
                    currentYear = maxYear;
                }
                renderDayCalendar(currentMonth, currentYear, type);
            });
            backMonthButton.addEventListener('click', function() {
                currentMonth--;
                if (currentMonth < 0) {
                    currentMonth = 11;
                    currentYear--;
                }
                renderDayCalendar(currentMonth, currentYear, type);
            });
            forwardMonthButton.addEventListener('click', function() {
                currentMonth++;
                if (currentMonth > 11) {
                    currentMonth = 0;
                    currentYear++;
                }
                renderDayCalendar(currentMonth, currentYear, type);
            });
        } else if (viewType === 'month') {
            backButton.addEventListener('click', function() {
                currentYear--;
                renderMonthCalendar(currentYear, type);
            });
            forwardButton.addEventListener('click', function() {
                currentYear++;
                renderMonthCalendar(currentYear, type);
            });
        } else if (viewType === 'year') {
            backButton.addEventListener('click', function() {
                currentYear -= 28;
                renderYearCalendar(currentYear, type);
            });
            forwardButton.addEventListener('click', function() {
                currentYear += 28;
                if (currentYear > maxYear) {
                    currentYear = maxYear;
                }
                renderYearCalendar(currentYear, type);
            });
        }

        // Initial render
        renderCalendar();
    }

    // Render Day Calendar
    function renderDayCalendar(month, year, type) {
      
        // debug print
 /*        console.log("renderDayCalendar type is:",type)
        console.log("month got parse to renderDayCalendar:",month)
        console.log("year got parse to renderDayCalendar:",year) */
        
        const containerId = `dayCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);
        
        // debug print
        console.log("container is:",container)

        const yearDisplayDay = container.querySelector('.calendar-year');
        const monthDisplayDay = container.querySelector('.calendar-month');
        const daysContainer = container.querySelector('.cal-days');
        const datesContainer = container.querySelector('.calendar-main');

        yearDisplayDay.textContent = `${year}`;
        monthDisplayDay.textContent = `${months[month]}`;

        // Calculate the first day of the month
        const firstDay = new Date(year, month, 1).getDay();
        const startingDay = (firstDay === 0) ? 6 : firstDay - 1;

        // Map days with conditional color for Sa and Su
        daysContainer.innerHTML = days.map((day, index) => {
            if (index === 5 || index === 6) { // 5 represents Saturday, 6 represents Sunday
                return `<div class="days weekend">${day}</div>`;
            } else {
                return `<div class="days">${day}</div>`;
            }
        }).join('');

        datesContainer.innerHTML = '';
        // Add empty divs for the days before the first day of the month
        for (let i = 0; i < startingDay; i++) {
            datesContainer.innerHTML += '<div></div>';
        }

        const lastDate = new Date(year, month + 1, 0).getDate();

        // Render days of the month
        for (let date = 1; date <= lastDate; date++) {
            const dateElement = document.createElement('div');
            dateElement.textContent = date;
            dateElement.classList.add('day'); // Only add 'day' class

            // Add class for Saturdays and Sundays
            const dayOfWeek = new Date(year, month, date).getDay();
            if (dayOfWeek === 6 || dayOfWeek === 0) {
                dateElement.classList.add('weekend');
            }

            dateElement.addEventListener('click', function () {
                // hide all calendar
                hideElements = [
                    dayCalendarStart, monthCalendarStart, yearCalendarStart,
                    dayCalendarEnd, monthCalendarEnd, yearCalendarEnd,
                    dayCalendarAlbumStart, monthCalendarAlbumStart, yearCalendarAlbumStart,
                    dayCalendarAlbumEnd, monthCalendarAlbumEnd, yearCalendarAlbumEnd
                ];
                toggleElementVisibility(hideElements,false);
                selectDate(year, month, date, type);
            });

            datesContainer.appendChild(dateElement);
        }
        
        
        // Toggle calendar visibility (month view)
        monthDisplayDay.addEventListener('click', function () {
            if (type ==="start" || type ==="albumStart"){
                hideElements = [
                    dayCalendarStart,yearCalendarStart,
                    dayCalendarAlbumStart,yearCalendarAlbumStart,
                ];
                toggleElementVisibility(hideElements,false);
                showElements = [ monthCalendarStart, monthCalendarAlbumStart];        
                toggleElementVisibility(showElements,true);
        
            } else if (type === "end"  || type==="albumEnd") {
                hideElements = [
                    dayCalendarEnd, yearCalendarEnd,
                    dayCalendarAlbumEnd, yearCalendarAlbumEnd
                ];
                toggleElementVisibility(hideElements,false);
                showElements = [monthCalendarEnd, monthCalendarAlbumEnd];
                toggleElementVisibility(showElements,true);
        
            }
            renderMonthCalendar(year,type);
        });
        
         // Toggle calendar visibility (year view)
        yearDisplayDay.addEventListener('click', function () {
            if(type ==="start" || type ==="albumStart"){
                hideElements = [
                    dayCalendarStart, monthCalendarStart, 
                    dayCalendarAlbumStart, monthCalendarAlbumStart,
                ];
                toggleElementVisibility(hideElements,false);
                showElements = [yearCalendarStart,yearCalendarAlbumStart];
                toggleElementVisibility(showElements,true);
                
                startYear = year -28; 
                if (year === maxYear) {
                    startYear === year;
                }
            } else if (type === "end" || type ==="albumEnd") {
                hideElements = [
                    dayCalendarEnd, monthCalendarEnd ,
                    dayCalendarAlbumEnd, monthCalendarAlbumEnd,
                ];
                toggleElementVisibility(hideElements,false);
                showElements = [yearCalendarEnd,yearCalendarAlbumEnd];
                toggleElementVisibility(showElements,true);
         
                startYear = year + 1; 
            }
            renderYearCalendar(startYear,type);
        });
    }

    // Render Month Calendar
    function renderMonthCalendar(year, type) {
         // debug print
        // console.log("renderMonthCalendar type is:",type)
         
        const containerId = `monthCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);

        const yearDisplayMonth = container.querySelector('.calendar-year');
        const monthsContainer = container.querySelector('.cal-months');

        yearDisplayMonth.textContent = `${year}`;

        monthsContainer.innerHTML = months.map((month, index) => {
            return `<div class="month" data-month="${index}">${month}</div>`;
        }).join('');

        // Add click event listeners to months
        monthsContainer.querySelectorAll('.month').forEach(monthElement => {
            monthElement.addEventListener('click', function () {
                const selectedMonth = parseInt(monthElement.dataset.month);
                renderDayCalendar(selectedMonth, year, type);
                document.getElementById(`dayCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`).classList.toggle('hidden', false);
                document.getElementById(`monthCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`).classList.toggle('hidden', true);
            });
        });
        
         // Toggle calendar visibility (year view)
        yearDisplayMonth.addEventListener('click', function () {
            if (type ==="start" || type ==="albumStart"){
                
                hideElements = [
                    dayCalendarStart, monthCalendarStart,
                    dayCalendarAlbumStart, monthCalendarAlbumStart
                ];
                toggleElementVisibility(hideElements,false)
                showElements= [yearCalendarStart, yearCalendarAlbumStart];
                toggleElementVisibility(showElements,true)
                
            } else if (type === "end"  || type==="albumEnd") {
                date
                hideElements = [
                    dayCalendarEnd, monthCalendarEnd,
                    dayCalendarAlbumEnd, monthCalendarAlbumEnd
                ];
                toggleElementVisibility(hideElements,false)
                showElements= [yearCalendarEnd,yearCalendarAlbumEnd];
                toggleElementVisibility(showElements,true)
            }
        
            renderYearCalendar(startYear,type);
        });
    }

    // Render Year Calendar
    function renderYearCalendar(startYear, type) {
        // debug print
        console.log("renderYearCalendar type is:",type)
        
        const containerId = `yearCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);

        const yearsContainer = container.querySelector('.cal-years');
        const endYear = Math.min(startYear + 27, maxYear); // Show a range of 28 years (4 rows x 7 columns)

        yearsContainer.innerHTML = '';
        for (let year = startYear; year <= endYear; year++) {
            const yearElement = document.createElement('div');
            yearElement.textContent = year;
            yearElement.classList.add('year');
            yearElement.addEventListener('click', function () {
                 currentYear = year;
                renderMonthCalendar(year, type);
                document.getElementById(`monthCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`).classList.toggle('hidden', false);
                document.getElementById(`yearCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`).classList.toggle('hidden', true);
            });
            yearsContainer.appendChild(yearElement);
        }
    }
     
    // Elements to check for outside click
    const calendarElements = [
        startDateInput, dayCalendarStart, monthCalendarStart, yearCalendarStart,
        albumStartDateInput, dayCalendarAlbumStart, monthCalendarAlbumStart, yearCalendarAlbumStart,
        endDateInput, dayCalendarEnd, monthCalendarEnd, yearCalendarEnd,
        albumEndDateInput, dayCalendarAlbumEnd, monthCalendarAlbumEnd, yearCalendarAlbumEnd
    ];

    // Close the calendar when clicking outside
    document.addEventListener('click', function (event) {
          // Elements to hide when clicking outside
        hideElements = [
            dayCalendarStart, monthCalendarStart, yearCalendarStart,
            dayCalendarEnd, monthCalendarEnd, yearCalendarEnd,
            dayCalendarAlbumStart, monthCalendarAlbumStart, yearCalendarAlbumStart,
            dayCalendarAlbumEnd, monthCalendarAlbumEnd, yearCalendarAlbumEnd
        ];
        
         // Elements to show when clicking outside
        showElements = [
            artistFilter,albumFilter, membersFilter, concertsFilter, submitFilter, resultContainer
        ];
        
        const isOutsideClick = calendarElements.every(element => !element.contains(event.target));
        if (isOutsideClick) {
            // Hide all calendar elements
            toggleElementVisibility(hideElements, false);

            // Show all filter and results containers
            toggleElementVisibility(showElements, true);
        }
        
        // set all calendar open to false
        isStartCalendarOpen = false;
        isEndCalendarOpen = false;
        isAlbumStartCalendarOpen = false;
        isAlbumEndCalendarOpen = false;
        
    });

    function selectDate(year, month, date, type) {
        // Debug print 
    /*     console.log("**********************")
        console.log("select date for type:", type) */
        
        const selectedDate = new Date(year,month,date);
        // Debug print  
        console.log("selected date from calendarPicker:", selectedDate);
        const currentDate = new Date();
    
        if (selectedDate > currentDate) {
            alert("Cannot select a future date for search.");
            return;
        }
    
        switch (type) {
            case 'start':
                const endDateValue = parseUKDate(endDateInput.value);
    
                if (endDateValue && selectedDate > endDateValue) {
                    alert("Start date cannot be later than end date.");
                    return; 
                }
    
                startDateInput.value = formatDateToUK(selectedDate);
                dayCalendarStart.classList.toggle('hidden', true);
                isStartCalendarOpen = false;
    
                // Debug print 
                console.log("Selected start year:", year);
                console.log("Selected start month:", month);
                break;
    
            case 'end': 
                const startDateValue = parseUKDate(startDateInput.value);
                
                console.log("start date value:",startDateValue)
    
                if (startDateValue && selectedDate < startDateValue) {
                    alert("End date cannot be earlier than start date.");
                    return; 
                }
                endDateInput.value = formatDateToUK(selectedDate);
                dayCalendarEnd.classList.toggle('hidden', true);
                isEndCalendarOpen = false;
    
                // Debug print 
                console.log("Selected end year:", year);
                console.log("Selected end month:", month);
    
                // Trigger change event to notify any listeners of the input change, in order to change the default endDate
                endDateInput.dispatchEvent(new Event('change'));
                break;
    
            case 'albumStart':
                const albumEndDateValue = parseUKDate(albumEndDateInput.value);
    
                if (albumEndDateValue && selectedDate > albumEndDateValue) {
                    alert("Album start date cannot be later than album end date.");
                    return; 
                }
    
                albumStartDateInput.value = formatDateToUK(selectedDate);
                dayCalendarAlbumStart.classList.toggle('hidden', true);
                isAlbumStartCalendarOpen = false;
    
                break;
    
            case 'albumEnd':
                const albumStartDateValue = parseUKDate(albumStartDateInput.value);
                
                if (albumStartDateValue && selectedDate < albumStartDateValue) {
                    alert("Album end date cannot be earlier than album start date.");
                    return; 
                }
    
                albumEndDateInput.value = formatDateToUK(selectedDate);
                dayCalendarAlbumEnd.classList.toggle('hidden', true);
                isAlbumEndCalendarOpen = false; 
    
              // Trigger change event to notify any listeners of the input change, in order to change the default endDate
                albumEndDateInput.dispatchEvent(new Event('change'));
                break;
        }
        
        // Make sure both sides of calender got closed
        if (!isStartCalendarOpen && !isEndCalendarOpen && !isAlbumStartCalendarOpen && !isAlbumEndCalendarOpen ) {
                
            // Show all filters and results container
            toggleElementVisibility(showElements, true);
        }
    }
   
});

// helper functions
 // Function to toggle visibility of elements
 function toggleElementVisibility(elements, isVisible) {
    elements.forEach(element => {
        if (element) {
            element.classList.toggle('hidden', !isVisible);
        }
    });
}

function formatDateToUK(date) {
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    return `${day}-${month}-${year}`;
}

function parseUKDate(dateString) {
    const [day, month, year] = dateString.split('-').map(Number);
    return new Date(year, month - 1, day);
}
    
export {formatDateToUK};