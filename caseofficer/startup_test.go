package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/core/host"
)

func ExampleStartupPing() {
	status := host.Ping(PkgPath)
	fmt.Printf("test: Ping() -> [status:%v]\n", status)

	//Output:
	//test: Ping() -> [status:OK]

}
