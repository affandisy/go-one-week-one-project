package services

type OpnameItemReq struct {
	ProductID uint `json:"product_id"`
	ActualQty int  `json:"actual_qty"`
}

type ProcessOpnameReq struct {
	Notes string          `json:"notes"`
	Items []OpnameItemReq `json:"items"`
}
