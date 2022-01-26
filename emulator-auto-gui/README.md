##static build

`go build -tags static --ldflags '-extldflags="-static"' -ldflags="-H windowsgui -w -s" -o game-auto.exe .`
