GO_BUILD_FLAGS=-mod vendor -trimpath
GO_BUILD_DBG_FLAGS = $(GO_BUILD_FLAGS) -gcflags "-N -l" -race
HOST=172.20.170.217

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

build_switcher:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/switcher fuxi/switcher

build_gs:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/gs fuxi/gs

build_map:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/map fuxi/map

build_robot:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/robot fuxi/robot

build_all: build_switcher build_gs build_map build_robot

run_switcher:
	./bin/switcher --pport 10001

run_gs:
	./bin/gs --pvid 1 --pport 10002

run_map:
	./bin/map --pvid 2 --pport 10003

run_robot:
	./bin/robot --pport 10004

stop_switcher:
	-ps aux | grep "[b]in/switcher" | awk '{print $$2}' | xargs kill -15

stop_gs:
	-ps aux | grep "[b]in/gs" | awk '{print $$2}' | xargs kill -15

stop_map:
	-ps aux | grep "[b]in/map" | awk '{print $$2}' | xargs kill -15

stop_robot:
	-ps aux | grep "[b]in/robot" | awk '{print $$2}' | xargs kill -15

stop_all: stop_robot stop_map stop_gs stop_switcher

img_switcher:
	docker build -f Dockerfile_switcher -t switcher .

img_gs:
	docker build --build-arg server=gs -t gs .

img_map:
	docker build --build-arg server=map -t map .

img_robot:
	docker build --build-arg server=robot -t robot .

img_all: img_switcher img_gs img_map img_robot

run_docker_switcher:
	docker ps -a -q --filter "name=switcher" -q | xargs docker rm
	docker run --name switcher \
		-p 8080:8080 \
		-p 8088:8088 \
		-e linker=$(HOST):8080 \
		-e provider=$(HOST):8088 \
		-e etcd=$(HOST):2379 \
		switcher

run_docker_gs:
	docker ps -a -q --filter "name=gs" -q | xargs docker rm
	docker run --name gs \
		-e etcd=$(HOST):2379 \
		-e svr=gs \
		-e args="--pvid=1" \
		gs

run_docker_map:
	docker ps -a -q --filter "name=map" -q | xargs docker rm
	docker run --name map \
		-e etcd=$(HOST):2379 \
		-e svr=map \
		-e args="--pvid=2" \
		map

run_docker_robot:
	docker ps -a -q --filter "name=robot" -q | xargs docker rm
	docker run --name robot \
		-e etcd=$(HOST):2379 \
		-e svr=robot \
		-e args="--num=1" \
		robot

run_docker_all: run_docker_switcher run_docker_gs run_docker_map run_docker_robot

stop_docker_switcher:
	docker ps -a -q --filter "name=switcher" -q | xargs docker stop

stop_docker_gs:
	docker ps -a -q --filter "name=gs" -q | xargs docker stop

stop_docker_map:
	docker ps -a -q --filter "name=map" -q | xargs docker stop

stop_docker_robot:
	docker ps -a -q --filter "name=robot" -q | xargs docker stop

stop_docker_all: stop_docker_robot stop_docker_map stop_docker_gs stop_docker_switcher

img_clean:
	docker images -a | grep none | awk '{ print $3; }' | xargs docker rmi

fmt:
	gofmt -s -w ${GOFILES}

vet:
	go vet ${VETPACKAGES}