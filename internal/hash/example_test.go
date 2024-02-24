package hasher

import "fmt"

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
