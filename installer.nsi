OutFile "build\installer\east-kishu-kot-setup.exe"
InstallDir "$PROGRAMFILES\EastKishuKOT"
Icon "ki-ho-multisize.ico"

; レジストリ登録用
InstallDirRegKey HKCU "Software\EastKishuKOT" "Install_Dir"

; アンインストーラー出力先
!define UNINSTALL_EXE "$INSTDIR\Uninstall.exe"

Section "Install"
  SetOutPath $INSTDIR

  ; EXE 本体
  File "build\bin\east-kishu-kot.exe"

  ; アイコン（ショートカット用）
  File "ki-ho-multisize.ico"

  ; アンインストーラーを生成
  WriteUninstaller "${UNINSTALL_EXE}"

  ; デスクトップショートカット作成
  CreateShortCut "$DESKTOP\EastKishuKOT.lnk" "$INSTDIR\east-kishu-kot.exe" "" "$INSTDIR\ki-ho-multisize.ico" 0

  ; アンインストーラーのショートカット（コントロールパネル登録用）
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\EastKishuKOT" "DisplayName" "東紀州KOT自動入力ツール"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\EastKishuKOT" "UninstallString" "${UNINSTALL_EXE}"
  WriteRegStr HKCU "Software\EastKishuKOT" "Install_Dir" "$INSTDIR"
SectionEnd

Section "Uninstall"
  ; EXE / ICO 削除
  Delete "$INSTDIR\east-kishu-kot.exe"
  Delete "$INSTDIR\ki-ho-multisize.ico"

  ; デスクトップショートカット削除
  Delete "$DESKTOP\EastKishuKOT.lnk"

  ; 自分自身を削除
  Delete "${UNINSTALL_EXE}"

  ; フォルダ削除
  RMDir "$INSTDIR"

  ; レジストリ掃除
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\EastKishuKOT"
  DeleteRegKey HKCU "Software\EastKishuKOT"
SectionEnd
