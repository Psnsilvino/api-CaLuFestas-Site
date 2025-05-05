package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                  primitive.ObjectID `json:"_id,omitempty"`
	Nome                string             `json:"nome"`
	Categoria           string             `json:"categoria"`
	Subcategoria           string             `json:"subcategoria"`
	Quantidade          int                `json:"quantidade"`
	QuantidadeEmLocacao int                `json:"quantidadeEmLocacao"`
	Preco               float64            `json:"preco"`
	Descricao           string             `json:"descricao"`
	Imagem              string             `json:"imagem"`
}
