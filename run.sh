#!/bin/bash
set -e

CurDir=$(cd $(dirname $0);pwd)

isMacOS=$(uname -a|grep -i Darwin|wc -l)
if [ $isMacOS -gt 0 ]; then
    echo "MacOS 请手动启动"
    exit 1
fi

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
    start|restart)
        start
        echo "started"
    ;;
    stop)
        stop
        echo "stopped"
    ;;
    *)
        echo "usage: ./run.sh {start|restart|stop}"
esac
