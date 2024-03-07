# Import the Active Directory module
Import-Module ActiveDirectory

# Function to check and delete disabled users older than 1 hour
function CheckAndDeleteDisabledUsers {
    # Define the time threshold for deletion (1 hour)
    $threshold = (Get-Date).AddHours(-1)
    #$threshold = (Get-Date).AddMinutes(-2)

    # Check for users disabled for more than 1 hour and delete them
    $disabledUsers = Get-ChildItem -Path "C:\Users\Administrator\Desktop\DisableTime\" -Filter *.txt
    foreach ($userFile in $disabledUsers) {
        $username = $userFile.BaseName
        $disableTimeString = Get-Content $userFile.FullName
        $disableTime = [DateTime]::ParseExact($disableTimeString, "M/d/yyyy h:mm:ss tt", $null)

        if ($disableTime -lt $threshold) {
            # Delete the user
            Write-Host "Deleting user $username as it has been disabled for more than 1 hour"
            Remove-ADUser -Identity $username -Confirm:$false

            # Remove the disable time file
            Remove-Item -Path $userFile.FullName -Force
        }
    }
}

# Infinite loop to continuously check and delete disabled users every 30 minutes
while ($true) {
    # Call the function to check and delete disabled users
    CheckAndDeleteDisabledUsers

    # Pause execution for 30 minutes
    Start-Sleep -Seconds (30 * 60)  # 30 minutes = 30 * 60 seconds
}
