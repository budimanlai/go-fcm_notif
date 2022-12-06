#!/bin/bash

APP="myapp"
REMOTEDIR="/home/apps/myapp"
LOGIN="dev@localhost"
PORT="22"

if [ -z "$VERSION" ]
then
    echo "Usage: VERSION=xx ${0}"
else
    GIT=`git log | head -n 1 | cut  -f 2 -d ' ' | head -c 8` 
    LOCALFILE="dist/$APP-$GIT"
    REMOTEFILE="dist/$APP-$VERSION-$GIT"
    pv ${LOCALFILE} | ssh $LOGIN -p $PORT "cd $REMOTEDIR; mkdir -p dist; cat > $REMOTEFILE; chmod +x $REMOTEFILE; rm -f $APP-$VERSION; ln -s $REMOTEFILE $APP-$VERSION"
fi