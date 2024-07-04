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
    //let startYear = Math.floor(currentYear / 28) * 28; 
    let startYear = currentYear -28; 
    let isYearDisabled = false;
    let isMonthDisabled = false;
    let isDayDisabled = false;
    let disableDate = null;
    let disableYear
    let disableMonth 
    let disableDay 

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
        
        // set status
        isStartCalendarOpen = true;
        isEndCalendarOpen = false
        let currentDay
        isDayDisabled = false;
        
        // check if startDate got input
        const startDateValue = startDateInput.value ? parseUKDate(startDateInput.value) : new Date();
        // check first if user select a date, if selected then convert the UKformat back to normal
        const endDateValue = endDateInput.value ? parseUKDate(endDateInput.value) : new Date();
       
        // if user input inside startDate then take that one
        if (startDateInput.value){
            currentYear = startDateValue.getFullYear();
            currentMonth = startDateValue.getMonth();
            currentDay = startDateValue.getDate();
        } else if (endDateInput.value){
            currentYear = endDateValue.getFullYear();
            currentMonth = endDateValue.getMonth();
            currentDay = endDateValue.getDate();
        }

        // If user has selected an end date, use that date
        if (endDateInput.value) {
            disableDate = endDateValue;
           
        } else{
            disableDate = new Date();
           
        }
        disableYear = disableDate.getFullYear();
        disableMonth = disableDate.getMonth();
        disableDay = disableDate.getDate();
       
        //debug print
        console.log("disable date is--------->",disableDate);
        console.log("disable Year is--------->",disableYear);
        console.log("disable Month is--------->",disableMonth);
        console.log("disable Day is--------->",disableDay);
        console.log("current Year is--------->",currentYear);
        console.log("current Month is--------->",currentMonth);
        console.log("current Day is--------->",currentDay);
        
        // Check if the start date is the first day of the month
        if (currentDay === 1) {
            currentMonth--;
            if (currentMonth < 0) { 
                currentMonth = 11; // Dec
                currentYear--;
            }
        }
        
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
        const endDateValue = endDateInput.value ? parseUKDate(endDateInput.value) : new Date();
        currentYear = endDateValue.getFullYear();
        currentMonth = endDateValue.getMonth();
        const currentDay = endDateValue.getDate();
        isDayDisabled = false;
        
        console.log("start date value is ====>>>>",startDateValue)
        
        // Set disableDate to the current date by default
        // if user input anything
        if (startDateInput.value){
            console.log("yes has value!!!!!!!!!!!!!!!!!!!!??????????????")
            disableDate = new Date();
            disableDate = startDateValue;
            disableYear = disableDate.getFullYear();
            disableMonth = disableDate.getMonth();
            disableDay = disableDate.getDate();
        } else{
            disableDate = null;
        }
       
        //debug print
        console.log("disable date is--------->",disableDate);
        console.log("disable Year is--------->",disableYear);
        console.log("disable Month is--------->",disableMonth);
        console.log("disable Day is--------->",disableDay);
        console.log("current Year is--------->",currentYear);
        console.log("current Month is--------->",currentMonth);
        console.log("current Day is--------->",currentDay);
        
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
        
        // toggle other calendars all hidden 
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
  
        // set status
        isAlbumStartCalendarOpen = true;
        isAlbumEndCalendarOpen = false
        let currentDay
        isDayDisabled = false;
        // in order to store currentDate for album
        let albumDate = new Date();
        let currentAlbumDay = albumDate.getDate();
        let currentAlbumMonth = albumDate.getMonth();
        let currentAlbumYear = albumDate.getFullYear();
        
        // check if startDate got input
        const albumStartDateValue = albumStartDateInput.value ? parseUKDate(albumStartDateInput.value) : new Date();
        
        // check first if user select a date, if selected then convert the UKformat back to normal
        const albumEndDateValue = albumEndDateInput.value ? parseUKDate(albumEndDateInput.value) : new Date();
        
        // if user input inside startDate then take that one
        if (albumStartDateInput.value){
            currentAlbumYear = albumStartDateValue.getFullYear();
            currentAlbumMonth = albumStartDateValue.getMonth();
            currentAlbumDay = albumStartDateValue.getDate();
            // asign currentDate back as global
            currentYear = currentAlbumYear;
            currentMonth = currentAlbumMonth;
        } else if (albumEndDateValue){
            currentAlbumYear = albumEndDateValue.getFullYear();
            currentAlbumMonth = albumEndDateValue.getMonth();
            currentAlbumDay = albumEndDateValue.getDate();
            currentYear = currentAlbumYear;
            currentMonth = currentAlbumMonth;
        }
      
        
      /*   currentYear = albumEndDateValue.getFullYear();
        currentMonth = albumEndDateValue.getMonth();
        let currentDay = albumEndDateValue.getDate(); */
         
         // If user has selected an end date, use that date
         if (albumEndDateInput.value) {
             disableDate = albumEndDateValue;
            
         } else {
            disableDate = new Date();
         }
         disableYear = disableDate.getFullYear();
         disableMonth = disableDate.getMonth();
         disableDay = disableDate.getDate();
      
         //debug print
         console.log("disable date is--------->",disableDate);
         console.log("disable Year is--------->",disableYear);
         console.log("disable Month is--------->",disableMonth);
         console.log("disable Day is--------->",disableDay);
         console.log("current Year is--------->",currentYear);
         console.log("current Month is--------->",currentMonth);
         console.log("current Day is--------->",currentDay);
        
        // Check if the start date is the first day of the month
    if (currentAlbumDay === 1) {
        currentAlbumMonth--;
        if (currentAlbumMonth < 0) { 
            currentAlbumMonth = 11; // Dec
            currentAlbumYear--;
        }
    }

    // Debug print
