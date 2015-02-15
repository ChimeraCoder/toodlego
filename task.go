package toodlego

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
	// Use Due() and Start() instead
	DueDate   ToodleTime `json:"duedate"`
	DueTime   ToodleTime `json:"duetime"`
	StartDate ToodleTime `json:"startdate"`
	StartTime ToodleTime `json:"starttime"`

	Length     int64  `json:"length"`
	TagsString string `json:"tag"`
	Parent     int64  `json:"parent"`
}

func (t Task) Due() *time.Time {
	return parseToodleDateTime(t.DueDate, t.DueTime)
}

func (t Task) Start() *time.Time {
	return parseToodleDateTime(t.StartDate, t.StartTime)
}

// ImplicitStart finds the start time for the task implied by the combination of
// the due time and the length of the task
// If either the length or the due time is not specified, ImplicitStart defaults to
// the explicit start time specified
func (t Task) ImplicitStart() *time.Time {
	// Default to using the time implied by the due time and length
	if t.Length != 0 && t.Due() != nil {
		result := t.Due().Add(-1 * time.Duration(t.Length) * time.Minute)
		return &result
	}
	return t.Start()
}

// IsChild is a convenience function that returns true if the task has a parent task
func (t Task) IsChild() bool {
	return t.Parent != 0
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
