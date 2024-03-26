# Path to the file you want to monitor for changes
$filePath = "S:\HR\newusers.csv"

# Path to store the hash
$hashFilePath = "C:\Users\Administrator\Desktop\csvhash.txt"

# Path to the script you want to run if the file changes
$scriptToRun = "C:\Users\Administrator\Desktop\Create_user_csv.ps1"

# Check if the hash file exists, if not create one
if (-not (Test-Path $hashFilePath)) {
    $initialHash = Get-FileHash -Path $filePath
    $initialHash | Out-File -FilePath $hashFilePath
}

# Read the initial hash from the hash file
$storedHash = Get-Content $hashFilePath

# Get the current hash of the file
$currentHash = Get-FileHash -Path $filePath

# Compare the stored hash with the current hash
if ($currentHash.Hash -ne $storedHash) {
    Write-Host "File has changed. Running script..."
    
    # Run the script
    & $scriptToRun
    
    # Update the stored hash with the new hash
    $currentHash.Hash | Out-File -FilePath $hashFilePath -Force
} else {
    Write-Host "File has not changed."
}
