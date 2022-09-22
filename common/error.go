package common

type Error struct {
	Status  int
	Message string
	Body    interface{}
}
