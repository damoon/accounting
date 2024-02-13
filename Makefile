
cli: .passwd .group .image
	@docker run -ti --rm \
		-v $(PWD)/.passwd:/etc/passwd \
		-v $(PWD)/.group:/etc/group \
		--user $(shell id -u):$(shell id -g) \
		-h dev-container \
		--network host \
		-v $(PWD):$(PWD) \
		-w $(PWD) \
		--entrypoint "bash" \
		$(shell cat .image)

pdf-to-txt:
	./bin/pdf-to-txt/run.sh

debug:
	go run ./bin/debug

.image: $(shell find ./docker -type f)
	docker build -f docker/Dockerfile ./docker
	@docker build -f docker/Dockerfile ./docker -q > .image

.passwd:
	echo "$(shell whoami):x:$(shell id -u):$(shell id -g):,,,:$(shell pwd):/bin/bash" > .passwd

.group:
	echo "$(shell id -g -n):x:$(shell id -g):" > .group

clean:
	rm .image .passwd .group
