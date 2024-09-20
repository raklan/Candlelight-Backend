@echo OFF

go test -v -coverprofile cover.out %1
go tool cover -html cover.out -o cover.html
rm cover.out
start cover.html