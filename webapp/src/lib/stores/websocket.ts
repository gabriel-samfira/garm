import { writable, get } from 'svelte/store';

// Event types that match the websocket API
export type EntityType = 
	| 'repository' 
	| 'organization' 
	| 'enterprise' 
	| 'pool' 
	| 'user' 
	| 'instance' 
	| 'job' 
	| 'controller' 
	| 'github_credentials' 
	| 'gitea_credentials' 
	| 'github_endpoint' 
	| 'scaleset';

export type Operation = 'create' | 'update' | 'delete';

export interface EventFilter {
	'entity-type': EntityType;
	operations: Operation[];
}

export interface FilterMessage {
	'send-everything'?: boolean;
	filters?: EventFilter[];
}

export interface WebSocketEvent {
	'entity-type': EntityType;
	operation: Operation;
	payload: any;
}

export interface WebSocketState {
	connected: boolean;
	connecting: boolean;
	error: string | null;
	lastEvent: WebSocketEvent | null;
}

// Create the websocket store
function createWebSocketStore() {
	const { subscribe, set, update } = writable<WebSocketState>({
		connected: false,
		connecting: false,
		error: null,
		lastEvent: null
	});

	let ws: WebSocket | null = null;
	let reconnectAttempts = 0;
	let maxReconnectAttempts = 50; // Increased for more persistent reconnection
	let baseReconnectInterval = 1000; // Base interval
	let reconnectInterval = 1000; // Current interval
	let maxReconnectInterval = 30000; // Max 30 seconds
	let reconnectTimeout: number | null = null;
	let pingInterval: number | null = null;
	let healthCheckInterval: number | null = null;
	let currentFilters: EventFilter[] = [];
	let manuallyDisconnected = false;
	let lastPongReceived = Date.now();

	// Event callbacks organized by entity type
	const eventCallbacks = new Map<EntityType, ((event: WebSocketEvent) => void)[]>();

	function getWebSocketUrl(): string {
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const host = window.location.host;
		return `${protocol}//${host}/api/v1/ws/events`;
	}

	function connect() {
		if (ws && (ws.readyState === WebSocket.CONNECTING || ws.readyState === WebSocket.OPEN)) {
			return;
		}

		manuallyDisconnected = false;
		update(state => ({ ...state, connecting: true, error: null }));

		try {
			const wsUrl = getWebSocketUrl();
			console.log(`[WebSocket] Connecting to ${wsUrl} (attempt ${reconnectAttempts + 1}/${maxReconnectAttempts})`);
			
			// Use cookie authentication - no need for Bearer token in protocol
			ws = new WebSocket(wsUrl);

			// Set connection timeout
			const connectionTimeout = setTimeout(() => {
				if (ws && ws.readyState === WebSocket.CONNECTING) {
					console.log('[WebSocket] Connection timeout');
					ws.close();
				}
			}, 10000); // 10 second timeout

			ws.onopen = () => {
				clearTimeout(connectionTimeout);
				console.log('[WebSocket] Connected to events endpoint');
				reconnectAttempts = 0;
				reconnectInterval = baseReconnectInterval;
				lastPongReceived = Date.now();
				
				update(state => ({ ...state, connected: true, connecting: false, error: null }));

				// Send current filters if any
				if (currentFilters.length > 0) {
					sendFilters(currentFilters);
				}

				// Setup ping to keep connection alive
				startHeartbeat();
			};

			ws.onmessage = (event) => {
				try {
					const data = JSON.parse(event.data);
					
					// Handle pong responses
					if (data.type === 'pong' || data.type === 'ping') {
						lastPongReceived = Date.now();
						return;
					}

					console.log('[WebSocket] Received event:', data);
					lastPongReceived = Date.now(); // Any message counts as activity

					// Update the store with the last event
					update(state => ({ ...state, lastEvent: data }));

					// Call registered callbacks for this entity type
					const callbacks = eventCallbacks.get(data['entity-type']) || [];
					callbacks.forEach(callback => {
						try {
							callback(data);
						} catch (err) {
							console.error('[WebSocket] Error in event callback:', err);
						}
					});
				} catch (err) {
					console.error('[WebSocket] Error parsing message:', err);
				}
			};

			ws.onclose = (event) => {
				clearTimeout(connectionTimeout);
				console.log('[WebSocket] Connection closed:', event.code, event.reason);
				cleanup();
				
				const wasManualDisconnect = event.code === 1000 && manuallyDisconnected;
				const errorMessage = event.code !== 1000 ? `Connection closed: ${event.reason || 'Unknown reason'}` : null;
				
				update(state => ({ 
					...state, 
					connected: false, 
					connecting: false,
					error: errorMessage
				}));

				// Attempt to reconnect unless it was explicitly a manual disconnect
				// This includes server restarts that result in clean closes (code 1000)
				if (!wasManualDisconnect) {
					console.log('[WebSocket] Connection lost, scheduling reconnect...');
					scheduleReconnect();
				}
			};

			ws.onerror = (error) => {
				clearTimeout(connectionTimeout);
				console.error('[WebSocket] Connection error:', error);
				cleanup();
				
				update(state => ({ 
					...state, 
					connected: false, 
					connecting: false, 
					error: 'WebSocket connection error' 
				}));

				// Schedule reconnect on error if not manually disconnected
				if (!manuallyDisconnected) {
					scheduleReconnect();
				}
			};

		} catch (err) {
			console.error('[WebSocket] Failed to create connection:', err);
			update(state => ({ 
				...state, 
				connected: false, 
				connecting: false, 
				error: err instanceof Error ? err.message : 'Failed to connect' 
			}));
		}
	}

	function startHeartbeat() {
		// Clear any existing intervals
		cleanup();

		// Setup ping to keep connection alive and detect dead connections
		pingInterval = window.setInterval(() => {
			if (ws && ws.readyState === WebSocket.OPEN) {
				// Check if we received a pong recently
				const timeSinceLastPong = Date.now() - lastPongReceived;
				if (timeSinceLastPong > 90000) { // No pong for 90 seconds
					console.log('[WebSocket] Connection appears dead (no pong), reconnecting...');
					ws.close();
					return;
				}

				// Send ping (note: browser WebSocket doesn't support ping frames, 
				// so we'll send a JSON ping message that the server should echo)
				try {
					ws.send(JSON.stringify({ type: 'ping', timestamp: Date.now() }));
				} catch (err) {
					console.error('[WebSocket] Failed to send ping:', err);
					ws.close();
				}
			}
		}, 30000); // Ping every 30 seconds

		// Health check interval to detect stale connections
		healthCheckInterval = window.setInterval(() => {
			if (ws && ws.readyState === WebSocket.OPEN) {
				// Additional health check - if we haven't received any data for too long,
				// assume the connection is stale
				const timeSinceLastPong = Date.now() - lastPongReceived;
				if (timeSinceLastPong > 120000) { // No activity for 2 minutes
					console.log('[WebSocket] Connection appears stale, reconnecting...');
					ws.close();
				}
			}
		}, 60000); // Check every minute
	}

	function cleanup() {
		if (pingInterval) {
			clearInterval(pingInterval);
			pingInterval = null;
		}
		if (healthCheckInterval) {
			clearInterval(healthCheckInterval);
			healthCheckInterval = null;
		}
	}

	function scheduleReconnect() {
		if (manuallyDisconnected) {
			return;
		}

		if (reconnectTimeout) {
			clearTimeout(reconnectTimeout);
		}

		reconnectAttempts++;
		
		// Reset attempts periodically to allow for long-term reconnection
		if (reconnectAttempts > maxReconnectAttempts) {
			console.log('[WebSocket] Max reconnect attempts reached, resetting counter for persistent reconnection');
			reconnectAttempts = 1;
			reconnectInterval = baseReconnectInterval;
		}

		const actualInterval = Math.min(reconnectInterval, maxReconnectInterval);
		console.log(`[WebSocket] Scheduling reconnect attempt ${reconnectAttempts} in ${actualInterval}ms`);

		reconnectTimeout = window.setTimeout(() => {
			if (!manuallyDisconnected) {
				connect();
				// Exponential backoff with jitter to avoid thundering herd
				const jitter = Math.random() * 1000; // 0-1 second jitter
				reconnectInterval = Math.min(reconnectInterval * 1.5 + jitter, maxReconnectInterval);
			}
		}, actualInterval);
	}

	function sendFilters(filters: EventFilter[]) {
		if (ws && ws.readyState === WebSocket.OPEN) {
			const message: FilterMessage = {
				'send-everything': false,
				filters: filters
			};
			console.log('[WebSocket] Sending filters:', message);
			ws.send(JSON.stringify(message));
			currentFilters = [...filters];
		}
	}

	function disconnect() {
		manuallyDisconnected = true;
		
		if (reconnectTimeout) {
			clearTimeout(reconnectTimeout);
			reconnectTimeout = null;
		}
		
		cleanup();

		if (ws) {
			ws.close(1000, 'Manual disconnect');
			ws = null;
		}

		// Clear all callbacks
		eventCallbacks.clear();
		currentFilters = [];

		update(state => ({ 
			...state, 
			connected: false, 
			connecting: false, 
			error: null,
			lastEvent: null
		}));
	}

	// Handle network connectivity changes
	function handleNetworkChange() {
		if (navigator.onLine && !manuallyDisconnected) {
			console.log('[WebSocket] Network came back online, attempting to reconnect');
			// Delay reconnection slightly to allow network to stabilize
			setTimeout(() => {
				if (!ws || ws.readyState === WebSocket.CLOSED || ws.readyState === WebSocket.CLOSING) {
					reconnectAttempts = 0; // Reset attempts on network recovery
					reconnectInterval = baseReconnectInterval;
					connect();
				}
			}, 2000);
		}
	}

	// Listen for network changes
	if (typeof window !== 'undefined') {
		window.addEventListener('online', handleNetworkChange);
		window.addEventListener('offline', () => {
			console.log('[WebSocket] Network went offline');
			update(state => ({ ...state, error: 'Network offline' }));
		});

		// Periodic check to ensure connection is maintained when there are active subscriptions
		setInterval(() => {
			// Only check if we have active subscriptions and are not manually disconnected
			if (!manuallyDisconnected && eventCallbacks.size > 0) {
				// If we should be connected but aren't, attempt to reconnect
				if (!ws || ws.readyState === WebSocket.CLOSED || ws.readyState === WebSocket.CLOSING) {
					console.log('[WebSocket] Periodic check: connection missing with active subscriptions, reconnecting...');
					connect();
				}
			}
		}, 10000); // Check every 10 seconds
	}

	// Subscribe to events for a specific entity type
	function subscribeToEntity(entityType: EntityType, operations: Operation[], callback: (event: WebSocketEvent) => void) {
		// Reset manual disconnect flag when new subscriptions are added
		manuallyDisconnected = false;
		
		// Add callback to the list for this entity type
		if (!eventCallbacks.has(entityType)) {
			eventCallbacks.set(entityType, []);
		}
		eventCallbacks.get(entityType)!.push(callback);

		// Add or update the filter for this entity type
		const existingFilterIndex = currentFilters.findIndex(f => f['entity-type'] === entityType);
		const newFilter: EventFilter = {
			'entity-type': entityType,
			operations: operations
		};

		if (existingFilterIndex >= 0) {
			// Merge operations with existing filter
			const existingOps = currentFilters[existingFilterIndex].operations;
			newFilter.operations = Array.from(new Set([...existingOps, ...operations]));
			currentFilters[existingFilterIndex] = newFilter;
		} else {
			currentFilters.push(newFilter);
		}

		// Send updated filters if connected
		if (ws && ws.readyState === WebSocket.OPEN) {
			sendFilters(currentFilters);
		}

		// Connect if not already connected or connecting
		if (!ws || ws.readyState === WebSocket.CLOSED || ws.readyState === WebSocket.CLOSING) {
			console.log('[WebSocket] New subscription added, ensuring connection...');
			connect();
		}

		// Return unsubscribe function
		return () => {
			const callbacks = eventCallbacks.get(entityType);
			if (callbacks) {
				const index = callbacks.indexOf(callback);
				if (index > -1) {
					callbacks.splice(index, 1);
				}
				
				// If no more callbacks for this entity type, remove the filter
				if (callbacks.length === 0) {
					eventCallbacks.delete(entityType);
					const filterIndex = currentFilters.findIndex(f => f['entity-type'] === entityType);
					if (filterIndex > -1) {
						currentFilters.splice(filterIndex, 1);
						if (ws && ws.readyState === WebSocket.OPEN) {
							sendFilters(currentFilters);
						}
					}
				}
			}
		};
	}

	return {
		subscribe,
		connect,
		disconnect,
		subscribeToEntity
	};
}

export const websocketStore = createWebSocketStore();