/*     console.log("year now is:", currentAlbumYear);
    console.log("month now after formatting is:", currentAlbumMonth); */
         
        renderDayCalendar(currentAlbumMonth, currentAlbumYear, "albumStart");
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
        const albumEndDateValue = albumEndDateInput.value ? parseUKDate(albumEndDateInput.value) : new Date();
        currentYear = albumEndDateValue.getFullYear();
        currentMonth = albumEndDateValue.getMonth();
        const currentDay =  albumEndDateValue.getDate();
        // if user input anything
        if (albumStartDateInput.value){
            disableDate = new Date();
            disableDate = albumStartDateValue;
            disableYear = disableDate.getFullYear();
            disableMonth = disableDate.getMonth();
            disableDay = disableDate.getDate();
        }else{
            disableDate = null;
        }
      
        isDayDisabled = false;
        //debug print
        console.log("disable date is--------->",disableDate);
        console.log("disable Year is--------->",disableYear);
        console.log("disable Month is--------->",disableMonth);
        console.log("disable Day is--------->",disableDay);
        
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
                } else if (currentYear === disableYear){
                   if (type ==='end'|| type ==='albumEnd'){
                        if (currentMonth < disableMonth){
                            currentMonth = disableMonth ;
                        }
                    }
                }                    
                renderDayCalendar(currentMonth, currentYear, type);
            });
            forwardYearButton.addEventListener('click', function() {
                currentYear++;
                if (currentYear > maxYear) {
                    currentYear = maxYear;
                } else if (currentYear === disableYear){
                    if (type ==='start'|| type ==='albumStart'){
                        if (currentMonth > disableMonth){
                            currentMonth = disableMonth ;
                        }
                    }
                
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
                if (currentYear > maxYear) {
                    currentYear = maxYear;
                }
                renderMonthCalendar(currentYear, type);
            });
        } else if (viewType === 'year') {
            backButton.addEventListener('click', function() {
                if (currentYear === disableYear+1) {
                    currentYear = disableYear;
                } else{
                    currentYear -= 28;
                }
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
        
        // reset every time when start
        isDayDisabled = false;
        isYearDisabled = false;
        isMonthDisabled = false;
      
        // debug print
        console.log("----> ** renderDayCalendar type is:",type)
        console.log("----> ** month got parsed:",month)
        console.log("----> ** year got parsed:",year) 
        console.log("----> ** disable Year is:",disableYear);
        console.log("----> ** disable month is:",disableMonth);
        console.log("----> ** disable day is:",disableDay);
        console.log("----> ** current Year is:",currentYear);
        console.log("----> ** current month is:",currentMonth); 
        
        const containerId = `dayCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);
        
        const yearDisplayDay = container.querySelector('.calendar-year');
        const monthDisplayDay = container.querySelector('.calendar-month');
        const daysContainer = container.querySelector('.cal-days');
        const datesContainer = container.querySelector('.calendar-main');
        
         // Check conditions to disable year or month
        if ((type === 'start' || type === 'albumStart') && (year >= disableYear ) ||
            (type === 'end' || type === 'albumEnd') && (year <= disableYear )) {
                isYearDisabled = true;
        }
        console.log("************* is year disable? **********",isYearDisabled)
        
         // If the year is disabled, add appropriate class and disable interaction
        const backYearButton = container.querySelector('.cal-btn.back-year');
        const forwardYearButton = container.querySelector('.cal-btn.front-year');

       // If the year is disabled, add appropriate class and disable interaction
        if (isYearDisabled) {
            // Disable interaction with year navigation buttons
            if (type ==='start' || type === 'albumStart'){
                if (year >= disableYear){
                    backYearButton.disabled = false;
                    forwardYearButton.disabled = true;
                } else{
                    forwardYearButton.disabled = true;
                    backYearButton.disabled = false;
                }
            } else if (type ==='end' || type ==='albumEnd'){
                if ( year <= disableYear){
                    backYearButton.disabled = true;
                    forwardYearButton.disabled = false;
                } else{
                    isYearDisabled = false;
                    backYearButton.disabled = false;
                    forwardYearButton.disabled = false;
                }
            }
        } else {
            backYearButton.disabled = false;
            forwardYearButton.disabled = false;
        } 

       yearDisplayDay.textContent = `${year}`;
       monthDisplayDay.textContent = `${months[month]}`;
        
       // month checking only if the year is same as disabledYear
        if (year === disableYear) {
            if ((type === 'start' || type === 'albumStart') && (month >= disableMonth) ||
                (type === 'end' || type === 'albumEnd') && (month <= disableMonth)) {
                isMonthDisabled = true;
            }
        }
        
        console.log("************* is month disable? **********",isMonthDisabled)
        
        const backMonthButton = container.querySelector('.cal-btn.back-month');
        const forwardMonthButton = container.querySelector('.cal-btn.front-month');
        
            // If the month is disabled, add appropriate class and disable interaction
        if (isMonthDisabled) {
            if (type ==='start' || type === 'albumStart'){
                backMonthButton.disabled = false;
                forwardMonthButton.disabled = true;
            } else if (type ==='end'|| type === 'albumEnd'){
                backMonthButton.disabled = true;
                forwardMonthButton.disabled = false;
            }
        } else {
            backMonthButton.disabled = false;
            forwardMonthButton.disabled = false;
        } 

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
            
          // Reset isDayDisabled for each date
            isDayDisabled = false;
            let checkDay = new Date(year, month, date);
            
            // debug print
             console.log("--> starting adding date inside calendar");
             console.log("--> parsing date :",date);
             console.log("--> disable date :",disableDay);
          /*    console.log("--> check date :",checkDay);
             console.log("--> disable date:", disableDate);  */
            
            const dateElement = document.createElement('div');
            if ((type === 'start' || type === 'albumStart') && checkDay >= disableDate ||
                (type === 'end' || type === 'albumEnd') && checkDay < disableDate ||
                checkDay > new Date()) {
                isDayDisabled = true;
            }
            console.log("************* is date disable? **********",isDayDisabled)
            
            if (isDayDisabled) {
                dateElement.textContent = date;
                dateElement.classList.add('disableDay');
            } else {
                dateElement.textContent = date;
                dateElement.classList.add('day'); // Only add 'day' class
            }

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
            if (type ==="start" || type ==="albumStart"){
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
              /*   console.log("************************")
                console.log("year is:",year);
                console.log("start year is:",startYear);
                console.log("max year is:",maxYear);
                console.log("************************") */
                if (year === maxYear) {
                    startYear === year;
                } else {
                    startYear = year + 1; 
                }
            }
            renderYearCalendar(startYear,type);
        });
    }

    // Render Month Calendar
    function renderMonthCalendar(year, type) {
          // reset
          isYearDisabled = false;
          
         // debug print
         console.log("----> ** renderMonthCalendar type is:",type)
         console.log("----> ** year got parsed:",year) 
         console.log("----> ** disable Year is:",disableYear);
         console.log("----> ** disable month is:",disableMonth);
         console.log("----> ** disable day is:",disableDay);
        
         console.log("----> ** current Year is:",currentYear);
         console.log("----> ** current month is:",currentMonth); 
        
        const containerId = `monthCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);

        const yearDisplayMonth = container.querySelector('.calendar-year');
        const monthsContainer = container.querySelector('.cal-months');
        
      
         // Check conditions to disable year or month
         if ((type === 'start' || type === 'albumStart') && (year >= disableYear ) ||
             (type === 'end' || type === 'albumEnd') && (year <= disableYear )) {
             isYearDisabled = true;
        }
        
        const backButton = container.querySelector('.cal-btn.back');
        const forwardButton = container.querySelector('.cal-btn.front');
            
        if (isYearDisabled){
            if  (type ==="start" || type ==="albumStart"){
                backButton.disabled = false;
                forwardButton.disabled = true;
            } else if (type === "end" || type ==="albumEnd") {
                backButton.disabled = true;
                forwardButton.disabled = false;
          }
        } else {
            backButton.disabled = false;
            forwardButton.disabled = false;
        }

        yearDisplayMonth.textContent = `${year}`;

        monthsContainer.innerHTML = months.map((month, index) => {
            
                // Reset isMonthDisabled for each month
                isMonthDisabled = false;
            
            if (year === disableYear){
                if ((type === 'start' || type === 'albumStart') && (index > disableMonth) ||
                    (type === 'end' || type === 'albumEnd') && (index < disableMonth)) {
                        isMonthDisabled = true;
                }
            }
            
            // Add disableMonth class if month is disabled
            const monthClass = isMonthDisabled ? 'month disableMonth' : 'month';

            return `<div class="${monthClass}" data-month="${index}">${month}</div>`;
        }).join('');

        // Add click event listeners to months
        monthsContainer.querySelectorAll('.month').forEach(monthElement => {
            monthElement.addEventListener('click', function () {
                const selectedMonth = parseInt(monthElement.dataset.month);
                // set currentMonth after click
                currentMonth = selectedMonth;
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
        
       // reset every time when start
         let isYearDisabled = false;
         
        // debug print
         console.log("----> ** renderYearCalendar type is:",type)
         console.log("----> ** startyear got parsed:",startYear) 
         console.log("----> ** disable Year is:",disableYear);
         console.log("----> ** disable month is:",disableMonth);
         console.log("----> ** disable day is:",disableDay);
         console.log("----> ** current Year is:",currentYear);
         console.log("----> ** current month is:",currentMonth); 
         
        const containerId = `yearCalendar${type.charAt(0).toUpperCase() + type.slice(1)}`;
        const container = document.getElementById(containerId);

        const yearsContainer = container.querySelector('.cal-years');
        const endYear = Math.min(startYear + 27, maxYear); // Show a range of 28 years (4 rows x 7 columns)
       
         // Check conditions to disable year or month
         if ((type === 'start' || type === 'albumStart') && (startYear >= disableYear ) ||
             (type === 'end' || type === 'albumEnd') && (startYear <= disableYear )) {
             isYearDisabled = true;
        }
        console.log("************* is year disable? **********",isYearDisabled)
        
        const backButton = container.querySelector('.cal-btn.back');
        const forwardButton = container.querySelector('.cal-btn.front');
            
         if (isYearDisabled){
            if  (type ==="start" || type ==="albumStart"){
                backButton.disabled = false;
                forwardButton.disabled = true;
            } else if (type === "end" || type ==="albumEnd") {
                backButton.disabled = true;
                forwardButton.disabled = false;
          }
        } else {
            backButton.disabled = false;
            forwardButton.disabled = false;
        }  
        
        yearsContainer.innerHTML = '';
        for (let year = startYear; year <= endYear; year++) {
            const yearElement = document.createElement('div');
            yearElement.textContent = year;
            yearElement.classList.add('year');
            
            // Add the disableYear class if the year should be disabled
            if ((type === 'start' || type === 'albumStart') && year > disableYear) {
                yearElement.classList.add('disableYear');
            } else if ((type === 'end' || type === 'albumEnd') && year < disableYear) {
                yearElement.classList.add('disableYear');
            }
        
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
    
                startDateInput.value = formatDateToUK(selectedDate);
                dayCalendarStart.classList.toggle('hidden', true);
                isStartCalendarOpen = false;
    
                // Debug print 
                console.log("Selected start year:", year);
                console.log("Selected start month:", month);
                break;
    
            case 'end': 
             
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
            
                albumStartDateInput.value = formatDateToUK(selectedDate);
                dayCalendarAlbumStart.classList.toggle('hidden', true);
                isAlbumStartCalendarOpen = false;
    
                break;
    
            case 'albumEnd':
    
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