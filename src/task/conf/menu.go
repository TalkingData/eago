package conf

import "eago/common/api/menu"

func NewMenu(conf *Conf) *menu.Menu {
	return menu.NewRootMenu(nil,
		[]*menu.Menu{
			{
				Uri:  "/task/results",
				Name: "结果-菜单",
			},
			{
				Uri:  "/task/tasks",
				Name: "任务-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_task/tasks",
						"任务-新增-按钮",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/task/schedules",
				Name: "计划任务-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
				Buttons: []*menu.Button{
					{
						"POST_task/schedules",
						"计划任务-新增-按钮",
						menu.NewPermIsRole(conf.Const.AdminRole),
					},
				},
			},
			{
				Uri:  "/task/works",
				Name: "执行器-菜单",
				Perm: menu.NewPermIsRole(conf.Const.AdminRole),
			},
		},
	)
}
