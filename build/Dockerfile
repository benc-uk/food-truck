# ================================================================================================
# === Stage 1: Build API server and frontend host
# ================================================================================================
FROM golang:1.18-alpine as go-build
WORKDIR /build
ARG GO_PACKAGE="github.com/benc-uk/food-truck/cmd/main"
ARG VERSION="1.0.0"
ARG BUILD_INFO="Not set"

# Install system dependencies
RUN apk update && apk add git gcc musl-dev

# Fetch and cache Go modules
COPY go.mod .
COPY go.sum .
# Copy in Go source files
COPY cmd/ ./cmd
COPY pkg/ ./pkg

# SQLite needs CGO_ENABLED=1
RUN CGO_ENABLED=1 GOOS=linux go build \
  -ldflags "-X main.version=$VERSION -X 'main.buildInfo=$BUILD_INFO'" \
  -o server ./cmd/...

# ================================================================================================
# === Stage 2: Bundle server and frontend into Docker image
# ================================================================================================
FROM golang:1.18-alpine
WORKDIR /app 

# Copy in database and frontend files
COPY web/client ./frontend
COPY data ./data

# Copy the Go server binary
COPY --from=go-build /build/server . 

EXPOSE 8080
# Start the backend passing the path to the database and the path to the frontend
CMD [ "./server", "/app/data/food-trucks.db", "./frontend"]