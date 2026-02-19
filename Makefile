COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))
BUILDENVVAR=CGO_ENABLED=0

BINARY_NAME=kubenexus-scheduler
DOCKER_IMAGE=kubenexus-scheduler
VERSION?=v0.1.0

.PHONY: all
all: generate test build

.PHONY: generate
generate:
	@echo "Generating code..."
	go run sigs.k8s.io/controller-tools/cmd/controller-gen@latest object:headerFile="hack/boilerplate.go.txt" paths="./pkg/apis/..."

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '-w' -o bin/$(BINARY_NAME) cmd/main.go

.PHONY: test
test:
	@echo "Running unit tests..."
	$(BUILDENVVAR) go test -v ./pkg/apis/... ./pkg/plugins/coscheduling/... ./pkg/plugins/resourcereservation/... ./pkg/workload/... ./pkg/utils/... ./pkg/scheduler/...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(BUILDENVVAR) go test -v -timeout=10m ./test/integration/...

.PHONY: test-e2e
test-e2e:
	@echo "Running E2E tests..."
	@echo "This will create a Kind cluster, run tests, and clean up"
	$(BUILDENVVAR) go test -v -timeout=30m ./test/e2e/...

.PHONY: test-benchmark
test-benchmark:
	@echo "Running performance benchmarks..."
	$(BUILDENVVAR) go test -bench=. -benchmem -benchtime=10s ./test/benchmark/...

.PHONY: test-all
test-all: test test-integration test-e2e
	@echo "All tests completed!"

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(BUILDENVVAR) go test -cover -coverprofile=coverage.out -covermode=atomic ./pkg/...
	@echo "Coverage report generated: coverage.out"
	@echo "View coverage: go tool cover -html=coverage.out"

.PHONY: test-coverage-html
test-coverage-html: test-coverage
	@echo "Opening coverage report in browser..."
	go tool cover -html=coverage.out

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	GOOS=linux GOARCH=amd64 $(BUILDENVVAR) go build -ldflags '-w' -o bin/$(BINARY_NAME)-linux cmd/main.go
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f pkg/apis/scheduling/v1alpha1/zz_generated.deepcopy.go

.PHONY: verify
verify: generate
	@echo "Verifying generated code is up to date..."
	git diff --exit-code pkg/apis/
