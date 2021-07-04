package model

import (
	"github.com/kamva/mgm/v3"
)

type News struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string `json:"title" faker:"sentence"`
	Content          string `json:"content" faker:"paragraph"`
	CreatedAt        string `json:"created_at" faker:"date"`
	Owner            string `json:"owner" faker:"name"`
}
