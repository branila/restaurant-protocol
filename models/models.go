package models

type Modifica struct {
	Tipo string // "+" per aggiungere, "-" per rimuovere
	Voce string
}

type Piatto struct {
	Nome      string
	Modifiche []Modifica
}

type Comanda struct {
	Numero   int
	Primo    *Piatto
	Secondo  *Piatto
	Contorno *Piatto
}

type Ordine struct {
	Tavolo  int
	Data    string
	Comande []Comanda
}
