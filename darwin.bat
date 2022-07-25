set PLUGIN=hpcr
set TF_BASE=build/terraform/plugins/www.ibm.com/local/%PLUGIN%
set VERSION=0.0.2

set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=arm64

set OUT=%TF_BASE%/%VERSION%/%GOOS%_%GOARCH%/terraform-provider-%PLUGIN%_v%VERSION%

go build -ldflags "-s -w" -trimpath -o %OUT%

