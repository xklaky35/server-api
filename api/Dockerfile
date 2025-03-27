FROM golang:1.23.4-alpine as build
WORKDIR /src


COPY go.mod go.sum /src/
RUN go mod download

COPY . ./
RUN go build -o /bin/server-api /src/

FROM alpine:latest
COPY --from=build /bin/server-api /bin/server-api 

ENV PORT=8080
EXPOSE 8080

COPY data/data.json /data/data.json

CMD ["/bin/server-api"]
