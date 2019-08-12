@set cli=poe-stash-cli.exe
@set server=poe-stash-server.exe

@echo %cli%
@echo %server%

@cd cmd\cli
@go build -o %cli%
@cd ..\..
@move cmd\cli\%cli% .

@cd cmd\server
@go build -o %server%
@cd ..\..
@move cmd\server\%server% .

@set TEMPDIR=temp738
@rmdir /s /q %TEMPDIR%
@mkdir %TEMPDIR%
@xcopy %cli% %TEMPDIR%
@xcopy %server% %TEMPDIR%
@mkdir %TEMPDIR%\data
@mkdir %TEMPDIR%\data\template
@xcopy /s data\template %TEMPDIR%\data\template

@echo Set fso = CreateObject("Scripting.FileSystemObject") > _zipup.vbs
@echo InputFolder = fso.GetAbsolutePathName(WScript.Arguments.Item(0)) >> _zipup.vbs
@echo ZipFile = fso.GetAbsolutePathName(WScript.Arguments.Item(1)) >> _zipup.vbs
@echo CreateObject("Scripting.FileSystemObject").CreateTextFile(ZipFile, True).Write "PK" ^& Chr(5) ^& Chr(6) ^& String(18, vbNullChar) >> _zipup.vbs
@echo Set objShell = CreateObject("Shell.Application") >> _zipup.vbs
@echo Set source = objShell.NameSpace(InputFolder).Items >> _zipup.vbs
@echo objShell.NameSpace(ZipFile).CopyHere(source) >> _zipup.vbs
@echo ' Keep script waiting until compression is done
@echo Do Until objShell.NameSpace( ZipFile ).Items.Count = objShell.NameSpace( InputFolder ).Items.Count >> _zipup.vbs
@echo     WScript.Sleep 200 >> _zipup.vbs
@echo Loop >> _zipup.vbs

@CScript _zipup.vbs %TEMPDIR% "poe-stash-windows-amd64.zip"
@del _zipup.vbs
@del %cli%
@del %server%
@rmdir /s /q %TEMPDIR%

@echo "zip generated"
