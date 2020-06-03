#!/bin/bash
set -e

CurDir=$(cd $(dirname $0);pwd)

start(){
  pid=$(ps -ef | grep '\/tbaoke' | grep -v grep | awk '{print $2}')
  if [ "x$pid" != "x" ]; then
    kill $pid
  fi
  timestamp=$(date +%s)
  sh -c "sed -i 's/main\.js.*\"/main.js?${timestamp}\"/' public/index.html"
  sh -c "sed -i 's/main\.css.*\"/main.css?${timestamp}\"/' public/index.html"
  nohup $CurDir/tbaoke >> $CurDir/logs/access.log 2>&1 &
}

stop(){
  pid=$(ps -ef | grep '\/tbaoke' | grep -v grep | awk '{print $2}')
  if [ "x$pid" != "x" ]; then
    kill $pid
  fi
}

case "$1" in
    start)
        start
        echo "started"
    ;;
    stop)
        stop
        echo "stopped"
    ;;
esac
