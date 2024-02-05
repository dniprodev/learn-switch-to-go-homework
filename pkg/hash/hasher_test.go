package hasher

import "testing"

func TestHello(t *testing.T) {
    want := false
    if got := CheckPasswordHash("password", "hash"); got != want {
        t.Errorf("CheckPasswordHash() = %v, want %v", got, want)
    }
}
