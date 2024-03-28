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

// Define a template for the HTML login form
var loginTpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
</head>
<body>
    <h1>Login</h1>
    <form method="post" action="/login">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required><br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required><br>
        <button type="submit">Login</button>
        <button type="button" onclick="forgotPassword()">Forgot Password</button>
    </form>

    <script>
        function forgotPassword() {
            // Redirect the user to the forgot password page or show a modal for password recovery
            // Example: window.location.href = "/forgot-password";
            alert("Redirecting to forgot password page...");
			window.location.href = "/request-reset";
        }
    </script>
</body>
</html>
`))

// Define a template for the HTML request application form
var requestAppTpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Request Application</title>
</head>
<body>
    <h1>Request Application</h1>
    <form method="post" action="/request">
        <label for="app-name">Application Name:</label>
        <input type="text" id="app-name" name="app-name" required><br>
        <button type="submit">Request Application</button>
    </form>
</body>
</html>
`))

var resetPasstpl = template.Must(template.New("").Parse(`
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

// Map to store session IDs
var sessions = make(map[string]struct{})

func main() {
	http.HandleFunc("/", showLogin)
	http.HandleFunc("/login", login)
	http.HandleFunc("/request-reset", showResetPassword)
	http.HandleFunc("/reset", resetPassword)
	http.HandleFunc("/request", requestApplicationHandler)
	http.ListenAndServe(":8080", nil)
}

func showLogin(w http.ResponseWriter, r *http.Request) {
	loginTpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	// Get the username and password from the form submission
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Call the PowerShell script to authenticate the user
	cmd := exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/Authenticate.ps1", "-Username", username, "-Password", password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
		return
	}

	// Check the output of the script
	if strings.TrimSpace(string(output)) == "Authentication successful" {
		// Generate a new session ID
		sessionID := generateRandomString(16)

		// Store the session ID
		sessions[sessionID] = struct{}{}

		// Set the session ID cookie in the response with security options
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   3600,                    // Example: Cookie expires in 1 hour
			HttpOnly: true,                    // Prevent JavaScript access to the cookie
			Secure:   true,                    // Cookie sent only over HTTPS
			SameSite: http.SameSiteStrictMode, // Limit cookie to same-site requests
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    username,
			Path:     "/",
			MaxAge:   3600,                    // Example: Cookie expires in 1 hour
			HttpOnly: true,                    // Prevent JavaScript access to the cookie
			Secure:   true,                    // Cookie sent only over HTTPS
			SameSite: http.SameSiteStrictMode, // Limit cookie to same-site requests
		})

		// If authentication succeeds, redirect to the home page or any other authorized page
		http.Redirect(w, r, "/request", http.StatusSeeOther)
		return
	}

	// If authentication fails, render the login form again with an error message
	loginTpl.Execute(w, "Invalid username or password")
}

func requestApplicationHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Verify session ID
	if _, ok := sessions[cookie.Value]; ok {
		// Check if the request is a POST method
		if r.Method == http.MethodPost {
			// Parse the form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
				return
			}

			// Get the application name from the form submission
			appName := r.Form.Get("app-name")

			// If the application name is provided, send an email
			if appName != "" {
				// Get the username from the cookie
				usernameCookie, err := r.Cookie("username")
				if err != nil {
					http.Error(w, "Failed to retrieve username from session", http.StatusInternalServerError)
					return
				}
				username := usernameCookie.Value

				// Send an email with the username and application name
				sendEmail(username, appName)

				// Respond to the user indicating that the application request was successful
				w.Write([]byte("Request successful"))
				return
			} else {
				// If the application name is not provided, render the request application form with an error message
				requestAppTpl.Execute(w, "Application name is required")
				return
			}
		}

		// Render the request application form
		requestAppTpl.Execute(w, nil)
		return
	}

	// If session ID is not valid or user is not logged in, redirect to login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func sendEmail(username string, appName string) {
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

	to := []string{"app.request@0x54.dev"}
	subject := "Subject: " + "Application request!" + "\r\n"
	body := "User " + username + " has requested the application: " + appName + "\r\n"
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

func showResetPassword(w http.ResponseWriter, r *http.Request) {
	resetPasstpl.Execute(w, nil)
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
	if email == "Please provide a username using the -Username flag." {
		// If the email address is not found, respond to the user with a 404 Not Found status
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate a random password of 8 characters
	password := generateRandomString(8)

	// Call the PowerShell script to reset the password
	cmd = exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/ResetPassword.ps1", "-Username", username, "-Password", password)

	// Send an email to the retrieved email address (replace with your email sending logic)
	sendEmailReset(email, password)

	// Respond to the user indicating that the password reset email has been sent
	w.Write([]byte("If user excists password reset email has been sent to your personal email address"))
}

func sendEmailReset(email string, password string) {
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
	body := "Your password has been reset. your new password is: " + password + "\r\n" + "You will need to change your password after login."
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

func generateResetKey() string {
	return generateRandomString(21) // You can adjust the length as needed
}
