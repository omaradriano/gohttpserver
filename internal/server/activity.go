package server

import (
	"fmt"
	// "sync"
	"time"
)

type Activity struct {
	AddedDate		time.Time 	`json:"addedDate"`
	Description 	string		`json:"description"`
	ID 				uint64		`json:"id"`
	Title			string		`json:"title"`
	LimitDate		time.Time	`json:"limitDate"`
	IsCompleted		bool		`json:"isCompleted"`
}

var ErrIDNotFound = fmt.Errorf("ID not found")

type Activities struct { //Este struct guarda una lista de Activity
	activities []Activity
}

/*
	* Agrega una actividad al "array" de actividades
	* devuelve el ID de la actividad agregada (como copia)
	* Toda la demas informacion se agrega en el request de la solicitud al usar la ruta
*/
func (c *Activities) InsertActivity(act Activity) uint64 { 
	act.ID = uint64(len(c.activities)) + 1 
	act.AddedDate = time.Now()
	limitDate := act.AddedDate
	act.LimitDate = limitDate.AddDate(0,0,5)
	act.IsCompleted = false
	c.activities = append(c.activities, act) 
	return act.ID 
}

func (c *Activities) GetActivity(id uint64) (Activity, error) {
	if id > uint64(len(c.activities)) {
		return Activity{}, ErrIDNotFound
	}

	for _, activity := range c.activities {
		if activity.ID == id {
			return activity, nil
		}
	}
	return Activity{}, ErrIDNotFound
}

func (c *Activities) EditActivity(act Activity) (uint64, error) {
	for i, activity := range c.activities {
		if activity.ID == act.ID {
			if act.Title != "" {
				c.activities[i].Title = act.Title
			}
			if act.Description != "" {
				c.activities[i].Description = act.Description
			}
			if !act.LimitDate.IsZero() {
				c.activities[i].LimitDate = act.LimitDate
			}
			c.activities[i].IsCompleted = act.IsCompleted
			return c.activities[i].ID, nil
		}
	}
	return 0, ErrIDNotFound
}

func (c *Activities) DeleteActivity(id uint64) error {
	if id < 1 {
		return fmt.Errorf("Id can't be lower than 0")
	}
	if len(c.activities) == 0 {
		return fmt.Errorf("There's no elements in the list")
	}

	toDeleteIndex := 0
	for index, activity := range c.activities {
		if activity.ID == id {
			toDeleteIndex = index
			break
		}
	}

	if toDeleteIndex != 0 {
		// fmt.Printf("Activity about to be deleted: %s", c.activities[toDeleteIndex].Title)
		c.activities = append(c.activities[:toDeleteIndex], c.activities[toDeleteIndex+1:]...)
		return nil
	}

	return fmt.Errorf("Activity doesn't exist")
}

// func (c *Activities) SearchActivity(id uint64) (uint64, error) {
	
// }


