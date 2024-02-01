# Go Hasher

A simple library to handle password hashing and checking.

## Installation

```bash
go get github.com/dniprodev/learn-switch-to-go-homework
```

## Usage

In your Go code, import the `hasher` package:

```golang
import "github.com/dniprodev/learn-switch-to-go-homework/hasher"
```

You can then use the function `HashPassword` to hash a password:

```golang
hashedPassword, err := hasher.HashPassword("myPassword")
if err != nil {
    // handle error
}
```

To check if a password matches a hash, use `CheckPasswordHash`:

```golang
match := hasher.CheckPasswordHash("myPassword", hashedPassword)
if match {
    fmt.Println("Password is correct")
} else {
    fmt.Println("Password is incorrect")
}
```

Please note:
1. Replace `"myPassword"` and `hashedPassword` with the actual password and hashed password you want to work with respectively.
2. You need to handle the error returned by `HashPassword` based on your program needs.
3. The `CheckPasswordHash` function returns a boolean to indicate if the password and the hash match.
