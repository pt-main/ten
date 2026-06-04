package ten

import "slices"

// Language represents the template language instance holding configuration,
// template sources, flags, and logical scope.
// Fields are unexported; use NewLanguage to create an instance.
type Language struct {
	flags          []string
	templateSource map[string]string
	Config         map[string]string
	logic_scope    map[string]bool
}

// NewLanguage creates a new Language instance with the provided flags,
// template source map, and configuration map.
// If flags is nil, an empty slice is used. If templateSource is nil,
// an empty map is used. If config is nil, default configuration is applied.
// Default configuration keys: "pre_plus", "post_plus", "raw", "string_start",
// "string_end", "ALPHA", "comment", "documentation".
// Missing keys in config are filled with defaults.
// The templateSource map is guaranteed to have a "nil" key with empty string.
func NewLanguage(
	flags []string,
	templateSource map[string]string,
	config map[string]string,
) (*Language, error) {
	default_config := map[string]string{
		"pre_plus": "pre+", "post_plus": "post+", "raw": "raw:",
		"string_start": "/ss(", "string_end": ")/se", "ALPHA": "",
		"comment": "//", "documentation": "has no docs",
	}
	if flags == nil {
		flags = make([]string, 0)
	}
	if templateSource == nil {
		templateSource = make(map[string]string)
	}
	if config == nil {
		config = default_config
	}
	default_keys := make([]string, 0, len(default_config))
	for k := range default_config {
		default_keys = append(default_keys, k)
	}
	config_keys := make([]string, 0, len(config))
	for k := range config {
		config_keys = append(config_keys, k)
	}
	for _, d_key := range default_keys {
		if !slices.Contains(config_keys, d_key) {
			config[d_key] = default_config[d_key]
		}
	}
	templateSource["nil"] = ""
	return &Language{
		flags:          flags,
		templateSource: templateSource,
		Config:         config,
		logic_scope:    map[string]bool{"false": false, "true": true},
	}, nil
}
