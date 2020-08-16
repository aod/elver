package config

var appName string

// SetAppName sets the global config app name which is used for retrieving
// specific configuration contents.
func SetAppName(name string) {
	appName = name
}
