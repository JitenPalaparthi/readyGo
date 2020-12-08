 #!/bin/sh

 # ~/go/bin/gox -os="linux darwin windows" -arch="amd64" -output="dist/readyGo_{{.OS}}_{{.Arch}}"

~/go/bin/gox -output="dist/readyGo_v0.0.3_{{.OS}}_{{.Arch}}"

# For all os 

 cd dist;gzip *