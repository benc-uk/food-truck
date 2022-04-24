# ================================================================================================
# === Stage 1: Build API server and frontend host
# ================================================================================================
FROM golang:1.18-alpine as go-build
WORKDIR /build
ARG GO_PACKAGE="github.com/benc-uk/food-truck/cmd/main"
ARG VERSION="1.0.0"
ARG BUILD_INFO="Dockerfile build"

# Install system dependencies
RUN apk update && apk add git gcc musl-dev

# Fetch and cache Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

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
# Can't use scratch image as we have built a binary with CGO_ENABLED
FROM golang:1.18-alpine 
WORKDIR /app 

# Copy in database and frontend files
COPY web/client ./frontend
COPY data ./data

# Copy the Go server binary
COPY --from=go-build /build/server . 

# Location of database file and frontend set as environment variables
ENV DATABASE_PATH "/app/data/food-trucks.db"
ENV FRONTEND_DIR "./frontend"
ENV PORT 8080

# This should *always* be set and overidden at runtime
ENV AZURE_MAPS_KEY ""

EXPOSE 8080

# Start the backend 
CMD [ "./server" ]