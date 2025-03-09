package env

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

type testAddr struct {
	Address string `env:"ADDRESS"`
	Port    int    `env:"PORT"`
}

type testConfig struct {
	Name  string `env:"NAME"`
	Age   int    `env:"AGE"`
	Email string `env:"EMAIL"`
	Addr  *testAddr
}

func Test_ReadEnv(t *testing.T) {

	expect := testConfig{
		Name:  "John Doe",
		Age:   30,
		Email: "johndoe@example.com",
		Addr:  &testAddr{Address: "127.0.0.1", Port: 8080},
	}

	os.Setenv("AGE", strconv.Itoa(expect.Age))
	os.Setenv("EMAIL", expect.Email)
	os.Setenv("ADDRESS", expect.Addr.Address)

	init := testConfig{
		Name: "John Doe",
		Age:  29,
		Addr: &testAddr{Port: 8080},
	}

	ReadEnv(&init)

	if !reflect.DeepEqual(expect, init) {
		t.Errorf("updateEnv: expect %v, got %v", expect, init)
	}
}
