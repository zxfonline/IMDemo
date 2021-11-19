@echo off
set GOPATH=%~dp0
set GO=go.exe

set GOOS=windows
set GOARCH=amd64

@REM go clean -cache

@REM echo fmt all src code
%GO% fmt ./src/...

@REM rd /S/Q  bin pkg runtime
rd /S/Q runtime

@set ROOT=%CD%

echo start build server

@cd .\src\im_main
%GO% install github.com/zxfonline/IMDemo/im_main
@cd %ROOT%

echo d|xcopy src\\runtime runtime /e /k

echo ok
pause