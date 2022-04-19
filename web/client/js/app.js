// ========================================================
// App and presentation layer for food truck app
// ========================================================

import 'https://atlas.microsoft.com/sdk/javascript/mapcontrol/2/atlas.min.js'
import { fetchTrucksNear } from './api-client.js'
import { AZURE_MAPS_KEY } from '../config.js'

let map = null
const defaultZoom = 15
var userDatasource = null
var truckDatasource = null
const radius = 600

// =============================================================================
// Initialize the map and application
// =============================================================================
function initApp() {
  map = new atlas.Map('truckMap', {
    view: 'Auto',

    authOptions: {
      authType: 'subscriptionKey',

      subscriptionKey: AZURE_MAPS_KEY,
    },
  })

  // Only when the map is ready, we can add the data
  map.events.add('ready', () => {
    // Load custom sprites
    map.imageSprite.add('truck', './img/food-truck.png')

    // Set up data sources
    truckDatasource = new atlas.source.DataSource('trucks', {
      cluster: false,
    })
    map.sources.add(truckDatasource)
    map.layers.add(
      new atlas.layer.SymbolLayer(truckDatasource, null, {
        iconOptions: {
          image: 'truck',
          size: 0.4,
          ignorePlacement: true,
          allowOverlap: true,
        },

        textOptions: {
          textField: ['to-string', ['get', 'name']],
          offset: [0, -3],
        },
      })
    )

    userDatasource = new atlas.source.DataSource()
    map.sources.add(userDatasource)
    map.layers.add(new atlas.layer.SymbolLayer(userDatasource))

    // Ask the browser to locate the user
    navigator.geolocation.getCurrentPosition(
      (position) => {
        refreshMap(position.coords.longitude, position.coords.latitude)
      },

      (err) => {
        alert('Unable to get location: ' + err.message)
      }
    )

    map.setCamera({
      zoom: defaultZoom,
    })
  })

  map.events.add('click', (evt) => {
    refreshMap(evt.position[0], evt.position[1])
  })
}

// =============================================================================
// Show nearby food trucks on the map
// =============================================================================
async function showTrucks(lat, long) {
  // Reset truck data
  truckDatasource.clear()

  // Call the API
  try {
    const trucks = await fetchTrucksNear(long, lat, radius)

    // Process the data
    for (const truck of trucks) {
      const point = new atlas.data.Point([truck.long, truck.lat])
      const feature = new atlas.data.Feature(point, {
        name: truck.name,
      })
      truckDatasource.add(feature)
    }
  } catch (err) {
    showError('API error: ' + err)
  }
}

// =============================================================================
// Button to recenter the user to their real location
// =============================================================================
document.querySelector('#locateUser').addEventListener('click', () => {
  navigator.geolocation.getCurrentPosition(
    (position) => {
      truckDatasource.clear()
      refreshMap(position.coords.longitude, position.coords.latitude)
    },

    (err) => {
      alert('Unable to get location: ' + err.message)
    }
  )
})

// =============================================================================
// Button to teleport to a place that sounds pretty terrible
// =============================================================================
document.querySelector('#sanFran').addEventListener('click', () => {
  truckDatasource.clear()
  refreshMap(-122.4205, 37.7758)
})

// =============================================================================
// Main function to refresh the map and show trucks
// =============================================================================
function refreshMap(lat, long) {
  userDatasource.clear()
  userDatasource.add(new atlas.data.Point([lat, long]))
  map.setCamera({
    center: [lat, long],
  })

  showTrucks(lat, long)
}

// =============================================================================
// Helper for showing errors
// =============================================================================
function showError(error) {
  const errDiv = document.querySelector('#error')
  errDiv.innerHTML = error
  errDiv.style.display = 'block'
}

// =============================================================================
// ENTRY POINT - It starts here
// =============================================================================
initApp()
