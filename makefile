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

export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjgyMzRjOGM1LTA1MDgtNDMwMS1iZmRmLWRhMWQ1MTUzOTlhMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlcy1nbyBzZXJ2aWNlIiwic3ViIjoiNGZjODAwZDQtMmMzZC00NWZjLWEwZmEtOTI2MzY0NGY2ZGU3IiwiZXhwIjoxNzc0Nzc0MzkxLCJuYmYiOjE3NzQ3NzA3OTEsIlJvbGVzIjpbIkFETUlOIl19.DMe04U6D2GMERZRCrnnU4LO7nIdn7OBvY2FNkyLxPb8gOi7FaQJS5mJJMcH8EowwrvovRtohwJE8aEkDWPdfpv7TUGhvfxE_ToX6Xo1R70vC1cudHA2BmgkYiyknFLfGbjYpoy45KW2hIAffQASvTb3VYcibX6eOfRAPv9G1owo2oS0Wk4VVMIipErNe3RRZz0z2ns_m8QYhYIajBLT9VUkVoBnbUbj2JcFwzri58-VZm6t0tL3wTZmYAdFn_efi9C15KDU1r7QzW6omZUnGEoDH6cFGojo4jyzuD1ZJkvHmSZCWjNgTpDlw9GNuyjRl7WuCwZ2kufDfc3x29ioSUw

curl-liveness:
	curl -i -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/liveness