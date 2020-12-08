 #!/bin/sh
 # git tag -a v0.0.3 -m "new release v0.0.3"

export GITHUB_TOKEN=10cbdbac9fb8b7739efb74b896db60489dc30b9a
export TAG=v0.0.3
~/go/bin/ghr -t $GITHUB_TOKEN -u JitenPalaparthi -r readyGo --replace --draft  $TAG dist/
