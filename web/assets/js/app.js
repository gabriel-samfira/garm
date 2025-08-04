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
            if (trigger.tagName === 'FORM') {
                showToast('Operation completed successfully', 'success');
                closeModal();
            }
        }
    });
    
    // Close modal when clicking outside
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('modal-backdrop')) {
            closeModal();
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
    
    // Handle custom HTMX triggers
    document.body.addEventListener('repositoryUpdated', function(evt) {
        closeModal();
        showToast('Repository updated successfully', 'success');
        // Manually refresh the repositories table to preserve polling
        const repoTable = document.getElementById('repositories-table');
        if (repoTable) {
            htmx.trigger(repoTable, 'refresh');
        }
    });
    
    document.body.addEventListener('repositoryCreated', function(evt) {
        closeModal();
        showToast('Repository created successfully', 'success');
        // Manually refresh the repositories table to preserve polling
        const repoTable = document.getElementById('repositories-table');
        if (repoTable) {
            htmx.trigger(repoTable, 'refresh');
        }
    });
    
    document.body.addEventListener('repositoryDeleted', function(evt) {
        showToast('Repository deleted successfully', 'success');
        // Manually refresh the repositories table to preserve polling
        const repoTable = document.getElementById('repositories-table');
        if (repoTable) {
            htmx.trigger(repoTable, 'refresh');
        }
    });
    
    document.body.addEventListener('repositoryCreateError', function(evt) {
        showToast('Failed to create repository', 'error');
    });
    
    document.body.addEventListener('repositoryUpdateError', function(evt) {
        showToast('Failed to update repository', 'error');
    });
    
    document.body.addEventListener('repositoryDeleteError', function(evt) {
        showToast('Failed to delete repository', 'error');
    });
    
    document.body.addEventListener('organizationUpdated', function(evt) {
        closeModal();
        showToast('Organization updated successfully', 'success');
        // Manually refresh the organizations table to preserve polling
        const orgTable = document.getElementById('organizations-table');
        if (orgTable) {
            htmx.trigger(orgTable, 'refresh');
        }
    });
    
    document.body.addEventListener('organizationCreated', function(evt) {
        closeModal();
        showToast('Organization created successfully', 'success');
        // Manually refresh the organizations table to preserve polling
        const orgTable = document.getElementById('organizations-table');
        if (orgTable) {
            htmx.trigger(orgTable, 'refresh');
        }
    });
    
    document.body.addEventListener('organizationDeleted', function(evt) {
        showToast('Organization deleted successfully', 'success');
        // Manually refresh the organizations table to preserve polling
        const orgTable = document.getElementById('organizations-table');
        if (orgTable) {
            htmx.trigger(orgTable, 'refresh');
        }
    });
    
    document.body.addEventListener('organizationCreateError', function(evt) {
        showToast('Failed to create organization', 'error');
    });
    
    document.body.addEventListener('organizationUpdateError', function(evt) {
        showToast('Failed to update organization', 'error');
    });
    
    document.body.addEventListener('organizationDeleteError', function(evt) {
        showToast('Failed to delete organization', 'error');
    });
    
    document.body.addEventListener('enterpriseUpdated', function(evt) {
        closeModal();
        showToast('Enterprise updated successfully', 'success');
        // Manually refresh the enterprises table to preserve polling
        const entTable = document.getElementById('enterprises-table');
        if (entTable) {
            htmx.trigger(entTable, 'refresh');
        }
    });
    
    document.body.addEventListener('enterpriseCreated', function(evt) {
        closeModal();
        showToast('Enterprise created successfully', 'success');
        // Manually refresh the enterprises table to preserve polling
        const entTable = document.getElementById('enterprises-table');
        if (entTable) {
            htmx.trigger(entTable, 'refresh');
        }
    });
    
    document.body.addEventListener('enterpriseDeleted', function(evt) {
        showToast('Enterprise deleted successfully', 'success');
        // Manually refresh the enterprises table to preserve polling
        const entTable = document.getElementById('enterprises-table');
        if (entTable) {
            htmx.trigger(entTable, 'refresh');
        }
    });
    
    document.body.addEventListener('enterpriseCreateError', function(evt) {
        showToast('Failed to create enterprise', 'error');
    });
    
    document.body.addEventListener('enterpriseUpdateError', function(evt) {
        showToast('Failed to update enterprise', 'error');
    });
    
    document.body.addEventListener('enterpriseDeleteError', function(evt) {
        showToast('Failed to delete enterprise', 'error');
    });
    
    document.body.addEventListener('poolUpdated', function(evt) {
        closeModal();
        showToast('Pool updated successfully', 'success');
        // Manually refresh the pools table to preserve polling
        const poolTable = document.getElementById('pools-table');
        if (poolTable) {
            htmx.trigger(poolTable, 'refresh');
        }
    });
    
    document.body.addEventListener('poolUpdateError', function(evt) {
        showToast('Failed to update pool', 'error');
    });
});

