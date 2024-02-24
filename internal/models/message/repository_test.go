package message

import "testing"

func TestRepository(t *testing.T) {
    want := 1
	var sut = Repository{}
	sut.save(Message{ text: "123"})
    if got := len(sut.findAll()); got != want {
        t.Errorf("findAll() = %v, want %v", got, want)
    }
}
