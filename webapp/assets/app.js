// GARM SPA - Matching HTMX Layout Exactly
class GarmSPA {
    constructor() {
        this.baseUrl = window.location.origin;
        this.authToken = null;
        this.currentUser = null;
        this.currentRoute = 'dashboard';
        this.sidebarOpen = false;
        
        // Navigation items matching HTMX version exactly  
        this.navItems = [
            {
                name: 'Dashboard',
                route: 'dashboard',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z" />`
            },
            {
                name: 'Repositories',
                route: 'repositories',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z" />`
            },
            {
                name: 'Organizations',
                route: 'organizations',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />`
            },
            {
                name: 'Enterprises',
                route: 'enterprises',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />`
            },
            {
                name: 'Pools',
                route: 'pools',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />`
            },
            {
                name: 'Instances',
                route: 'instances',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />`
            },
            {
                name: 'Scale Sets',
                route: 'scalesets',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />`
            }
        ];

        this.configNavItems = [
            {
                name: 'Credentials',
                route: 'credentials',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />`
            },
            {
                name: 'Endpoints',
                route: 'endpoints',
                icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />`
            }
        ];
        
        this.init();
    }

    async init() {
        this.setupEventListeners();
        this.loadAuthState();
        this.initDarkMode();
        
        if (this.authToken) {
            try {
                await this.checkAuth();
                this.showMainApp();
                this.buildNavigation();
                this.navigate(this.currentRoute);
            } catch (error) {
                this.logout();
            }
        } else {
            this.showLogin();
        }
    }

    setupEventListeners() {
        // Login form
        document.getElementById('login-form').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleLogin();
        });

        // Theme toggles
        document.getElementById('theme-toggle-desktop').addEventListener('click', () => this.toggleDarkMode());
        document.getElementById('theme-toggle-mobile').addEventListener('click', () => this.toggleDarkMode());

        // Mobile menu
        document.getElementById('open-mobile-menu').addEventListener('click', () => this.openMobileMenu());
        document.getElementById('close-mobile-menu').addEventListener('click', () => this.closeMobileMenu());
        document.getElementById('overlay-backdrop').addEventListener('click', () => this.closeMobileMenu());

        // Handle browser back/forward
        window.addEventListener('popstate', (e) => {
            if (e.state && e.state.route) {
                this.navigate(e.state.route, false);
            }
        });
    }

    loadAuthState() {
        this.authToken = localStorage.getItem('garm_token');
        this.currentUser = localStorage.getItem('garm_user');
    }

    initDarkMode() {
        const isDarkMode = localStorage.getItem('darkMode') === 'true';
        
        if (isDarkMode) {
            document.documentElement.classList.add('dark');
            this.updateDarkModeIcons(true);
        } else {
            this.updateDarkModeIcons(false);
        }
    }

    updateDarkModeIcons(isDark) {
        const sunIconsDesktop = document.getElementById('sun-icon-desktop');
        const moonIconDesktop = document.getElementById('moon-icon-desktop');
        const sunIconsMobile = document.getElementById('sun-icon-mobile');
        const moonIconMobile = document.getElementById('moon-icon-mobile');
        
        if (isDark) {
            sunIconsDesktop.classList.remove('hidden');
            moonIconDesktop.classList.add('hidden');
            sunIconsMobile.classList.remove('hidden');
            moonIconMobile.classList.add('hidden');
        } else {
            sunIconsDesktop.classList.add('hidden');
            moonIconDesktop.classList.remove('hidden');
            sunIconsMobile.classList.add('hidden');
            moonIconMobile.classList.remove('hidden');
        }
    }

    toggleDarkMode() {
        const isDarkMode = document.documentElement.classList.contains('dark');
        
        if (isDarkMode) {
            document.documentElement.classList.remove('dark');
            localStorage.setItem('darkMode', 'false');
            this.updateDarkModeIcons(false);
        } else {
            document.documentElement.classList.add('dark');
            localStorage.setItem('darkMode', 'true');
            this.updateDarkModeIcons(true);
        }
    }

    openMobileMenu() {
        this.sidebarOpen = true;
        document.getElementById('mobile-overlay').classList.remove('hidden');
    }

    closeMobileMenu() {
        this.sidebarOpen = false;
        document.getElementById('mobile-overlay').classList.add('hidden');
    }

    buildNavigation() {
        const desktopNav = document.getElementById('desktop-nav');
        const mobileNav = document.getElementById('mobile-nav');
        
        // Build desktop navigation
        let desktopNavHTML = '';
        
        // Main navigation items
        this.navItems.forEach(item => {
            const isActive = this.currentRoute === item.route;
            desktopNavHTML += `
                <a href="#${item.route}" class="nav-link group flex items-center px-2 py-2 text-sm font-medium rounded-md ${
                    isActive 
                        ? 'bg-gray-100 text-gray-900 dark:bg-gray-700 dark:text-white' 
                        : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white'
                }" data-route="${item.route}">
                    <svg class="mr-3 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        ${item.icon}
                    </svg>
                    ${item.name}
                </a>
            `;
        });

        // Configuration section
        desktopNavHTML += `
            <div class="border-t border-gray-200 dark:border-gray-600 mt-4 pt-4">
        `;
        
        this.configNavItems.forEach(item => {
            const isActive = this.currentRoute === item.route;
            desktopNavHTML += `
                <a href="#${item.route}" class="nav-link group flex items-center px-2 py-2 text-sm font-medium rounded-md ${
                    isActive 
                        ? 'bg-gray-100 text-gray-900 dark:bg-gray-700 dark:text-white' 
                        : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white'
                }" data-route="${item.route}">
                    <svg class="mr-3 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        ${item.icon}
                    </svg>
                    ${item.name}
                </a>
            `;
        });

        desktopNavHTML += `
            </div>
            <div class="border-t border-gray-200 dark:border-gray-600 mt-4 pt-4">
                <button id="logout-button" class="group flex items-center px-2 py-2 text-sm font-medium rounded-md w-full text-left text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white">
                    <svg class="mr-3 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    Logout
                </button>
            </div>
        `;

        desktopNav.innerHTML = desktopNavHTML;

        // Mobile navigation (similar structure but with different spacing)
        const mobileHeader = `
            <div class="flex-shrink-0 flex items-center justify-between px-4">
                <div class="flex items-center">
                    <img src="/assets/garm-light.svg" alt="GARM Logo" class="h-20 w-20 dark:hidden">
                    <img src="/assets/garm-dark.svg" alt="GARM Logo" class="h-20 w-20 hidden dark:block">
                    <h1 class="ml-3 text-xl font-bold text-gray-900 dark:text-white">GARM</h1>
                </div>
                <button id="theme-toggle-mobile-sidebar" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors" title="Toggle theme">
                    <svg id="sun-icon-mobile-sidebar" class="w-5 h-5 text-black hover:text-gray-800 hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
                    </svg>
                    <svg id="moon-icon-mobile-sidebar" class="w-5 h-5 text-yellow-400 hover:text-yellow-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
                    </svg>
                </button>
            </div>
            <nav class="mt-5 px-2 space-y-1">
        `;

        let mobileNavHTML = mobileHeader;
        
        [...this.navItems, ...this.configNavItems].forEach(item => {
            const isActive = this.currentRoute === item.route;
            mobileNavHTML += `
                <a href="#${item.route}" class="nav-link group flex items-center px-2 py-2 text-base font-medium rounded-md ${
                    isActive 
                        ? 'bg-gray-100 text-gray-900 dark:bg-gray-700 dark:text-white' 
                        : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white'
                }" data-route="${item.route}">
                    <svg class="mr-4 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        ${item.icon}
                    </svg>
                    ${item.name}
                </a>
            `;
        });

        mobileNavHTML += `
                <button id="logout-button-mobile" class="group flex items-center px-2 py-2 text-base font-medium rounded-md w-full text-left text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-white">
                    <svg class="mr-4 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    Logout
                </button>
            </nav>
        `;

        mobileNav.innerHTML = mobileNavHTML;

        // Attach navigation event listeners
        document.querySelectorAll('.nav-link').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const route = e.currentTarget.getAttribute('data-route');
                this.navigate(route);
                this.closeMobileMenu();
            });
        });

        // Attach logout event listeners
        document.getElementById('logout-button').addEventListener('click', () => this.logout());
        const logoutMobile = document.getElementById('logout-button-mobile');
        if (logoutMobile) {
            logoutMobile.addEventListener('click', () => this.logout());
        }

        // Mobile sidebar theme toggle
        const mobileSidebarToggle = document.getElementById('theme-toggle-mobile-sidebar');
        if (mobileSidebarToggle) {
            mobileSidebarToggle.addEventListener('click', () => this.toggleDarkMode());
        }
    }

    async handleLogin() {
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const button = document.getElementById('login-button');
        const errorDiv = document.getElementById('login-error');
        const errorMessage = document.getElementById('login-error-message');

        button.disabled = true;
        button.textContent = 'Signing in...';
        errorDiv.classList.add('hidden');

        try {
            const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                const data = await response.json();
                this.authToken = data.token;
                this.currentUser = username;
                
                localStorage.setItem('garm_token', this.authToken);
                localStorage.setItem('garm_user', username);
                
                this.showMainApp();
                this.buildNavigation();
                this.navigate('dashboard');
            } else {
                const errorData = await response.json();
                errorMessage.textContent = errorData.error || 'Login failed';
                errorDiv.classList.remove('hidden');
            }
        } catch (error) {
            errorMessage.textContent = 'Network error. Please try again.';
            errorDiv.classList.remove('hidden');
        } finally {
            button.disabled = false;
            button.textContent = 'Sign in';
        }
    }

    async checkAuth() {
        const response = await fetch(`${this.baseUrl}/api/v1/controller-info`, {
            headers: {
                'Authorization': `Bearer ${this.authToken}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Authentication failed');
        }
    }

    logout() {
        this.authToken = null;
        this.currentUser = null;
        localStorage.removeItem('garm_token');
        localStorage.removeItem('garm_user');
        this.showLogin();
    }

    showLogin() {
        document.getElementById('loading').classList.add('hidden');
        document.getElementById('main-app').classList.add('hidden');
        document.getElementById('login-page').classList.remove('hidden');
    }

    showMainApp() {
        document.getElementById('loading').classList.add('hidden');
        document.getElementById('login-page').classList.add('hidden');
        document.getElementById('main-app').classList.remove('hidden');
    }

    navigate(route, pushState = true) {
        this.currentRoute = route;
        
        if (pushState) {
            history.pushState({ route }, '', `#${route}`);
        }

        // Update page title
        const pageTitles = {
            'dashboard': 'Dashboard',
            'repositories': 'Repositories',
            'organizations': 'Organizations',
            'enterprises': 'Enterprises',
            'pools': 'Pools',
            'instances': 'Instances',
            'scalesets': 'Scale Sets',
            'credentials': 'Credentials',
            'endpoints': 'Endpoints'
        };
        
        document.getElementById('page-title').textContent = pageTitles[route] || 'Dashboard';

        // Rebuild navigation to update active states
        this.buildNavigation();

        // Load page content
        this.loadPageContent(route);
    }

    async loadPageContent(route) {
        const content = document.getElementById('content');
        content.innerHTML = '<div class="text-center py-8"><div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div></div>';

        try {
            switch (route) {
                case 'dashboard':
                    await this.loadDashboard();
                    break;
                case 'repositories':
                    await this.loadRepositories();
                    break;
                case 'organizations':
                    await this.loadOrganizations();
                    break;
                case 'enterprises':
                    await this.loadEnterprises();
                    break;
                case 'pools':
                    await this.loadPools();
                    break;
                case 'instances':
                    await this.loadInstances();
                    break;
                case 'scalesets':
                    await this.loadScaleSets();
                    break;
                case 'credentials':
                    await this.loadCredentials();
                    break;
                case 'endpoints':
                    await this.loadEndpoints();
                    break;
                default:
                    content.innerHTML = '<div class="text-center py-8"><h1 class="text-2xl font-bold text-gray-900 dark:text-white">Page Not Found</h1></div>';
            }
        } catch (error) {
            content.innerHTML = `
                <div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
                    <div class="flex">
                        <div class="flex-shrink-0">
                            <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                            </svg>
                        </div>
                        <div class="ml-3">
                            <h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading page</h3>
                            <p class="mt-2 text-sm text-red-700 dark:text-red-300">${error.message}</p>
                        </div>
                    </div>
                </div>
            `;
        }
    }

    async apiRequest(endpoint, options = {}) {
        const url = `${this.baseUrl}/api/v1${endpoint}`;
        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${this.authToken}`,
                ...options.headers
            }
        });

        if (!response.ok) {
            throw new Error(`API request failed: ${response.statusText}`);
        }

        return response.json();
    }

    async loadDashboard() {
        const [repos, orgs, pools, instances] = await Promise.all([
            this.apiRequest('/repositories').catch(() => []),
            this.apiRequest('/organizations').catch(() => []),
            this.apiRequest('/pools').catch(() => []),
            this.apiRequest('/instances').catch(() => [])
        ]);

        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
                <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 rounded-md bg-blue-500 text-white flex items-center justify-center">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z" />
                                    </svg>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Repositories</dt>
                                    <dd class="text-lg font-medium text-gray-900 dark:text-white">${repos.length}</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 rounded-md bg-green-500 text-white flex items-center justify-center">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                                    </svg>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Organizations</dt>
                                    <dd class="text-lg font-medium text-gray-900 dark:text-white">${orgs.length}</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 rounded-md bg-purple-500 text-white flex items-center justify-center">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                    </svg>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Pools</dt>
                                    <dd class="text-lg font-medium text-gray-900 dark:text-white">${pools.length}</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 rounded-md bg-yellow-500 text-white flex items-center justify-center">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
                                    </svg>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Instances</dt>
                                    <dd class="text-lg font-medium text-gray-900 dark:text-white">${instances.length}</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="mt-8 bg-white dark:bg-gray-800 shadow rounded-lg">
                <div class="p-6">
                    <h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">GARM SPA vs HTMX Comparison</h3>
                    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                        Side-by-side comparison of implementation approaches
                    </p>
                    
                    <div class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2">
                        <div class="border border-gray-300 dark:border-gray-600 rounded-lg p-4">
                            <h4 class="text-sm font-medium text-gray-900 dark:text-white">HTMX Version Features</h4>
                            <ul class="mt-2 text-sm text-gray-600 dark:text-gray-400 space-y-1">
                                <li>✅ Server-rendered templates</li>
                                <li>✅ Smaller JavaScript bundle</li>
                                <li>✅ Template-based development</li>
                                <li>✅ Progressive enhancement</li>
                                <li><a href="/web/" class="text-blue-600 hover:text-blue-800">→ View HTMX Version</a></li>
                            </ul>
                        </div>
                        
                        <div class="border border-gray-300 dark:border-gray-600 rounded-lg p-4">
                            <h4 class="text-sm font-medium text-gray-900 dark:text-white">SPA Version Features</h4>
                            <ul class="mt-2 text-sm text-gray-600 dark:text-gray-400 space-y-1">
                                <li>✅ Client-side routing</li>
                                <li>✅ Instant navigation</li>
                                <li>✅ No npm dependencies</li>
                                <li>✅ Embedded in Go binary</li>
                                <li><span class="text-green-600">✅ Current Version</span></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    async loadRepositories() {
        const repos = await this.apiRequest('/repositories');
        const content = document.getElementById('content');
        
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
                ${repos.length === 0 ? `
                    <div class="p-6 text-center">
                        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z" />
                        </svg>
                        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">No repositories found</p>
                    </div>
                ` : `
                    <div class="overflow-x-auto">
                        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                            <thead class="bg-gray-50 dark:bg-gray-700">
                                <tr>
                                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
                                    <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Endpoint</th>
                                    <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Credential</th>
                                    <th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
                                    <th class="px-3 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Actions</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                                ${repos.map(repo => `
                                    <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div>
                                                <a href="#" class="text-sm font-medium text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 hover:underline">${repo.owner}/${repo.name}</a>
                                            </div>
                                        </td>
                                        <td class="px-3 py-4 whitespace-nowrap">
                                            <div class="flex items-center space-x-2">
                                                <div class="flex items-center text-gray-600 dark:text-gray-400" title="${this.getForgeType(repo.endpoint)}">
                                                    ${this.getRepositoryForgeIcon(repo.endpoint)}
                                                </div>
                                                <span class="text-sm text-gray-500 dark:text-gray-400">${repo.endpoint?.name || ''}</span>
                                            </div>
                                        </td>
                                        <td class="px-3 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">${repo.credentials_name}</td>
                                        <td class="px-3 py-4 whitespace-nowrap">
                                            <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full ${repo.pool_manager_status?.is_running ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'}">
                                                ${repo.pool_manager_status?.is_running ? 'Running' : 'Stopped'}
                                            </span>
                                        </td>
                                        <td class="px-3 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div class="flex justify-end space-x-2">
                                                <button class="text-green-600 dark:text-green-400 hover:text-green-900 dark:hover:text-green-300">Edit</button>
                                                <button class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300">Delete</button>
                                            </div>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                `}
            </div>
        `;
    }

    async loadOrganizations() {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
                <p class="text-gray-600 dark:text-gray-400">Organizations page content would be implemented here using the same API endpoints as the HTMX version.</p>
            </div>
        `;
    }

    async loadEnterprises() {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
                <p class="text-gray-600 dark:text-gray-400">Enterprises page content would be implemented here using the same API endpoints as the HTMX version.</p>
            </div>
        `;
    }

    async loadPools() {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
                <p class="text-gray-600 dark:text-gray-400">Pools page content would be implemented here using the same API endpoints as the HTMX version.</p>
            </div>
        `;
    }

    async loadInstances() {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
                <p class="text-gray-600 dark:text-gray-400">Instances page content would be implemented here using the same API endpoints as the HTMX version.</p>
            </div>
        `;
    }

    async loadScaleSets() {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
                <p class="text-gray-600 dark:text-gray-400">Scale Sets page content would be implemented here using the same API endpoints as the HTMX version.</p>
            </div>
        `;
    }

    async loadCredentials() {
        try {
            // Fetch credentials and endpoints from both GitHub and Gitea
            const [githubCredentials, giteaCredentials, githubEndpoints, giteaEndpoints] = await Promise.all([
                this.apiRequest('/github/credentials').catch(() => []),
                this.apiRequest('/gitea/credentials').catch(() => []),
                this.apiRequest('/github/endpoints').catch(() => []),
                this.apiRequest('/gitea/endpoints').catch(() => [])
            ]);
            
            // Combine credentials and endpoints with forge type information
            const allCredentials = [
                ...githubCredentials.map(cred => ({ ...cred, endpoint_type: 'github' })),
                ...giteaCredentials.map(cred => ({ ...cred, endpoint_type: 'gitea' }))
            ];
            
            const allEndpoints = [
                ...githubEndpoints.map(ep => ({ ...ep, endpoint_type: 'github' })),
                ...giteaEndpoints.map(ep => ({ ...ep, endpoint_type: 'gitea' }))
            ];
            
            this.renderCredentials(allCredentials, allEndpoints);
        } catch (error) {
            console.error('Error loading credentials:', error);
            this.renderCredentialsError(error.message);
        }
    }

    async loadEndpoints() {
        try {
            // Fetch both GitHub and Gitea endpoints
            const [githubEndpoints, giteaEndpoints] = await Promise.all([
                this.apiRequest('/github/endpoints').catch(() => []),
                this.apiRequest('/gitea/endpoints').catch(() => [])
            ]);
            
            // Combine endpoints with forge type information
            const allEndpoints = [
                ...githubEndpoints.map(ep => ({ ...ep, endpoint_type: 'github' })),
                ...giteaEndpoints.map(ep => ({ ...ep, endpoint_type: 'gitea' }))
            ];
            
            this.renderEndpoints(allEndpoints);
        } catch (error) {
            console.error('Error loading endpoints:', error);
            this.renderEndpointsError(error.message);
        }
    }

    renderEndpoints(endpoints) {
        const content = document.getElementById('content');
        
        content.innerHTML = `
            <div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
                <div class="p-6 border-b border-gray-200 dark:border-gray-700">
                    <div class="flex items-center justify-between">
                        <h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">Endpoints</h3>
                        <button id="create-endpoint-btn" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                            <svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                            </svg>
                            Create Endpoint
                        </button>
                    </div>
                </div>
                
                ${endpoints.length === 0 ? `
                    <div class="p-6 text-center">
                        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">No endpoints found</p>
                        <p class="text-xs text-gray-500 dark:text-gray-500">Create your first endpoint to get started</p>
                    </div>
                ` : `
                    <div class="overflow-x-auto">
                        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                            <thead class="bg-gray-50 dark:bg-gray-700">
                                <tr>
                                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
                                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Description</th>
                                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">API URL</th>
                                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Type</th>
                                    <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Actions</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                                ${endpoints.map(endpoint => `
                                    <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="text-sm font-medium text-gray-900 dark:text-white">${endpoint.name}</div>
                                        </td>
                                        <td class="px-6 py-4">
                                            <div class="text-sm text-gray-900 dark:text-white">${endpoint.description || '-'}</div>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="text-sm text-gray-900 dark:text-white">${endpoint.api_base_url}</div>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="flex items-center">
                                                ${this.getForgeIcon(endpoint.endpoint_type)}
                                                <span class="text-sm text-gray-900 dark:text-white">${endpoint.endpoint_type === 'github' ? 'GitHub' : 'Gitea'}</span>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <button onclick="garmApp.editEndpoint('${endpoint.name}', '${endpoint.endpoint_type}')" class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 mr-3">
                                                Edit
                                            </button>
                                            <button onclick="garmApp.deleteEndpoint('${endpoint.name}', '${endpoint.endpoint_type}')" class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300">
                                                Delete
                                            </button>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                `}
            </div>
        `;
        
        // Attach event listeners
        document.getElementById('create-endpoint-btn')?.addEventListener('click', () => this.showCreateEndpointModal());
    }

    getForgeIcon(forgeType) {
        if (forgeType === 'github') {
            return `<svg class="w-5 h-5 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-5 h-5 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg>`;
        } else {
            return `<svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
        }
    }

    getRepositoryForgeIcon(endpoint) {
        if (!endpoint || !endpoint.endpoint_type) {
            return '';
        }
        
        if (endpoint.endpoint_type === 'gitea') {
            return `<svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
        } else {
            // GitHub
            return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
        }
    }

    getForgeType(endpoint) {
        if (!endpoint) return '';
        return endpoint.endpoint_type === 'gitea' ? 'Gitea' : 'GitHub';
    }

    renderEndpointsError(message) {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
                <div class="flex">
                    <div class="flex-shrink-0">
                        <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                        </svg>
                    </div>
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading endpoints</h3>
                        <p class="mt-2 text-sm text-red-700 dark:text-red-300">${message}</p>
                    </div>
                </div>
            </div>
        `;
    }

    renderCredentials(credentials, endpoints) {
        const content = document.getElementById('content');
        
        content.innerHTML = `
            <div class="space-y-6">
                <!-- Header -->
                <div class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Credentials</h1>
                        <p class="mt-2 text-sm text-gray-700 dark:text-gray-300">
                            Manage authentication credentials for your GitHub and Gitea endpoints.
                        </p>
                    </div>
                    <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
                        <button
                            id="create-credentials-btn"
                            type="button"
                            class="block rounded-md bg-blue-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                        >
                            Add Credentials
                        </button>
                    </div>
                </div>

                <!-- Table -->
                <div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
                    ${credentials.length === 0 ? `
                        <div class="px-4 py-8 text-center">
                            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
                            </svg>
                            <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No credentials</h3>
                            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Get started by creating new credentials.</p>
                        </div>
                    ` : `
                        <div class="overflow-x-auto">
                            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                                <thead class="bg-gray-50 dark:bg-gray-700">
                                    <tr>
                                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
                                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Description</th>
                                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Endpoint</th>
                                        <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Auth Type</th>
                                        <th scope="col" class="relative px-6 py-3"><span class="sr-only">Actions</span></th>
                                    </tr>
                                </thead>
                                <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                                    ${credentials.map(credential => `
                                        <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                                            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                                                ${credential.name}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                                                ${credential.description}
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                                                <div class="flex items-center">
                                                    ${this.getForgeIcon(endpoints.find(e => e.name === credential.endpoint_name)?.endpoint_type || '')}
                                                    <span class="ml-2">${credential.endpoint_name}</span>
                                                </div>
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                                                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full ${credential.auth_type === 'pat' ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'}">
                                                    ${credential.auth_type === 'pat' ? 'PAT' : 'App'}
                                                </span>
                                            </td>
                                            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                <div class="flex justify-end space-x-2">
                                                    <button
                                                        onclick="garmApp.editCredentials('${credential.name}')"
                                                        class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                                                    >
                                                        Edit
                                                    </button>
                                                    <button
                                                        onclick="garmApp.deleteCredentials('${credential.name}')"
                                                        class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                                                    >
                                                        Delete
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        </div>
                    `}
                </div>
            </div>
        `;

        // Add event listeners
        document.getElementById('create-credentials-btn')?.addEventListener('click', () => this.showCreateCredentialsModal());
    }

    renderCredentialsError(message) {
        const content = document.getElementById('content');
        content.innerHTML = `
            <div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
                <div class="flex">
                    <div class="flex-shrink-0">
                        <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                        </svg>
                    </div>
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading credentials</h3>
                        <p class="mt-2 text-sm text-red-700 dark:text-red-300">${message}</p>
                    </div>
                </div>
            </div>
        `;
    }

    showCreateEndpointModal() {
        // Implementation will be added here
        alert('Create endpoint modal - to be implemented');
    }

    async editEndpoint(name, type) {
        // Implementation will be added here
        alert(`Edit ${type} endpoint: ${name} - to be implemented`);
    }

    async deleteEndpoint(name, endpointType) {
        if (!confirm(`Are you sure you want to delete the endpoint "${name}"? This action cannot be undone.`)) {
            return;
        }
        
        try {
            const endpoint = endpointType === 'github' ? 'github' : 'gitea';
            await this.apiRequest(`/${endpoint}/endpoints/${name}`, {
                method: 'DELETE'
            });
            
            this.loadEndpoints(); // Refresh the list
        } catch (error) {
            alert(`Error deleting endpoint: ${error.message}`);
        }
    }

    showCreateCredentialsModal() {
        // Implementation will be added here
        alert('Create credentials modal - to be implemented');
    }

    async editCredentials(name) {
        // Implementation will be added here
        alert(`Edit credentials: ${name} - to be implemented`);
    }

    async deleteCredentials(name) {
        if (!confirm(`Are you sure you want to delete the credentials "${name}"? This action cannot be undone.`)) {
            return;
        }
        
        try {
            await this.apiRequest(`/credentials/${name}`, {
                method: 'DELETE'
            });
            
            this.loadCredentials(); // Refresh the list
        } catch (error) {
            alert(`Error deleting credentials: ${error.message}`);
        }
    }
}

// Initialize the application when DOM is ready
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new GarmSPA();
});