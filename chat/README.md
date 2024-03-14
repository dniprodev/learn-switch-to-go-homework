# switch-to-go-homework module which is lets-go-chat indeed 😅

lets-go-chat is a small Go module with multiple packages.

## Go Hasher

A simple library to handle password hashing and checking.

### Installation

```bash
go get github.com/dniprodev/learn-switch-to-go-homework/chat
```

### Usage

In your Go code, import the `hasher` package:

```golang
import "github.com/dniprodev/learn-switch-to-go-homework/chat/hasher"
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

Sure, here's a condensed version:

**Local Deployment Guide:**

1. Install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/).

1. Clone the repo: `git clone <URL of your repository>`

1. Move to the project directory: `cd path/to/project`

1. Update the `.env` file with the appropriate values for your PostgreSQL and MongoDB setup.

1. Build and Start the app: 
   
   ```bash
   docker-compose up --build
   ```

1. To validate that the application is running, navigate to `http://localhost:${HTTP_PORT}` in your browser (replace `${HTTP_PORT}` with your chosen port from `.env`).

1. To shut down the app and delete its resources, run:

   ```bash
   docker-compose down
   ```
Make sure to replace all placeholders (like `<URL of your repository>` and `path/to/project`) with actual values.
