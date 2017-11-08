package doc

import yaml "gopkg.in/yaml.v2"

type DocumentVersionFragment struct {
	Version string `yaml:"version" json:"version"`
}

type DocumentDomain struct {
	Key   string `yaml:"key" json:"key"`
	Title string `yaml:"title" json:"title"`
}

type DocumentErrorDescriptionFragment struct {
	Friendly  string `yaml:"friendly" json:"friendly"`
	Technical string `yaml:"technical" json:"technical"`
}

type DocumentErrorDescription DocumentErrorDescriptionFragment

func (dd *DocumentErrorDescription) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var f DocumentErrorDescriptionFragment
	if err := unmarshal(&f); err != nil {
		var s string
		if err := unmarshal(&s); err != nil {
			return err
		}

		*dd = DocumentErrorDescription{
			Friendly:  s,
			Technical: s,
		}
	} else {
		*dd = DocumentErrorDescription(f)
	}

	return nil
}

type DocumentErrorArgumentFragment struct {
	Type        string      `yaml:"type" json:"type"`
	Description string      `yaml:"description" json:"description"`
	Validators  []string    `yaml:"validators" json:"validators"`
	Default     interface{} `yaml:"default" json:"default"`
}

type DocumentErrorArgument DocumentErrorArgumentFragment

func (dea *DocumentErrorArgument) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var f DocumentErrorArgumentFragment
	if err := unmarshal(&f); err != nil {
		return err
	}

	*dea = DocumentErrorArgument(f)

	if dea.Validators == nil {
		dea.Validators = []string{}
	}

	if dea.Type == "" {
		dea.Type = "string"
	}

	return nil
}

func (dea *DocumentErrorArgument) IsOptional() bool {
	return dea != nil && dea.Default != nil
}

type DocumentErrorArgumentItem struct {
	Name     string
	Argument *DocumentErrorArgument
}

type DocumentErrorArguments map[string]*DocumentErrorArgument

type DocumentErrorFragment struct {
	Title            string                      `json:"title"`
	Description      *DocumentErrorDescription   `json:"description"`
	Arguments        DocumentErrorArguments      `json:"arguments"`
	OrderedArguments []DocumentErrorArgumentItem `json:"-"`
}

type DocumentError DocumentErrorFragment

func (de *DocumentError) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var f DocumentErrorFragment
	if err := unmarshal(&f); err != nil {
		return err
	}

	*de = DocumentError(f)

	if de.Arguments == nil {
		de.Arguments = DocumentErrorArguments{}
	}

	// Extract the iteration order.
	var fi struct {
		Title       string                    `yaml:"title"`
		Description *DocumentErrorDescription `yaml:"description"`
		ArgumentsIt yaml.MapSlice             `yaml:"arguments"`
	}
	if err := unmarshal(&fi); err != nil {
		return err
	}

	for _, item := range fi.ArgumentsIt {
		name := item.Key.(string)

		de.OrderedArguments = append(de.OrderedArguments, DocumentErrorArgumentItem{
			Name:     name,
			Argument: de.Arguments[name],
		})
	}

	return nil
}

type DocumentErrors map[string]DocumentError

type DocumentSection struct {
	Title  string         `yaml:"title" json:"title"`
	Errors DocumentErrors `yaml:"errors" json:"errors"`
}

type DocumentSections map[string]DocumentSection

type Document struct {
	Version  int              `yaml:"version" json:"version"`
	Domain   DocumentDomain   `yaml:"domain" json:"domain"`
	Sections DocumentSections `yaml:"sections" json:"sections"`
}
