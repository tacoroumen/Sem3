# Import the Active Directory module
Import-Module ActiveDirectory

# Read CSV file
$data = Import-Csv -Path "C:\Users\Administrator\Desktop\users.csv"

# Iterate over each row
foreach ($row in $data) {
    # Extracting information from CSV
    $username = $row.Username

    # Disable the user
    Write-Host "Disabling user: $username"
    Set-ADUser -Identity $username -Enabled $false

    # Add timestamp when user was disabled
    $disableTimePath = "C:\Users\Administrator\Desktop\DisableTime\$username.txt"
    New-Item -ItemType File -Path $disableTimePath -Force
    $timestamp = Get-Date -Format "M/d/yyyy h:mm:ss tt"
    $timestamp | Out-File -FilePath $disableTimePath

    # Remove the user from specified groups
    Remove-ADGroupMember -Identity "taco users" -Members $username -Confirm:$false
    Remove-ADGroupMember -Identity $row.GroupName -Members $username -Confirm:$false

    # Invoke ansible-playbook to update configurations
    Invoke-Expression -Command "ssh troumen@10.0.0.3 'ansible-playbook /home/troumen/ansible/delete_Desktop.yml --vault-password-file /home/troumen/ansible/vault_file.txt -e vm_name=$username'"
}
