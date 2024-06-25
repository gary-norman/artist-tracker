// Ensure to add your Mapbox access token
mapboxgl.accessToken = 'pk.eyJ1IjoibG9yZXdvcmxkIiwiYSI6ImNsd3FseDNsbDAzZjMyanF2czh3Mmt4eTgifQ.-_bXsAv_SR1bpcmvOSpDuA';

// Create the map instance
const map = new mapboxgl.Map({
    container: 'map', // container ID
    style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8', // map style
    center: [-98.5795, 39.8283], // starting position [lng, lat]
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

function expandDates(dateString) {
    let tags = '';
    const dateArray = dateString.split(", ");
    for (const i of dateArray) {
        tags += formatDate(parseDate(i)) + ", "
    }
    tags = tags.replace(/,\s*$/, '');
    return tags;
}


async function loadGeoJSONForArtist() {
    const artistName = getArtistNameFromURL();
    if (!artistName) {
        console.error('No artist name found in URL.');
        return;
    }

    
    // Fetch artist ID based on artist name
    const artistID = await fetchArtistID(artistName);
    if (!artistID) {
        console.error('No artist ID found for artist name:', artistName);
        return;
    }

    const extras = [0,2,4,5,6,8,11,13,14,15,16,18,19,20,21,23,26,32,34,35,40,42,48]

    let geoJSONPath
    const fileNB = artistID - 1;
    if (extras.includes(fileNB)) {
        geoJSONPath = `/db/mapbox/${fileNB}.geojson`;
    } else {
        geoJSONPath = `/db/mapbox_std/${fileNB}.geojson`;
    }

    try {
        const response = await fetch(geoJSONPath);
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }

        const geojson = await response.json();
        console.log('GeoJSON Data:', geojson);

        // add markers to map
        for (const feature of geojson.features) {
            // create an HTML element for each feature
            const el = document.createElement('div');
            el.className = 'marker';

            // make a marker for each feature and add it to the map
            new mapboxgl.Marker(el)
                .setLngLat(feature.geometry.coordinates)
                .setPopup(
                    new mapboxgl.Popup({ offset: 20 }) // add popups
                        .setHTML(
                            `<p class="p--bold">${feature.properties.title}</p>
                             <p class="small justify">${expandDates(feature.properties.date)}</p>
                             <p class="small">${feature.properties.eventAddress}</p>`
                        )
                )
                .addTo(map);
        }

    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

let mapClick = document.querySelectorAll(".artistPageTourdate");
mapClick.forEach(tourdate => {
    tourdate.addEventListener('click', () => map.flyTo({
        center: [-73.99156, 40.74971],
        zoom: 16,
        pitch: 60,
        essential: true // this animation is considered essential with respect to prefers-reduced-motion
    }));
});

// // document.querySelectorAll('[id^="location"]')
//     document.querySelectorAll('[class="artistPageTourdate"]').addEventListener('click', () => {
//     // document.getElementById('tempclick').addEventListener('click', () => {
//     // Fly to a random location
//     map.flyTo({
//         center: [-73.99156, 40.74971],
//         zoom: 16,
//         pitch: 60,
//         essential: true // this animation is considered essential with respect to prefers-reduced-motion
//     });
// });
    


// Gary
// Fetch GeoJSON data from a local file
/* fetch('/db/mapbox_std/48.geojson')
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(geojson => {
        console.log('GeoJSON Data:', geojson);

        // add markers to map
        for (const feature of geojson.features) {
            // create an HTML element for each feature
            const el = document.createElement('div');
            el.className = 'marker';

            // make a marker for each feature and add it to the map
            new mapboxgl.Marker(el)
                .setLngLat(feature.geometry.coordinates)
                .setPopup(
                    new mapboxgl.Popup({ offset: 20 }) // add popups
                        .setHTML(
                            `<p class="p--bold">${feature.properties.title}</p>
                             <p class="small">${formatDate(parseDate(feature.properties.date))}</p>
                             <p class="small">${feature.properties.eventAddress}</p>`
                        )
                )
                .addTo(map);
        }

        // // Add the fetched GeoJSON data to the map as a source
        // map.on('load', () => {
        //     map.addSource('geojson-data', {
        //         type: 'geojson',
        //         data: geojson
        //     });
        //
        //     // Add a layer to use the GeoJSON data
        //     // map.addLayer({
        //     //     id: 'geojson-layer',
        //     //     type: 'circle   ',
        //     //     source: 'geojson-data',
        //     //     paint: {
        //     //         'circle-radius': 6,
        //     //         'circle-color': '#B42222'
        //     //     }
        //     // });
        //     //
        //     // // Optionally, add popups or other interactions here
        //     // map.on('click', 'geojson-layer', (e) => {
        //     //     const coordinates = e.features[0].geometry.coordinates.slice();
        //     //     const description = e.features[0].properties.description;
        //     //
        //     //     while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
        //     //         coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
        //     //     }
        //     //
        //     //     new mapboxgl.Popup()
        //     //         .setLngLat(coordinates)
        //     //         .setHTML(description)
        //     //         .addTo(map);
        //     // });
        //
        //     // Change the cursor to a pointer when the mouse is over the places layer.
        //     map.on('mouseenter', 'geojson-layer', () => {
        //         map.getCanvas().style.cursor = 'pointer';
        //     });
        //
        //     // Change it back to a pointer when it leaves.
        //     map.on('mouseleave', 'geojson-layer', () => {
        //         map.getCanvas().style.cursor = '';
        //     });
        // });
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
 */

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

document.addEventListener('DOMContentLoaded', loadGeoJSONForArtist);