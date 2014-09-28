package toodledo

import "time"

type TaskResponse struct {
	Meta  TaskResponseMeta
	Tasks []Task
}

type TaskResponseMeta struct {
	Num   int64 `json:"num"`
	Total int64 `json:"total"`
}

type ToodleTime int64

type Task struct {
	Completed int64  `json:"completed"`
	Folder    int64  `json:"folder"`
	ID        int64  `json:"id"`
	Modified  int64  `json:"modified"`
	Priority  int    `json:"priority"`
	Star      int    `json:"star"`
	Title     string `json:"title"`

	// Optional fields

	// Date and time are ToodleTimes and should not be used directly
	DueDate   ToodleTime `json:"duedate"`
	DueTime   ToodleTime `json:"duetime"`
	StartDate ToodleTime `json:"startdate"`
	StartTime ToodleTime `json:"starttime"`

	Length     int64  `json:"length"`
	TagsString string `json:"tags"`
	Parent     int64  `json:"parent"`
}

func (t Task) Due() *time.Time {
	return parseToodleDateTime(t.DueDate, t.DueTime)
}

func (t Task) Start() *time.Time {
	return parseToodleDateTime(t.StartDate, t.StartTime)
}

func parseToodleDateTime(tdate, ttime ToodleTime) *time.Time {
	if tdate == 0 && ttime == 0 {
		return nil
	}
	//The day and clock time (hour/min/sec) are correct when read in GMT and interpreted as local time
	year, month, day := time.Unix(int64(tdate), 0).UTC().Date()
	hour, min, sec := time.Unix(int64(ttime), 0).UTC().Clock()

	//For some reason it gives the year as 1974 sometimes
	if year < 2010 {
		year, _, _ = time.Now().Date()
	}

	due := time.Date(year, month, day, hour, min, sec, 0, time.Local)
	return &due
}
