SALES_VERSION := 0.0.1
AUTH_VERSION := 0.0.1

run-sales:
	cd ./apis/services/sales && \
	go build -ldflags " \
		-X main.build=$(SALES_VERSION) \
		-X main.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	" && \
	./sales

run-auth:
	cd ./apis/services/auth && \
	go build -ldflags " \
		-X main.build=$(AUTH_VERSION) \
		-X main.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	" && \
	./auth

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

export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjgyMzRjOGM1LTA1MDgtNDMwMS1iZmRmLWRhMWQ1MTUzOTlhMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlcy1nbyBzZXJ2aWNlIiwic3ViIjoiNGZjODAwZDQtMmMzZC00NWZjLWEwZmEtOTI2MzY0NGY2ZGU3IiwiZXhwIjoxNzc0NzkwNzc2LCJuYmYiOjE3NzQ3ODcxNzYsIlJvbGVzIjpbIkFETUlOIl19.B5EV8cCkdywWbNcjjzU_Sn9SZHP96FOF26ZOHMVIOUQecS2Zj2Dl6VDqkiFZCT6aSK-l5OsNQFf0bSPgBeQMsQRD2jCARZy9cvx6dXnpXz_oP3BLEWbFO8Ldx75DoWKShnhR4inirO0_cbPisDd3wX1CNN1209whx5VAZcPyloRo7YD_oMQPeYHqQWyEwEIWyWk5aRzsjNkarT1rSLMlWvgjLWP1nkK4lNw66xamO97VpVOUAMTj49m_wicfvZbPLwbtrW2trG3OTYg4CpbWDwmLIzgY03H9Bk0FoMgr1iC_xcBejT4tDjuk3tnZ6-wp8Amn432vX-1JxbaTCaye4g

curl-liveness:
	curl -i -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/liveness

gen-token:
	go run cmd/rsa-kpg/main.go