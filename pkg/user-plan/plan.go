package userplan

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// UserPlan contains quota and feature of an user plan
type UserPlan struct {
	Name            string `yaml:"name" json:"name"`
	TotalArticles   uint   `yaml:"total_articles" json:"total_articles"`
	TotalCategories uint   `yaml:"total_categories" json:"total_categories"`
}

// UserPlans contains all user plans by name
type UserPlans struct {
	Plans []UserPlan `yaml:"user_plans"`
}

// GetPlan return a plan by its name and fallback to first plan if missing
func (p UserPlans) GetPlan(name string) (result *UserPlan) {
	if len(p.Plans) == 0 {
		return nil
	}
	for _, plan := range p.Plans {
		if plan.Name == name {
			return &plan
		}
	}
	// Fallback to first plan
	plan := p.Plans[0]
	return &plan
}

// NewUserPlans create new user plans from definition file
func NewUserPlans(filename string) (UserPlans, error) {
	plans := UserPlans{
		Plans: []UserPlan{},
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
