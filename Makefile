setup:
	mkdir -p bin

build_windows:
	docker build -t vss_windows_image_64 -f Dockerfile.windows --progress=plain . && \
	docker run --name "vss_windows_container_64" vss_windows_image_64 && \
	docker cp vss_windows_container_64:/go/bin/vss.exe ./bin/vss_x86_64.exe && \
	docker rm vss_windows_container_64 && \
	docker image rm vss_windows_image_64