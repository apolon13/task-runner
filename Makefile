build-for-lin:
	GOOS=linux go build -o bin/task-runner
build-all:
	make build-for-lin