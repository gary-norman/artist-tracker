document.addEventListener('DOMContentLoaded', loadGeoJSONForArtist);
mapBlack = document.getElementById("map-container")


mapboxgl.accessToken = 'pk.eyJ1IjoibG9yZXdvcmxkIiwiYSI6ImNsemlqbHpreDBmem0ya3IxcjBkMDltZGYifQ.Pt0OW2VelQK5-o1EDti5aQ'
let mapProject = 'globe'
// toggle map projection type
document.getElementById('mapProject').addEventListener('click', async() => {
    if (mapProject === 'globe' ) {
        mapProject = 'mercator'
    } else {
        mapProject = 'globe'
    }
    console.log('mapProject:' + mapProject);
    mapBlack.classList.add("concerts-map-trans");
    loadGeoJSONForArtist().then(r => {});
    await wait(1000)
    mapBlack.classList.remove("concerts-map-trans");
})

// this function returns a promise that resolves after n milliseconds
const wait = (n) => new Promise((resolve) => setTimeout(resolve, n));


function parseDate(dateStr) {
    const [day, month, year] = dateStr.split('-');
    return new Date(`${year}-${month}-${day}`);
}

function parseDates(dateString) {
    // Split the input string by commas and trim any extra whitespace from each date and return the array
    return dateString.split(',').map(date => date.trim());
}

function formatDate(date) {
    const options = { day: '2-digit', month: 'short', year: 'numeric' };
    return date.toLocaleDateString('en-GB', options);
}

function expandDates(dateString) {
    let tags = '';
    const dateArray = dateString.split(", ");
    for (const i of dateArray) {
        tags += formatDate(parseDate(i)) + ", "
    }
    tags = tags.replace(/,\s*$/, '');
    return tags;
}

// Function to generate HTML for all dates
function generateDatesHTML(dates) {
    return dates.map(singleDate => `
    <p class="pic date">${singleDate}</p>
  `).join('');
}

