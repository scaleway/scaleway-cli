package inference

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	product_catalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
)

type ListServerlessModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelCustom struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	OwnedBy          string    `json:"owned_by"`
	ConsumptionMode  string    `json:"consumption_mode"`
	InputTokenPrice  float64   `json:"input_token_price"`
	OutputTokenPrice float64   `json:"output_token_price"`
	Tasks            []string  `json:"tasks"`
	Reasoning        bool      `json:"reasoning"`
}

func listServerlessModels() *core.Command {
	return &core.Command{
		Short:     `List available serverless models`,
		Long:      `List all available serverless models using the OpenAI-compatible API endpoint. Prices per million of tokens.`,
		Namespace: "inference",
		Resource:  "model",
		Verb:      "list-serverless",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
		Run: func(ctx context.Context, argsI any) (any, error) {
			scwClient := core.ExtractClient(ctx)

			secret, ok := scwClient.GetSecretKey()
			if !ok {
				return nil, errors.New("secret key not found")
			}

			resp, err := fetchModels(ctx, secret)
			if err != nil {
				return nil, err
			}

			catalogApi := product_catalog.NewPublicCatalogAPI(scwClient)
			catalogRes, err := catalogApi.ListPublicCatalogProducts(
				&product_catalog.PublicCatalogAPIListPublicCatalogProductsRequest{
					ProductTypes: []product_catalog.ListPublicCatalogProductsRequestProductType{
						product_catalog.ListPublicCatalogProductsRequestProductTypeGenerativeAPIs,
					},
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch product catalog: %w", err)
			}

			// Build a map of information per model and consumption mode
			// Key: "modelName:consumptionMode" (e.g., "glm-5.2:realtime")
			type ModelInfo struct {
				InputPrice  float64
				OutputPrice float64
				Tasks       []string
				Reasoning   bool
			}
			infoByModelAndMode := make(map[string]ModelInfo)

			for _, product := range catalogRes.Products {
				if product.Properties == nil || product.Properties.GenerativeAPIs == nil {
					continue
				}
				if product.Price == nil || product.Price.RetailPrice == nil {
					continue
				}

				modelName := product.Product
				genAPIs := product.Properties.GenerativeAPIs
				consumptionMode := string(genAPIs.ConsumptionMode)
				price := product.Price.RetailPrice.ToFloat()

				// Create a unique key per model + consumption mode
				key := modelName + ":" + consumptionMode

				info := infoByModelAndMode[key]

				info.Reasoning = genAPIs.Reasoning

				// Accumulate unique tasks
				taskSet := make(map[string]bool)
				for _, t := range info.Tasks {
					taskSet[t] = true
				}
				for _, t := range genAPIs.Tasks {
					if !taskSet[string(t)] {
						info.Tasks = append(info.Tasks, string(t))
						taskSet[string(t)] = true
					}
				}

				if genAPIs.TokenType == product_catalog.PublicCatalogProductPropertiesGenerativeAPIsTokenTypeInputToken {
					info.InputPrice = price * 1_000
				} else if genAPIs.TokenType == product_catalog.PublicCatalogProductPropertiesGenerativeAPIsTokenTypeOutputToken {
					info.OutputPrice = price * 1_000
				}
				infoByModelAndMode[key] = info
			}

			// For each model from the OpenAI API, create entries for each available consumption mode
			var models []ModelCustom
			for _, m := range resp.Data {
				// Check if we have realtime info
				realtimeKey := m.ID + ":realtime"
				if info, ok := infoByModelAndMode[realtimeKey]; ok {
					models = append(models, ModelCustom{
						ID:               m.ID,
						CreatedAt:        time.Unix(m.Created, 0),
						OwnedBy:          m.OwnedBy,
						ConsumptionMode:  "realtime",
						InputTokenPrice:  info.InputPrice,
						OutputTokenPrice: info.OutputPrice,
						Tasks:            info.Tasks,
						Reasoning:        info.Reasoning,
					})
				}

				// Check if we have batch info
				batchKey := m.ID + ":batch"
				if info, ok := infoByModelAndMode[batchKey]; ok {
					models = append(models, ModelCustom{
						ID:               m.ID,
						CreatedAt:        time.Unix(m.Created, 0),
						OwnedBy:          m.OwnedBy,
						ConsumptionMode:  "batch",
						InputTokenPrice:  info.InputPrice,
						OutputTokenPrice: info.OutputPrice,
						Tasks:            info.Tasks,
						Reasoning:        info.Reasoning,
					})
				}

				// If no specific mode is found, add an entry without mode
				if _, hasRealtime := infoByModelAndMode[realtimeKey]; !hasRealtime {
					if _, hasBatch := infoByModelAndMode[batchKey]; !hasBatch {
						// No mode found, add with default info (if available)
						info := infoByModelAndMode[m.ID]
						models = append(models, ModelCustom{
							ID:               m.ID,
							CreatedAt:        time.Unix(m.Created, 0),
							OwnedBy:          m.OwnedBy,
							ConsumptionMode:  "",
							InputTokenPrice:  info.InputPrice,
							OutputTokenPrice: info.OutputPrice,
							Tasks:            info.Tasks,
							Reasoning:        info.Reasoning,
						})
					}
				}
			}

			return models, nil
		},
	}
}

func fetchModels(ctx context.Context, secret string) (*ListServerlessModelsResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.scaleway.ai/v1/models",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var listResp ListServerlessModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &listResp, nil
}
