set PATH=%CD%\wix\tools;%PATH%
go-msi make -p wix.json -a amd64 -m nuv.msi -l LICENSE --version 0.3.0 --src templates