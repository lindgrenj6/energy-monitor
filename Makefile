all: build

build:
	go build -o energy-monitor

clean:
	go clean

push:
	podman build . -t quay.io/jlindgren/energy-monitor
	podman push quay.io/jlindgren/energy-monitor
