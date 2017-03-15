package client

import (
	"errors"
	"math/rand"

	"github.com/Leon2012/goconfd/libs/ringhash"
)

type Selector interface {
	Select(values []string) (string, error)
}

type DefaultSelector struct {
}

func (s *DefaultSelector) Select(values []string) (string, error) {
	max := len(values)
	idx := rand.Intn(max)
	return values[idx], nil
}

type RingSelector struct {
	Key string
}

func NewRingSelector(k string) *RingSelector {
	return &RingSelector{
		Key: k,
	}
}

func (s *RingSelector) Select(values []string) (string, error) {
	if s.Key == "" {
		return "", errors.New("no default key")
	}
	max := len(values)
	ring := ringhash.New(max, nil)
	ring.Add(values...)
	str := ring.Get(s.Key)
	return str, nil
}
