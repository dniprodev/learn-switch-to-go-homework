package hasher

import "testing"

func TestHello(t *testing.T) {
    want := true
    if got := CheckPasswordHash("password", "hash"); got != want {
        t.Errorf("CheckPasswordHash() = %q, want %q", got, want)
    }
}