// Function to load GeoJSON data for the artist based on artist name in URL
async function loadGeoJSONForArtist() {
    const artistName = getArtistNameFromURL();
    if (!artistName) {
        console.error('No artist name found in URL.');
        return;
    }

    try {
        // Fetch artist ID based on artist name
        const artistID = await fetchArtistID(artistName);
        if (!artistID) {
            console.error('No artist ID found for artist name:', artistName);
            return;
        }

        // Determine GeoJSON file path based on artist ID
       /* let geoJSONPath;
        const extras = [0, 2, 4, 5, 6, 8, 11, 13, 14, 15, 16, 18, 19, 20, 21, 23, 26, 32, 34, 35, 40, 42, 48];
        const fileNB = artistID - 1;
        if (extras.includes(fileNB)) {
            geoJSONPath = `/db/mapbox/${fileNB}.geojson`;
        } else {
            geoJSONPath = `/db/mapbox_std/${fileNB}.geojson`;
        }*/

        const fileNB = artistID - 1;
        const geoJSONPath = `/db/mapbox_std/${fileNB}.geojson`;

        // Fetch GeoJSON data
        const response = await fetch(geoJSONPath);
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        const geojson = await response.json();
        const aveLoc = geojson.aveCoords
        console.log('GeoJSON Data:', geojson);
        console.log('Average location:', aveLoc)

        // Create the map instance
        const map = new mapboxgl.Map({
            container: 'map', // container ID
            style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8', // map style
            projection: mapProject,
            // center: [-98.5795, 39.8283], // starting position [lng, lat]
            center: aveLoc, // starting position [lng, lat]
            zoom: 1 // starting zoom
        });

        map.on('style.load', () => {
            map.setFog({
                "range": [0.8, 9],
                "color": "#471e50",
                "horizon-blend": .2,
                "high-color": "#07173e",
                "space-color": "#000000",
                "star-intensity": 0.35
            });
        });

        // Assuming `geojson` is your GeoJSON data object
        geojson.features.forEach((feature, index) => {
            const el = document.createElement('div');
            el.className = 'marker';

            // Generate the HTML for the dates
            const datesHTML = generateDatesHTML(parseDates(expandDates(feature.properties.date)));

            // Conditionally include the address if there is more than one date
            const addressHTML = parseDates(expandDates(feature.properties.date)).length < 2
                ? `<p class="small">${feature.properties.eventAddress}</p>`
                : '';

            const title = feature.properties.title.replace(" at ", " in ")
            // Create Mapbox Marker for each feature
            new mapboxgl.Marker(el)
                .setLngLat(feature.geometry.coordinates)
                .setPopup(
                    new mapboxgl.Popup({ offset: 20 })
                        .setHTML(`
          <p class="p--bold">${title}</p>
          <div class="content go-across-md scroll">
            ${datesHTML}
            ${addressHTML}
          </div>
        `)
                )
                .addTo(map);

            // Store the index as a data attribute on the marker element for click events
            el.dataset.index = index;
        });

        // reset map position
        // document.getElementById('kaartReset').addEventListener('click', () => {
        document.getElementById('resetMap').addEventListener('click', async() => {
            map.flyTo({
                center: aveLoc,
                zoom: 1,
                pitch: 0,
                essential: true // animation considered essential for accessibility
            });
        })

        // Select all tour date elements
        const mapClick = document.querySelectorAll(".artistPageTourdate");
        // Add click event listener to each tour date element
        mapClick.forEach(tourdate => {
            tourdate.addEventListener('click', () => {
                // Get the concert ID (date string) from the clicked tourdate element
                const concertId = tourdate.dataset.index;
                
                // debug print
                console.log("Concert ID:", concertId);

                // Ensure geojson is defined and features exist
                if (geojson && geojson.features) {
                    // Find the concert data from `geojson` based on concertId
                    const concert = geojson.features.find(feature => feature.properties.date.includes(concertId));

                    if (concert) {
                        
                        // debug print
                        console.log("got concert");
                        console.log("Feature Date:", concert.properties.date);
                        
                        const coordinates = concert.geometry.coordinates;
                        
                        // debug print
                        console.log("Coordinates:", coordinates);

                        // Fly to the coordinates of the feature with the fetched concert ID
                        map.flyTo({
                            center: coordinates,
                            zoom: 16,
                            pitch: 60,
                            essential: true // animation considered essential for accessibility
                        });
                    } else {
                        console.log(`Concert data not found for concert ID: ${concertId}`);
                    }
                } else {
                    console.log("GeoJSON data not loaded or features are missing.");
                }
            });
        });

    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

function getArtistNameFromURL() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.has('name') ? decodeURIComponent(urlParams.get('name')) : null;
}

async function fetchArtistID(artistName) {
    
    try {
        const response = await fetch(`/artist/id?name=${encodeURIComponent(artistName)}`);
        if (!response.ok) {
            throw new Error('Failed to fetch artist ID');
        }
        const data = await response.json();
        console.log('Artist ID:', data.artistId);
        return data.artistId; 
    } catch (error) {
        console.error('Error fetching artist ID:', error);
        throw error;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const globeIcon = document.querySelector('.globe');
    const mapIcon = document.querySelector('.kaart');
    const parent = document.getElementById('mapProject');
    let switchButtonText = document.querySelector('#mapProject .button-text');
    let resetButtonText = document.querySelector('#resetMap .button-text');
    const mapControls = document.getElementById('map-controls');
    const mapControlsContainer = document.querySelector('.mapControls');

    parent.addEventListener('click', () => {
        toggleMapView();
    });

    //initialise button text on load
    toggleButtonText();
    toggleControlsView();

    window.addEventListener('resize', () => {
        toggleButtonText();
        toggleControlsView();
    });


    async function toggleMapView() {
        if (globeIcon.classList.contains('hide-icon')){
            globeIcon.classList.remove('hide-icon');
            mapIcon.classList.add('hide-icon');
            toggleButtonText();

        } else {
            mapIcon.classList.remove('hide-icon');
            globeIcon.classList.add('hide-icon');
            toggleButtonText();
        }
    }

    function toggleButtonText() {
        let isGlobe = false;

        if (globeIcon.classList.contains('hide-icon')) {
            isGlobe = true;
        }
        if (window.innerWidth < 380) {
            if (isGlobe) {
                switchButtonText.textContent = '2D view';
            } else {
                switchButtonText.textContent = '3D view';
            }
            resetButtonText.textContent = 'Reset';
        } else if (window.innerWidth < 650) {
            if (isGlobe) {
                switchButtonText.textContent = 'Switch to 2D';
            } else {
                switchButtonText.textContent = 'Switch to 3D';
            }
            resetButtonText.textContent = 'Reset map';
        } else {
            if (isGlobe) {
                switchButtonText.textContent = 'Switch to 2D map';
            } else {
                switchButtonText.textContent = 'Switch to 3D map';
            }
            resetButtonText.textContent = 'Reset map position';
        }
    }

    function toggleControlsView(){
        if (window.innerWidth < 550) {
            mapControls.classList.add('stretch');
            mapControlsContainer.classList.remove('space');
        } else {
            mapControls.classList.remove('stretch');
            mapControlsContainer.classList.add('space');
        }
    }
});