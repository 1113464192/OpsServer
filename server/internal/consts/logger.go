package consts

const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

const (
	GetHostInfoCmd = `systemDisk=$(df -Th | awk '{if ($NF=="/")print$(NF-2)}' | grep -Eo "[0-9]+")
				 dataDisk=$(df -Th | awk '{if ($NF=="/data")print$(NF-2)}' | grep -Eo "[0-9]+")
				 if [[ -z ${dataDisk} ]];then
				 	dataDisk="-1"
				 fi
				 mem=$(free -m | awk '/Mem/{print $NF}')
				 iowait=$(iostat | awk '/avg-cpu:/ {getline; print $(NF-2)}')
				 idle=$(iostat | awk '/avg-cpu:/ {getline; print $(NF)}')
				 load=$(uptime | awk -F"[, ]+" '{print $(NF-1)}')
				 echo "$systemDisk $dataDisk $mem $iowait $idle $load" | awk '{print $1,$2,$3,$4,$5,$6}'`
)
