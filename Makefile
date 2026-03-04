generate:
	mkdir -p pkg/api
	oapi-codegen -generate types,chi-server -package api -o pkg/api/api.gen.go api.yaml
build: generate
	go build -o app main.go