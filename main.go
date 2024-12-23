// go mod init
// go get github.com/golang-jwt/jwt/v5
// go get github.com/urfave/cli/v2
package main
import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/urfave/cli/v2"
)
var secretKey = []byte("your-secret-key") // Replace with a secure key
var users = []string{}                    // Stores users in-memory
// GenerateJWT generates a JWT token for a specific user
func GenerateJWT(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expiration: 1 hour
        "iat":      time.Now().Unix(),                    // Issued at time
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}
// AddUser adds a new user to the user list
func AddUser(username string) {
    users = append(users, username)
    fmt.Printf("User '%s' added successfully.\n", username)
}
// ListUsers displays the list of users
func ListUsers() {
    if len(users) == 0 {
        fmt.Println("No users available. Add a new user first.")
        return
    }
    fmt.Println("Existing users:")
    for i, user := range users {
        fmt.Printf("[%d] %s\n", i+1, user)
    }
}
// Prompt for input from the user
func promptInput(prompt string) string {
    fmt.Print(prompt)
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    return strings.TrimSpace(scanner.Text())
}
func main() {
    app := &cli.App{
        Name:  "jwtgen",
        Usage: "A multi-step CLI to manage users and generate JWT tokens",
        Action: func(c *cli.Context) error {
            for {
                fmt.Println("Select an option:")
                fmt.Println("[1] Create new user")
                fmt.Println("[2] Select existing user")
                choice := promptInput("Enter your choice: ")
                switch choice {
                case "1":
                    // Create a new user
                    username := promptInput("Enter username: ")
                    AddUser(username)
                case "2":
                    // Select an existing user
                    if len(users) == 0 {
                        fmt.Println("No users available. Please add a user first.")
                        continue
                    }
                    ListUsers()
                    userChoice := promptInput("Select a user by number to get JWT token: ")
                    selectedIndex := -1
                    fmt.Sscanf(userChoice, "%d", &selectedIndex)
                    if selectedIndex < 1 || selectedIndex > len(users) {
                        fmt.Println("Invalid selection. Please try again.")
                        continue
                    }
                    selectedUser := users[selectedIndex-1]
                    fmt.Printf("You selected user: %s\n", selectedUser)
                    // Generate token
                    fmt.Println("Generating JWT token...")
                    token, err := GenerateJWT(selectedUser)
                    if err != nil {
                        return err
                    }
                    fmt.Println("JWT Token generated successfully:")
                    fmt.Println(token)
                default:
                    fmt.Println("Invalid choice. Please try again.")
                }
            }
        },
    }
    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
