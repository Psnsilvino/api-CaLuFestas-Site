package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                  primitive.ObjectID `json:"_id,omitempty"`
	Nome                string             `json:"nome"`
	Categoria           string             `json:"categoria"`
	Quantidade          int                `json:"quantidade"`
	QuantidadeEmLocacao int                `json:"quantidadeEmLocacao"`
	Preco               float64            `json:"preco"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
}

type ProductResponse struct {
    ID                  string  `json:"_id"`
    Nome                string  `json:"nome"`
    Categoria           string  `json:"categoria"`
    Quantidade          int     `json:"quantidade"`
    QuantidadeEmLocacao int     `json:"quantidadeEmLocacao"`
    Preco               float64 `json:"preco"`
    CreatedAt           time.Time `json:"createdAt"`
    UpdatedAt           time.Time `json:"updatedAt"`
}

