
test:
	 @ echo "-> Start tests ..."
	 @ richgo test -v ./... -coverprofile=cover-all.out
	 @ richgo tool cover -func=cover-all.out
	 @ echo "-> Done!!!"
.PHONY: test

lint:
	@ echo "-> Running linters ..."
	@ golangci-lint run --exclude-use-default=false
	@ echo "-> Done!"
.PHONY: lint

ci: test lint