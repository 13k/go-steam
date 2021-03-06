package trade

import (
	"time"

	"github.com/13k/go-steam/economy/inventory"
	"github.com/13k/go-steam/economy/trade/api"
)

type Slot uint

func (t *Trade) action(status *api.Result, err error) error {
	if err != nil {
		return err
	}

	return t.onStatus(status)
}

// Poll returns the next batch of events to process. These can be queued from calls to methods like
// `AddItem` or, if there are no queued events, from a new HTTP request to Steam's API (blocking!).
// If the latter is the case, this method may also sleep before the request to conform to the
// polling interval of the official Steam client.
func (t *Trade) Poll() ([]interface{}, error) {
	if t.queuedEvents != nil {
		return t.Events(), nil
	}

	if d := time.Since(t.lastPoll); d < pollTimeout {
		time.Sleep(pollTimeout - d)
	}

	t.lastPoll = time.Now()

	if err := t.action(t.api.GetStatus()); err != nil {
		return nil, err
	}

	return t.Events(), nil
}

func (t *Trade) GetTheirInventory(contextID uint64, appID uint32) (*inventory.Inventory, error) {
	return inventory.GetFullInventory(func(start uint) (*inventory.PartialInventory, error) {
		return t.api.GetForeignInventory(contextID, appID, start)
	})
}

func (t *Trade) GetOwnInventory(contextID uint64, appID uint32) (*inventory.Inventory, error) {
	return t.api.GetOwnInventory(contextID, appID)
}

func (t *Trade) GetMain() (*api.Main, error) {
	return t.api.GetMain()
}

func (t *Trade) AddItem(slot Slot, item *Item) error {
	return t.action(t.api.AddItem(uint(slot), item.AssetID, item.ContextID, item.AppID))
}

func (t *Trade) RemoveItem(slot Slot, item *Item) error {
	return t.action(t.api.RemoveItem(uint(slot), item.AssetID, item.ContextID, item.AppID))
}

func (t *Trade) Chat(message string) error {
	return t.action(t.api.Chat(message))
}

func (t *Trade) SetCurrency(amount uint, currency *Currency) error {
	return t.action(t.api.SetCurrency(amount, currency.CurrencyID, currency.ContextID, currency.AppID))
}

func (t *Trade) SetReady(ready bool) error {
	return t.action(t.api.SetReady(ready))
}

// Confirm may only be called after a successful `SetReady(true)`.
func (t *Trade) Confirm() error {
	return t.action(t.api.Confirm())
}

func (t *Trade) Cancel() error {
	return t.action(t.api.Cancel())
}
