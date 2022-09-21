package reflect

import (
	"reflect"
	"testing"
)

type User struct {
	Name        string `json:"name"`
	Country     string `json:"my_country"`
	Age         string
	CapitalCity string
}

func TestGetJSONFieldNames(t *testing.T) {
	got := GetJSONFieldNames(User{})
	want := []string{"name", "my_country", "age", "capital_city"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
