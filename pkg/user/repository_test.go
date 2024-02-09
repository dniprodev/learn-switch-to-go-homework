package user

import "testing"

func TestRepository(t *testing.T) {
    want := User{id: "id", name: "name", password: "password"}
	var sut = Repository{}
	sut.save(User{id: "id", name: "name", password: "password"})
    if got, ok := sut.findByUsername("name"); got.id != want.id || !ok {
        t.Errorf("findAll() = %v, want %v", got, want)
    }
}
