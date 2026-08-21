package domain

import "context"

type Item struct {
	ID, InstanceID, Path string
	Attempts             int
	State                string
}
type Repository interface {
	Put(context.Context, Item) error
	Pending(context.Context, int) ([]Item, error)
	Mark(context.Context, string, string) error
}
type Files interface {
	Write(context.Context, string, []byte) (string, error)
	Read(context.Context, string) ([]byte, error)
	Remove(context.Context, string) error
}
