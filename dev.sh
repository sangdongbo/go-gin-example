#!/bin/bash

echo "选择热重载工具："
echo "1. Air (功能丰富)"
echo "2. Fresh (最快)"
echo "3. Reflex (轻量)"
echo "4. 原生 Go run (最简单)"

read -p "请选择 (1-4): " choice

case $choice in
    1)
        echo "启动 Air..."
        if ! command -v air &> /dev/null; then
            echo "安装 Air..."
            go install github.com/cosmtrek/air@latest
        fi
        air
        ;;
    2)
        echo "启动 Fresh..."
        if ! command -v fresh &> /dev/null; then
            echo "安装 Fresh..."
            go install github.com/gravityblast/fresh@latest
        fi
        fresh
        ;;
    3)
        echo "启动 Reflex..."
        if ! command -v reflex &> /dev/null; then
            echo "安装 Reflex..."
            go install github.com/cespare/reflex@latest
        fi
        reflex -c reflex.conf
        ;;
    4)
        echo "使用原生 Go run..."
        echo "注意：需要手动重启"
        go run main.go
        ;;
    *)
        echo "无效选择，使用 Air..."
        air
        ;;
esac