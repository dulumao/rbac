package rbac

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNew(t *testing.T) {
	pg := New()

	pg.Users("matt", "super administrator", "administrator")
	pg.Users("tom", "administrator")

	superAdministratorRole := NewRole("super administrator")
	administratorRole := NewRole("administrator")

	c1 := NewController("environment", []*Action{
		{Name: "index"},
		{Name: "delete"},
		{Name: "new"},
		{Name: "read"},
	})

	c2 := NewController("member", []*Action{
		{Name: "index"},
	})

	m := NewModule("platform", "system")
	m.SetControllers(c1, c2)

	superAdministratorRole.SetModules(m)

	pg.SetRoles(superAdministratorRole, administratorRole)

	println("----------")
	spew.Dump(pg)
	println("----------")
	spew.Dump(pg.UserRole("matt"))
	spew.Dump(pg.RoleUsers())
	spew.Dump(pg.Can("matt", "platform_system", "environment", "read"))
	spew.Dump(pg.Can("matt", []string{"platform", "system"}, "environment", "read"))
}

func TestError(t *testing.T) {
	pg := New()

	pg.Users("matt", "administrator")

	role := NewRole("administrator")

	dashboard := NewModule("dashboard")
	system := NewModule("platform", "system")
	member := NewModule("platform", "member")

	c0 := NewController("dashboard", []*Action{
		{Name: "index"},
	})

	c1 := NewController("environment", []*Action{
		{Name: "index"},
	})

	c2 := NewController("role", []*Action{
		{Name: "index"},
	})

	dashboard.SetControllers(c0)
	system.SetControllers(c1)
	member.SetControllers(c2)
	role.SetModules(dashboard, system, member)
	pg.SetRoles(role)

	spew.Dump(pg)
	spew.Dump(pg.Can("matt", []string{"platform", "system"}, "environment", "index"))
	spew.Dump(pg.CanModule("matt", []string{"platform", "system"}))
	spew.Dump(pg.Can("matt", []string{"platform", "member"}, "role", "index"))
}
