package userplan

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// UserPlan contains quota and feature of an user plan
type UserPlan struct {
	Name                 string `json:"name"`
	TotalArticles        uint   `yaml:"total_articles" json:"total_articles"`
	TotalCategories      uint   `yaml:"total_categories" json:"total_categories"`
	TotalRules           uint   `yaml:"total_rules" json:"total_rules"`
	TotalAPIKeys         uint   `yaml:"total_api_keys" json:"total_api_keys"`
	TotalArchiveServices uint   `yaml:"total_archive_services" json:"total_archive_services"`
}

// UserPlans contains all user plans by name
type UserPlans struct {
	Plans map[string]UserPlan `yaml:"user_plans"`
}

// GetPlan return a plan by its name and fallback to default plan if missing
// Returns nil if default plan is also missing
func (p UserPlans) GetPlan(name string) (result *UserPlan) {
	if plan, ok := p.Plans[name]; ok {
		plan.Name = name
		result = &plan
	} else if plan, ok := p.Plans["default"]; ok {
		plan.Name = "default"
		result = &plan
	}
	return
}

// GetPlans return all plans
func (p UserPlans) GetPlans() []UserPlan {
	plans := []UserPlan{}
	for name, plan := range p.Plans {
		plan.Name = name
		plans = append(plans, plan)
	}
	return plans
}

// NewUserPlans create new user plans from definition file
func NewUserPlans(filename string) (UserPlans, error) {
	plans := UserPlans{
		Plans: make(map[string]UserPlan),
	}
	if filename == "" {
		return plans, nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("unable to load user plans from file (%s): %s", filename, err.Error())
		return plans, err
	}
	err = yaml.Unmarshal(data, &plans)
	if err != nil {
		err = fmt.Errorf("unable to read user plans definition: %s", err.Error())
		return plans, err
	}
	return plans, nil
}
