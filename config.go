package main

type Account struct {
	AccountId string `yaml:"AccountId"`
	RoleName  string `yaml:"RoleName"`
}

type Group struct {
	ManagedPolicies []string `yaml:"ManagedPolicies"`
}

type Policy struct {
	Statements []Statement `yaml:"Statements"`
}

type Role struct {
	Policies []string `yaml:"Policies"`
}

type Statement struct {
	Actions   []string `yaml:"Actions"`
	Effect    string   `yaml:"Effect"`
	Resources []string `yaml:"Resources"`
}

type User struct {
	Groups []string `yaml:"Groups"`
}

type Variables struct {
	BackendOrganization string `yaml:"BackendOrganization"`
	MainAccountId       string `yaml:"MainAccountId"`
	MainProfile         string `yaml:"MainProfile"`
	Region              string `yaml:"Region"`
}

type OrganisationRolesConfig struct {
	Accounts  map[string]Account `yaml:"Accounts"`
	Groups    map[string]Group   `yaml:"Groups"`
	Policies  map[string]Policy  `yaml:"Policies"`
	Roles     map[string]Role    `yaml:"Roles"`
	Users     map[string]User    `yaml:"Users"`
	Variables Variables          `yaml:"Variables"`
}
