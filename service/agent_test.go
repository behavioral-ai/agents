package service

import (
	"fmt"
	"github.com/behavioral-ai/core/core"
)

func ExampleNewAgent() {
	a := NewAgent(core.Origin{Region: "us-central"}, nil, nil)

	fmt.Printf("test: NewAgent() -> [uri:%v]\n", a)

	//Output:
	//test: NewAgent() -> [uri:service:us-central..]

}
