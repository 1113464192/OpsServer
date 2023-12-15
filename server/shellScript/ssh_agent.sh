#!/bin/bash
cd $(dirname $0) &&\
eval $(ssh-agent -a ${agent_sock_path}) &&\
./ssh_agent.exp &&\
lsof ${agent_sock_path} | grep -v "COMMAND" |awk '{printf "%s%d", "Success\nPID:\n", $2}' &&\
/bin/rm ${id_key_path}