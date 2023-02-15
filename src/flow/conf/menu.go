package conf

import (
	"eago/common/api/menu"
)

func NewMenu(conf *Conf) *menu.Menu {
	newInsBtn := &menu.Button{
		"POST_flow/instantiate",
		"发起流程-按钮",
		nil,
	}

	return menu.NewRootMenu(nil,
		[]*menu.Menu{
			{
				Uri:     "/flow/instances/todo",
				Name:    "我的流程-代办-菜单",
				Buttons: []*menu.Button{newInsBtn},
			},
			{
				Uri:     "/flow/instances/done",
				Name:    "我的流程-已办-菜单",
				Buttons: []*menu.Button{newInsBtn},
			},
			{
				Uri:     "/flow/instances/my",
				Name:    "我的流程-我发起的-菜单",
				Buttons: []*menu.Button{newInsBtn},
			},
			{
				Uri:  "/flow/instances",
				Name: "流程实例-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
			},
			{
				Uri:  "/flow/categories",
				Name: "类别-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_flow/categories",
						"类别-新增-按钮",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/flow/triggers",
				Name: "触发器-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_flow/triggers",
						"类别-新增-触发器",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/flow/nodes",
				Name: "节点-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_flow/nodes",
						"节点-新增-触发器",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/flow/forms",
				Name: "表单-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_flow/forms",
						"表单-新增-触发器",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/flow/flows",
				Name: "流程-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_flow/flows",
						"流程-新增-触发器",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
		},
	)
}
