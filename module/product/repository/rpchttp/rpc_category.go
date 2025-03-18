package rpchttp

import (
	"bytes"
	"context"
	"ecommerce/common"
	"ecommerce/module/product/domain/query"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type rpcFindCategoriesByIDs struct {
  url string
}

func NewRPCFindCategoriesByIDs(url string) *rpcFindCategoriesByIDs {return &rpcFindCategoriesByIDs{url}}

func (c *rpcFindCategoriesByIDs) FindCategoriesByIDs(ctx context.Context, categoryIDs []common.UUID) ([]query.CategoryDTO, error) {

  // Tạo payload chứa danh sách category IDs
    payloadData := struct {
        IDs []common.UUID `json:"ids"`
    }{
        IDs: categoryIDs,
    }

    // Chuyển payload sang dạng JSON
    payloadBytes, err := json.Marshal(payloadData)
    if err != nil {
        return nil, fmt.Errorf("error marshalling payload: %w", err)
    }

    findCateByIDsURL := fmt.Sprintf("%s/query-categories-ids", c.url)

    // Tạo request mới với phương thức POST và payload
    req, err := http.NewRequest("POST", findCateByIDsURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }
    // Đính kèm context vào request
    req = req.WithContext(ctx)
    req.Header.Add("Content-Type", "application/json")

    // Tạo HTTP client và thực hiện request
    client := &http.Client{
      Timeout: time.Second * 10,
    }
    res, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }
    defer res.Body.Close()

    // Kiểm tra mã trạng thái HTTP
    if res.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(res.Body)
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(body))
    }

    // Đọc và giải mã kết quả trả về
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    var responseData struct {
      Data []query.CategoryDTO `json:"data"`
    }

    err = json.Unmarshal(body, &responseData)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling response: %w", err)
    }

    return responseData.Data, nil
}