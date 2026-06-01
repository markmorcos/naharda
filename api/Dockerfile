# Multi-stage build → distroless static (project.md §6, §9.6).
FROM golang:1.24-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /naharda-api ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /naharda-api /naharda-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/naharda-api"]
