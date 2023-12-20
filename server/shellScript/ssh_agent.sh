#!/bin/bash
cd $(dirname $0) &&\
eval $(ssh-agent -a ${agent_sock_path}) &&\
if [[ -n ${id_key_passphrase} ]]
then
    ./ssh_agent.exp
else
     timeout 5 ssh-add ${id_key_path}
fi &&\
&&\
lsof ${agent_sock_path} | grep -v "COMMAND" | awk '{print $2}' | sort | uniq |  awk '{printf "%s%d\n", "Success\nPID:\n", $1}' &&\
/bin/rm ${id_key_path}