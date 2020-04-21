package inventory

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// PartialInventory is a partial inventory as sent by the Steam API.
type PartialInventory struct {
	Inventory

	Success   bool
	Error     string
	More      bool
	MoreStart MoreStart `json:"more_start"`
}

type MoreStart uint

func (m *MoreStart) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("false")) {
		return nil
	}

	return json.Unmarshal(data, (*uint)(m))
}

func PerformRequest(client *http.Client, req *http.Request) (*PartialInventory, error) {
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	inv := &PartialInventory{}

	if err = json.NewDecoder(resp.Body).Decode(inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// TODO: use only one iterator function and indicate the first page with start = 0
// type PartialInventoryFetcher func(start uint) (*PartialInventory, error)

func GetFullInventory(
	getFirst func() (*PartialInventory, error),
	getNext func(start uint) (*PartialInventory, error),
) (*Inventory, error) {
	first, err := getFirst()
	if err != nil {
		return nil, err
	}
	if !first.Success {
		return nil, errors.New("GetFullInventory API call failed: " + first.Error)
	}

	result := &first.Inventory
	var next *PartialInventory
	for latest := first; latest.More; latest = next {
		next, err := getNext(uint(latest.MoreStart))
		if err != nil {
			return nil, err
		}
		if !next.Success {
			return nil, errors.New("GetFullInventory API call failed: " + next.Error)
		}

		result = Merge(result, &next.Inventory)
	}

	return result, nil
}

// Merges the given Inventory into a single Inventory.
// The given slice must have at least one element. The first element of the slice is used
// and modified.
func Merge(p ...*Inventory) *Inventory {
	inv := p[0]
	for idx, i := range p {
		if idx == 0 {
			continue
		}

		for key, value := range i.Items {
			inv.Items[key] = value
		}
		for key, value := range i.Descriptions {
			inv.Descriptions[key] = value
		}
		for key, value := range i.Currencies {
			inv.Currencies[key] = value
		}
	}

	return inv
}
