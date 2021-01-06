FROM golang:1.15-alpine AS build

LABEL maintener="Maximilian Zinke <me@mxzinke.dev>"
WORKDIR /app

COPY . .

RUN go get
RUN go build -o build_artifact .

FROM alpine:latest

COPY --from=build /app/build_artifact ./binary

ENTRYPOINT [ "./binary" ]