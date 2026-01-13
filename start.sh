# /bin/bash

echo "Start Robot Control System"

echo "Building artifact"
go build -o /home/undead/go/src/github.com/undeadpelmen/webrobot-robot/build/webrobot-robot /home/undead/go/src/github.com/undeadpelmen/webrobot-robot


echo "Starting app"
/home/undead/go/src/github.com/undeadpelmen/webrobot-robot/build/webrobot-robot
