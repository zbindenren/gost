package configuration

import (
	"reflect"
	"testing"
)

func TestMarshalJsonSubnet(t *testing.T) {
	c := &Configuration{
		Private:  true,
		Token:    "123",
		Username: "Ant Man",
	}

	err := c.Save()
	if err != nil {
		t.Fatalf("could not save configuration: %s", err)
	}
	c2, err := LoadConfiguration()
	if err != nil {
		t.Fatalf("could not load configuration: %s", err)
	}
	if !reflect.DeepEqual(c, c2) {
		t.Errorf("configurations are different but should not: %v --- %v", c, c2)
	}
}
