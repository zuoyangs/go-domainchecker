#!/bin/bash  
  
# 检查进程是否在运行  
is_process_running() {  
    if ps -ef |grep go-domainchecker |grep -v grep > /dev/null; then  
        return 0  # 进程正在运行，返回成功  
    else  
        return 1  # 进程没有运行，返回失败  
    fi  
}  
  
# 启动逻辑  
start() {  
    cd "$(dirname "$0")"  
    nohup ./go-domainchecker &  
    sleep 1  
  
    if is_process_running; then  
        echo "go-domainchecker 启动成功"  
    else  
        echo "启动失败，请检查"  
    fi  
}  
  
# 停止逻辑  
stop() {  
    if is_process_running; then  
        kill -9 $(ps -ef |grep go-domainchecker |grep -v grep | awk '{print $2}')  
        echo "go-domainchecker 已停止"  
    else  
        echo "没有找到名为 go-domainchecker 的进程"  
    fi  
}

# 根据命令行参数执行相应的操作  
case "$1" in  
    start)  
        start  
        ;;  
    stop)  
        stop  
        ;;  
    *)  
        echo "Usage: $0 {start|stop}"  
        exit 1  
esac
