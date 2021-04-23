# !/bin/sh
#  # git tag -a v0.0.3 -m "new release v0.0.3"
# export GITHUB_TOKEN="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"
# export TAG=$(git tag --no-merged)
# ~/go/bin/ghr -t $GITHUB_TOKEN -u JitenPalaparthi -r readyGo --replace --draft  $TAG dist/

# // git describe --abbrev=0 --tags


#!/bin/bash
echo "Github.Com personal token: $1"

#export GITHUB_TOKEN="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"

export GITHUB_TOKEN=$1

RES=$(git show-ref --tags)
if [ -z "$RES" ]; then
    NEW_TAG="v0.0.1"
else
    LATEST_TAG=$(git describe --tags --abbrev=0)
    IFS='.' read -r -a array <<< ${LATEST_TAG:1}
    one=${array[0]}
    two=${array[1]}
    three=${array[2]}

    if [ "$three" == "9" ]; then
        if [ "$two" == "9" ]; then
            three=0
            two=0
            ((one++))
        else
            ((two++))
            three=0
        fi
    elif [ "$two" == "9" ] && [ "$three" == "9" ]; then
        ((one++))
        two=0
    else
        ((three++))
    fi

    NEW_TAG="v${one}.${two}.${three}"
fi

echo $NEW_TAG

git tag -a $NEW_TAG -m "new release $NEW_TAG"

#export GITHUB_TOKEN="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"

goreleaser