// Toast notification system using existing toast elements
function showToast(message, type = 'info', duration = 5000) {
    let toastEl, messageEl;
    
    if (type === 'success') {
        toastEl = document.getElementById('success-toast');
        messageEl = document.getElementById('success-message');
    } else if (type === 'error') {
        toastEl = document.getElementById('error-toast');
        messageEl = document.getElementById('error-message');
    }
    
    if (toastEl && messageEl) {
        messageEl.textContent = message;
        toastEl.classList.remove('hidden');
        
        // Auto-hide after duration
        setTimeout(() => {
            toastEl.classList.add('hidden');
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

// Logout functionality
function logout() {
    console.log('Logout initiated');
    
    // Clear token from local storage (if using local storage authentication)
    localStorage.removeItem('garm_token');
    
    // Also clear any other auth-related items from storage
    localStorage.removeItem('garm_user');
    sessionStorage.removeItem('garm_token');
    sessionStorage.removeItem('garm_user');
    
    // Make a request to the logout endpoint to clear the server-side cookie
    fetch('/web/api/auth/logout', {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        console.log('Logout response:', response.status);
        // Always redirect to login regardless of response
        window.location.href = '/web/login';
    })
    .catch(error => {
        console.error('Logout error:', error);
        // Even if logout request fails, redirect to login since we cleared local storage
        window.location.href = '/web/login';
    });
}

console.log('GARM Web Interface loaded');

// Delete confirmation functions for unified entity template
function showDeleteRepositoryConfirm(id, name) {
    if (confirm(`Are you sure you want to delete the repository "${name}"? This action cannot be undone.`)) {
        // Use HTMX to make the DELETE request
        htmx.ajax('DELETE', `/web/api/repositories/${id}`, {
            target: 'body',
            swap: 'none',
            headers: {
                'HX-Request': 'true'
            }
        }).then(() => {
            showToast('Repository deleted successfully', 'success');
            // Redirect to repositories list
            window.location.href = '/web/repositories';
        }).catch(() => {
            showToast('Failed to delete repository', 'error');
        });
    }
}

function showDeleteOrganizationConfirm(id, name) {
    if (confirm(`Are you sure you want to delete the organization "${name}"? This action cannot be undone.`)) {
        // Use HTMX to make the DELETE request
        htmx.ajax('DELETE', `/web/api/organizations/${id}`, {
            target: 'body',
            swap: 'none',
            headers: {
                'HX-Request': 'true'
            }
        }).then(() => {
            showToast('Organization deleted successfully', 'success');
            // Redirect to organizations list
            window.location.href = '/web/organizations';
        }).catch(() => {
            showToast('Failed to delete organization', 'error');
        });
    }
}

function showDeleteEnterpriseConfirm(id, name) {
    if (confirm(`Are you sure you want to delete the enterprise "${name}"? This action cannot be undone.`)) {
        // Use HTMX to make the DELETE request
        htmx.ajax('DELETE', `/web/api/enterprises/${id}`, {
            target: 'body',
            swap: 'none',
            headers: {
                'HX-Request': 'true'
            }
        }).then(() => {
            showToast('Enterprise deleted successfully', 'success');
            // Redirect to enterprises list
            window.location.href = '/web/enterprises';
        }).catch(() => {
            showToast('Failed to delete enterprise', 'error');
        });
    }
}