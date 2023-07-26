# Execute on powershell 7 (pwsh)
# Remove service
#sc.exe delete rccmd
Remove-Service -Name rccmd

# Create service
New-Service -Name "rccmd" -BinaryPathName "C:\Users\next\rccmdServer.exe -debug -src 192.168.1.142"
