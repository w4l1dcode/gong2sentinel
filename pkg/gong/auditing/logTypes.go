package auditing

// LogTypeStructMap maps log types to their corresponding structs
var LogTypeStructMap = map[string]interface{}{
	"AccessLog":                  AccessLog{},
	"UserActivityLog":            UserActivityLog{},
	"UserCallPlay":               UserCallPlay{},
	"ExternallySharedCallAccess": ExternallySharedCallAccess{},
	"ExternallySharedCallPlay":   ExternallySharedCallPlay{},
}
