package credential

import (
	"testing"
)

func TestMockStore_SetAndGet(t *testing.T) {
	store := NewMockStore()

	if err := store.Set("svc", "key1", "value1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := store.Get("svc", "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != "value1" {
		t.Errorf("Get = %q, want %q", got, "value1")
	}
}

func TestMockStore_SetOverwrite(t *testing.T) {
	store := NewMockStore()

	store.Set("svc", "key1", "old")
	store.Set("svc", "key1", "new")

	got, err := store.Get("svc", "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != "new" {
		t.Errorf("Get = %q, want %q", got, "new")
	}
}

func TestMockStore_Delete(t *testing.T) {
	store := NewMockStore()

	store.Set("svc", "key1", "value1")

	if err := store.Delete("svc", "key1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := store.Get("svc", "key1")
	if err == nil {
		t.Fatal("expected error after Delete, got nil")
	}
}

func TestMockStore_GetNonExistentKey(t *testing.T) {
	store := NewMockStore()

	_, err := store.Get("svc", "missing")
	if err == nil {
		t.Fatal("expected error for non-existent key, got nil")
	}
}

func TestMockStore_GetNonExistentService(t *testing.T) {
	store := NewMockStore()

	_, err := store.Get("nosvc", "key")
	if err == nil {
		t.Fatal("expected error for non-existent service, got nil")
	}
}

func TestMockStore_DeleteNonExistentKey(t *testing.T) {
	store := NewMockStore()

	err := store.Delete("svc", "missing")
	if err == nil {
		t.Fatal("expected error for deleting non-existent key, got nil")
	}
}

func TestMockStore_DeleteNonExistentService(t *testing.T) {
	store := NewMockStore()

	err := store.Delete("nosvc", "key")
	if err == nil {
		t.Fatal("expected error for deleting from non-existent service, got nil")
	}
}

func TestMockStore_MultipleServices(t *testing.T) {
	store := NewMockStore()

	store.Set("svc1", "key", "val1")
	store.Set("svc2", "key", "val2")

	got1, _ := store.Get("svc1", "key")
	got2, _ := store.Get("svc2", "key")

	if got1 != "val1" {
		t.Errorf("svc1 Get = %q, want %q", got1, "val1")
	}
	if got2 != "val2" {
		t.Errorf("svc2 Get = %q, want %q", got2, "val2")
	}
}
