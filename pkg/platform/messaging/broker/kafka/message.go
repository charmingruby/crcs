package kafka

type Message struct {
	ID    string
	Topic string
	Key   []byte
	Value []byte
}
