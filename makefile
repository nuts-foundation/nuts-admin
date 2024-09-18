.PHONY: dev

run-generators: gen-api

install-tools:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

gen-api:
	oapi-codegen -generate server,types -package api api/api.yaml > api/generated.go

dev:
	make -j3 watch run-nuts-node run-api

dev-no-api:
	make -j3 watch run-nuts-node

watch:
	npm install
	npm run watch

run-nuts-node:
	docker compose pull
	docker compose up --wait

run-api:
	go run . live --configfile=deploy/admin.config.yaml

docker:
	docker build -t nutsfoundation/nuts-admin:main .
