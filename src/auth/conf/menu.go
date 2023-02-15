package conf

import "eago/common/api/menu"

func NewMenu(conf *Conf) *menu.Menu {
	return menu.NewRootMenu(nil,
		[]*menu.Menu{
			{
				Uri:  "/auth/products",
				Name: "产品线-菜单",
				Buttons: []*menu.Button{
					{
						"POST_auth/products",
						"产品线-新增-按钮",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/auth/departments",
				Name: "部门-菜单",
				Buttons: []*menu.Button{
					{
						"GET_auth/departments/tree",
						"产品线-查看树结构-按钮",
						nil,
					},
					{
						"GET_auth/departments",
						"产品线-列表方式查看-按钮",
						nil,
					},
				},
			},
			{
				Uri:  "/auth/groups",
				Name: "组-菜单",
				Buttons: []*menu.Button{
					{
						"POST_auth/groups",
						"组-新增-按钮",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/auth/roles",
				Name: "角色-菜单",
				Buttons: []*menu.Button{
					{
						"POST_auth/roles",
						"角色-新增-按钮",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/auth/users",
				Name: "用户-菜单",
			},
		},
	)
}
