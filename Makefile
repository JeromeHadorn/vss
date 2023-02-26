setup:
	mkdir -p bin

build_win32:
	docker build -t vss_windows_image_32 -f Dockerfile.win32 --progress=plain . && \
	docker run --name "vss_windows_container_32" vss_windows_image_32 && \
	docker cp vss_windows_container_32:/go/bin/vss.exe ./bin/vss_x86_32.exe && \
	docker rm vss_windows_container_32 && \
	docker image rm vss_windows_image_32

build_win64:
	docker build -t vss_windows_image_64 -f Dockerfile.win64 --progress=plain . && \
	docker run --name "vss_windows_container_64" vss_windows_image_64 && \
	docker cp vss_windows_container_64:/go/bin/vss.exe ./bin/vss_x86_64.exe && \
	docker rm vss_windows_container_64 && \
	docker image rm vss_windows_image_64

build:
	CGO_ENABLED="1" \
	GOOS="windows" \
    go build -o bin/ cmd/example/*
