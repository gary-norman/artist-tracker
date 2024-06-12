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

// Fetch GeoJSON data from a local file
fetch('/db/mapbox/0.geojson')
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
                    new mapboxgl.Popup({ offset: 25 }) // add popups
                        .setHTML(
                            `<p class="p--bold">${feature.properties.title}</p>
                             <p class="small"> {{ $dateParts := formatDate feature.properties.date }} {{ $dateParts.Day }} {{ $dateParts.Month }} {{ $dateParts.Year }}</p>
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
