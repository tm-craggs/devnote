package config

type TemplateConfig struct {
	Project     string        `yaml:"project"`
	Date        DateConfig    `yaml:"date"`
	Time        TimeConfig    `yaml:"time"`
	Timestamp   TimeConfig    `yaml:"timestamp"`
	Commits     CommitsConfig `yaml:"commits"`
	Branch      string        `yaml:"branch"`
	AuthorName  string        `yaml:"authorName"`
	AuthorEmail string        `yaml:"authorEmail"`
}

type DateConfig struct {
	Format   string `yaml:"format"`
	Timezone string `yaml:"timezone"`
}

type TimeConfig struct {
	Format   string `yaml:"format"`
	Timezone string `yaml:"timezone"`
}

type CommitsConfig struct {
	Author          string  `yaml:"author"`
	Since           string  `yaml:"since"`
	Until           string  `yaml:"until"`
	Path            *string `yaml:"path"` // pointer for null values
	N               *int    `yaml:"n"`    // pointer for null values
	IncludeMerges   bool    `yaml:"include_merges"`
	Branch          string  `yaml:"branch"`
	MessageContains *string `yaml:"message_contains"` // pointer for null values
}
