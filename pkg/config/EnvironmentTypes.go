package config

type EnvironmentTypes struct {
	Release string
	Debug   string
}

var (
	EnvTypes = EnvironmentTypes{
		Release: "release",
		Debug:   "local",
	}
)
