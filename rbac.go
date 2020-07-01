package rbac

import (
	"errors"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
)

var SPLIT = "_"
var SUB = "_"

func SetSPLIT(v string) {
	SPLIT = v
}

func SetSUB(v string) {
	SUB = v
}

type RBAC struct {
	groups     []*Role
	users      map[string]string
	ruleCached []string
	l          sync.Mutex
}

type Role struct {
	Name    string
	Modules []*Module
}

func New() *RBAC {
	return &RBAC{
		users:      make(map[string]string),
		ruleCached: make([]string, 0),
	}
}

type Module struct {
	ID          uint
	Level       int
	Name        []string
	Controllers []*Controller
}

type Controller struct {
	ID      uint
	Level   int
	Name    string
	Actions []*Action
}

type Action struct {
	ID    uint
	Level int
	Name  string
}

func NewRole(name string) *Role {
	return &Role{
		Name: strcase.ToKebab(name),
	}
}

func (r *Role) SetModules(m ...*Module) {
	r.Modules = append(r.Modules, m...)
}

func NewModule(v ...string) *Module {
	return &Module{
		Name: v,
	}
}

func NewController(v string, actions []*Action) *Controller {
	return &Controller{
		Name:    v,
		Actions: actions,
	}
}

func (m *Module) SetControllers(c ...*Controller) {
	m.Controllers = append(m.Controllers, c...)
}

func (rbac *RBAC) SetRoles(g ...*Role) {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	rbac.groups = append(rbac.groups, g...)

	for _, v := range g {
		for _, m := range v.Modules {
			for _, c := range m.Controllers {
				for _, a := range c.Actions {
					var ruleCached = strings.Join([]string{v.Name, strings.Join(m.Name, SUB), c.Name, a.Name}, SPLIT)
					rbac.ruleCached = append(rbac.ruleCached, ruleCached)
				}
			}
		}
	}
}

func (rbac *RBAC) Users(username string, groupName string) {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	rbac.users[username] = groupName
}

func (rbac *RBAC) Can(username string, module interface{}, controller string, action string) bool {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	var groupName string
	var isFoundGroup bool
	var moduleName string

	if groupName, isFoundGroup = rbac.users[username]; !isFoundGroup {
		return false
	}

	switch m := module.(type) {
	case []string:
		moduleName = strings.Join(m, SUB)
	case string:
		moduleName = m
	default:
		panic(errors.New("module type error"))
	}

	var ruleCached = strings.Join([]string{groupName, moduleName, controller, action}, SPLIT)

	for _, c := range rbac.ruleCached {
		if c == ruleCached {
			return true
		}
	}

	return false
}
