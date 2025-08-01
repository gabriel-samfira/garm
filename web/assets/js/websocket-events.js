/**
 * GARM WebSocket Events Client
 * Connects to the GARM WebSocket events endpoint and handles real-time updates
 */

class GarmWebSocketClient {
    constructor() {
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000; // Start with 1 second
        this.maxReconnectDelay = 30000; // Max 30 seconds
        this.handlers = new Map();
        this.isAuthenticated = false;
        this.entityFilters = new Set();
    }

    /**
     * Connect to the WebSocket events endpoint
     */
    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/events/`;
        
        console.log('Connecting to GARM WebSocket:', wsUrl);
        
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = (event) => {
            console.log('WebSocket connected');
            this.reconnectAttempts = 0;
            this.reconnectDelay = 1000;
            this.sendAuthIfNeeded();
        };
        
        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleMessage(data);
            } catch (error) {
                console.error('Failed to parse WebSocket message:', error, event.data);
            }
        };
        
        this.ws.onclose = (event) => {
            console.log('WebSocket connection closed:', event.code, event.reason);
            this.isAuthenticated = false;
            this.scheduleReconnect();
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    /**
     * Handle incoming WebSocket messages
     */
    handleMessage(data) {
        // Check if this is an authentication response or an event
        if (data.authenticated !== undefined) {
            this.isAuthenticated = data.authenticated;
            if (this.isAuthenticated) {
                console.log('WebSocket authenticated successfully');
                this.sendEventFilters();
            } else {
                console.error('WebSocket authentication failed');
            }
            return;
        }

        // Handle entity events
        if (data['entity-type'] && data.operation && data.payload) {
            const entityType = data['entity-type'];
            const operation = data.operation;
            const payload = data.payload;
            
            console.log('Received WebSocket event:', entityType, operation, payload);
            
            // Call registered handlers
            const handlerKey = `${entityType}:${operation}`;
            if (this.handlers.has(handlerKey)) {
                this.handlers.get(handlerKey).forEach(handler => {
                    try {
                        handler(payload, operation, entityType);
                    } catch (error) {
                        console.error('Handler error:', error);
                    }
                });
            }
            
            // Call generic handlers
            const genericKey = `${entityType}:*`;
            if (this.handlers.has(genericKey)) {
                this.handlers.get(genericKey).forEach(handler => {
                    try {
                        handler(payload, operation, entityType);
                    } catch (error) {
                        console.error('Generic handler error:', error);
                    }
                });
            }
        }
    }

    /**
     * Send authentication token if available
     */
    sendAuthIfNeeded() {
        // Try to get JWT token from localStorage or cookies
        const token = this.getAuthToken();
        if (token) {
            this.sendMessage({
                type: 'auth',
                token: token
            });
        } else {
            // For now, we'll try to establish the connection anyway
            // The WebSocket endpoint might use session authentication
            this.isAuthenticated = true;
            this.sendEventFilters();
        }
    }

    /**
     * Get authentication token from storage
     */
    getAuthToken() {
        // Try localStorage first
        let token = localStorage.getItem('garm_token');
        if (token) return token;
        
        // Try sessionStorage
        token = sessionStorage.getItem('garm_token');
        if (token) return token;
        
        // Try to extract from cookies
        const cookies = document.cookie.split(';');
        for (let cookie of cookies) {
            const [name, value] = cookie.trim().split('=');
            if (name === 'garm_token' || name === 'jwt_token') {
                return value;
            }
        }
        
        return null;
    }

    /**
     * Send event filters to subscribe to specific events
     */
    sendEventFilters() {
        if (this.entityFilters.size === 0) {
            // Subscribe to all events by default for now
            this.sendMessage({
                "send-everything": true
            });
        } else {
            // Send specific filters
            const filters = Array.from(this.entityFilters).map(filter => {
                const [entityType, operations] = filter.split(':');
                return {
                    "entity-type": entityType,
                    "operations": operations === '*' ? ["create", "update", "delete"] : operations.split(',')
                };
            });
            
            this.sendMessage({
                "send-everything": false,
                "filters": filters
            });
        }
    }

    /**
     * Send a message to the WebSocket
     */
    sendMessage(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn('WebSocket not ready, cannot send message:', message);
        }
    }

    /**
     * Subscribe to events for a specific entity type and operation
     */
    subscribe(entityType, operation, handler) {
        const key = `${entityType}:${operation}`;
        
        if (!this.handlers.has(key)) {
            this.handlers.set(key, new Set());
        }
        
        this.handlers.get(key).add(handler);
        
        // Add to entity filters
        this.entityFilters.add(`${entityType}:${operation}`);
        
        // Update filters if already connected
        if (this.isAuthenticated) {
            this.sendEventFilters();
        }
        
        console.log(`Subscribed to ${entityType}:${operation} events`);
    }

    /**
     * Unsubscribe from events
     */
    unsubscribe(entityType, operation, handler) {
        const key = `${entityType}:${operation}`;
        
        if (this.handlers.has(key)) {
            this.handlers.get(key).delete(handler);
            
            if (this.handlers.get(key).size === 0) {
                this.handlers.delete(key);
                this.entityFilters.delete(`${entityType}:${operation}`);
            }
        }
    }

    /**
     * Schedule reconnection with exponential backoff
     */
    scheduleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('Max reconnection attempts reached');
            return;
        }
        
        this.reconnectAttempts++;
        
        console.log(`Reconnecting in ${this.reconnectDelay}ms (attempt ${this.reconnectAttempts})`);
        
        setTimeout(() => {
            this.connect();
        }, this.reconnectDelay);
        
        // Exponential backoff
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, this.maxReconnectDelay);
    }

    /**
     * Disconnect from WebSocket
     */
    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
        this.isAuthenticated = false;
        this.handlers.clear();
        this.entityFilters.clear();
    }

    /**
     * Get connection status
     */
    isConnected() {
        return this.ws && this.ws.readyState === WebSocket.OPEN && this.isAuthenticated;
    }
}

// Create global instance
window.garmWebSocket = new GarmWebSocketClient();

// Auto-connect when the page loads
document.addEventListener('DOMContentLoaded', function() {
    window.garmWebSocket.connect();
});

// Reconnect when the page becomes visible (handles tab switching)
document.addEventListener('visibilitychange', function() {
    if (!document.hidden && !window.garmWebSocket.isConnected()) {
        window.garmWebSocket.connect();
    }
});

// Helper function to update entity events in real-time
function updateEntityEvents(entityId, entityType) {
    // Subscribe to events for this specific entity type
    window.garmWebSocket.subscribe(entityType, '*', function(payload, operation) {
        // Check if this event is for the current entity
        if (payload.id === entityId) {
            // Trigger refresh of events section
            const eventsContainer = document.getElementById(`${entityType}-events-${entityId}`);
            if (eventsContainer) {
                // Use HTMX to refresh the events section
                htmx.trigger(eventsContainer, 'refreshEvents');
            }
        }
    });
}

// Export for use in templates
window.updateEntityEvents = updateEntityEvents;