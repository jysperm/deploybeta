package datasource

import "testing"

func TestMakeRedisImage(t *testing.T) {
	if err := MakeRedisImage(); err != nil {
		t.Fatal(err)
	}
}

func TestMakeMongoImage(t *testing.T) {
	if err := MakeMognoImage(); err != nil {
		t.Fatal(err)
	}
}
