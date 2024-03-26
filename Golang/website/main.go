package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"
)

type EmailConfig struct {
	Mailadress  string `json:"mailadress"`
	Wachtwoord  string `json:"wachtwoord"`
	Smtp_server string `json:"smtp_server"`
	Smtp_poort  int    `json:"smtp_poort"`
}

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

// Function to generate a random string of given length
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time

	// Characters allowed in the random string
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*-_"

	// Create a byte slice to store the random string
	randomString := make([]byte, length)

	// Generate random characters and append them to the byte slice
	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return
	return string(randomString)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	// Get the username from the form submission
	username := r.FormValue("username")

	// Call the PowerShell script to get the email address
	cmd := exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/GetEmail.ps1", "-Username", username)
	emailBytes, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to retrieve email address", http.StatusInternalServerError)
		return
	}

	// Trim any leading or trailing whitespace characters, including carriage returns and line feeds
	email := strings.TrimSpace(string(emailBytes))

	// Generate a random password of 8 characters
	password := generateRandomString(8)

	// Call the PowerShell script to reset the password
	cmd = exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/ResetPassword.ps1", "-Username", username, "-Password", password)

	// Send an email to the retrieved email address (replace with your email sending logic)
	sendEmail(email, password)

	// Respond to the user indicating that the password reset email has been sent
	w.Write([]byte("Password reset email has been sent to your personal email address"))
}

// Dummy function to simulate sending an email
func sendEmail(email string, password string) {
	// Replace this function with your actual email sending logic
	// For example, you can use a third-party library like sendgrid-go or gomail
	// This is just a placeholder
	println("Sending email to:", email)

	// Open and read the email configuration from a JSON file
	file, err := os.Open("mail.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := EmailConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	auth := smtp.PlainAuth("", config.Mailadress, config.Wachtwoord, config.Smtp_server)

	to := []string{email}
	subject := "Subject: " + "Password reset!" + "\r\n"
	body := "Your password has been reset. your new password is: " + password + "\r\n" + "Please change your password after login."
	msg := []byte(subject +
		"\r\n" +
		body)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.Smtp_server, config.Smtp_poort), auth, config.Mailadress, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email sent successfully!")
}
