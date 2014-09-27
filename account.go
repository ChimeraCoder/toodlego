package toodledo

type Account struct {
	Alias            string `json:"alias"`
	DateFormat       int    `json:"dateformat"`
	Email            string `json:"email"`
	HideMonths       int    `json:"hidemonths"`
	HotListduedate   int    `json:"hotlistduedate"`
	HotListpriority  int    `json:"hotlistpriority"`
	HotListstar      int    `json:"hotliststar"`
	HotListstatus    int    `json:"hotliststatus"`
	LastDeleteNote   int    `json:"lastdelete_note"`
	LastDeleteTask   int    `json:"lastdelete_task"`
	LastEditContext  int    `json:"lastedit_context"`
	LastEditFolder   int    `json:"lastedit_folder"`
	LastEditGoal     int    `json:"lastedit_goal"`
	LastEditList     int    `json:"lastedit_list"`
	LastEditLocation int    `json:"lastedit_location"`
	LastEditNote     int    `json:"lastedit_note"`
	LastEditOutline  int    `json:"lastedit_outline"`
	LastEditTask     int    `json:"lastedit_task"`
	Pro              int    `json:"pro"`
	ShowTabNums      int    `json:"showtabnums"`
	Timezone         int    `json:"timezone"`
	Userid           string `json:"userid"`
}
