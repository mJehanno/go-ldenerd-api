FROM golang
COPY . /go/src/go-ldener-api
ENTRYPOINT /go/src/go-ldener-api
EXPOSE 8000