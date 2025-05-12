package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]reflect.Value
}

func NewContainer() *Container {
	return &Container{constructors: make(map[string]reflect.Value)}
}

// RegisterType - зарегистрировать конструктор по созданию типа
func (c *Container) RegisterType(name string, constructor interface{}) {
	val := reflect.ValueOf(constructor)

	if val.Kind() != reflect.Func {
		panic("constructor must be a func")
	}
	c.constructors[name] = val
}

// Resolve - создать объект с использованием конструктора
func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, ok := c.constructors[name]
	if !ok {
		return nil, fmt.Errorf("no constructor for %q", name)
	}

	if constructor.Type().NumIn() != 0 {
		return nil, errors.New("constructor must not accept arguments")
	}

	results := constructor.Call(nil)
	if len(results) != 1 {
		return nil, errors.New("constructor must return result")
	}

	return results[0].Interface(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
