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
	m := NewModule("platform", "system")

	c2 := NewController("member", []*Action{
		{Name: "index"},
	})

	m.SetControllers(c2)
	role.SetModules(m)
	pg.SetRoles(role)

	spew.Dump(pg)
	spew.Dump(pg.Can("matt", []string{"platform", "system"}, "member", "index"))
	spew.Dump(pg.CanModule("matt", []string{"platform", "system"}))
	spew.Dump(pg.CanController("matt", []string{"platform", "system"}, "member"))
}
