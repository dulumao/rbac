package rbac

type IRule interface {
	GetParentID() uint
	GetID() uint
	GetLevel() int
	GetName() string
}

func GetModuleTrees(rules []IRule) []*Module {
	var modules []*Module
	var parents []IRule

	for _, r := range rules {
		if r.GetParentID() == 0 {
			parents = append(parents, r)
		}
	}

	// 检测 module 的下级
	for _, p := range parents {
		var children int

		for _, r := range rules {
			if r.GetParentID() == p.GetID() && p.GetLevel() == r.GetLevel() {
				modules = append(modules, &Module{
					ID:    r.GetID(),
					Level: r.GetLevel(),
					Name:  []string{p.GetName(), r.GetName()},
				})

				children++
			}
		}

		if children == 0 {
			modules = append(modules, &Module{
				ID:    p.GetID(),
				Level: p.GetLevel(),
				Name:  []string{p.GetName()},
			})
		}
	}

	// level 2
	for _, m := range modules {
		for _, r := range rules {
			if r.GetParentID() == m.ID {
				var c = &Controller{
					ID:    r.GetID(),
					Level: r.GetLevel(),
					Name:  r.GetName(),
				}

				m.Controllers = append(m.Controllers, c)
			}
		}
	}

	for _, m := range modules {
		// level 3
		for _, c := range m.Controllers {
			for _, r := range rules {
				if r.GetParentID() == c.ID /*&& r.GetLevel() == 3*/ {
					c.Actions = append(c.Actions, &Action{
						ID:    r.GetID(),
						Name:  r.GetName(),
						Level: r.GetLevel(),
					})
				}
			}
		}
	}

	return modules
}
