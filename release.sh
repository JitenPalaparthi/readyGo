#!/bin/bash
echo "Github.Com personal token: $1"
# Write a logic to validate the token. If the token is expired or not validated message and stop the further process.

# If there are any changes. Tagging will fail so. Test in the local repo that no changes are pending for commit. If there are any 
# Stop to procedd further

export GITHUB_TOKEN=$1

RES=$(git show-ref --tags)
UNCOMITTED=$(git status --porcelain)
if [ -z "$RES" ]; then
    NEW_TAG="v0.0.1"
else
if [ -z "$UNCOMITTED" ]; then
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

    echo $NEW_TAG

    # Write tag validation before running goreleaser.

    git tag -a $NEW_TAG -m "new release $NEW_TAG"

    rm -rf dist/

    goreleaser
    else 
echo "One or more changes are uncomitted;commit or stash them and try again"
fi
fi

 # git ls-remote --tags --refs --sort="v:refname" git://github.com/jitenpalaparthi/readygo.git | tail -n1 | sed 's/.*\///'

