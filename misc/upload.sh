 #!/bin/sh
 # git tag -a v0.0.3 -m "new release v0.0.3"
export GITHUB_TOKEN="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"
export TAG=v0.0.3
~/go/bin/ghr -t $GITHUB_TOKEN -u JitenPalaparthi -r readyGo --replace --draft  $TAG dist/