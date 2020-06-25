package service

import (
	"database/sql"
	"errors"
	"github.com/Alex27Khalupka/Go-course-task/pkg/model"
	"log"
)

// func GetGroups gets all groups with tasks assigned to them from database
func GetGroups(db *sql.DB) []model.Groups {
	rows, err := db.Query("SELECT groups.id, groups.title FROM groups")
	if err != nil {
		log.Fatal(err)
	}

	var groups []model.Groups

	for rows.Next() {
		var group model.Groups
		if err := rows.Scan(&group.ID, &group.Title); err != nil {
			log.Fatal(err)
		}
		for _, task := range getGroupTasks(db, group.ID) {
			for _, timeFrame := range getTaskTimeFrames(db, task.ID) {
				task.TimeFrames = append(task.TimeFrames, timeFrame)
			}
			group.Tasks = append(group.Tasks, task)
		}
		groups = append(groups, group)
	}

	return groups
}

// func GetGroups gets all tasks with time frames assigned to them from database
func GetTasks(db *sql.DB) []model.Tasks {

	rows, err := db.Query("SELECT tasks.id, tasks.title, tasks.group_id FROM tasks")
	if err != nil {
		log.Fatal(err)
	}

	var tasks []model.Tasks
	for rows.Next() {
		var task model.Tasks
		if err := rows.Scan(&task.ID, &task.Title, &task.GroupId); err != nil {
			log.Fatal(err)
		}
		for _, item := range getTaskTimeFrames(db, task.ID) {
			task.TimeFrames = append(task.TimeFrames, item)
		}
		tasks = append(tasks, task)
	}

	return tasks
}

// func PostGroups inserts group from request to database
// returns created group
func PostGroups(db *sql.DB, group model.Groups) (model.Groups, error) {
	if err := db.QueryRow("INSERT INTO groups (title) VALUES ($1) RETURNING id", group.Title).Scan(&group.ID); err != nil {
		return group, errors.New("group with this title already exists")
	}

	return group, nil
}

// func PostTAsks inserts task from request to database
// returns created task
// return error if groups id doesn't exist
func PostTasks(db *sql.DB, task model.Tasks) (model.Tasks, error) {
	if err := db.QueryRow("INSERT INTO tasks (title, group_id) VALUES ($1, $2) RETURNING id", task.Title, task.GroupId).
		Scan(&task.ID); err != nil {
		return task, errors.New("group with this id doesn't exist")
	}

	return task, nil
}

// func PostTimeFrames inserts time frame from request to database
// returns created time frame
// returns error if tasks id doesn't exist
func PostTimeFrames(db *sql.DB, timeFrame model.TimeFrames) (model.TimeFrames, error) {
	if err := db.QueryRow("INSERT INTO time_frames (task_id, from_time, to_time) VALUES ($1, $2, $3) RETURNING id",
		timeFrame.TaskId, timeFrame.From, timeFrame.To).
		Scan(&timeFrame.ID); err != nil {
		return timeFrame, errors.New("task with this id doesn't exist")
	}

	return timeFrame, nil
}

// func PutGroup updates group with given id
// returns updated group with tasks assigned to it
// return error if groups id doesn't exists
func PutGroups(db *sql.DB, id string, group model.Groups) (model.Groups, error) {

	var newGroup model.Groups

	_, err := db.Query("UPDATE groups SET title = $1 WHERE id = $2", group.Title, id)
	if err != nil {
		return newGroup, errors.New("wrong id type")
	}

	rows, err := db.Query("SELECT id, title FROM groups WHERE id = $1", id)
	if err != nil {
		return newGroup, err
	}

	for rows.Next() {
		if err := rows.Scan(&newGroup.ID, &newGroup.Title); err != nil {
			return newGroup, err
		}
	}

	if newGroup.ID == "" {
		return newGroup, errors.New("group with this id doesn't exist")
	}

	for _, task := range getGroupTasks(db, newGroup.ID) {
		for _, timeFrame := range getTaskTimeFrames(db, task.ID) {
			task.TimeFrames = append(task.TimeFrames, timeFrame)
		}
		newGroup.Tasks = append(newGroup.Tasks, task)
	}

	return newGroup, nil
}

// func PutTask updates group with given id
// returns updated task with time frames assigned to it
// return error if tasks or groups id doesn't exists
func PutTasks(db *sql.DB, id string, task model.Tasks) (model.Tasks, error) {

	var newTask model.Tasks

	_, err := db.Query("UPDATE tasks SET title = $1, group_id = $2 WHERE id = $3", task.Title, task.GroupId, id)
	if err != nil {
		return newTask, errors.New("group with this id doesn't exist")
	}

	rows, err := db.Query("SELECT id, title, group_id FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&newTask.ID, &newTask.Title, &newTask.GroupId); err != nil {
			log.Fatal(err)
		}
	}
	if newTask.ID == "" {
		return newTask, errors.New("task with this id doesn't exist")
	}

	for _, timeFrame := range getTaskTimeFrames(db, newTask.ID) {
		newTask.TimeFrames = append(newTask.TimeFrames, timeFrame)
	}

	return newTask, nil
}

// func DeleteGroups deletes group with given id from database
// with all assigned tasks to it
// returns error if id isn't uuid
func DeleteGroups(db *sql.DB, id string) error {
	rows, err := db.Query("SELECT id FROM tasks WHERE group_id = $1", id)
	if err != nil {
		return errors.New("id type isn't uuid")
	}

	var taskID string

	for rows.Next() {
		if err := rows.Scan(&taskID); err != nil {
			return err
		}
		if err = DeleteTasks(db, taskID); err != nil {
			return err
		}
	}

	_, err = db.Query("DELETE FROM groups WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// func DeleteTasks deletes task with given id from database
// with all assigned time frames to it
// returns error if id isn't uuid
func DeleteTasks(db *sql.DB, id string) error {
	_, err := db.Query("DELETE FROM time_frames WHERE task_id = $1", id)
	if err != nil {
		return errors.New("id type isn't uuid")
	}

	_, err = db.Query("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// func DeleteTimeFrames deletes time frames with given id from database
// with all assigned time frames to it
// returns error if id isn't uuid
func DeleteTimeFrames(db *sql.DB, id string) error {
	_, err := db.Query("DELETE FROM time_frames WHERE id = $1", id)
	if err != nil {
		return errors.New("id type isn't uuid")
	}
	return nil
}

// func getGroupTasks returns all tasks assigned to given group
func getGroupTasks(db *sql.DB, id string) []model.Tasks {
	rows, err := db.Query("SELECT tasks.id, tasks.title, tasks.group_id FROM tasks WHERE group_id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	var tasks []model.Tasks
	for rows.Next() {
		var task model.Tasks
		if err := rows.Scan(&task.ID, &task.Title, &task.GroupId); err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	return tasks

}

// func getTaskTimeFrames returns all time frames assigned to given task
func getTaskTimeFrames(db *sql.DB, id string) []model.TimeFrames {
	var timeFrames []model.TimeFrames

	rows, err := db.Query("SELECT time_frames.id, time_frames.from_time, time_frames.to_time, time_frames.task_id FROM time_frames WHERE time_frames.task_id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var timeFrame model.TimeFrames
		if err := rows.Scan(&timeFrame.ID, &timeFrame.From, &timeFrame.To, &timeFrame.TaskId); err != nil {
			log.Fatal(err)
		}
		timeFrames = append(timeFrames, timeFrame)
	}
	return timeFrames
}
