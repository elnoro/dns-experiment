.PHONY: build-image
build-image:
	podman build . -t foxylockv2

.PHONY: run-image
run-image:
	podman run -p 53:1053 --env-file=.env --rm -v /etc/foxylock:/etc/foxylock foxylockv2 -conf /etc/foxylock/Corefile

.PHONY: test
test:
	go test ./... -timeout=30s -race

coverage:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html
