package service

import (
	"fmt"
	"github.com/advanced-go/common/core"
)

func ExampleNewAgent() {
	a := NewAgent(core.Origin{Region: "us-central"}, nil)

	fmt.Printf("test: NewAgent() -> [uri:%v]\n", a)

	//Output:
	//test: NewAgent() -> [uri:service:us-central..]

}
