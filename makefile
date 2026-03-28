SALES_VERSION := 0.0.1

run-sales:
	cd ./apis/services/sales && \
	go build -ldflags " \
		-X main.build=$(SALES_VERSION) \
		-X main.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	" && \
	./sales

help-sales:
	cd ./apis/services/sales && \
	go run main.go --help

build-sales:
	docker build \
		-f ./zarf/docker/dockerfile.sales \
		-t vishalgovind/sales \
		--build-arg BUILD_REF=$(SALES_VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

curl-test:
	curl -i -X GET http://localhost:3000/test

dev-apply-sales:
	kubectl apply -f ./zarf/k8s/sales/sales.yaml

dev-restart:
	kubectl rollout restart deployment sales

dev-logs:
	kubectl logs --selector app=sales --all-containers=true --tail=100 --max-log-requests=6

dev-status:
	watch kubectl get pods -o wide --all-namespaces --show-labels

# generate RSA private key: openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# generate RSA public key: openssl rsa -pubout -in private.pem -out public.pem
# generate RSA key pair: go run cmd/rsa-kpg/main.go