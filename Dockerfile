FROM golang

COPY . /go/src/github.com/CatWantsMeow/gtetris/
WORKDIR /go/src/github.com/CatWantsMeow/gtetris/

RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build . && \
    mkdir -p /go/bin && \
    mv -v gtetris /go/bin/

FROM scratch
COPY --from=0 /go/bin/gtetris /gtetris
CMD ["/gtetris"]
