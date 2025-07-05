@echo off
setlocal

echo === きーほくんアイコンを注入中... ===

:: manifest.json をコピー（上書き）
copy /Y meta\manifest.json build\bin\manifest.json >nul

echo === Wailsビルド開始 ===
wails build -clean

:: ビルド失敗時は中止
if errorlevel 1 (
  echo ビルドに失敗しました。アイコンは注入されていません。
  pause
  exit /b 1
)

:: Resource Hacker のパス
set RH="C:\Program Files (x86)\Resource Hacker\ResourceHacker.exe"
:: .ico ファイルパス
set ICO=ki-ho-multisize.ico
:: 対象 exe
set EXE=build\bin\east-kishu-kot.exe
:: 出力 exe（同一）
set OUT=build\bin\east-kishu-kot.exe

:: Resource Hacker の存在確認
if not exist %RH% (
  echo ResourceHacker.exe が見つかりません： %RH%
  pause
  exit /b 1
)

:: .ico ファイルの存在確認
if not exist %ICO% (
  echo アイコンファイルが見つかりません： %ICO%
  pause
  exit /b 1
)

echo === Resource Hackerでアイコンを注入中... ===

:: MAINICON（一般的なアプリ用アイコン）
%RH% -open %EXE% -save %OUT% -action addoverwrite -res %ICO% -mask ICONGROUP,MAINICON,

:: タスクバー・ピン止め等のため ICON 単体も注入
%RH% -open %EXE% -save %OUT% -action addoverwrite -res %ICO% -mask ICON,1,

:: --- inject-icon.rc → inject-icon.res に変換（GoRC が必要） ---
set GORC="C:\Program Files (x86)\GoRC\GoRC.exe"
if not exist %GORC% (
  echo GoRC が見つかりません： %GORC%
  pause
  exit /b 1
)

%GORC% /r inject-icon.rc

if not exist inject-icon.res (
  echo inject-icon.res の生成に失敗しました。
  pause
  exit /b 1
)

:: --- Resource Hacker で .res を注入 ---
%RH% -open %EXE% -save %OUT% -action addoverwrite -res inject-icon.res

if errorlevel 1 (
  echo アイコン注入に失敗しました。
  pause
  exit /b 1
) else (
  echo アイコンが exe に注入されました！（MAINICON + ICON）
)

echo.
echo === NSISインストーラを作成中... ===

:: 出力フォルダ作成
mkdir build\installer >nul 2>nul

:: NSIS 実行パス
set MAKENSIS="C:\Program Files (x86)\NSIS\makensis.exe"

if not exist %MAKENSIS% (
  echo NSISが見つかりません。インストーラは作成されません。
  pause
  exit /b 1
)

%MAKENSIS% installer.nsi

if errorlevel 1 (
  echo インストーラ作成に失敗しました。
) else (
  echo インストーラが build\installer フォルダに作成されました！
)

echo.
echo === 全処理完了！お疲れさまでした ===
pause
endlocal
