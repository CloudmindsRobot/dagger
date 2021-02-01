package models

type GroupYAML struct {
	Groups []RuleYAML `yaml:"groups"`
}

type RuleYAML struct {
	Name  int         `yaml:"name"`
	Rules []AlertYAML `yaml:"rules"`
}

type AlertYAML struct {
	Alert       interface{}            `yaml:"alert"`
	Expr        interface{}            `yaml:"expr"`
	For         int                    `yaml:"for"`
	Labels      map[string]interface{} `yaml:"labels"`
	Annotations map[string]interface{} `yaml:"annotations"`
}
