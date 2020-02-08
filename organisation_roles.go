package main

import (
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (extension *ExtensionAwsOrganisationRoles) Exec(globalConfig schema.TerralessGlobalConfig, data schema.TerralessData) error {
	logger.Info(fmt.Sprintf("[%s] Executing", ExtensionName))

	config := readConfig()

	logger.Debug(fmt.Sprintf("AWS Organisation Roles configuration: %v\n", config))

	// Render Main Account Terraform Files
	buffer := generateMainAccountUsers(config)
	groups := generateMainAccountGroups(config)
	buffer.Write(groups.Bytes())
	support.WriteToFile("main-account.tf", buffer)

	// Render Sub Account Terraform Files
	buffer = generateAccountRoles(config)
	support.WriteToFile("account/main.tf", buffer)

	// Render Module Call Terraform Files
	buffer = generateModuleCall(config)
	support.WriteToFile("account-modules.tf", buffer)

	// Render Module Call Terraform Files
	buffer = generateMainTemplate(config)
	support.WriteToFile("main.tf", buffer)

	// Generate Deploy File
	buffer = generateDeployScript(config)
	support.WriteToFile("deploy.sh", buffer)

	logger.Info(fmt.Sprintf("AWS Organisation Roles finished writting terraform files\n"))

	return nil
}

func readConfig() *OrganisationRolesConfig {
	yamlBytes, err := ioutil.ReadFile("organisation-roles.yml")
	config := &OrganisationRolesConfig{}

	if err != nil {
		logrus.Warnf("Could not parse organisation roles config! Error: %s\n", err)
		return nil
	}

	if err := yaml.Unmarshal(yamlBytes, config); err != nil {
		logrus.Fatalf("Could not parse project config! Error: %s\n", err)
		return nil
	}

	logger.Info(fmt.Sprintf("AWS Organisation Roles: %v\n", config))

	return config
}
