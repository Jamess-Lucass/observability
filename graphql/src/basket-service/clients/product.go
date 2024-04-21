package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type GraphQLResponse[T any] struct {
	Data *T `json:"data"`
}

type Product struct {
	Id *uuid.UUID `json:"id"`
}

// TODO: Better error handling
func GetProduct(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
	query Product($productId: UUID!) {
		product(id: $productId) {
		  id
		}
	  }
    `

	variables := map[string]uuid.UUID{
		"productId": id,
	}

	request := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/graphql", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var response GraphQLResponse[Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	if response.Data.Id != nil {
		return true, nil
	}

	return false, nil
}
