package context

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	getAllEnvVariables    = os.Environ
	getAllDotEnvVariables = godotenv.Read
)

type Environment map[string]string

func prepareEnvironment() map[string]string {
	env := make(Environment)
	env.loadDotEnvFile()
	env.loadEnvironment()

	return env
}

func (env *Environment) loadEnvironment() {
	for _, envVar := range getAllEnvVariables() {
		values := strings.SplitN(envVar, "=", 2)
		key := values[0]
		value := values[1]

		(*env)[key] = value
	}
}

func (env *Environment) loadDotEnvFile() {
	myEnv, err := getAllDotEnvVariables()
	if err != nil {
		return
	}

	for key, value := range myEnv {
		(*env)[key] = value
	}
}
