package steamid

import (
	"testing"
)

// TestSteamID3 tests a steamid3 format
func TestSteamID3(t *testing.T) {
	sid := "[U:1:69038686]"
	id, err := NewId(sid)
	if err != nil {
		t.Fatal(err.Error())
	}
	if id.ToUint64() != uint64(76561198029304414) {
		t.Fatalf("%d != 76561198029304414", id.ToUint64())
	}
}

// TestSteamID64 tests a steamid64 format
func TestSteamID64(t *testing.T) {
	sid := "76561198029304414"
	id, err := NewId(sid)
	if err != nil {
		t.Fatal(err.Error())
	}
	if id.ToUint64() != uint64(76561198029304414) {
		t.Fatalf("%d != 76561198029304414", id.ToUint64())
	}
}

// TestSteamID32 tests a steamid32 format
func TestSteamID32(t *testing.T) {
	sid := "69038686"
	id, err := NewId(sid)
	if err != nil {
		t.Fatal(err.Error())
	}
	if id.ToUint64() != uint64(76561198029304414) {
		t.Fatalf("%d != 76561198029304414", id.ToUint64())
	}
}
