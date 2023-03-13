package articleconst

const (
	StatusValueAvailable   string = "available"
	StatusValueOutOfStock  string = "outofstock"
	StatusValueUnavailable string = "unavailable"

	StatusClassAvailable   string = "stock-row-available"
	StatusClassOutOfStock  string = "stock-row-outofstock"
	StatusClassUnavailable string = "stock-row-unavailable"
	StatusClassError       string = "stock-row-error"

	StatusLabelAvailable   string = "Disponible"
	StatusLabelOutOfStock  string = "Epuisé"
	StatusLabelUnavailable string = "Indisponible"
	StatusLabelError       string = "Status inconnu"

	FilterValueAll string = ""
	FilterValueDes string = "DES:"
	FilterValueRef string = "REF:"
	FilterValueCat string = "CAT:"

	FilterLabelAll string = "Tout"
	FilterLabelDes string = "Désignation"
	FilterLabelRef string = "Référence"
	FilterLabelCat string = "Catégorie"

	UnitLiter       string = "l"
	UnitMeter       string = "m"
	UnitLinearMeter string = "ml"
	UnitPiece       string = "unit"
	UnitSquareMeter string = "m²"
	UnitRoll        string = "rouleau"
)
