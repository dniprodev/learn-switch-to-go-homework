package user

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCustomerRepository(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join(".", "testdata", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	customerRepo, err := NewRepository(ctx, connStr)
	assert.NoError(t, err)

	err = customerRepo.Save(ctx, User{
		ID:       "123",
		Name:     "Henry",
		Password: "123",
	})
	assert.NoError(t, err)

	customer, err := customerRepo.FindByUsername(ctx, "Henry")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "123", customer.ID)
	assert.Equal(t, "Henry", customer.Name)
	assert.Equal(t, "123", customer.Password)
}
