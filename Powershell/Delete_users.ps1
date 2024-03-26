# Import the Active Directory module
Import-Module ActiveDirectory


$logFilePath = "C:\Users\Administrator\Desktop\log.txt"

# Define the time threshold for deletion (1 hour)
$threshold = (Get-Date).AddHours(-1)
#$threshold = (Get-Date).AddMinutes(-2)

# Setting home folder and profile path
$homeFolder = "\\WIN-ULE5N582EEV\Homefolder\$username"
$profilePath = "\\WIN-ULE5N582EEV\User-profiles\$username"

# Check for users disabled for more than 1 hour and delete them
$disabledUsers = Get-ChildItem -Path "C:\Users\Administrator\Desktop\DisableTime" -Filter *.txt
foreach ($userFile in $disabledUsers) {
    $username = $userFile.BaseName
    $timestamp = Get-Date -Format "M/d/yyyy HH:mm:ss"
    $disableTimeString = Get-Content $userFile.FullName
    $disableTime = [DateTime]::ParseExact($disableTimeString, "M/d/yyyy h:mm:ss tt", $null)

    if ($disableTime -lt $threshold) {
        # Delete the user
        $deleteMessage = "Deleting user $username as it has been disabled for more than 1 hour"
        Write-Host $deleteMessage
        "[" + $timestamp + "] " + $deleteMessage | Out-File -FilePath $logFilePath -Append

        # Remove the disable time file
        Remove-Item -Path $userFile.FullName -Force
        Remove-Item -Path $homeFolder -Force
        Remove-Item -Path $profilePath -Force
    }
}

