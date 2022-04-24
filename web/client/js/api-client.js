// ========================================================
// Simple API Client for Food Trucks
// ========================================================

const API_ENDPOINT = '/'

const HEADERS = {
  'Content-Type': 'application/json',
}

// ========================================================
// Call to the find 5 trucks near point API
// ========================================================
export async function getTrucksNear(lat, long) {
  return await apiCall(`trucks/${lat}/${long}`)
}

// ========================================================
// Call to the find trucks in radius of point API
// ========================================================
export async function getTrucksInRadius(lat, long, radius) {
  return await apiCall(`trucks/${lat}/${long}?radius=${radius}`)
}

// ========================================================
// Used to fetch config from the server
// ========================================================
export async function getConfig() {
  return await apiCall(`config`)
}

// ========================================================
// Used to fetch status from the server
// ========================================================
export async function getStatus() {
  return await apiCall(`status`)
}

// ========================================================
// Wrapper around fetch, private and not exported
// ========================================================
async function apiCall(apiPath, method = 'get', data = null) {
  const resp = await fetch(`${API_ENDPOINT}${apiPath}`, {
    method,
    headers: HEADERS,
    body: data ? JSON.stringify(data) : null,
  })

  if (!resp.ok) {
    throw 'Failed to fetch or update truck data: ' + resp.status
  }

  const respData = await resp.json()
  return respData
}
