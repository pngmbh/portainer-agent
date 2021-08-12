IMG=r.planetary-quantum.com/quantum-public/portainer-agent

.PHONY: build
build: dist
	GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build --installsuffix cgo --ldflags '-s' "cmd/agent/main.go"
	mv main dist/agent
	docker build -t "$(IMG):latest" -f build/linux/Dockerfile .

.PHONY: push
push: guard-VERSION
	docker tag "$(IMG):latest" "$(IMG):$(VERSION)"
	docker push "$(IMG):$(VERSION)"
	docker push "$(IMG):latest"

dist:
	mkdir -p dist

guard-%:
	@ if [ "${${*}}" = "" ]; then \
        echo "Environment variable $* not set"; \
        exit 1; \
    fi
