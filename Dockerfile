FROM golang:1.19 as builder
WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o energy-monitor

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /build/energy-monitor /energy-monitor
ENTRYPOINT ["/energy-monitor"]
