package message

type Repository struct {
	messages []Message
}

func (r *Repository)save(message Message) {
	r.messages = append(r.messages, message)
}

func (r *Repository)findAll() []Message {
	return r.messages
}