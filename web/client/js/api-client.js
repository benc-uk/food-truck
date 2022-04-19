// ========================================================
// API Client for Food Trucks
// ========================================================

let API_ENDPOINT = '/'

if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
  API_ENDPOINT = 'http://localhost:8080/'
}

const HEADERS = {
  'Content-Type': 'application/json',
}

export async function fetchTrucksNear(lat, long, radius) {
  return await apiCall(`trucks/${lat}/${long}`) //?radius=${radius}`)
}

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
