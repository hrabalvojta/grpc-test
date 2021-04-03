# Build stage
FROM golang:alpine AS build-env
ADD . /src/grpc-test
ENV CGO_ENABLED=0
RUN cd /src/grpc-test && go build -o /app

# Production stage
FROM scratch
COPY --from=build-env /app /

ENTRYPOINT ["/app"]
