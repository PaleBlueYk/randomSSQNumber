SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o randSSQ main.go
copy config\config.toml bin\config.toml
xcopy source bin\source /e