package anrok

// Address 客户地址或发货地址（字段均可选，但不可传空字符串；见 Anrok API 文档）
type Address struct {
	Country    string `json:"country,omitempty"`
	Line1      string `json:"line1,omitempty"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
}

// CustomerTaxId 客户税号
type CustomerTaxId struct {
	Value string `json:"value,omitempty"`
}

// TransactionLineItem 交易行项目
type TransactionLineItem struct {
	ID                    string `json:"id"`
	ProductExternalID     string `json:"productExternalId"`
	Amount                int    `json:"amount"`
	IsTaxIncludedInAmount *bool  `json:"isTaxIncludedInAmount,omitempty"`
	Quantity              string `json:"quantity,omitempty"`
}

// CreateOrUpdateTransactionRequest 对应 POST /v1/seller/transactions/createOrUpdate
type CreateOrUpdateTransactionRequest struct {
	LineItems          []TransactionLineItem `json:"lineItems"`
	CurrencyCode       string              `json:"currencyCode"`
	CustomerAddress    Address             `json:"customerAddress"`
	CustomerName       string              `json:"customerName,omitempty"`
	CustomerTaxIds     []CustomerTaxId     `json:"customerTaxIds,omitempty"`
	ShipFromAddress    *Address            `json:"shipFromAddress,omitempty"`
	AccountingDate     string              `json:"accountingDate,omitempty"`
	AccountingTime     string              `json:"accountingTime,omitempty"`
	AccountingTimeZone string              `json:"accountingTimeZone,omitempty"`
	TaxDate            string              `json:"taxDate,omitempty"`
	CustomerID         string              `json:"customerId,omitempty"`
	ID                 string              `json:"id"`
}

// CreateEphemeralTransactionRequest 对应 POST /v1/seller/transactions/createEphemeral（无 id）
type CreateEphemeralTransactionRequest struct {
	LineItems          []TransactionLineItem `json:"lineItems"`
	CurrencyCode       string              `json:"currencyCode"`
	CustomerAddress    Address             `json:"customerAddress"`
	CustomerName       string              `json:"customerName,omitempty"`
	CustomerTaxIds     []CustomerTaxId     `json:"customerTaxIds,omitempty"`
	ShipFromAddress    *Address            `json:"shipFromAddress,omitempty"`
	AccountingDate     string              `json:"accountingDate,omitempty"`
	AccountingTime     string              `json:"accountingTime,omitempty"`
	AccountingTimeZone string              `json:"accountingTimeZone,omitempty"`
	TaxDate            string              `json:"taxDate,omitempty"`
	CustomerID         string              `json:"customerId,omitempty"`
}

// NotTaxedReason 未征税原因
type NotTaxedReason struct {
	Type string `json:"type"`
}

// TransactionTaxDetail 管辖区内单项税明细
type TransactionTaxDetail struct {
	TaxName       string `json:"taxName"`
	TaxableAmount string `json:"taxableAmount"`
	TaxAmount     string `json:"taxAmount"`
	TaxRate       string `json:"taxRate"`
}

// TransactionJurisLine 行项目下的管辖区分解
type TransactionJurisLine struct {
	Name           string                 `json:"name"`
	Taxes          []TransactionTaxDetail `json:"taxes"`
	NotTaxedReason *NotTaxedReason        `json:"notTaxedReason"`
}

// TransactionLineItemResult 响应中的行项目税额结果
type TransactionLineItemResult struct {
	ID                 string                 `json:"id"`
	TaxAmountToCollect int64                  `json:"taxAmountToCollect"`
	PreTaxAmount       string                 `json:"preTaxAmount"`
	Jurises            []TransactionJurisLine `json:"jurises"`
}

// JurisSummary 管辖区汇总
type JurisSummary struct {
	Name            string           `json:"name"`
	NotTaxedReasons []NotTaxedReason `json:"notTaxedReasons"`
}

// CreateOrUpdateTransactionResponse createOrUpdate 成功响应（含 version）
type CreateOrUpdateTransactionResponse struct {
	Version            int                       `json:"version"`
	TaxAmountToCollect int64                     `json:"taxAmountToCollect"`
	LineItems          []TransactionLineItemResult `json:"lineItems"`
	PreTaxAmount       string                    `json:"preTaxAmount"`
	JurisSummaries     []JurisSummary            `json:"jurisSummaries"`
}

// CreateEphemeralTransactionResponse createEphemeral 成功响应（文档示例中无 version）
type CreateEphemeralTransactionResponse struct {
	TaxAmountToCollect int64                       `json:"taxAmountToCollect"`
	LineItems          []TransactionLineItemResult `json:"lineItems"`
	PreTaxAmount       string                      `json:"preTaxAmount"`
	JurisSummaries     []JurisSummary              `json:"jurisSummaries"`
}
