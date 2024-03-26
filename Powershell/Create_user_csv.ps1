# Import the Active Directory module
Import-Module ActiveDirectory

# Read CSV file
$data = Import-Csv -Path "S:\HR\newusers.csv"

# Iterate over each row
foreach ($row in $data) {
    # Extracting information from CSV
    $firstName = $row.FirstName
    $lastName = $row.LastName
    $username = $row.Username
    $groupName = $row.GroupName

    # Define password (ensure it meets complexity requirements)
    $password = ConvertTo-SecureString -String "Welcome1!" -AsPlainText -Force

    # Constructing the full name
    $fullName = $firstName + " " + $lastName

    # Setting home folder and profile path
    $homeFolder = "\\WIN-ULE5N582EEV\Homefolder\$username"
    $profilePath = "\\WIN-ULE5N582EEV\User-profiles\$username"

    # Determine OU based on group name
    $ou = switch ($groupName) {
        "IT" { "OU=IT,DC=tacoroumen,DC=local" }
        "Finance" { "OU=Finance,DC=tacoroumen,DC=local" }
        "HR" { "OU=HR,DC=tacoroumen,DC=local" }
    }

    # Creating the user in specified OU
    Write-Host "Creating user: $fullName with username: $username"
    New-ADUser `
        -Name $fullname `
        -GivenName $firstName `
        -Surname $lastName `
        -SamAccountName $username `
        -UserPrincipalName ($username + "@tacoroumen.local") `
        -AccountPassword $password `
        -Enabled $true `
        -HomeDirectory $homeFolder `
        -ProfilePath $profilePath `
        -HomeDrive "H:" `
        -ChangePasswordAtLogon:$true `
        -Path $ou


    New-Item $homeFolder -Type Directory
    New-Item $profilePath -Type Directory

    $homeACL = Get-Acl $homeFolder
    $homeRuleUser = New-Object System.Security.AccessControl.FileSystemAccessRule("$username","FullControl","ContainerInherit,ObjectInherit","None","Allow")
    $homeRuleAdmin = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrator","FullControl","ContainerInherit,ObjectInherit","None","Allow")
    $homeACL.AddAccessRule($homeRuleUser)
    $homeACL.AddAccessRule($homeRuleAdmin)
    Set-Acl -Path $homeFolder -AclObject $homeACL

    $profileACL = Get-Acl $profilePath
    $profileRuleUser = New-Object System.Security.AccessControl.FileSystemAccessRule("$username","FullControl","ContainerInherit,ObjectInherit","None","Allow")
    $profileRuleAdmin = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrator","FullControl","ContainerInherit,ObjectInherit","None","Allow")
    $profileACL.AddAccessRule($profileRuleUser)
    $profileACL.AddAccessRule($profileRuleAdmin)
    Set-Acl -Path $profilePath -AclObject $profileACL

    # Adding the user to the specified group
    Add-ADGroupMember -Identity "taco users" -Members $username
    Add-ADGroupMember -Identity $groupName -Members $username
    Add-ADGroupMember -Identity "Roaming user profiles users and computers" -Members $username

    # Invoke ansible-playbook
    Invoke-Expression -Command "ssh troumen@10.0.0.3 'ansible-playbook /home/troumen/ansible/create_Desktop.yml --vault-password-file /home/troumen/ansible/vault_file.txt -e vm_name=$username'"
}
