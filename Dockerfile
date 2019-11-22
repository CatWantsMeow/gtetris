FROM golang

COPY . /go/src/github.com/CatWantsMeow/tetris/
WORKDIR /go/src/github.com/CatWantsMeow/tetris/

RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build . && \
    mkdir -p /go/bin && \
    mv -v tetris /go/bin/

FROM scratch
COPY --from=0 /go/bin/tetris /tetris
CMD ["/tetris"]
