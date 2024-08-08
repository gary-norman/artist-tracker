import {parseDate,formatDate,expandDates,parseDates,generateDatesHTML,wait} from "./map_all.js"

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

  // Create the map instance
  const map = new mapboxgl.Map({
      container: 'map',
      style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8',
      center: [-98.5795, 39.8283],
      zoom: 1
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

  // Load GeoJSON data for the artist
  await loadGeoJSONForArtist(map);
}

async function loadGeoJSONForArtist(map) {
  const artistName = getArtistNameFromURL();
  if (!artistName) {
      console.error('No artist name found in URL.');
      return;
  }

  let artistID;
  try {
      artistID = await fetchArtistID(artistName);
  } catch (error) {
      console.error('Error fetching artist ID:', error);
      return;
  }

  const fileNB = artistID - 1;
  const geoJSONPath = `/db/mapbox_std/${fileNB}.geojson`;

  let geojson;
  try {
      const response = await fetch(geoJSONPath);
      if (!response.ok) throw new Error('Network response was not ok');
      geojson = await response.json();
  } catch (error) {
      console.error('Error fetching GeoJSON data:', error);
      return;
  }

  const aveLoc = geojson.aveCoords;
  console.log('GeoJSON Data:', geojson);
  console.log('Average location:', aveLoc);

  map.setCenter(aveLoc);
  map.setZoom(1);

  geojson.features.forEach((feature, index) => {
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

      el.dataset.index = index;
  });

  document.getElementById('resetMap').addEventListener('click', () => {
      map.flyTo({
          center: aveLoc,
          zoom: 1,
          pitch: 0,
          essential: true
      });
  });

  document.querySelectorAll(".artistPageTourdate").forEach(tourdate => {
      tourdate.addEventListener('click', () => {
          const concertId = tourdate.dataset.index;
          const concert = geojson.features.find(feature => feature.properties.date.includes(concertId));

          if (concert) {
              const coordinates = concert.geometry.coordinates;
              map.flyTo({
                  center: coordinates,
                  zoom: 16,
                  pitch: 60,
                  essential: true
              });
          } else {
              console.log(`Concert data not found for concert ID: ${concertId}`);
          }
      });
  });
}

function setupEventListeners() {
  const globeIcon = document.querySelector('.globe');
  const mapIcon = document.querySelector('.kaart');
  const parent = document.getElementById('mapProject');
  const mapBlack = document.getElementById('map-container');
  let switchButtonText = document.querySelector('#mapProject .button-text');
  let resetButtonText = document.querySelector('#resetMap .button-text');
  const mapControls = document.getElementById('map-controls');
  const mapControlsContainer = document.querySelector('.mapControls');
  let mapProject = 'globe';

  parent.addEventListener('click', async () => {
      mapProject = mapProject === 'globe' ? 'mercator' : 'globe';
      console.log('mapProject:', mapProject);
      mapBlack.classList.add("concerts-map-trans");
      await loadGeoJSONForArtist(map); // Reload map with updated projection
      await wait(1000);
      mapBlack.classList.remove("concerts-map-trans");
  });

  window.addEventListener('resize', () => {
      toggleButtonText();
      toggleControlsView();
  });

  function toggleButtonText() {
      let isGlobe = globeIcon && globeIcon.classList.contains('hide-icon');
      if (window.innerWidth < 380) {
          switchButtonText.textContent = isGlobe ? '2D view' : '3D view';
          resetButtonText.textContent = 'Reset';
      } else if (window.innerWidth < 650) {
          switchButtonText.textContent = isGlobe ? 'Switch to 2D' : 'Switch to 3D';
          resetButtonText.textContent = 'Reset map';
      } else {
          switchButtonText.textContent = isGlobe ? 'Switch to 2D map' : 'Switch to 3D map';
          resetButtonText.textContent = 'Reset map position';
      }
  }

  function toggleControlsView() {
      if (window.innerWidth < 550) {
          mapControls.classList.add('stretch');
          mapControlsContainer.classList.remove('space');
      } else {
          mapControls.classList.remove('stretch');
          mapControlsContainer.classList.add('space');
      }
  }
}

// Utility functions
function getArtistNameFromURL() {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.has('name') ? decodeURIComponent(urlParams.get('name')) : null;
}

async function fetchArtistID(artistName) {
  try {
      const response = await fetch(`/artist/id?name=${encodeURIComponent(artistName)}`);
      if (!response.ok) throw new Error('Failed to fetch artist ID');
      const data = await response.json();
      console.log('Artist ID:', data.artistId);
      return data.artistId;
  } catch (error) {
      console.error('Error fetching artist ID:', error);
      throw error;
  }
}