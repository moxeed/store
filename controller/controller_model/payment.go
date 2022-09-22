package controller_model

type OpenTerminalModel struct {
	OrderCode uint `query:"orderCode"`
}

type VerifyModel struct {
	Authority string `query:"Authority"`
}
