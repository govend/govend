package filters

// ExceptionList is a list of exception packages to skip.
var ExceptionList = []string{
	"appengine",
	"appengine/*",
	"appengine_internal",
	"appengine_internal/*",
}
