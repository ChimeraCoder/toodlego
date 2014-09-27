package toodledo

type TaskResponse struct {
	Meta  TaskResponseMeta
	Tasks []Task
}

type TaskResponseMeta struct {
	Num   int64 `json:"num"`
	Total int64 `json:"total"`
}

type Task struct {
	Completed int64  `json:"completed"`
	Folder    int64  `json:"folder"`
	ID        int64  `json:"id"`
	Modified  int64  `json:"modified"`
	Priority  int    `json:"priority"`
	Star      int    `json:"star"`
	Title     string `json:"title"`

	// Optional fields
	DueDate int64 `json:"duedate"`
	DueTime int64 `json:"duetime"`
}
