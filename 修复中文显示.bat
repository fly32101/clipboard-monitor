@echo off
chcp 65001 >nul
title 修复中文显示问题

echo ========================================
echo           中文显示修复工具
echo ========================================
echo.

echo 正在检测系统配置...
echo.

REM 检查是否已启用 UTF-8 支持
reg query "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Nls\CodePage" /v ACP | find "65001" >nul
if %errorlevel% == 0 (
    echo ✅ 系统已启用 UTF-8 支持
    goto :run_program
)

echo ⚠️  系统未启用 UTF-8 支持
echo.
echo 选择修复方案:
echo 1. 启用系统 UTF-8 支持 (推荐，需要重启)
echo 2. 仅为此程序设置环境变量 (临时)
echo 3. 直接运行程序
echo.
set /p choice=请选择 (1-3): 

if "%choice%"=="1" goto :enable_utf8
if "%choice%"=="2" goto :set_env
if "%choice%"=="3" goto :run_program

:enable_utf8
echo.
echo 正在启用系统 UTF-8 支持...
echo 注意: 此操作需要管理员权限和重启电脑
echo.
pause
echo.

REM 尝试启用 UTF-8 支持
reg add "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Nls\CodePage" /v ACP /t REG_SZ /d 65001 /f >nul 2>&1
reg add "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Nls\CodePage" /v OEMCP /t REG_SZ /d 65001 /f >nul 2>&1

if %errorlevel% == 0 (
    echo ✅ UTF-8 支持已启用
    echo 请重启电脑后再运行程序
    pause
    exit
) else (
    echo ❌ 需要管理员权限，请以管理员身份运行此脚本
    echo 或手动在设置中启用 UTF-8 支持
    pause
    goto :set_env
)

:set_env
echo.
echo 设置临时环境变量...
set LANG=zh_CN.UTF-8
set LC_ALL=zh_CN.UTF-8
set LC_CTYPE=zh_CN.UTF-8
set FYNE_FONT=
set FYNE_THEME=light
echo ✅ 环境变量已设置

:run_program
echo.
echo 启动剪贴板监听器...
if exist "clipboard-monitor.exe" (
    clipboard-monitor.exe
) else (
    echo ❌ 找不到 clipboard-monitor.exe
    echo 请先编译程序或将此脚本放在程序同一目录
    pause
)

echo.
echo 程序已退出
pause
