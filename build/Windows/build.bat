@echo off
setlocal

set "output_path=%~1"
rem Write to .\bin if no custom path specified
if "%output_path%"=="" set "output_path=.\bin\"
rem Ensure the path ends with '/'
if not "%output_path:~-1%"=="\" set "output_path=%output_path%\"
go build -o "%output_path%sentinel.exe"

endlocal
