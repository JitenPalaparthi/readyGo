 #!/bin/sh
 # git tag -a v0.0.3 -m "new release v0.0.3"

#export GITHUB_TOKEN="10cbdbac9fb8b7739efb74b896db60489dc30b9a"
export GITHUB_TOKEN="74e5670ae4c9271b5964d322a6e941c3aff98c21"
export TAG=v0.0.3
~/go/bin/ghr -t $GITHUB_TOKEN -u JitenPalaparthi -r readyGo --replace --draft  $TAG dist/