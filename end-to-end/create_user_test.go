package user_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func cleanUserDB(t *testing.T) {
	err := godotenv.Load(".e2e.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Connection string
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := pgx.Connect(context.Background(), psqlconn)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "TRUNCATE TABLE users")
	if err != nil {
		t.Fatalf("Failed to clean database: %v", err)
	}
}

// Test cases for POST /user
func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		payload       string
		expStatusCode int
	}{
		{
			name: "Valid User",
			payload: `{
				"userName": "TestUser",
				"password": "testpassword"
			}`,
			expStatusCode: http.StatusOK,
		},
		{
			name: "Missing Username",
			payload: `{
				"password": "testpassword"
			}`,
			expStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cleanUserDB(t) // Clean the database before each test

			req, _ := http.NewRequest("POST", "http://localhost:8080/user", strings.NewReader(tc.payload))
			resp, err := http.DefaultClient.Do(req)

			require.NoError(t, err)                             // Assert that there was no error making request
			require.Equal(t, tc.expStatusCode, resp.StatusCode) // Assert the status code is what we expect
		})
	}
}
