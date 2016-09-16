package domain

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

type MockConnection struct {}

type MockDbRow struct {}

func (m *MockConnection) QueryRow(query string, args ...interface{}) DbRow {
  fmt.Println("MockConnection/QueryRow")
  return &MockDbRow{}
}

func (m *MockDbRow) Scan(dest ...interface{}) error {
  fmt.Println("MockDbRow/Scan")
  if b, ok := dest[0].(*string); ok {
    *b = "123"
  }
  if b, ok := dest[1].(*string); ok {
    *b = "boo foo"
  }
  if b, ok := dest[2].(*string); ok {
    *b = "+61432606721"
  }  
  return nil
}

func TestGetUserById(t *testing.T) {
  user, err := GetUserById(&MockConnection{}, "123")
  assert.Nil(t, err)
  assert.Equal(t, &UserEntity{Id: "123", Name: "boo foo", Phone: "+61432606721"}, user)
}