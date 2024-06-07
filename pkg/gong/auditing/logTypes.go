package auditing

// Map log types to their corresponding structs
var logTypeStructMap = map[string]interface{}{
	"AccessLog":                  AccessLog{},
	"UserActivityLog":            UserActivityLog{},
	"UserCallPlay":               UserCallPlay{},
	"ExternallySharedCallAccess": ExternallySharedCallAccess{},
	"ExternallySharedCallPlay":   ExternallySharedCallPlay{},
}
