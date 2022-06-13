$ErrorActionPreference = "Stop"
echo "=== building... ==="

mkdir bin -ErrorAction SilentlyContinue
go build -o bin ./...
go clean

echo "=== build done ==="
