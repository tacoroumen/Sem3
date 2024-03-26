package main

import (
	"html/template"
	"net/http"
	"os/exec"
)

// Define a template for the HTML form
var tpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Password Reset</title>
</head>
<body>
    <h1>Password Reset</h1>
    <form method="post" action="/reset">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required>
        <button type="submit">Reset Password</button>
    </form>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", showForm)
	http.HandleFunc("/reset", resetPassword)
	http.ListenAndServe(":8080", nil)
}

func showForm(w http.ResponseWriter, r *http.Request) {
	// Render the HTML form
	tpl.Execute(w, nil)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	// Get the username from the form submission
	username := r.FormValue("username")

	// Call the PowerShell script to get the email address
	cmd := exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/GetEmail.ps1", "-Username", username)
	email, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to retrieve email address", http.StatusInternalServerError)
		return
	}

	// Send an email to the retrieved email address (replace with your email sending logic)
	sendEmail(string(email))

	// Respond to the user indicating that the password reset email has been sent
	w.Write([]byte("Password reset email has been sent to the provided email address"))
}

// Dummy function to simulate sending an email
func sendEmail(email string) {
	// Replace this function with your actual email sending logic
	// For example, you can use a third-party library like sendgrid-go or gomail
	// This is just a placeholder
	println("Sending email to:", email)

}
