// Package hasher provides functionality for calculating and checking hash for string
package hasher

import "fmt"

// CheckPasswordHash checks does password has provided hash
func CheckPasswordHash(password, hash string) bool {
	return false
}

// HashPassword calculates hash for given password
func HashPassword(password string) (string, error) {
	return "", nil
}

// ExampleUsage shows how to use CheckPasswordHash and HashPassword functions.
func ExampleUsage() {
	hash, err := HashPassword("password")
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

	result := CheckPasswordHash("password", hash)

	fmt.Printf("Result: %v", result)
}
