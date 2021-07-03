package model

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type News struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string    `json:"title" faker:"sentence"`
	Content          string    `json:"content" faker:"paragraph"`
	CreatedAt        time.Time `json:"created_at" faker:"date"`
	Owner            string    `json:"owner" faker:"name"`
}
