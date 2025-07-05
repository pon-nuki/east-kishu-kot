OutFile "build\installer\east-kishu-kot-setup.exe"
InstallDir "$PROGRAMFILES\EastKishuKOT"
Icon "ki-ho-multisize.ico"

; ���W�X�g���o�^�p
InstallDirRegKey HKCU "Software\EastKishuKOT" "Install_Dir"

; �A���C���X�g�[���[�o�͐�
!define UNINSTALL_EXE "$INSTDIR\Uninstall.exe"

Section "Install"
  SetOutPath $INSTDIR

  ; EXE �{��
  File "build\bin\east-kishu-kot.exe"

  ; �A�C�R���i�V���[�g�J�b�g�p�j
  File "ki-ho-multisize.ico"

  ; �A���C���X�g�[���[�𐶐�
  WriteUninstaller "${UNINSTALL_EXE}"

  ; �f�X�N�g�b�v�V���[�g�J�b�g�쐬
  CreateShortCut "$DESKTOP\EastKishuKOT.lnk" "$INSTDIR\east-kishu-kot.exe" "" "$INSTDIR\ki-ho-multisize.ico" 0

  ; �A���C���X�g�[���[�̃V���[�g�J�b�g�i�R���g���[���p�l���o�^�p�j
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\EastKishuKOT" "DisplayName" "���I�BKOT�������̓c�[��"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\EastKishuKOT" "UninstallString" "${UNINSTALL_EXE}"
  WriteRegStr HKCU "Software\EastKishuKOT" "Install_Dir" "$INSTDIR"
SectionEnd

Section "Uninstall"
  ; EXE / ICO �폜
  Delete "$INSTDIR\east-kishu-kot.exe"
  Delete "$INSTDIR\ki-ho-multisize.ico"

  ; �f�X�N�g�b�v�V���[�g�J�b�g�폜
  Delete "$DESKTOP\EastKishuKOT.lnk"

  ; �������g���폜
  Delete "${UNINSTALL_EXE}"

  ; �t�H���_�폜
  RMDir "$INSTDIR"

  ; ���W�X�g���|��
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\EastKishuKOT"
  DeleteRegKey HKCU "Software\EastKishuKOT"
SectionEnd
