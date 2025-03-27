FROM golang:alpine as build
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download

COPY . ./
RUN go build -o /bin/apiUpdater /src/

FROM scratch
COPY --from=build /bin/apiUpdater /bin/apiUpdater 

CMD ["/bin/apiUpdater"]
