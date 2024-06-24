// Ensure to add your Mapbox access token
mapboxgl.accessToken = 'pk.eyJ1IjoibG9yZXdvcmxkIiwiYSI6ImNsd3FseDNsbDAzZjMyanF2czh3Mmt4eTgifQ.-_bXsAv_SR1bpcmvOSpDuA';

// Create the map instance
const map = new mapboxgl.Map({
    container: 'map', // container ID
    style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8', // map style
    center: [-98.5795, 39.8283], // starting position [lng, lat]
    zoom: 3 // starting zoom
});

map.on('style.load', () => {
    map.setFog({
        "range": [0.8, 9],
        "color": "#471e50",
        "horizon-blend": .2,
        "high-color": "#07173e",
        "space-color": "#000000",
        "star-intensity": 0.35
    }); //
});

function parseDate(dateStr) {
    const [day, month, year] = dateStr.split('-');
    return new Date(`${year}-${month}-${day}`);
}

function formatDate(date) {
    const options = { day: '2-digit', month: 'short', year: 'numeric' };
    return date.toLocaleDateString('en-GB', options);
}

// Generate GeoJSON file paths dynamically
const geojsonFiles = [];
const numFiles = 51; // Number of GeoJSON files

for (let i = 0; i < numFiles; i++) {
    geojsonFiles.push(`/db/mapbox_std/${i}.geojson`);
}

// Fetch all GeoJSON files and process them
Promise.all(geojsonFiles.map(file => fetch(file)
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    })
))
.then(geojsons => {
    geojsons.forEach(geojson => {
        console.log('GeoJSON Data:', geojson);

        // add markers to map
        for (const feature of geojson.features) {
            const el = document.createElement('div');
            el.className = 'marker';

            new mapboxgl.Marker(el)
                .setLngLat(feature.geometry.coordinates)
                .setPopup(
                    new mapboxgl.Popup({ offset: 20 })
                        .setHTML(
                            `<div class="p--bold flexrow"><a href="/artist?name=${feature.properties.artist}">${feature.properties.artist}</a>
                             <b>live in ${feature.properties.eventAddress}</b></div>
                             <p class="small">${formatDate(parseDate(feature.properties.date))}</p>
                             <p class="small">${feature.properties.eventAddress}</p>`
                        )
                )
                .addTo(map);
        }
    });
})
.catch(error => {
    console.error('There has been a problem with your fetch operation:', error);
});
