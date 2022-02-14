GO=go
GOCOVER=$(GO) tool cover

.PHONY: test/cover
test/cover:
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out