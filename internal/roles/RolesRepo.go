package roles

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

type Targets struct {
	Targets map[string]map[string]*Role `yaml:"targets"`
}

type Role struct {
	MaxExpSeconds int      `yaml:"max-exp"`
	Scopes        []string `yaml:"scopes"`
}

type TargetRole struct {
	Name          string
	RoleName      string
	MaxExpSeconds int
	Scopes        []string
}

type RolesRepo interface {
	GetTarget(name string, role string) (*TargetRole, error)
}

func (t *Targets) GetTarget(name string, role string) (*TargetRole, error) {
	target, ok := t.Targets[name]
	if ok {
		r, ok := target[role]
		if ok {
			return &TargetRole{
				Name:          name,
				RoleName:      role,
				MaxExpSeconds: r.MaxExpSeconds,
				Scopes:        r.Scopes,
			}, nil
		} else {
			return nil, errors.New(fmt.Sprintf("role %s doesn't exist", role))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("target %s doesn't exist", name))
	}
}

func GetRepo(file string) (RolesRepo, error) {
	v := &Targets{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(bytes))
	if err := yaml.Unmarshal(bytes, v); err != nil {
		return nil, err
	}
	fmt.Println(v)
	return v, nil
}
