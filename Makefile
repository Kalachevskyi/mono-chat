
.PHONY: test
test:
	 @ echo "-> Start tests ..."
	 @ richgo test -v ./... -coverprofile=cover-all.out
	 @ richgo tool cover -func=cover-all.out
	 @ echo "-> Done!!!"