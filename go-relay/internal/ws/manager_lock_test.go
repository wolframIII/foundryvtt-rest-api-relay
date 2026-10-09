package ws

import "testing"

func TestRemoveClientCallbackRunsOutsideManagerLock(t *testing.T) {
	assertCallbackCanRead := func(t *testing.T, remove func(*ClientManager, *Client)) {
		t.Helper()

		manager := NewClientManager(nil, "test")
		client := &Client{id: "client-1", apiKey: "test-key"}
		manager.clients[client.id] = client
		manager.tokenGroups[client.apiKey] = map[string]bool{client.id: true}

		callbackLockWasFree := false
		manager.OnClientRemoved = func(string, string) {
			callbackLockWasFree = manager.mu.TryRLock()
			if callbackLockWasFree {
				manager.mu.RUnlock()
			}
		}

		remove(manager, client)
		if !callbackLockWasFree {
			t.Fatal("OnClientRemoved ran while the client-manager lock was held")
		}
	}

	t.Run("RemoveClient", func(t *testing.T) {
		assertCallbackCanRead(t, func(manager *ClientManager, client *Client) {
			manager.RemoveClient(client.id)
		})
	})

	t.Run("RemoveClientIfMatch", func(t *testing.T) {
		assertCallbackCanRead(t, func(manager *ClientManager, client *Client) {
			manager.RemoveClientIfMatch(client.id, client)
		})
	})
}
