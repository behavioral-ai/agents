package module

const (
	Domain    = "github/behavioral-ai/agents"
	RouteName = "resiliency"
	Version   = "2.2.2"
	Ver1      = "v1"
	Ver2      = "v2"
)

// Configuration keys used on startup for map values
const (
	PackageNameUserKey     = "user"    // type:text, empty:false
	PackageNamePasswordKey = "pswd"    // type:text, empty:false
	PackageNameRetriesKey  = "retries" // type:int, default:-1, range:0-100
)
