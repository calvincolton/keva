FROM golang:1.19 as build

COPY . /go/src/keva

WORKDIR /go/src/keva

RUN go get github.com/gorilla/mux github.com/lib/pq

RUN CGO_ENABLED=0 GOOS=linux go build -o keva

FROM scratch as image

COPY --from=build /go/src/keva/keva .

EXPOSE 8080

CMD ["/keva"]