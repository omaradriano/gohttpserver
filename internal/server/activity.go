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
	c.activities = append(c.activities, act) 
	return act.ID 
}

func (c *Activities) GetActivity(id uint64) (Activity, error) {
	if id > uint64(len(c.activities)) {
		return Activity{}, ErrIDNotFound
	}

	return c.activities[id-1], nil
}


