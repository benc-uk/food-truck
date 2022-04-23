# 🍖🤤 Food Truck Engineering Challenge 🍔🚚

This repo is the output of my work on the [Food Truck - Take Home Engineering Challenge](https://github.com/timfpark/take-home-engineering-challenge)

## 📝 Design & Architecture Notes

Here is a outline of some of the key decisions, trade offs, limitations & choices to this point in the project

### Components & Features

- Backend REST API server & host for frontend JS app. This is written in Go.
  > Note: Co-hosting like this is preferred to prevent CORS issues, when we don't have a layer-7 routing service in front
- Frontend JS application, written in vanilla JS
  - Frontend uses Azure Maps for visualization
- Dockerfile to build the service + frontend as a container
- Auto generated OpenAPI / Swagger support, with SwaggerUI
- Middleware for CORS, logging, health checks and Prometheus metrics
- Makefile to assist with local development and CI/CD, includes:
  - Linting with golangci-lint
  - Local hot reload with cosmtrek/air
  - Build, tag & push container images
- Unit & integration tests
- Runnable images are hosted [here on GitHub Container Registry](https://github.com/benc-uk/food-truck/pkgs/container/food-truck)
- Bicep template for deployment to Azure running as a Azure Container App
- Dev Container with all required tooling, to support local development
- [Performance and black box REST API tests using k6](./tests/)

### Project Structure

- The repo adopts the "[Standard Go Project Layout](https://github.com/golang-standards/project-layout)" which can be unfamiliar to those that have not worked with Go. A full breakdown of the repo structure is provided below
- In addition several of my own open source projects were used to bring many reuseable assets, boilerplate and working source code:
  - My own project starter template: [benc-uk/project-starter](https://github.com/benc-uk/project-starter) 
  - Go REST API for bookstore: [benc-uk/go-rest-books](https://github.com/benc-uk/go-rest-books) 
  - Bicep code was borrowed from:  [benc-uk/bicep-iac](https://github.com/benc-uk/bicep-iac) and [benc-uk/chatr](https://github.com/benc-uk/chatr)
  - Performance tests were taken from other projects and reports use [my k6 report generator (benc-uk/k6-reporter)](https://github.com/benc-uk/k6-reporter)
- The breaking up of the API code into many packages, interfaces and multiple source files was probably overkill at this point (with single endpoint etc) However it provided a common set of abstractions (routing, services, middleware, data-handlers) and also importantly testability

### Data

- SQLite DB was used rather than a CSV file, this was sourced from here https://san-francisco.datasettes.com/food-trucks. A full blown database service was felt to be overkill for the assignment but being able to use SQL queries represented a reasonable degree of "realism".
- SQLite is very fast & powerful, but it's a extremely bad choice for a backend service with CRUD, so this decision was purely tactical and doesn't represent best/recommended practice.
- There are some duplicates in the database, with the same name, lat & long but different IDs. This can result in less than 5 trucks being shown on the map even when 5 or more are present in the data returned from the API.

### Front End

- Purposefully kept very clean/simple, no framework required. Vue.js or React would be overkill at this point. [Alpine.js](https://alpinejs.dev/) was considered but with no reactive components in the app, even this minimal library would be too much.
- Vanilla "modern" JS was used, with ES5 modules allowing code structuring and fetch for API calls.

### API

- Classic REST API, using HTTP.
- Only queries & GETs are used with a single `/trucks` endpoint.
- See [api/spec.yaml](./api/spec.yaml) for the OpenAPI description of the API, this is auto generated.
- See the [API docs for further details](./api/)

### Limitations, Known Issues & Backlog

- The query for finding nearby trucks is _extremely_ sub-optimal and a borderline hack. Switching to a database service with spatial support like Cosmos, PostgreSQL or Azure SQL should be the highest priority.
- Fix the leaky/poor abstraction in the database spec & improve the unit testing method.
- ~~Auth key to Azure Maps should be fetched with API, not baked into frontend code.~~ DONE!
- GitHub Actions for CI & CD
  - ~~Automate builds~~ DONE!
  - ~~Automate tests~~ DONE!
- Rate limiting on the API (should use a upstream traffic gateway, e.g. ingress controller in Kubernetes NGINX/Envoy or Azure service like App Gateway).
- Auth in front of the API (likewise this should be handled by the gateway to do JWT validation etc).
- Sem ver for images and releases.
- Add CLI tool
- ~~Add end to end API & performance tests, k6.io is my tool of choice for this, or Postman/Newman.~~
- Switch to RFC 7807 (Problem Details) for API errors https://datatracker.ietf.org/doc/html/rfc7807
- Consider switching to dependency injection but should weigh up the pros & cons.
- Add rest of the application functionality ;)

The rest of the readme follows in a format similar to one I use on my many open source projects on GitHub

---

# 🚚 Food Truck Application

![](https://img.shields.io/github/license/benc-uk/food-truck)
![](https://img.shields.io/github/last-commit/benc-uk/food-truck)
![](https://img.shields.io/github/release/benc-uk/food-truck)

![](https://img.shields.io/github/last-commit/benc-uk/food-truck)
![](https://img.shields.io/github/workflow/status/benc-uk/food-truck/CI%20Build?label=ci-build)

# 🏃‍♂️ Getting Started

The makefile is the main starting point for working with this project, simply calling `make` will provide this help text. If you are using a system without access to make, I suggest using the provided dev container.

```text
help                 💬 This help message :)
lint                 🌟 Lint & format, will not fix but sets exit code on error
lint-fix             🔍 Lint & format, will try to fix errors and modify code
image                📦 Build container image from Dockerfile
push                 📤 Push container image to registry
build                🔨 Run a local build without a container
run                  🏃 Run server & frontend host, with hot reload for local dev
install-tools        🔮 Install dev tools
generate             🔬 Generate Swagger / OpenAPI spec
test                 🧪 Run unit and integration tests
test-perf            📈 Run performance tests
deploy               🚀 Deploy to Azure using Bicep & Azure CLI
```

# 🚀 Installing / Deploying

Deploy to Azure using `make deploy` this will deploy to Azure Container Apps (NOTE: only certain regions presently supported) plus the Azure Maps account. Set `AZURE_IMAGE` variable when running `make deploy` in order to deploy your own image, otherwise `ghcr.io/benc-uk/food-truck:latest` will be used.

Deployment is done through Bicep with a set of template & modules. You will need Azure CLI with the Bicep add-on installed.

# 🏃 Running Locally

Before running locally you will need to deploy Azure Maps account in Azure and obtain the shared access key. This can easily be done with the Azure CLI, e.g.

```bash
RES_GROUP=__CHANGE_ME__
az maps account create --name food-truck-maps --resource-group $RES_GROUP --kind Gen2 --sku G2
echo Access key is: $(az maps account keys list --name food-truck-maps --resource-group $RES_GROUP --query primaryKey -o tsv)
```

## 📦 Running as a container

- PRE-REQS: Docker engine installed locally and Docker CLI
- Build the images locally `make image`
  - Be sure to override and set your own `IMAGE_REG` and `IMAGE_REPO` e.g. `make image IMAGE_REPO=myfoodtruck`
- ALTERNATIVELY: Run the public image directly using

```bash
docker run --rm -it -p 8080:8080 \
  -e AZURE_MAPS_KEY=__YOUR_MAPS_KEY__ \
  ghcr.io/benc-uk/food-truck:latest
```

📝 NOTE! When running in a container the frontend is served from a different path, `/app/` so to access it locally use http://localhost:8080/app/ (Note the trailing slash!)

## 💻 Running using Go

If you do not have Go installed, then open the repo the provided dev container, this has all the tools you will need

Copy `.env.sample` to `.env` and edit the file, setting **AZURE_MAPS_KEY** to the correct value.

Now run the server and backend API using:

```bash
make run
```

Access the app using http://localhost:8080/app/

# 🧪 Running Tests

Run the integration tests and unit tests with:

```bash
make test
```

# 🏗️ Architecture

Standard SPA (Single Page Application) style frontend with REST backend

```text
                   ┌───────────────────────────────────┐
                   │             CONTAINER             │
 ┌───────────┐     │ ┌─────────────┐    ┌────────────┐ │
 │           │HTTP │ │             ├────►            │ │
 │ FRONTEND  ├─────┼─►  REST API   │    │  SQLITE DB │ │
 │ (BROWSER) │     │ │             ◄────┤            │ │
 └───────────┘     │ └─────────────┘    └────────────┘ │
                   └───────────────────────────────────┘
```

# 🔧 Configuration

The following env vars are used by the backend server. `.env` (aka dotenv) is looked for and loaded when working locally

| Setting / Variable | Purpose                                                  | Default                 |
| ------------------ | -------------------------------------------------------- | ----------------------- |
| AZURE_MAPS_KEY     | Azure Maps shared access key                             | _None_                  |
| PORT               | Port the server will listen on                           | 8080                    |
| DATABASE_PATH      | Path to the database file                                | "./data/food-trucks.db" |
| FRONTEND_DIR       | Where the frontend HTML/JS is located for static serving | "./web/client"          |

# 📂 Repository Structure

A brief description of the top-level directories of this project is as follows:

```text
├── api          - API spec and docs
├── build        - Build assets, Dockerfiles etc
├── cmd          - Go source for executables, main server is here
├── data         - Application data, SQLite db file(s)
├── deploy       - Deployment assets, Bicep templates to deploy to Azure
├── pkg          - Go source packages
│   ├── api      - Base API common
│   ├── data     - Data layer for calling SQLite
│   └── trucks   - Truck API and service
├── scripts      - Some helper scripts
├── tests        - End to end performance tests
└── web
    └── client   - The application frontend source code
```

# 🌐 REST API

See the [API documentation](./api/) for full information about the food truck API(s)  

To aid development work and testing a REST API file is provided `api/test-api.http` this requires the [REST Client extension for VSCode](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) and allows for making test calls to the API from within VSCode

## Other Endpoints

- `/config` - Used by the frontend to get configuration, i.e. the Azure Maps key
- `/metrics` - Metrics in Prometheus format for observability
- `/status` - Simple status API
- `/health` - Heath check for use with Kubernetes & load balancer probes
- `/swagger` - Swagger UI

# 🪵 Change Log

See [complete change log](./CHANGELOG.md)

# ⚖️ License

This project uses the MIT software license. See [full license file](./LICENSE)
