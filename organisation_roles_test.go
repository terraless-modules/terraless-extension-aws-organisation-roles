package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

type TestLogger struct {
	logs map[string][]string
}
func (testLogger *TestLogger) Debug(msg string, args ...interface{}) {
	testLogger.logs["debug"] = append(testLogger.logs["debug"], msg)
}
func (testLogger *TestLogger) Error(msg string, args ...interface{}) {
	testLogger.logs["error"] = append(testLogger.logs["error"], msg)
}
func (testLogger *TestLogger) Info(msg string, args ...interface{}) {
	testLogger.logs["info"] = append(testLogger.logs["info"], msg)
}
func (testLogger *TestLogger) Warn(msg string, args ...interface{}) {
	testLogger.logs["warn"] = append(testLogger.logs["warn"], msg)
}
func (testLogger *TestLogger) Trace(msg string, args ...interface{}) {
	testLogger.logs["trace"] = append(testLogger.logs["warn"], msg)
}
func (testLogger *TestLogger) IsDebug() bool {
	return true
}
func (testLogger *TestLogger) IsError() bool {
	return true
}
func (testLogger *TestLogger) IsInfo() bool {
	return true
}
func (testLogger *TestLogger) IsTrace() bool {
	return true
}
func (testLogger *TestLogger) IsWarn() bool {
	return true
}
func (testLogger *TestLogger) Named(name string) hclog.Logger {
	return &TestLogger{}
}
func (testLogger *TestLogger) ResetNamed(name string) hclog.Logger {
	return &TestLogger{}
}
func (testLogger *TestLogger) SetLevel(level hclog.Level) {}
func (testLogger *TestLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	buffer := bytes.Buffer{}
	return log.New(&buffer, "", 0)
}
func (testLogger *TestLogger) With(args ...interface{}) hclog.Logger {
	return &TestLogger{}
}
func (testLogger *TestLogger) Messages() map[string][]string {
	return testLogger.logs
}

var testLogger TestLogger
func TestMain(m *testing.M) {
	testLogger = TestLogger{
		logs: map[string][]string{},
	}
	logger = &testLogger

	os.Exit(m.Run())
}

func TestOrganisationRoles_Exec(t *testing.T) {
	testLogger = TestLogger{
		logs: map[string][]string{},
	}
	logger = &testLogger

	// given
	globalConfig := schema.TerralessGlobalConfig{
		Teams: []schema.TerralessTeam{
			{
				Name: "Odania",
				Data: map[string]string{
					"baseAccountId": "my-account-id",
				},
				Providers: []schema.TerralessProvider{
					{
						Type: "dummy",
						Name: "dummy-provider",
						Data: map[string]string{
							"accountId": "account-id-1",
							"color": "color1",
						},
					},
					{
						Type: "aws",
						Name: "aws-provider",
						Data: map[string]string{
							"accountId": "account-id-2",
							"color": "color2",
						},
						Roles: []string{
							"admin",
							"developer",
						},
					},
				},
			},
		},
	}
	terralessData := schema.TerralessData{}

	// when
	extension := ExtensionAwsOrganisationRoles{}
	err := extension.Exec(globalConfig, terralessData)

	// then
	assert.Equal(t, nil, err)

	assert.NotNil(t, retrieveConfigurationString(testLogger.Messages()["info"], "AWS Organisation Roles finished writting terraform files"))
	assert.NotNil(t, retrieveConfigurationString(testLogger.Messages()["debug"], "AWS Organisation Roles configuration"))
}

func retrieveConfigurationString(logs []string, prefix string) *string {
	for _, line := range logs {
		if strings.HasPrefix(line, prefix) {
			return &line
		}
	}

	logrus.Fatalf("Failed to detect 'AWS Organisation Roles configuration'")
	return nil
}
