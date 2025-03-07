#!/bin/bash

SERVER_NAME="MygateBot"

Start_Bot() {

  # nohup ./$SERVER_NAME >/dev/null 2> nohup.out &
  nohup ./"${SERVER_NAME}" > server.log  2>&1 &
  echo "MygateBot Server is started."
  >server.log
  tail -10f server.log
}

Stop_Bot() {
  # kill $1 if it exists.
  PID_LIST=$(ps -ef | grep $SERVER_NAME | grep -v grep | awk '{printf "%s ", $2}')
  for PID in $PID_LIST; do
    if kill -9 $PID; then
      echo "Process $one($PID) was stopped at " $(date)
      echo "MygateBot Server is stoped."
    fi
  done
}

Status_Bot() {
  PID_NUM=$(ps -ef | grep $SERVER_NAME | grep -v grep | wc -l)
  if [ $PID_NUM -gt 0 ]; then
    {
      echo "MygateBot server is started."
    }
  else
    {
      echo "MygateBot server is stoped."
    }
  fi
}

case "$1" in
'start')
  Start_Bot
  ;;
'stop')
  Stop_Bot
  ;;
'restart')
  Stop_Bot
  Start_Bot
  ;;
'status')
  Status_Bot
  ;;
*)
  echo "Usage: $0 {start|stop}"
  echo "  start : To start the application of MygateBot"
  echo "  stop  : To stop the application of MygateBot"
  echo "  restart  : To restart the application of MygateBot"
  echo "  status  : To view status the application of MygateBot"
  RETVAL=1
  ;;
esac

exit 0