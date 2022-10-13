package pubsub

type PubSub interface {
	Connect(connectionString string) error
}
