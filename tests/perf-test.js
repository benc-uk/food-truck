import http from 'k6/http'
import { group, check, sleep } from 'k6'
import { Trend } from 'k6/metrics'

// Adds reporting
import { htmlReport } from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js'
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js'

// Top level test parameters & defaults
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'
const API_PATH = __ENV.API_PATH || '/'
const SLEEP_DURATION = parseFloat(__ENV.TEST_SLEEP) || 0.1

// Execution options for k6 run
export const options = {
  scenarios: {
    constant_load: {
      // See https://k6.io/docs/using-k6/scenarios/#executors
      executor: 'constant-vus',

      // common scenario configuration
      startTime: '0s',
      gracefulStop: '5s',

      // executor-specific configuration
      vus: __ENV.TEST_VUS || 10,
      duration: __ENV.TEST_DURATION || '10s',
    },
  },

  // These will need adjusting
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<5000'],
  },
}

// Force the tests to run as a single iteration, no scenario
if (__ENV.TEST_ONE_ITER == 'true') {
  delete options.scenarios
  options.vus = 1
  options.iterations = 1
}

let getTrucksSanFran = new Trend('getTrucksSanFran')
let getTrucksLondon = new Trend('getTrucksLondon')
let getTrucksRadius = new Trend('getTrucksRadius')

export default function () {
  const endpoint = BASE_URL + API_PATH
  const headers = {}

  //
  // Tests for the truck API
  //
  group('Trucks API', () => {
    //
    // Find five trucks in San Fran
    //
    let resp = http.get(`${endpoint}trucks/37.7758/-122.4205`, { headers })

    check(resp, {
      getTrucksSanFran_returns_200: (r) => r.status === 200,
      getTrucksSanFran_five_results: (r) => r.json().length >= 5,
    })

    getTrucksSanFran.add(resp.timings.duration)

    sleep(SLEEP_DURATION)

    //
    // Find five trucks near London
    //
    resp = http.get(`${endpoint}trucks/51.403278/0.056169`, { headers })

    check(resp, {
      getTrucksLondon_returns_200: (r) => r.status === 200,
      getTrucksLondon_five_results: (r) => r.json().length >= 5,
    })

    getTrucksLondon.add(resp.timings.duration)

    sleep(SLEEP_DURATION)

    //
    // Find any trucks in radius of point
    //
    resp = http.get(`${endpoint}trucks/37.7758/-122.4205?radius=400`, { headers })

    check(resp, {
      getTrucksRadius_returns_200: (r) => r.status === 200,
      getTrucksRadius_has_results: (r) => r.json().length >= 0,
    })

    getTrucksRadius.add(resp.timings.duration)

    sleep(SLEEP_DURATION)
  })
}

// Generates a HTML report
export function handleSummary(data) {
  return {
    'output/load-test-summary.html': htmlReport(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  }
}
