package rbac

type IRule interface {
	GetParentID() uint
	GetID() uint
	GetLevel() int
	GetName() string
}

func GetModuleTrees(rules []IRule) []*Module {
	var modules []*Module
	var parents []*Module

	for _, r := range rules {
		if r.GetParentID() == 0 {
			parents = append(parents, &Module{
				ID:    r.GetID(),
				Level: r.GetLevel(),
				Name:  []string{r.GetName()},
			})
		}
	}

	// level 1
	for _, p := range parents {
		for _, r := range rules {
			if r.GetParentID() == p.ID && p.Level == r.GetLevel() {
				//m.ID = r.GetID()
				//m.Name = append(m.Name, r.GetName())
				modules = append(modules, &Module{
					ID:    r.GetID(),
					Level: r.GetLevel(),
					Name:  []string{p.Name[0], r.GetName()},
				})
			}
		}
	}

	for _, m := range modules {
		// level 2
		for _, r := range rules {
			if r.GetParentID() == m.ID {
				var c = &Controller{
					ID:   r.GetID(),
					Name: r.GetName(),
				}

				m.Controllers = append(m.Controllers, c)
			}
		}
	}

	for _, m := range modules {
		// level 3
		for _, c := range m.Controllers {
			for _, r := range rules {
				if r.GetParentID() == c.ID {
					c.Actions = append(c.Actions, &Action{
						Name: r.GetName(),
					})
				}
			}
		}
	}

	return modules
}

type Tree struct {
	Name     string
	ID       uint
	level    int
	Children []*Tree
}

func GetRuleTrees(rules []IRule, parentID uint) []*Tree {
	var t []*Tree

	for _, r := range rules {
		if r.GetParentID() == parentID {
			t = append(t, &Tree{
				Name:     r.GetName(),
				ID:       r.GetID(),
				Children: GetRuleTrees(rules, r.GetID()),
				level:    r.GetLevel(),
			})
		}
	}

	return t
}
