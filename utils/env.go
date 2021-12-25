package utils

import "os"

func GetEnvVar(name string, defaultValue string) string {
	value, present := os.LookupEnv(name)
	if !present {
		value = defaultValue
	}
	return value
}

func HasEnvVar(name string) bool {
	_, found := os.LookupEnv(name)
	return found
}
