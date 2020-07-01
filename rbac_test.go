package rbac

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNew(t *testing.T) {
	//SetSPLIT(".")
	//SetSUB(".")

	pg := New()

	pg.Users("matt", "super administrator")
	pg.Users("matt", "administrator")

	role := NewRole("super administrator")
	m := NewModule("platform", "system")

	c1 := NewController("environment", []*Action{
		{Name: "index"},
		{Name: "delete"},
		{Name: "new"},
		{Name: "read"},
	})

	c2 := NewController("member", []*Action{
		{Name: "index"},
	})

	m.SetControllers(c1, c2)
	role.SetModules(m)
	pg.SetRoles(role)

	spew.Dump(pg)
	spew.Dump(pg.Can("matt", "platform_system", "environment", "read"))
	spew.Dump(pg.Can("matt", []string{"platform", "system"}, "environment", "read"))
}
