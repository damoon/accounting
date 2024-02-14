
cli: .cli/passwd .cli/group .cli/image
	@docker run -ti --rm \
		-v $(PWD)/.cli/passwd:/etc/passwd \
		-v $(PWD)/.cli/group:/etc/group \
		--user $(shell id -u):$(shell id -g) \
		-h dev-container \
		--network host \
		-v $(PWD):$(PWD) \
		-w $(PWD) \
		--entrypoint "bash" \
		$(shell cat .cli/image)

pdf-to-txt:
	./bin/pdf-to-txt/run.sh

debug:
	go run ./bin/debug

.cli/image: $(shell find ./docker -type f)
	docker build -f docker/Dockerfile ./docker
	@docker build -f docker/Dockerfile ./docker -q > .cli/image

.cli/passwd:
	echo "$(shell whoami):x:$(shell id -u):$(shell id -g):,,,:$(shell pwd):/bin/bash" > .cli/passwd

.cli/group:
	echo "$(shell id -g -n):x:$(shell id -g):" > .cli/group

clean:
	rm -f .cli/image .cli/passwd .cli/group
