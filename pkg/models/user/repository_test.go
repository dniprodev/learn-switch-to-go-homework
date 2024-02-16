package user

import (
	"testing"
)

func TestRepository(t *testing.T) {
	want := User{ID: "id", Name: "name", Password: "password"}
	var sut = Repository{}
	sut.Save(User{ID: "id", Name: "name", Password: "password"})
	if got, ok := sut.FindByUsername("name"); got.Name != want.Name || got.Password != want.Password || !ok {
		t.Errorf("findAll() = %v, want %v", got, want)
	}
}
