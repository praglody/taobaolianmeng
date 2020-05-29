#!/bin/bash
set -e

CurDir=$(cd $(dirname $0);pwd)

start(){
  pid=$(ps -ef | grep '\/tbaoke' | grep -v grep | awk '{print $2}')
  if [ "x$pid" != "x" ]; then
    kill $pid
  fi
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
