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
- Runnable images are hosted here on GitHub
- Bicep template for deployment to Azure running as a Azure Container App

### Project Structure

- The repo adopts the "[Standard Go Project Layout](https://github.com/golang-standards/project-layout)" which can be unfamiliar to those that have not worked with Go. A full breakdown of the repo structure is provided below
- In addition my [own project starter template](https://github.com/benc-uk/project-starter) and [Go REST API for bookstore](https://github.com/benc-uk/go-rest-books) were used to bring a lot of reuseable assets, boilerplate and working source code.
- The breaking up of the API code into many packages and multiple source files was probably overkill at this point in the project but it provided a common set of abstractions (controllers, services, middleware, data-handlers) and also testability

### Data

- SQLite DB was used rather than a CSV file, this was sourced from here https://san-francisco.datasettes.com/food-trucks. A full blown database service was felt to be overkill for the assignment but being able to use SQL queries represented 
- SQLite is extremely fast & powerful, but it's a very bad choice for a real backend, this decision was purely tactical.
- There are some duplicates in the database, with the same name, lat & long but different IDs. This can result in less than 5 trucks being shown on the map even when 5 or more are present in the data returned from the API

### Front End

- Purposefully kept very clean/simple, no framewor required. Vue.js or React would be overkill at this point.
- Vanilla "modern" JS with ES5 modules and fetch.

### API

- Standard REST API pattern was used, using standard HTTP. 
- Only queries & GETs are used with a single `/trucks` endpoint.
- See [api/spec.yaml](./api/spec.yaml) for the OpenAPI description of the API, this is auto generated.
- See the [API docs for further details](./api/)

### Limitations, Known Issues & Backlog

- The query for finding nearby trucks is _extremely_ suboptimal and a borderline hack. Switching to a database service with spatial support like Cosmos, PostgreSQL or Azure SQL should be the highest priority
- Auth key to Azure Maps should be fetched with API, not baked into frontend code.
- GitHub Actions for CI & CD
  - Automate tests
- Rate limiting on the API (should use a upstream traffic gateway, e.g. ingress controller in Kubernetes NGINX/Envoy or Azure service like App Gateway)
- Auth in front of the API (likewise this should be handled by the gateway to do JWT validation etc)
- Sem ver for images and releases
- Fix the unit tests with the database mocked/stubbed
- Switch to RFC 7807 (Problem Details) for API errors https://datatracker.ietf.org/doc/html/rfc7807

The rest of the readme follows in a format similar to one I use on my many open source projects on GitHub

---

# 🚚 Food Truck Application

![](https://img.shields.io/github/license/benc-uk/food-truck)
![](https://img.shields.io/github/last-commit/benc-uk/food-truck)
![](https://img.shields.io/github/release/benc-uk/food-truck)

![](https://img.shields.io/github/checks-status/benc-uk/food-truck/main)
![](https://img.shields.io/github/workflow/status/benc-uk/food-truck/CI%20Build?label=ci-build)

# 🏃‍♂️ Getting Started

The makefile is the main starting point for working with this repo, simply calling `make` will provide this help text

```text
help                 💬 This help message :)
lint                 🌟 Lint & format, will not fix but sets exit code on error
lint-fix             🔍 Lint & format, will try to fix errors and modify code
image                📦 Build container image from Dockerfile
push                 📤 Push container image to registry
build                🔨 Run a local build without a container
run                  🏃 Run backend server, with hot reload, for local development
run-frontend         💻 Run frontend, with hot reload, for local development
install-tools        🔮 Install dev tools
generate             🔬 Generate Swagger / OpenAPI spec
test                 🥽 Run unit and integration tests
```

# 🚀 Installing / Deploying

- Deploy to Azure using `make deploy`
- 💥 NOTE: There is currently a limitation where the key for Azure maps is baked into the frontend code. This is a high priority backlog item. Update the key in web/client/config.js, build & push the image before deploying

# 📦 Running as container

- Build the images locally `make image`
  - Be sure to override and set your own `IMAGE_REG` and `IMAGE_REPO` e.g. `make image IMAGE_REG=foo`
- Run the public image using `docker run --rm -it -p 8080:8080 ghcr.io/benc-uk/food-truck:latest`

📝 NOTE! When running in a container the frontend is served from a different path, `/app/` so to access it locally use http://localhost:8080/app/

# 💻 Running locally

Run the server using:

```bash
make run
```

In a separate terminal session, run:

```bash
make run-frontend
```

Access the app using http://localhost:3000/

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

Details of any configuration files, environmental variables, command line parameters, etc.

For services

| Setting / Variable | Purpose                                                  | Default                 |
| ------------------ | -------------------------------------------------------- | ----------------------- |
| PORT               | Port the server will listen on.                          | 8080                    |
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
│   ├── api      - Base API
│   ├── data     - Data layer for calling SQLite
│   └── trucks   - Truck API and service
├── scripts      -
├── tests
└── web
    └── client
```

# 🌐 API

See the [API documentation](./api/) for full infomration about the API(s).

# 🪵 Change Log

See [complete change log](./CHANGELOG.md)

# ⚖️ License

This project uses the MIT software license. See [full license file](./LICENSE)
