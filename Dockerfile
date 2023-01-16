FROM golang:1.17 as builder
ARG server
ENV CGO_ENABLED 0
WORKDIR /go/src/fuxi
COPY . /go/src/fuxi
RUN go build -mod vendor -trimpath -o bin/$server fuxi/$server

FROM alpine
ARG server
ENV svr=""
ENV etcd=""
ENV pport=0
ENV args=""
WORKDIR /go/src/fuxi
COPY --from=builder /go/src/fuxi/bin/$server ./
ENTRYPOINT ["sh", "-c", "./${svr} --etcd ${etcd} --pport ${pport} ${args}"]