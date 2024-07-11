document.addEventListener('DOMContentLoaded', loadGeoJSONForAllLocations);

// Ensure to add your Mapbox access token
mapboxgl.accessToken = 'pk.eyJ1IjoibG9yZXdvcmxkIiwiYSI6ImNsd3FseDNsbDAzZjMyanF2czh3Mmt4eTgifQ.-_bXsAv_SR1bpcmvOSpDuA';

// Map container
const mapBlack = document.getElementById("map-container");

let mapProject = 'globe';

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

// This function returns a promise that resolves after n milliseconds
const wait = (n) => new Promise((resolve) => setTimeout(resolve, n));

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

// Load GeoJSON data for all locations
async function loadGeoJSONForAllLocations() {
    const geojsonFiles = [];
    const numFiles = 51; // Number of GeoJSON files

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

        console.log('GeoJSON Data:', geojsons);

        const map = new mapboxgl.Map({
            container: 'map', // container ID
            style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8', // map style
            projection: mapProject,
            center: [-98.5795, 39.8283], // starting position [lng, lat]
            zoom: 1 // starting zoom
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

        // reset map position
        document.getElementById('kaartReset').addEventListener('click', () => {
            map.flyTo({
                center: [-98.5795, 39.8283],
                zoom: 1,
                pitch: 0,
                essential: true // animation considered essential for accessibility
            });
        });

        // Select all tour date elements
        const mapClick = document.querySelectorAll(".artistPageTourdate");
        mapClick.forEach(tourdate => {
            tourdate.addEventListener('click', () => {
                const concertId = tourdate.dataset.index;

                console.log("Concert ID:", concertId);

                const concert = geojsons.flatMap(geojson => geojson.features)
                    .find(feature => feature.properties.date.includes(concertId));

                if (concert) {
                    const coordinates = concert.geometry.coordinates;

                    console.log("Coordinates:", coordinates);

                    map.flyTo({
                        center: coordinates,
                        zoom: 16,
                        pitch: 60,
                        essential: true // animation considered essential for accessibility
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


// // Ensure to add your Mapbox access token
// mapboxgl.accessToken = 'pk.eyJ1IjoibG9yZXdvcmxkIiwiYSI6ImNsd3FseDNsbDAzZjMyanF2czh3Mmt4eTgifQ.-_bXsAv_SR1bpcmvOSpDuA';
//
// // Create the map instance
// const map = new mapboxgl.Map({
//     container: 'map', // container ID
//     style: 'mapbox://styles/loreworld/clx6fy3dp01w001pnegho7ud8', // map style
//     center: [-98.5795, 39.8283], // starting position [lng, lat]
//     zoom: 1 // starting zoom
// });
//
// map.on('style.load', () => {
//     map.setFog({
//         "range": [0.8, 9],
//         "color": "#471e50",
//         "horizon-blend": .2,
//         "high-color": "#07173e",
//         "space-color": "#000000",
//         "star-intensity": 0.35
//     }); //
// });
//
// // At low zooms, complete a revolution every two minutes.
// const secondsPerRevolution = 90;
// // Above zoom level 5, do not rotate.
// const maxSpinZoom = 5;
// // Rotate at intermediate speeds between zoom levels 3 and 5.
// const slowSpinZoom = 3;
//
// let userInteracting = false;
// const spinEnabled = true;
//
// function spinGlobe() {
//     const zoom = map.getZoom();
//     if (spinEnabled && !userInteracting && zoom < maxSpinZoom) {
//         let distancePerSecond = 360 / secondsPerRevolution;
//         if (zoom > slowSpinZoom) {
//             // Slow spinning at higher zooms
//             const zoomDif =
//                 (maxSpinZoom - zoom) / (maxSpinZoom - slowSpinZoom);
//             distancePerSecond *= zoomDif;
//         }
//         const center = map.getCenter();
//         center.lng -= distancePerSecond;
//         // Smoothly animate the map over one second.
//         // When this animation is complete, it calls a 'moveend' event.
//         map.easeTo({ center, duration: 1000, easing: (n) => n });
//     }
// }
//
// // Pause spinning on interaction
// map.on('mousedown', () => {
//     userInteracting = true;
// });
// map.on('dragstart', () => {
//     userInteracting = true;
// });
//
// // When animation is complete, start spinning if there is no ongoing interaction
// map.on('moveend', () => {
//     spinGlobe();
// });
//
// spinGlobe();
//
// function parseDate(dateStr) {
//     const [day, month, year] = dateStr.split('-');
//     return new Date(`${year}-${month}-${day}`);
// }
//
// function formatDate(date) {
//     const options = { day: '2-digit', month: 'short', year: 'numeric' };
//     return date.toLocaleDateString('en-GB', options);
// }
//
// // Generate GeoJSON file paths dynamically
// const geojsonFiles = [];
// const numFiles = 51; // Number of GeoJSON files
//
// for (let i = 0; i < numFiles; i++) {
//     geojsonFiles.push(`/db/mapbox_std/${i}.geojson`);
// }
//
// // Fetch all GeoJSON files and process them
// Promise.all(geojsonFiles.map(file => fetch(file)
//     .then(response => {
//         if (!response.ok) {
//             throw new Error('Network response was not ok ' + response.statusText);
//         }
//         return response.json();
//     })
//     .catch(error => {
//         console.error('There has been a problem with your fetch operation:', error);
//     })
// ))
//     .then(geojsons => {
//         geojsons.forEach(geojson => {
//             console.log('GeoJSON Data:', geojson);
//
//             // add markers to map
//             for (const feature of geojson.features) {
//                 const el = document.createElement('div');
//                 el.className = 'marker';
//
//                 new mapboxgl.Marker(el)
//                     .setLngLat(feature.geometry.coordinates)
//                     .setPopup(
//                         new mapboxgl.Popup({ offset: 20 })
//                             .setHTML(
//                                 `<div class="p--bold flexrow"><a href="/artist?name=${feature.properties.artist}">${feature.properties.artist}</a>
//                              <b>live in ${feature.properties.eventAddress}</b></div>
//                              <p class="small">${formatDate(parseDate(feature.properties.date))}</p>
//                              <p class="small">${feature.properties.eventAddress}</p>`
//                             )
//                     )
//                     .addTo(map);
//             }
//         });
//     })
//     .catch(error => {
//         console.error('There has been a problem with your fetch operation:', error);
//     });