build-for-win:
	GOOS=windows go build -o bin/task-runner.exe
build-for-lin:
	GOOS=linux go build -o bin/task-runner
build-all:
	make build-for-win
	make build-for-lin