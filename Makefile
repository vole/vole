all:
	go run src/vole/vole.go

assets:
	go-bindata -prefix "static" -o="src/lib/assets/main.go" -pkg="assets" -debug="false" -ignore=".DS_Store" static/...
