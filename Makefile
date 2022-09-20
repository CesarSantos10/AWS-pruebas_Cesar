.PHONY: build

build:
	sam build

build-SumAgeFn:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(ARTIFACTS_DIR)/handler functions/sum-age-fn/main.go

build-UpdateStudentCompleteFn:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(ARTIFACTS_DIR)/handler functions/update-student-complete/main.go