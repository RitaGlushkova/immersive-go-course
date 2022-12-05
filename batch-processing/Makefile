SHELL=/bin/bash
.PHONY: run test develop

outputs:
	mkdir -p outputs

run: outputs
	docker build --target run -t run .
	docker run \
		--env-file docker_env \
		--mount type=bind,source="$$(pwd)/inputs",target=/inputs \
		--mount type=bind,source="$$(pwd)/outputs",target=/outputs \
		--mount type=bind,source="$$(echo $$HOME)/.aws",target=/root/.aws \
		--rm run

test: outputs
	docker build --target test -t test .
	docker run \
		--rm test

develop: outputs
	docker build --target develop -t develop .
	docker run -it \
		--env-file docker_env \
		--mount type=bind,source="$$(pwd)",target=/app \
		--mount type=bind,source="/tmp",target=/tmp \
		--mount type=bind,source="$$(echo $$HOME)/.aws",target=/root/.aws \
		--mount type=bind,source="$$(pwd)/inputs",target=/inputs \
		--mount type=bind,source="$$(pwd)/outputs",target=/outputs \
		--rm develop