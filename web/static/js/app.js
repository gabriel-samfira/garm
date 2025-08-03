// GARM Web Interface JavaScript

// Global HTMX configuration
document.addEventListener('DOMContentLoaded', function() {
    // Configure HTMX
    htmx.config.globalViewTransitions = true;
    htmx.config.defaultSwapStyle = 'innerHTML';
    htmx.config.scrollBehavior = 'smooth';
    
    // Global error handler
    document.body.addEventListener('htmx:responseError', function(evt) {
        showToast('Request failed: ' + evt.detail.xhr.status, 'error');
    });
    
    // Global success handler for form submissions
    document.body.addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.xhr.status >= 200 && evt.detail.xhr.status < 300) {
            const trigger = evt.detail.elt;
            const response = evt.detail.xhr;
            
            // Check if response has custom HX-Trigger with showToast (skip global message)
            const hxTrigger = response.getResponseHeader('HX-Trigger');
            const hasCustomToast = hxTrigger && hxTrigger.includes('showToast');
            
            console.log('HTMX afterRequest:', {
                status: evt.detail.xhr.status,
                tagName: trigger.tagName,
                hxTrigger: hxTrigger,
                hasCustomToast: hasCustomToast
            });
            
            if (trigger.tagName === 'FORM' && !hasCustomToast) {
                showToast('Operation completed successfully', 'success');
                closeModal();
            }
        }
    });
    
    // Handle custom HTMX triggers
    document.body.addEventListener('repositoryUpdated', function(evt) {
        closeModal();
    });
    
    document.body.addEventListener('organizationUpdated', function(evt) {
        closeModal();
    });
    
    document.body.addEventListener('enterpriseUpdated', function(evt) {
        closeModal();
    });
    
    document.body.addEventListener('poolUpdated', function(evt) {
        showToast('Pool updated successfully', 'success', 3000);
        closeModal();
    });
    
    // Handle toast notifications from HTMX triggers
    document.body.addEventListener('showToast', function(evt) {
        console.log('showToast event received:', evt.detail);
        if (evt.detail && evt.detail.message) {
            showToast(evt.detail.message, evt.detail.type || 'info', evt.detail.duration || 5000);
        }
    });
    
    // Close modal when clicking outside
    document.addEventListener('click', function(e) {
        try {
            if (e.target && e.target.classList && e.target.classList.contains('modal-backdrop')) {
                closeModal();
            }
        } catch (error) {
            console.error('Error in modal click handler:', error, e.target);
        }
    });
    
    // Handle escape key to close modals
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            closeModal();
        }
    });
    
    // Auto-hide flash messages
    setTimeout(function() {
        const flash = document.querySelector('.flash-message');
        if (flash) {
            hideElement(flash);
        }
    }, 5000);
    
    // Initialize tooltips if any
    initTooltips();
});

// Toast notification system
function showToast(message, type = 'info', duration = 5000) {
    // Remove existing error toasts if showing a new one (only one error at a time)
    if (type === 'error') {
        const existingErrors = document.querySelectorAll('.toast-error');
        existingErrors.forEach(toast => toast.remove());
    }
    
    // Remove any existing toasts if showing a success (success replaces errors)
    if (type === 'success') {
        const existingToasts = document.querySelectorAll('.toast');
        existingToasts.forEach(toast => toast.remove());
    }
    
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    // For error toasts, make them dismissable with duration = 0 meaning no auto-dismiss
    const showCloseButton = type === 'error' || duration === 0;
    
    toast.innerHTML = `
        <div class="flex items-center space-x-2">
            <span>${message}</span>
            ${showCloseButton ? '<button onclick="this.parentElement.parentElement.remove()" class="ml-2 text-current hover:opacity-75">×</button>' : ''}
        </div>
    `;
    
    document.body.appendChild(toast);
    
    // Auto-remove after duration (but only if duration > 0)
    if (duration > 0) {
        setTimeout(() => {
            if (toast.parentElement) {
                toast.remove();
            }
        }, duration);
    }
}

// Modal management
function closeModal() {
    const modalContainer = document.getElementById('modal-container');
    if (modalContainer) {
        modalContainer.innerHTML = '';
    }
}

// Safe modal close function for inline onclick handlers
function safeCloseModal(element) {
    try {
        const modal = element.closest('.fixed');
        if (modal) {
            modal.remove();
        }
    } catch (error) {
        console.error('Error closing modal:', error);
    }
}

