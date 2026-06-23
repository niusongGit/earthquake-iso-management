#!/bin/bash

# 查找并杀掉所有匹配 ./earthquake-iso-management 的旧进程
PIDS=$(ps -ef | grep "./earthquake-iso-management" | grep -v grep | awk '{print $2}')
if [ -n "$PIDS" ]; then
    echo "找到旧进程 PID: $PIDS，正在终止..."
    # 使用 kill -9 强制终止（谨慎使用，生产环境建议先尝试普通 kill）
    kill -9 $PIDS
    sleep 1  # 等待进程彻底退出
else
    echo "未找到旧进程"
fi

# 删除旧的日志文件
if [ -f "debug.log" ]; then
    rm -f debug.log
    echo "已删除旧日志文件 debug.log"
fi

mv earthquake-iso-management-new earthquake-iso-management

# 确保执行权限
chmod +x earthquake-iso-management


# 后台启动 AirGo（假设子命令为 start，日志重定向到 airgo.log）
nohup ./earthquake-iso-management > debug.log 2>&1 &
echo "地震国际标准管理平台 已启动，新进程 PID: $!"
