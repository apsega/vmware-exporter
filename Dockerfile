# Build the application from source
FROM golang:1.21.1 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /vmware-exporter

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /vmware-exporter /vmware-exporter

USER nonroot:nonroot

ENTRYPOINT ["/vmware-exporter"]