function openModal(content) {
    const modalContainer = document.getElementById('modal-container');
    if (modalContainer) {
        modalContainer.innerHTML = content;
    }
}

// Utility functions
function hideElement(element) {
    element.style.transition = 'opacity 0.5s';
    element.style.opacity = '0';
    setTimeout(() => element.remove(), 500);
}

function showElement(element) {
    element.style.transition = 'opacity 0.5s';
    element.style.opacity = '1';
}

// Confirmation dialogs
function confirmDelete(message = 'Are you sure you want to delete this item?') {
    return confirm(message);
}

// Initialize tooltips (if using a tooltip library)
function initTooltips() {
    const tooltipElements = document.querySelectorAll('[data-tooltip]');
    tooltipElements.forEach(element => {
        element.addEventListener('mouseenter', showTooltip);
        element.addEventListener('mouseleave', hideTooltip);
    });
}

function showTooltip(e) {
    const text = e.target.getAttribute('data-tooltip');
    if (!text) return;
    
    const tooltip = document.createElement('div');
    tooltip.className = 'tooltip';
    tooltip.textContent = text;
    tooltip.style.cssText = `
        position: absolute;
        background: #1f2937;
        color: white;
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 12px;
        z-index: 1000;
        pointer-events: none;
    `;
    
    document.body.appendChild(tooltip);
    
    const rect = e.target.getBoundingClientRect();
    tooltip.style.left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2) + 'px';
    tooltip.style.top = rect.top - tooltip.offsetHeight - 5 + 'px';
    
    e.target._tooltip = tooltip;
}

function hideTooltip(e) {
    if (e.target._tooltip) {
        e.target._tooltip.remove();
        delete e.target._tooltip;
    }
}

// Form validation helpers
function validateForm(form) {
    const requiredFields = form.querySelectorAll('[required]');
    let isValid = true;
    
    requiredFields.forEach(field => {
        if (!field.value.trim()) {
            field.classList.add('border-red-500');
            isValid = false;
        } else {
            field.classList.remove('border-red-500');
        }
    });
    
    return isValid;
}

// Auto-refresh functionality
function startAutoRefresh(selector, interval = 30000) {
    setInterval(() => {
        const element = document.querySelector(selector);
        if (element && element.hasAttribute('hx-get')) {
            htmx.trigger(element, 'refresh');
        }
    }, interval);
}

// Search/filter functionality
function filterTable(input, tableId) {
    const filter = input.value.toLowerCase();
    const table = document.getElementById(tableId);
    const rows = table.getElementsByTagName('tr');
    
    for (let i = 1; i < rows.length; i++) { // Skip header row
        const row = rows[i];
        const cells = row.getElementsByTagName('td');
        let found = false;
        
        for (let j = 0; j < cells.length; j++) {
            if (cells[j].textContent.toLowerCase().includes(filter)) {
                found = true;
                break;
            }
        }
        
        row.style.display = found ? '' : 'none';
    }
}

// Copy to clipboard functionality
function copyToClipboard(text) {
    if (navigator.clipboard) {
        navigator.clipboard.writeText(text).then(() => {
            showToast('Copied to clipboard', 'success');
        });
    } else {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        try {
            document.execCommand('copy');
            showToast('Copied to clipboard', 'success');
        } catch (err) {
            showToast('Failed to copy to clipboard', 'error');
        }
        document.body.removeChild(textArea);
    }
}

// Status color helpers
function getStatusClass(status) {
    switch (status.toLowerCase()) {
        case 'running':
        case 'active':
        case 'online':
            return 'status-running';
        case 'stopped':
        case 'failed':
        case 'terminated':
        case 'offline':
            return 'status-stopped';
        case 'pending':
        case 'installing':
            return 'status-pending';
        default:
            return 'status-unknown';
    }
}

// Export functionality (for future use)
function exportData(data, filename, type = 'json') {
    const blob = new Blob([JSON.stringify(data, null, 2)], { 
        type: type === 'json' ? 'application/json' : 'text/csv' 
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

// Keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl/Cmd + K for search (future feature)
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const searchInput = document.querySelector('[data-search]');
        if (searchInput) {
            searchInput.focus();
        }
    }
    
    // Ctrl/Cmd + N for new item (if on list pages)
    if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
        e.preventDefault();
        const addButton = document.querySelector('[data-add-button]');
        if (addButton) {
            addButton.click();
        }
    }
});

console.log('GARM Web Interface loaded');