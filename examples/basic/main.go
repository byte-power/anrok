// Basic 演示 Anrok SDK 的临时计税（CreateEphemeral）与落库交易（CreateOrUpdate）。
//
// 运行前设置环境变量 ANROK_API_KEY，然后在模块根目录执行：
//
//	go run ./examples/basic
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/byte-power/anrok"
)

func main() {
	apiKey := os.Getenv("ANROK_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置环境变量 ANROK_API_KEY（Anrok 控制台中的 API key）")
	}

	client := anrok.NewClient(apiKey, nil)

	lineItems := []anrok.TransactionLineItem{
		{
			ID:                "item-1",
			ProductExternalID: "test_oneff_product_id",
			Amount:            15000,
			Quantity:          "1",
		},
		{
			ID:                "item-2",
			ProductExternalID: "test_product_id_sub",
			Amount:            31000,
			Quantity:          "12.3",
		},
	}
	incl := true
	lineItems[1].IsTaxIncludedInAmount = &incl

	customerAddr := anrok.Address{
		Country:    "us",
		Line1:      "1450 Cherokee St",
		City:       "Denver",
		Region:     "CO",
		PostalCode: "80204",
	}

	shipFrom := &anrok.Address{
		Country:    "us",
		Line1:      "230 S LaSalle St",
		City:       "Chicago",
		Region:     "IL",
		PostalCode: "60604",
	}

	now := time.Now().UTC().Format(time.RFC3339)

	// 1) 临时交易：只算税，不在 Anrok 中保存（适合开票前预览）
	ephemeralReq := anrok.CreateEphemeralTransactionRequest{
		CurrencyCode:       "usd",
		AccountingTime:     now,
		AccountingTimeZone: "UTC",
		LineItems:          lineItems,
		CustomerAddress:    customerAddr,
		ShipFromAddress:    shipFrom,
		CustomerName:       "Example Customer",
	}

	ephemeral, err := client.CreateEphemeralTransaction(ephemeralReq)
	if err != nil {
		handleError("CreateEphemeralTransaction", err)
	}
	printJSON("CreateEphemeralTransaction 响应摘要", map[string]any{
		"taxAmountToCollect": ephemeral.TaxAmountToCollect,
		"preTaxAmount":       ephemeral.PreTaxAmount,
		"lineItemCount":      len(ephemeral.LineItems),
		"raw":                ephemeral,
	})

	// 2) 创建或更新交易：写入 Anrok（示例使用带时间戳的唯一 id，避免重复运行冲突）
	txID := fmt.Sprintf("example-sdk-%d", time.Now().UnixNano())
	createReq := anrok.CreateOrUpdateTransactionRequest{
		ID:                 txID,
		CurrencyCode:       "usd",
		AccountingTime:     now,
		AccountingTimeZone: "UTC",
		LineItems:          lineItems,
		CustomerAddress:    customerAddr,
		ShipFromAddress:    shipFrom,
		CustomerName:       "Example Customer",
	}

	saved, err := client.CreateOrUpdateTransaction(createReq)
	if err != nil {
		handleError("CreateOrUpdateTransaction", err)
	}
	printJSON("CreateOrUpdateTransaction 响应摘要", map[string]any{
		"version":            saved.Version,
		"taxAmountToCollect": saved.TaxAmountToCollect,
		"preTaxAmount":       saved.PreTaxAmount,
		"transactionId":      txID,
		"raw":                saved,
	})
}

func handleError(op string, err error) {
	var rateLimitErr *anrok.RateLimitError
	if errors.As(err, &rateLimitErr) {
		log.Fatalf("%s: rate limited, retry after %d seconds", op, rateLimitErr.RetryAfter)
	}

	var typedErr *anrok.TypedError
	if errors.As(err, &typedErr) {
		switch {
		case typedErr.IsType(anrok.ErrTaxDateTooFarInFuture):
			log.Fatalf("%s: tax date is too far in the future", op)
		case typedErr.IsType(anrok.ErrCustomerAddressCouldNotResolve):
			log.Fatalf("%s: customer address could not be resolved, please check address fields", op)
		case typedErr.IsType(anrok.ErrProductExternalIdUnknown):
			log.Fatalf("%s: unknown product ID, please add it in Anrok first", op)
		default:
			log.Fatalf("%s: API error %s (HTTP %d)", op, typedErr.Type, typedErr.StatusCode)
		}
	}

	var apiErr *anrok.APIError
	if errors.As(err, &apiErr) {
		log.Fatalf("%s: HTTP %d: %s", op, apiErr.StatusCode, apiErr.Body)
	}

	log.Fatalf("%s: %v", op, err)
}

func printJSON(title string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("%s: marshal: %v", title, err)
	}
	fmt.Printf("\n=== %s ===\n%s\n", title, b)
}
