package context

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareEnvironment_NoEnvSet_EnvIsEmpty(t *testing.T) {
	setEnvVarsFun(map[string]string{})
	setDotEnvVarsFun(map[string]string{})

	env := prepareEnvironment()

	assert.Empty(t, env)
}

func TestPrepareEnvironment_OnlyEnvIsSet_EnvContainsItems(t *testing.T) {
	setEnvVarsFun(map[string]string{
		"myEnv":  "myVar",
		"myEnv2": "myVar2",
	})
	setDotEnvVarsFun(map[string]string{})

	env := prepareEnvironment()

	assert.Equal(t, "myVar", env["myEnv"])
	assert.Equal(t, "myVar2", env["myEnv2"])
}

func TestPrepareEnvironment_OnlyDotEnvIsSet_EnvContainsItems(t *testing.T) {
	setEnvVarsFun(map[string]string{})
	setDotEnvVarsFun(map[string]string{
		"myEnv":  "myVar",
		"myEnv2": "myVar2"})

	env := prepareEnvironment()

	assert.Equal(t, "myVar", env["myEnv"])
	assert.Equal(t, "myVar2", env["myEnv2"])
}

func TestPrepareEnvironment_BothSourcesAreSet_EnvContainsMergedItems(t *testing.T) {
	setEnvVarsFun(map[string]string{
		"myEnv": "myVar",
	})
	setDotEnvVarsFun(map[string]string{
		"myEnv2": "myVar2",
	})

	env := prepareEnvironment()

	assert.Equal(t, "myVar", env["myEnv"])
	assert.Equal(t, "myVar2", env["myEnv2"])
}

func TestPrepareEnvironment_BothSourcesContainsSameKey_EnvTakesPrecedence(t *testing.T) {
	setEnvVarsFun(map[string]string{
		"myEnv": "fromEnv",
	})
	setDotEnvVarsFun(map[string]string{
		"myEnv": "fromDotEnv",
	})

	env := prepareEnvironment()

	assert.Equal(t, "fromEnv", env["myEnv"])
}

func setEnvVarsFun(mapping map[string]string) func() {
	originalGetEnvVarsFunc := getAllEnvVariables

	var env []string
	for key, value := range mapping {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	getAllEnvVariables = func() []string { return env }

	return func() {
		getAllEnvVariables = originalGetEnvVarsFunc
	}
}

func setDotEnvVarsFun(mapping map[string]string) func() {
	originalGetDotEnvVarsFunc := getAllDotEnvVariables

	getAllDotEnvVariables = func(...string) (map[string]string, error) { return mapping, nil }

	return func() {
		getAllDotEnvVariables = originalGetDotEnvVarsFunc
	}
}
