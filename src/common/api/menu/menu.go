package menu

import (
	"github.com/gin-gonic/gin"
)

type Menu struct {
	Uri      string    `json:"uri"`
	Name     string    `json:"name"`
	Perm     *itemPerm `json:"-"`
	Buttons  []*Button `json:"buttons"`
	SubMenus []*Menu   `json:"sub_menus"`
}

func NewRootMenu(perm *itemPerm, subMenu []*Menu) *Menu {
	newM := &Menu{
		Perm: perm,
	}
	newM.AddSubMenu(subMenu...)

	return newM
}

// AddButton 追加按钮
func (m *Menu) AddButton(bts ...*Button) {
	if m.Buttons == nil {
		m.Buttons = make([]*Button, len(bts))
	}
	for idx, bt := range bts {
		m.Buttons[idx] = bt
	}
}

// AddSubMenu 追加子菜单
func (m *Menu) AddSubMenu(subMenu ...*Menu) {
	if m.SubMenus == nil {
		m.SubMenus = make([]*Menu, len(subMenu))
	}
	for idx, sMenu := range subMenu {
		m.SubMenus[idx] = sMenu
	}
}

// ListMenusByContext 根据gin.Context中的当前用户携带的权限信息列出菜单
func (m *Menu) ListMenusByContext(c *gin.Context) []*Menu {
	return recursionMenu(m, c)
}

func recursionMenu(m *Menu, c *gin.Context) []*Menu {
	mArr := make([]*Menu, 0)
	if m == nil {
		return mArr
	}

	for _, sub := range m.SubMenus {
		if sub.Perm != nil {
			if ok, _ := sub.Perm.hasPerm(c); !ok {
				continue
			}
		}

		newM := &Menu{
			Uri:      sub.Uri,
			Name:     sub.Name,
			Buttons:  make([]*Button, 0),
			SubMenus: recursionMenu(sub, c),
		}
		for _, b := range sub.Buttons {
			if b.Perm != nil {
				if ok, _ := b.Perm.hasPerm(c); !ok {
					continue
				}
			}
			newM.Buttons = append(newM.Buttons, &Button{b.Id, b.Name, nil})
		}

		mArr = append(mArr, newM)
	}

	return mArr
}
