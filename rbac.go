package rbac

import (
	"errors"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
)

var SPLIT = "_"
var SUB = "_"
var GOD = "__$$GOD$$__"

func SetSPLIT(v string) {
	SPLIT = v
}

func SetSUB(v string) {
	SUB = v
}

type RBAC struct {
	roles      []*Role
	users      map[string][]string
	ruleCached []string
	l          sync.Mutex
	// 上帝用户
	god string
}

type Role struct {
	Name    string
	Modules []*Module
}

func New() *RBAC {
	return &RBAC{
		users:      make(map[string][]string),
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

func (rbac *RBAC) SetGod(username string) {
	rbac.god = username
}

func (rbac *RBAC) IsGod(username string) bool {
	if rbac.god == username {
		return true
	}

	return false
}

func (rbac *RBAC) Users(username string, roles ...string) {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	if _, ok := rbac.users[username]; ok {
		rbac.users[username] = append(rbac.users[username], roles...)
	} else {
		rbac.users[username] = roles
	}
}

func (rbac *RBAC) UserRole(username string) (bool, []string) {
	var roles []string
	var isFoundRole bool

	if rbac.IsGod(username) {
		return true, []string{GOD}
	}

	if roles, isFoundRole = rbac.users[username]; !isFoundRole {
		return false, roles
	}

	return true, roles
}

func (rbac *RBAC) RoleUsers() map[string][]string {
	var roles = make(map[string][]string)

	for username, userRoles := range rbac.users {
		for _, r := range userRoles {
			if _, ok := roles[r]; !ok {
				roles[r] = make([]string, 0)
			}

			roles[r] = append(roles[r], username)
		}
	}

	return roles
}

func (rbac *RBAC) getModule(username string, module interface{}) (bool, []string, string) {
	var roles []string
	var isFoundRole bool
	var moduleName string

	if isFoundRole, roles = rbac.UserRole(username); !isFoundRole {
		return false, roles, ""
	}

	switch m := module.(type) {
	case []string:
		moduleName = strings.Join(m, SUB)
	case string:
		moduleName = m
	default:
		panic(errors.New("module type error"))
	}

	return true, roles, moduleName
}

func (rbac *RBAC) Can(username string, module interface{}, controller string, action string) bool {
	rbac.l.Lock()
	defer rbac.l.Unlock()

	var roles []string
	var isFoundRole bool
	var moduleName string

	if isFoundRole, roles, moduleName = rbac.getModule(username, module); !isFoundRole {
		return false
	}

	if roles[0] == GOD {
		return true
	}

	for _, r := range roles {
		var data = strings.Join([]string{strcase.ToKebab(r), moduleName, controller, action}, SPLIT)

		data = strings.TrimRight(data, SPLIT)

		for _, c := range rbac.ruleCached {
			if strings.HasPrefix(c, data) {
				return true
			}
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
