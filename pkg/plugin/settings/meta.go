package settings

const (
	name        = "kind"
	description = "A UI for kind"
	icon        = "flask"
)

// GetName returns the plugin name
func GetName() string {
	return name
}

// GetDescription returns the plugin description
func GetDescription() string {
	return description
}
