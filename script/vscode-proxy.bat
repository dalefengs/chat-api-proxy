@echo off
setlocal

set host=https://proxy.fungs.cn

set extensions_dir=%userprofile%\.vscode\extensions
if not exist "%extensions_dir%" (
  echo ERROR: VSCode extensions directory not found!
  pause
  exit /b 1
)

for /f "tokens=*" %%a in ('dir /b /ad "%extensions_dir%" ^| findstr /r /c:"^github\.copilot-[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*$"') do (
  set copilot_dir=%%a
  echo find copilot extension: %%a
  goto :found
)

echo ERROR: Copilot extension not found!
pause
exit /b 1

:found
set copilot_dir=%extensions_dir%\%copilot_dir%
echo find copilot_dir %copilot_dir%
set extension_file=%copilot_dir%\dist\extension.js
if not exist "%extension_file%" (
  echo ERROR: Copilot extension entry file not found!
  pause
  exit /b 1
)
echo please be patient...

set tmp_file=%copilot_dir%\dist\extension.js.tmp
powershell -Command "(Get-Content '%extension_file%') -replace 'Gu.Utils.joinPath\(r,\"/copilot_internal/v2/token\"\).toString\(\)', '\"%host%/copilot_internal/v2/token\"' | Out-File -encoding ASCII '%tmp_file%'"

move /y "%tmp_file%" "%extension_file%" > nul


echo copilot extension configuration complete.

echo copilot chat extension configuring...

for /f "tokens=*" %%a in ('dir /b /ad "%extensions_dir%" ^| findstr /r /c:"^github\.copilot-chat-[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*$"') do (
  set copilot_chat_dir=%%a
  echo find copilot chat extension: %%a
  goto :foundChat
)

echo ERROR: Copilot Chat extension not found!
pause
exit /b 1


:foundChat
set copilot_chat_dir=%extensions_dir%\%copilot_chat_dir%
echo find copilot_chat_dir %copilot_chat_dir%
set extension_chat_file=%copilot_chat_dir%\dist\extension.js
if not exist "%extension_chat_file%" (
  echo ERROR: Copilot chat extension entry file not found!
  pause
  exit /b 1
)

set tmp_chat_file=%copilot_chat_dir%\dist\extension.js.tmp
powershell -Command "(Get-Content '%extension_chat_file%') -replace '\"https://api.githubcopilot.com\"', '\"%host%\"' | Out-File -encoding ASCII '%tmp_chat_file%'"

move /y "%tmp_chat_file%" "%extension_chat_file%" > nul



echo done. please restart your vscode.
pause
