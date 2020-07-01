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
	roles      []*Role
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

	rbac.roles = append(rbac.roles, g...)

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

func (rbac *RBAC) Users(username string, roleName string) {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	rbac.users[username] = roleName
}

func (rbac *RBAC) UserRole(username string) (bool, string) {
	var roleName string
	var isFoundRole bool

	if roleName, isFoundRole = rbac.users[username]; !isFoundRole {
		return false, ""
	}

	return true, roleName
}

func (rbac *RBAC) RoleUsers() map[string][]string {
	var roles = make(map[string][]string)

	for username, r := range rbac.users {
		if _, ok := roles[r]; !ok {
			roles[r] = make([]string, 0)
		}

		roles[r] = append(roles[r], username)
	}

	return roles
}

func (rbac *RBAC) getModule(username string, module interface{}) (bool, string, string) {
	var roleName string
	var isFoundRole bool
	var moduleName string

	if isFoundRole, roleName = rbac.UserRole(username); !isFoundRole {
		return false, "", ""
	}

	switch m := module.(type) {
	case []string:
		moduleName = strings.Join(m, SUB)
	case string:
		moduleName = m
	default:
		panic(errors.New("module type error"))
	}

	return true, roleName, moduleName
}

func (rbac *RBAC) Can(username string, module interface{}, controller string, action string) bool {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	var roleName string
	var isFoundRole bool
	var moduleName string

	if isFoundRole, roleName, moduleName = rbac.getModule(username, module); !isFoundRole {
		return false
	}

	roleName = strcase.ToKebab(roleName)

	var ruleCached = strings.Join([]string{roleName, moduleName, controller, action}, SPLIT)

	ruleCached = strings.TrimRight(ruleCached, SPLIT)

	for _, c := range rbac.ruleCached {
		if strings.HasPrefix(c, ruleCached) {
			return true
		}
	}

	return false
}

func (rbac *RBAC) CanModule(username string, module interface{}) bool {
	return rbac.Can(username, module, "", "")
}

func (rbac *RBAC) CanController(username string, module interface{}, controller string) bool {
	return rbac.Can(username, module, controller, "")
}
