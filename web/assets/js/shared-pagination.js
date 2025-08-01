// Shared pagination and search functionality for entity lists
function createEntityPagination(config) {
    const {
        entityType,         // 'repositories', 'organizations', 'enterprises', 'pools'
        tableId,           // ID of the table body element
        searchInputId,     // ID of the search input
        apiEndpoint,       // API endpoint for the entity
        searchPlaceholder  // Placeholder text for search input
    } = config;
    
    let currentPage = 1;
    let hasNextPage = false;

    function changePage(direction) {
        const newPage = currentPage + direction;
        if (newPage < 1) return;
        
        currentPage = newPage;
        document.getElementById('current-page').value = currentPage;
        document.getElementById('page-numbers').textContent = `Page ${currentPage}`;
        
        // Make the HTMX request
        htmx.trigger(`#${tableId}`, 'refresh');
        
        updatePaginationButtons();
    }

    function updatePaginationButtons() {
        const prevBtn = document.getElementById('prev-page');
        const nextBtn = document.getElementById('next-page');
        
        if (prevBtn) prevBtn.disabled = currentPage <= 1;
        if (nextBtn) nextBtn.disabled = !hasNextPage;
    }

    // Reset to page 1 when searching
    function setupSearchHandler() {
        const searchInput = document.getElementById(searchInputId);
        if (searchInput) {
            searchInput.addEventListener('input', function() {
                currentPage = 1;
                document.getElementById('current-page').value = 1;
                document.getElementById('page-numbers').textContent = 'Page 1';
            });
        }
    }

    // Listen for HTMX responses to update pagination info
    function setupPaginationListener() {
        document.addEventListener('htmx:afterRequest', function(evt) {
            if (evt.detail.target && evt.detail.target.id === tableId) {
                // Check if we got results to determine if there might be a next page
                const tableBody = document.getElementById(tableId);
                const rows = tableBody.querySelectorAll('tr:not([data-has-more])');
                const hasMoreMarker = tableBody.querySelector('tr[data-has-more]');
                const perPage = parseInt(document.querySelector('[name="per_page"]').value) || 25;
                
                // Check if there are more pages based on the hidden marker
                hasNextPage = hasMoreMarker !== null;
                updatePaginationButtons();
                
                // Update pagination info
                const info = document.getElementById('pagination-info');
                if (info && rows.length > 0 && !rows[0].textContent.includes(`No ${entityType} found`)) {
                    const start = (currentPage - 1) * perPage + 1;
                    const end = start + rows.length - 1;
                    info.textContent = `Showing ${start}-${end} ${entityType}`;
                } else if (info) {
                    info.textContent = `No ${entityType} found`;
                }
            }
        });
    }

    // Setup pagination controls
    function setupPaginationControls() {
        // Add event listeners to pagination buttons
        document.addEventListener('click', function(evt) {
            if (evt.target.id === 'prev-page') {
                evt.preventDefault();
                changePage(-1);
            } else if (evt.target.id === 'next-page') {
                evt.preventDefault();
                changePage(1);
            }
        });
    }

    // Initialize the pagination system
    function init() {
        setupSearchHandler();
        setupPaginationListener();
        setupPaginationControls();
        
        // Update the tbody to trigger on refresh
        const tableBody = document.getElementById(tableId);
        if (tableBody) {
            tableBody.setAttribute('hx-trigger', 'load, every 30s, refresh');
        }
    }

    // Public API
    return {
        init: init,
        changePage: changePage,
        getCurrentPage: () => currentPage,
        hasNextPage: () => hasNextPage
    };
}

// Global function for inline onclick handlers
function changePage(direction) {
    if (window.entityPagination) {
        window.entityPagination.changePage(direction);
    }
}