param (
    [Parameter(Mandatory=$true)]
    [string]$Username,

    [Parameter(Mandatory=$true)]
    [string]$Password
)

# Convert plain text password to SecureString
$newPassword = ConvertTo-SecureString -String $password -AsPlainText -Force

# Reset user password and force password change at next logon
Set-ADAccountPassword -Identity $username -NewPassword $newPassword -Reset

Write-Host "Password for $username has been reset successfully. The user will be prompted to change the password at next login."
