package services

import (
	"errors"
)

type Database interface {
	InsertMessage(json []byte) error
	PopMessage() error
	GetAllMessages() ([][]byte, error)
}

type Node struct {
	value []byte
	next  *Node
}

type Stack struct {
	top           *Node
	currentLength int
	DB            Database
}

func (s *Stack) Push(value []byte) error {
	newNode := &Node{value: value, next: s.top}
	s.top = newNode
	err := s.DB.InsertMessage(value)
	if err != nil {
		return err
	}
	s.currentLength++
	return nil
}

func (s *Stack) Pop() ([]byte, error) {
	if s.top == nil {
		return nil, errors.New("stack is empty")
	}
	value := s.top.value
	s.top = s.top.next
	err := s.DB.PopMessage()
	if err != nil {
		return nil, err
	}
	s.currentLength--
	return value, nil
}

func (s *Stack) LoadStackIfExist() error {
	messages, err := s.DB.GetAllMessages()
	if err != nil {
		return err
	}

	for _, node := range messages {
		newNode := &Node{value: node, next: s.top}
		s.top = newNode
		s.currentLength++
	}

	return nil
}
