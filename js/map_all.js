document.addEventListener('DOMContentLoaded', async () => {
 
  try {
    await initializeMap();
    setupEventListeners();
  } catch (error) {
    console.error('Error during initialization:', error);
  }
  
});

async function initializeMap() {
     // Fetch and set the Mapbox access token
    try {
        const response = await fetch('/map_token');
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();
        mapboxgl.accessToken = data.MAPBOXGL_ACCESS_TOKEN;
    } catch (error) {
        console.error('Error fetching Mapbox token:', error);
        return; 
    }
    
    // Map container
    let mapBlack = document.getElementById("map-container");

    let mapProject = 'globe';
    
    const map = new mapboxgl.Map({
      container: 'map',
      style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8',
      center: [-98.5795, 39.8283],
      zoom: 1,
      projection: mapProject
    });
  
    map.on('style.load', () => {
      map.setFog({
        "range": [0.8, 9],
        "color": "#471e50",
        "horizon-blend": 0.2,
        "high-color": "#07173e",
        "space-color": "#000000",
        "star-intensity": 0.35
      });
    });
  
    await loadGeoJSONForAllLocations(map);
  
    // reset map position
    document.getElementById('resetMap').addEventListener('click', () => {
      map.flyTo({
        center: [-98.5795, 39.8283],
        zoom: 1,
        pitch: 0,
        essential: true // animation considered essential for accessibility
      });
    });
    
    // Toggle map projection type
    document.getElementById('mapProject').addEventListener('click', async () => {
        if (mapProject === 'globe') {
            mapProject = 'mercator';
        } else {
            mapProject = 'globe';
        }
        console.log('mapProject:' + mapProject);
        mapBlack.classList.add("concerts-map-trans");
        loadGeoJSONForAllLocations().then(r => { });
        await wait(1000);
        mapBlack.classList.remove("concerts-map-trans");
    });
  }


// Load GeoJSON data for all locations
async function loadGeoJSONForAllLocations(map) {
  const geojsonFiles = [];
  const numFiles = 51;

  for (let i = 0; i < numFiles; i++) {
      geojsonFiles.push(`/db/mapbox_std/${i}.geojson`);
  }

  try {
      const geojsons = await Promise.all(geojsonFiles.map(file => fetch(file).then(response => {
          if (!response.ok) {
              throw new Error('Network response was not ok ' + response.statusText);
          }
          return response.json();
      })));

      geojsons.forEach(geojson => {
          geojson.features.forEach(feature => {
              const el = document.createElement('div');
              el.className = 'marker';

              const datesHTML = generateDatesHTML(parseDates(expandDates(feature.properties.date)));
              const addressHTML = parseDates(expandDates(feature.properties.date)).length < 2
                  ? `<p class="small">${feature.properties.eventAddress}</p>`
                  : '';

              const title = feature.properties.title.replace(" at ", " in ");
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

              el.dataset.index = feature.properties.date;
          });
      });
      
      // Select all tour date elements
      const mapClick = document.querySelectorAll(".artistPageTourdate");
      mapClick.forEach(tourdate => {
          tourdate.addEventListener('click', () => {
              const concertId = tourdate.dataset.index;

              const concert = geojsons.flatMap(geojson => geojson.features)
                  .find(feature => feature.properties.date.includes(concertId));

              if (concert) {
                  const coordinates = concert.geometry.coordinates;
                  
                  console.log("Coordinates:", coordinates);
                  
                  map.flyTo({
                      center: coordinates,
                      zoom: 16,
                      pitch: 60,
                      essential: true  // animation considered essential for accessibility
                  });
              } else {
                  console.log(`Concert data not found for concert ID: ${concertId}`);
              }
          });
      });

  } catch (error) {
      console.error('There has been a problem with your fetch operation:', error);
  }
}

function setupEventListeners() {
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
      if (globeIcon.classList.contains('hide-icon')) {
        globeIcon.classList.remove('hide-icon');
        mapIcon.classList.add('hide-icon');
        toggleButtonText();
      } else {
        mapIcon.classList.remove('hide-icon');
        globeIcon.classList.add('hide-icon');
        toggleButtonText();
      }
      // Reinitialize the map with the new projection type
      await initializeMap();
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
}

// Helper functions for date handling
function parseDate(dateStr) {
    const [day, month, year] = dateStr.split('-');
    return new Date(`${year}-${month}-${day}`);
}

function formatDate(date) {
    const options = { day: '2-digit', month: 'short', year: 'numeric' };
    return date.toLocaleDateString('en-GB', options);
}

function expandDates(dateString) {
    let tags = '';
    const dateArray = dateString.split(", ");
    for (const i of dateArray) {
        tags += formatDate(parseDate(i)) + ", ";
    }
    tags = tags.replace(/,\s*$/, '');
    return tags;
}

function parseDates(dateString) {
    return dateString.split(',').map(date => date.trim());
}

// Function to generate HTML for all dates
function generateDatesHTML(dates) {
    return dates.map(singleDate => `
        <p class="pic date">${singleDate}</p>
    `).join('');
}

// This function returns a promise that resolves after n milliseconds
function wait(n) {
    return new Promise((resolve) => setTimeout(resolve, n));
}

export {parseDate,formatDate,expandDates,parseDates,generateDatesHTML,wait}