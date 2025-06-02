package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Nome       string             `json:"nome" bson:"nome"`
	Preco      float64            `json:"preco" bson:"preco"`
	Quantidade int                `json:"quantidade" bson:"quantidade"`
}

type Locacao struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Nome         string             `json:"nome" bson:"nome"`
	Endereco     string             `json:"endereco" bson:"endereco"`
	DataEntrega  string             `json:"data_entrega" bson:"data_entrega"`
	DataRetirada string             `json:"data_retirada" bson:"data_retirada"`
	Pagamento    string             `json:"pagamento" bson:"pagamento"`
	Total        float64            `json:"total" bson:"total"`
	Items        []Item             `json:"items" bson:"items"`
}
