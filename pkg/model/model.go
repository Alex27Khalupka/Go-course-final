package model

import "time"

type Groups struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Tasks []Tasks `json:"tasks"`
}

type Tasks struct {
	ID         string       `json:"id"`
	Title      string       `json:"title"`
	GroupId    string       `json:"group_id"`
	TimeFrames []TimeFrames `json:"time_frames"`
}

type TimeFrames struct {
	ID     string    `json:"id"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	TaskId string    `json:"task_id"`
}

type ResponseGroups struct {
	Groups []Groups `json:"groups"`
}

type ResponseTasks struct {
	Tasks []Tasks `json:"tasks"`
}
