GO_BUILD_FLAGS=-mod vendor -trimpath
GO_BUILD_DBG_FLAGS = $(GO_BUILD_FLAGS) -gcflags "-N -l" -race

build_switcher:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/switcher fuxi/switcher

build_gs:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/gs fuxi/gs

build_map:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/map fuxi/map

build_robot:
	go build $(GO_BUILD_DBG_FLAGS) -o ./bin/robot fuxi/robot

build_all:build_switcher build_gs build_map build_robot

run_switcher:
	./bin/switcher

run_gs:
	./bin/gs --pvid 1

run_map:
	./bin/map --pvid 2

run_robot:
	./bin/robot

stop_switcher:
	-ps aux | grep "[b]in/switcher" | awk '{print $$2}' | xargs kill -15

stop_gs:
	-ps aux | grep "[b]in/gs" | awk '{print $$2}' | xargs kill -15

stop_map:
	-ps aux | grep "[b]in/map" | awk '{print $$2}' | xargs kill -15

stop_robot:
	-ps aux | grep "[b]in/robot" | awk '{print $$2}' | xargs kill -15

stop_all:stop_robot stop_map stop_gs stop_switcher