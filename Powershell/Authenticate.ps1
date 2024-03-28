# Define parameters for username and password
param (
    [Parameter(Mandatory=$true)]
    [string]$Username,

    [Parameter(Mandatory=$true)]
    [string]$Password
)

# Create a credential object from the provided username and password
$credentials = New-Object System.Management.Automation.PSCredential ($Username, (ConvertTo-SecureString -String $Password -AsPlainText -Force))

# Validate the credentials against Active Directory
try {
    $user = Get-ADUser -Credential $credentials -Identity $Username -ErrorAction Stop
    Write-Host "Authentication successful"
} catch {
    Write-Host "Authentication failed. Please check your username and password." -ForegroundColor Red
}
