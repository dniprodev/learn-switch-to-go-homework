package message

import "sync"

type Repository struct {
	messages []Message
	lock     sync.Mutex
}

func (r *Repository) save(message Message) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.messages = append(r.messages, message)
}

func (r *Repository) findAll() []Message {
	r.lock.Lock()
	defer r.lock.Unlock()

	// TODO: Q: is it ok to return it that way?
	result := make([]Message, len(r.messages))
	copy(result, r.messages)
	return result
}
