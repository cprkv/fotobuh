$ErrorActionPreference = "Stop"

./build.ps1

echo "=== running... ==="
&"./bin/$args"
