@echo off
setlocal

echo === ���[�ق���A�C�R���𒍓���... ===

:: manifest.json ���R�s�[�i�㏑���j
copy /Y meta\manifest.json build\bin\manifest.json >nul

echo === Wails�r���h�J�n ===
wails build -clean

:: �r���h���s���͒��~
if errorlevel 1 (
  echo �r���h�Ɏ��s���܂����B�A�C�R���͒�������Ă��܂���B
  pause
  exit /b 1
)

:: Resource Hacker �̃p�X
set RH="C:\Program Files (x86)\Resource Hacker\ResourceHacker.exe"
:: .ico �t�@�C���p�X
set ICO=ki-ho-multisize.ico
:: �Ώ� exe
set EXE=build\bin\east-kishu-kot.exe
:: �o�� exe�i����j
set OUT=build\bin\east-kishu-kot.exe

:: Resource Hacker �̑��݊m�F
if not exist %RH% (
  echo ResourceHacker.exe ��������܂���F %RH%
  pause
  exit /b 1
)

:: .ico �t�@�C���̑��݊m�F
if not exist %ICO% (
  echo �A�C�R���t�@�C����������܂���F %ICO%
  pause
  exit /b 1
)

echo === Resource Hacker�ŃA�C�R���𒍓���... ===

:: MAINICON�i��ʓI�ȃA�v���p�A�C�R���j
%RH% -open %EXE% -save %OUT% -action addoverwrite -res %ICO% -mask ICONGROUP,MAINICON,

:: �^�X�N�o�[�E�s���~�ߓ��̂��� ICON �P�̂�����
%RH% -open %EXE% -save %OUT% -action addoverwrite -res %ICO% -mask ICON,1,

:: --- inject-icon.rc �� inject-icon.res �ɕϊ��iGoRC ���K�v�j ---
set GORC="C:\Program Files (x86)\GoRC\GoRC.exe"
if not exist %GORC% (
  echo GoRC ��������܂���F %GORC%
  pause
  exit /b 1
)

%GORC% /r inject-icon.rc

if not exist inject-icon.res (
  echo inject-icon.res �̐����Ɏ��s���܂����B
  pause
  exit /b 1
)

:: --- Resource Hacker �� .res �𒍓� ---
%RH% -open %EXE% -save %OUT% -action addoverwrite -res inject-icon.res

if errorlevel 1 (
  echo �A�C�R�������Ɏ��s���܂����B
  pause
  exit /b 1
) else (
  echo �A�C�R���� exe �ɒ�������܂����I�iMAINICON + ICON�j
)

echo.
echo === NSIS�C���X�g�[�����쐬��... ===

:: �o�̓t�H���_�쐬
mkdir build\installer >nul 2>nul

:: NSIS ���s�p�X
set MAKENSIS="C:\Program Files (x86)\NSIS\makensis.exe"

if not exist %MAKENSIS% (
  echo NSIS��������܂���B�C���X�g�[���͍쐬����܂���B
  pause
  exit /b 1
)

%MAKENSIS% installer.nsi

if errorlevel 1 (
  echo �C���X�g�[���쐬�Ɏ��s���܂����B
) else (
  echo �C���X�g�[���� build\installer �t�H���_�ɍ쐬����܂����I
)

echo.
echo === �S���������I����ꂳ�܂ł��� ===
pause
endlocal
