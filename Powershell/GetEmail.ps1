param(
    [string]$Username
)

if (-not $Username) {
    Write-Host "Please provide a username using the -Username flag."
    exit
}

# Get the user from the local domain controller
$user = Get-ADUser -Identity $Username -Properties EmailAddress

if (-not $user) {
    Write-Host "User '$Username' not found."
    exit
}

$email = $user.EmailAddress
if ($email) {
    Write-Output "$email"
} else {
    Write-Output "USer does not have an email address."
}
