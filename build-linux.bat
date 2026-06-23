@echo off
chcp 65001 >nul
echo ========================================
echo  地震国际标准管理平台 - Linux打包脚本
echo ========================================
echo.

echo [1/3] 构建前端...
cd web
call npm run build
if %errorlevel% neq 0 (
    echo 前端构建失败！
    pause
    exit /b 1
)
cd ..
echo 前端构建完成。
echo.

echo [2/3] 交叉编译Linux二进制文件...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o earthquake-iso-management-new .
if %errorlevel% neq 0 (
    echo Go编译失败！
    pause
    exit /b 1
)
echo.

echo [3/3] 清理环境变量...
set GOOS=
set GOARCH=
set CGO_ENABLED=
echo.

echo ========================================
echo  打包完成！
echo  输出文件: earthquake-iso-management-new
echo  部署方式: 上传至Linux服务器后执行
echo    chmod +x earthquake-iso-management-new
echo    ./earthquake-iso-management-new
echo ========================================
pause
