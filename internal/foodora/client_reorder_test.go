package foodora

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientCustomerAddresses(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method=%s", r.Method)
		}
		if r.URL.Path != "/customers/addresses" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"data":{"items":[{"id":"1","is_selected":true}]}}`))
	}))
	t.Cleanup(srv.Close)

	c, err := New(Options{
		BaseURL:     srv.URL + "/",
		AccessToken: "tok",
		UserAgent:   "ua",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	resp, err := c.CustomerAddresses(context.Background())
	if err != nil {
		t.Fatalf("CustomerAddresses: %v", err)
	}
	if resp.Status != 200 {
		t.Fatalf("status=%d", resp.Status)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Items[0]["id"] != "1" {
		t.Fatalf("unexpected items: %#v", resp.Data.Items)
	}
}

func TestClientOrderReorder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s", r.Method)
		}
		if r.URL.Path != "/orders/abc/reorder" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			t.Fatalf("content-type=%q", ct)
		}

		var got ReorderRequestBody
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.ReorderTime != "2025-12-20T01:02:03+0100" {
			t.Fatalf("reorder_time=%q", got.ReorderTime)
		}
		if got.Address == nil || got.Address["id"] != "addr1" {
			t.Fatalf("address=%#v", got.Address)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "status": 200,
  "data": {
    "vendor_id": 1,
    "vendor_code": "v",
    "vendor_info": {"name":"Test Vendor","vertical":"restaurants","time_zone":"Europe/Vienna"},
    "cart": {
      "total_value": 12.3,
      "vendor_cart": [
        {"products":[{"name":"Burger","variation_name":"Cheese","quantity":1,"total_price":12.3,"price":12.3,"is_available":true,"toppings":[]}]}
      ]
    }
  }
}`))
	}))
	t.Cleanup(srv.Close)

	c, err := New(Options{
		BaseURL:     srv.URL + "/",
		AccessToken: "tok",
		UserAgent:   "ua",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	resp, err := c.OrderReorder(context.Background(), "abc", ReorderRequestBody{
		Address:     map[string]any{"id": "addr1"},
		ReorderTime: "2025-12-20T01:02:03+0100",
	})
	if err != nil {
		t.Fatalf("OrderReorder: %v", err)
	}
	if resp.Status != 200 {
		t.Fatalf("status=%d", resp.Status)
	}
	if resp.Data.VendorCode != "v" || resp.Data.VendorInfo == nil || resp.Data.VendorInfo.Name != "Test Vendor" {
		t.Fatalf("unexpected vendor: %#v", resp.Data)
	}
	if resp.Data.Cart.TotalValue != 12.3 {
		t.Fatalf("cart total=%v", resp.Data.Cart.TotalValue)
	}
	if len(resp.Data.Cart.VendorCart) != 1 || len(resp.Data.Cart.VendorCart[0].Products) != 1 {
		t.Fatalf("unexpected products: %#v", resp.Data.Cart.VendorCart)
	}
}
