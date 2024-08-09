package models

type Account struct {
	AccountNumber     int     `json:"account_number"`
	AccountBalance    float64 `json:"account_balance"`
	AccountDrawdown   float64 `json:"account_drawdown"`
	AccountEquity     float64 `json:"account_equity"`
	AccountMargin     float64 `json:"account_margin"`
	AccountFreeMargin float64 `json:"account_free_margin"`
}

type Symbol struct {
	TotalDD  int     `json:"total_dd"`
	SellDD   int     `json:"sell_dd"`
	BuyDD    int     `json:"buy_dd"`
	NetLot   float64 `json:"net_lot"`
	TotalLot float64 `json:"total_lot"`
	SellLot  float64 `json:"sell_lot"`
	BuyLot   float64 `json:"buy_lot"`
	NetOC    int     `json:"net_oc"`
	TotalOC  int     `json:"total_oc"`
	SellOC   int     `json:"sell_oc"`
	BuyOC    int     `json:"buy_oc"`
	Profit   float64 `json:"profit"`
}

type Total struct {
	TotalDD      int     `json:"total_dd"`
	SellDD       int     `json:"sell_dd"`
	BuyDD        int     `json:"buy_dd"`
	NetLot       float64 `json:"net_lot"`
	TotalLot     float64 `json:"total_lot"`
	TotalSellLot float64 `json:"total_sell_lot"`
	TotalBuyLot  float64 `json:"total_buy_lot"`
	NetOC        int     `json:"net_oc"`
	TotalOC      int     `json:"total_oc"`
	SellOC       int     `json:"sell_oc"`
	BuyOC        int     `json:"buy_oc"`
	Profit       float64 `json:"profit"`
}

type Message struct {
	Account Account           `json:"account"`
	Total   Total             `json:"total"`
	Symbols map[string]Symbol `json:"symbols"`
}
