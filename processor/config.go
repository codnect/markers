package processor

type Config struct {
	Version    string      `json:"version"`
	Parameters []Parameter `json:"parameters"`
	Overrides  []Override  `json:"overrides"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type Override struct {
	Package    string      `json:"package"`
	Version    string      `json:"version"`
	Parameters []Parameter `json:"parameters"`
}
