all: build

build:
	go build -o energy-monitor

run:
	go run .

clean:
	go clean

push:
	podman build . -t quay.io/jlindgren/energy-monitor
	podman push quay.io/jlindgren/energy-monitor
