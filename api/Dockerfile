FROM golang:alpine as build
WORKDIR /src


COPY go.mod go.sum /src/
RUN go mod download

COPY . ./
RUN go build -o /bin/server-api /src/

FROM alpine:latest
COPY --from=build /bin/server-api /bin/server-api 

ENV PORT=:3000
EXPOSE 3000

#COPY ./data/data.json /welcome-page-api/data.json

CMD ["/bin/server-api"]
