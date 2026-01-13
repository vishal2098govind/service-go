SALES_VERSION := 0.0.1

run_sales:
	cd ./apis/services/sales && \
	go build -ldflags " \
		-X main.build=$(SALES_VERSION) \
		-X main.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	" && \
	./sales


build_sales:
	docker build \
		-f ./zarf/docker/dockerfile.sales \
		-t vishalgovind/sales \
		--build-arg BUILD_REF=$(SALES_VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

dev-status:
	watch kubectl get pods -o wide --all-namespaces --show-labels