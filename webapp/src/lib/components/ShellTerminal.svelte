<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { createShellConnection, type ShellConnection } from '$lib/utils/shell';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import '@xterm/xterm/css/xterm.css';

  export let runnerName: string;
  export let onClose: () => void;

  interface TabSession {
    id: string;
    title: string;
    connection: ShellConnection | null;
    terminal: Terminal | null;
    fitAddon: FitAddon | null;
    isInitialized: boolean;
    isConnecting: boolean;
    isConnected: boolean;
    error: string;
  }

  let tabs: TabSession[] = [];
  let activeTabId: string = '';
  let tabCounter = 1;

  let terminalContainer: HTMLDivElement;
  let isMaximized = false;
  let isResizing = false;
  let isDragging = false;
  let initialMouseX = 0;
  let initialMouseY = 0;
  let initialWidth = 0;
  let initialHeight = 0;
  let initialLeft = 0;
  let initialTop = 0;
  let windowedWidth = 800;
  let windowedHeight = 500;
  let windowedLeft = 0;
  let windowedTop = 0;
  let hasRestoredState = false;

  // Computed values for active tab
  $: activeTab = tabs.find(tab => tab.id === activeTabId);
  $: connection = activeTab?.connection || null;
  $: terminal = activeTab?.terminal || null;
  $: isConnecting = activeTab?.isConnecting || false;
  $: isConnected = activeTab?.isConnected || false;
  $: error = activeTab?.error || '';

  function createNewTab(): string {
    const id = `tab-${Date.now()}-${tabCounter}`;
    const newTab: TabSession = {
      id,
      title: `Shell ${tabCounter}`,
      connection: null,
      terminal: null,
      fitAddon: null,
      isInitialized: false,
      isConnecting: true,
      isConnected: false,
      error: ''
    };
    
    tabs = [...tabs, newTab];
    tabCounter++;
    return id;
  }

  function switchToTab(tabId: string) {
    if (activeTabId !== tabId) {
      console.log(`switchToTab: Switching from ${activeTabId} to ${tabId}`);
      activeTabId = tabId;
      
      // Only focus and fit - NEVER reinitialize
      setTimeout(() => {
        const tab = tabs.find(t => t.id === tabId);
        console.log(`switchToTab: Tab ${tabId} - isInitialized=${tab?.isInitialized}, hasTerminal=${!!tab?.terminal}, isConnected=${tab?.isConnected}`);
        if (tab?.terminal && tab.isInitialized && tab.isConnected && tab.fitAddon) {
          console.log(`switchToTab: Fitting terminal for tab ${tabId} (now visible)`);
          // Fit the terminal now that it's visible and can calculate dimensions correctly
          fitTerminalForTab(tab);
          tab.terminal.focus();
        }
      }, 10);
    }
  }

  function closeTab(tabId: string) {
    const tabIndex = tabs.findIndex(tab => tab.id === tabId);
    if (tabIndex === -1) return;

    const tab = tabs[tabIndex];
    
    // Send close message and dispose resources
    cleanupTab(tab);

    // Remove tab from array
    tabs = tabs.filter(t => t.id !== tabId);

    // Switch to another tab if this was active
    if (activeTabId === tabId) {
      if (tabs.length > 0) {
        // Switch to the previous tab, or first tab if we closed the first one
        const newActiveIndex = Math.max(0, Math.min(tabIndex - 1, tabs.length - 1));
        activeTabId = tabs[newActiveIndex]?.id || '';
      } else {
        // No tabs left, close the terminal
        onClose();
      }
    }
  }

  function cleanupTab(tab: TabSession) {
    // Close connection (this should send the close shell message)
    if (tab.connection) {
      tab.connection.close();
    }
    // Dispose terminal
    if (tab.terminal) {
      tab.terminal.dispose();
    }
  }

  async function createConnection(tabId: string) {
    const tabIndex = tabs.findIndex(tab => tab.id === tabId);
    if (tabIndex === -1) return;

    try {
      const newConnection = await createShellConnection(
        runnerName,
        (data: Uint8Array) => handleData(tabId, data),
        () => handleReady(tabId),
        () => handleExit(tabId),
        (errorMsg: string) => handleError(tabId, errorMsg)
      );
      
      tabs[tabIndex].connection = newConnection;
      tabs = [...tabs];
    } catch (err) {
      tabs[tabIndex].error = err instanceof Error ? err.message : 'Failed to connect';
      tabs[tabIndex].isConnecting = false;
      tabs = [...tabs];
    }
  }

  function initializeTerminal(tabId: string) {
    const tabIndex = tabs.findIndex(t => t.id === tabId);
    if (tabIndex === -1) return;
    
    const tab = tabs[tabIndex];
    
    // NEVER reinitialize a terminal that's already been initialized
    if (!tab.terminal || tab.isInitialized) {
      console.log(`initializeTerminal: Skipping tab ${tabId}, isInitialized=${tab.isInitialized}, hasTerminal=${!!tab.terminal}`);
      return;
    }
    
    console.log(`initializeTerminal: Initializing tab ${tabId} for the first time`);

    // Create FitAddon for this tab if it doesn't exist
    if (!tab.fitAddon) {
      tab.fitAddon = new FitAddon();
      tab.terminal.loadAddon(tab.fitAddon);
    }

    // Find the terminal element for this tab
    const terminalElement = document.querySelector(`.terminal-tab[data-tab-id="${tabId}"]`) as HTMLDivElement;
    if (!terminalElement) return;

    // Open terminal in its dedicated element - THIS SHOULD ONLY HAPPEN ONCE
    tab.terminal.open(terminalElement);
    
    if (tab.fitAddon) {
      tab.fitAddon.fit();
    }

    // Handle terminal input
    tab.terminal.onData((data) => {
      if (tab.connection && tab.isConnected) {
        const encoder = new TextEncoder();
        tab.connection.sendData(encoder.encode(data));
      }
    });

    // Handle terminal resize
    tab.terminal.onResize(({ cols, rows }) => {
      if (tab.connection && tab.isConnected) {
        tab.connection.resize(cols, rows);
      }
    });

    // Mark as initialized so this never happens again
    tab.isInitialized = true;

    // Update tabs array
    tabs[tabIndex] = tab;
    tabs = [...tabs];

    // Focus the terminal if it's the active tab
    if (tabId === activeTabId) {
      tab.terminal.focus();
    }
  }

  function fitTerminal() {
    if (!activeTab?.terminal || !activeTab.fitAddon) return;

    activeTab.fitAddon.fit();
    
    if (connection && isConnected) {
      connection.resize(activeTab.terminal.cols, activeTab.terminal.rows);
    }
  }

  function fitTerminalForTab(tab: TabSession) {
    if (!tab.terminal || !tab.fitAddon || !tab.connection || !tab.isConnected) {
      console.log(`fitTerminalForTab: Skipping tab ${tab.id} - missing requirements`);
      return;
    }

    // Only fit if the terminal is properly initialized (has element)
    if (!tab.terminal.element) {
      console.log(`fitTerminalForTab: Skipping tab ${tab.id} - no terminal element`);
      return;
    }

    console.log(`fitTerminalForTab: Fitting tab ${tab.id}, cols=${tab.terminal.cols}, rows=${tab.terminal.rows}`);
    tab.fitAddon.fit();
    tab.connection.resize(tab.terminal.cols, tab.terminal.rows);
  }

  function fitAllTerminals() {
    console.log(`fitAllTerminals: Fitting ${tabs.length} tabs`);
    
    // First, fit the visible terminal to get the correct dimensions
    const activeTab = tabs.find(tab => tab.id === activeTabId);
    let targetCols = 80, targetRows = 24; // defaults
    
    if (activeTab?.terminal && activeTab.fitAddon && activeTab.connection && activeTab.isConnected && activeTab.terminal.element) {
      console.log(`fitAllTerminals: Fitting visible tab ${activeTab.id} to get dimensions`);
      activeTab.fitAddon.fit();
      targetCols = activeTab.terminal.cols;
      targetRows = activeTab.terminal.rows;
      activeTab.connection.resize(targetCols, targetRows);
      console.log(`fitAllTerminals: Target dimensions from visible terminal: ${targetCols}x${targetRows}`);
    }
    
    // Then apply the same dimensions to all other terminals
    tabs.forEach(tab => {
      if (tab.terminal && tab.fitAddon && tab.connection && tab.isConnected && tab.terminal.element && tab.id !== activeTabId) {
        console.log(`fitAllTerminals: Applying ${targetCols}x${targetRows} to hidden tab ${tab.id}`);
        // Set the terminal size to match the visible one
        tab.terminal.resize(targetCols, targetRows);
        tab.connection.resize(targetCols, targetRows);
      }
    });
  }

  // Solarized Dark theme
  const solarizedDark = {
    background: '#002b36',
    foreground: '#839496',
    cursor: '#93a1a1',
    black: '#073642',
    red: '#dc322f',
    green: '#859900',
    yellow: '#b58900',
    blue: '#268bd2',
    magenta: '#d33682',
    cyan: '#2aa198',
    white: '#eee8d5',
    brightBlack: '#586e75',
    brightRed: '#cb4b16',
    brightGreen: '#859900',
    brightYellow: '#b58900',
    brightBlue: '#268bd2',
    brightMagenta: '#d33682',
    brightCyan: '#2aa198',
    brightWhite: '#fdf6e3'
  };

  // Solarized Light theme
  const solarizedLight = {
    background: '#fdf6e3',
    foreground: '#657b83',
    cursor: '#586e75',
    black: '#073642',
    red: '#dc322f',
    green: '#859900',
    yellow: '#b58900',
    blue: '#268bd2',
    magenta: '#d33682',
    cyan: '#2aa198',
    white: '#eee8d5',
    brightBlack: '#002b36',
    brightRed: '#cb4b16',
    brightGreen: '#859900',
    brightYellow: '#b58900',
    brightBlue: '#268bd2',
    brightMagenta: '#d33682',
    brightCyan: '#2aa198',
    brightWhite: '#657b83'
  };

  function updateTerminalTheme() {
    const isDarkMode = document.documentElement.classList.contains('dark');
    const theme = isDarkMode ? solarizedDark : solarizedLight;
    
    // Update theme for all terminals
    tabs.forEach(tab => {
      if (tab.terminal) {
        tab.terminal.options.theme = theme;
        // Re-fit terminal after theme change
        if (tab.fitAddon) {
          setTimeout(() => tab.fitAddon?.fit(), 0);
        }
      }
    });
  }

  onMount(() => {
    // Detect initial theme from document class
    const isDarkMode = document.documentElement.classList.contains('dark');
    const theme = isDarkMode ? solarizedDark : solarizedLight;

    // Create the first tab
    const firstTabId = createNewTab();
    activeTabId = firstTabId;

    // Create terminal and fitAddon for first tab
    const newTerminal = new Terminal({
      cursorBlink: true,
      theme,
      fontSize: 13,
      fontFamily: 'Monaco, "Menlo", "Ubuntu Mono", monospace',
      allowTransparency: true
    });

    const newFitAddon = new FitAddon();
    newTerminal.loadAddon(newFitAddon);
    
    const tabIndex = tabs.findIndex(t => t.id === firstTabId);
    if (tabIndex !== -1) {
      tabs[tabIndex].terminal = newTerminal;
      tabs[tabIndex].fitAddon = newFitAddon;
      tabs = [...tabs];
      
      // Initialize the terminal immediately
      setTimeout(() => {
        initializeTerminal(firstTabId);
      }, 10);
    }

    // Create the connection for first tab
    createConnection(firstTabId);

    // Handle window resize - only fit the currently visible terminal
    function onWindowResize() {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        fitAllTerminals(); // This now only fits the active terminal
      }, 100);
    }

    // Listen for theme changes by observing the document class
    const observer = new MutationObserver(() => updateTerminalTheme());
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class']
    });

    window.addEventListener('resize', onWindowResize);
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    document.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('resize', onWindowResize);
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
      document.removeEventListener('keydown', handleKeyDown);
      observer.disconnect();
    };
  });

  onDestroy(() => {
    // Clean up all tabs - this will send close messages for each connection
    tabs.forEach(tab => {
      cleanupTab(tab);
    });
  });

  function handleData(tabId: string, data: Uint8Array) {
    const tab = tabs.find(t => t.id === tabId);
    if (tab?.terminal) {
      const text = new TextDecoder().decode(data);
      tab.terminal.write(text);
    }
  }

  function handleReady(tabId: string) {
    const tabIndex = tabs.findIndex(t => t.id === tabId);
    if (tabIndex === -1) return;

    tabs[tabIndex].isConnecting = false;
    tabs[tabIndex].isConnected = true;
    tabs = [...tabs];

    // Send resize message immediately after receiving ShellReadyMessage
    const tab = tabs[tabIndex];
    if (tab.connection && tab.terminal && tab.fitAddon) {
      fitTerminalForTab(tab);
    }
  }

  function handleExit(tabId: string) {
    const tabIndex = tabs.findIndex(t => t.id === tabId);
    if (tabIndex === -1) return;

    tabs[tabIndex].isConnected = false;
    tabs = [...tabs];
    
    if (tabs[tabIndex].terminal) {
      tabs[tabIndex].terminal.write('\r\n[Shell session ended]');
    }
  }

  function handleError(tabId: string, errorMsg: string) {
    const tabIndex = tabs.findIndex(t => t.id === tabId);
    if (tabIndex === -1) return;

    tabs[tabIndex].error = errorMsg;
    tabs[tabIndex].isConnecting = false;
    tabs[tabIndex].isConnected = false;
    tabs = [...tabs];
  }

  // Handle window resize
  let resizeTimeout: NodeJS.Timeout;

  // Remove automatic fitting on tab switch - let it happen only when needed

  // Function to add a new tab
  function addNewTab() {
    const isDarkMode = document.documentElement.classList.contains('dark');
    const theme = isDarkMode ? solarizedDark : solarizedLight;
    
    const newTabId = createNewTab();
    
    // Create terminal and fitAddon for new tab
    const newTerminal = new Terminal({
      cursorBlink: true,
      theme,
      fontSize: 13,
      fontFamily: 'Monaco, "Menlo", "Ubuntu Mono", monospace',
      allowTransparency: true
    });

    const newFitAddon = new FitAddon();
    newTerminal.loadAddon(newFitAddon);
    
    const tabIndex = tabs.findIndex(t => t.id === newTabId);
    if (tabIndex !== -1) {
      tabs[tabIndex].terminal = newTerminal;
      tabs[tabIndex].fitAddon = newFitAddon;
      tabs = [...tabs];
      
      // Initialize the terminal immediately  
      setTimeout(() => {
        initializeTerminal(newTabId);
      }, 10);
    }

    // Switch to new tab and create connection
    activeTabId = newTabId;
    createConnection(newTabId);
  }

  function toggleMaximize() {
    if (!isMaximized) {
      // Store windowed dimensions and position before going fullscreen
      saveWindowedState();
      terminalContainer.requestFullscreen?.();
    } else {
      // Apply windowed state immediately before exiting fullscreen
      hasRestoredState = true;
      restoreWindowedState();
      document.exitFullscreen?.();
    }
  }

  function saveWindowedState() {
    const rect = terminalContainer.getBoundingClientRect();
    windowedWidth = rect.width;
    windowedHeight = rect.height;
    windowedLeft = rect.left;
    windowedTop = rect.top;
  }

  function restoreWindowedState() {
    terminalContainer.style.position = 'absolute';
    terminalContainer.style.width = `${windowedWidth}px`;
    terminalContainer.style.height = `${windowedHeight}px`;
    terminalContainer.style.left = `${windowedLeft}px`;
    terminalContainer.style.top = `${windowedTop}px`;
    terminalContainer.style.margin = '0';
    terminalContainer.style.zIndex = '1000';
  }

  // Handle ESC key to preemptively restore state
  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape' && isMaximized && !hasRestoredState) {
      // User pressed ESC in fullscreen, preemptively restore state
      hasRestoredState = true;
      restoreWindowedState();
    }
  }

  // Handle fullscreen change events
  function handleFullscreenChange() {
    const wasMaximized = isMaximized;
    isMaximized = !!document.fullscreenElement;

    // Restore windowed state when exiting fullscreen (only if not already restored)
    if (wasMaximized && !isMaximized) {
      if (!hasRestoredState) {
        // ESC key was pressed, restore state immediately
        restoreWindowedState();
      }
      hasRestoredState = false; // Reset flag
      
      // Terminal fit with minimal delay - fit all terminals after fullscreen changes
      setTimeout(() => {
        fitAllTerminals();
      }, 10);
    } else if (!wasMaximized && isMaximized) {
      // Clear positioning when entering fullscreen
      terminalContainer.style.position = '';
      terminalContainer.style.width = '';
      terminalContainer.style.height = '';
      terminalContainer.style.left = '';
      terminalContainer.style.top = '';
      terminalContainer.style.margin = '';
      terminalContainer.style.zIndex = '';
      hasRestoredState = false; // Reset flag
      
      setTimeout(() => {
        fitAllTerminals();
      }, 10);
    }
  }

  // Manual resize functionality
  function startResize(event: MouseEvent) {
    if (isMaximized) return; // Don't allow resize in fullscreen
    
    isResizing = true;
    initialMouseX = event.clientX;
    initialMouseY = event.clientY;
    initialWidth = terminalContainer.offsetWidth;
    initialHeight = terminalContainer.offsetHeight;
    
    event.preventDefault();
    document.addEventListener('mousemove', handleResize);
    document.addEventListener('mouseup', stopResize);
    document.body.style.userSelect = 'none';
    document.body.style.cursor = 'nw-resize';
  }

  function handleResize(event: MouseEvent) {
    if (!isResizing) return;
    
    const deltaX = event.clientX - initialMouseX;
    const deltaY = event.clientY - initialMouseY;
    
    // Calculate new dimensions with viewport constraints
    // Leave some padding (32px) from viewport edges for visibility
    const maxWidth = window.innerWidth - 32;
    const maxHeight = window.innerHeight - 32;
    
    const newWidth = Math.max(300, Math.min(maxWidth, initialWidth + deltaX));
    const newHeight = Math.max(200, Math.min(maxHeight, initialHeight + deltaY));
    
    terminalContainer.style.width = `${newWidth}px`;
    terminalContainer.style.height = `${newHeight}px`;
    
    // Debounced terminal resize - fit all terminals during manual resize
    clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(() => {
      fitAllTerminals();
    }, 50);
  }

  function stopResize() {
    isResizing = false;
    document.removeEventListener('mousemove', handleResize);
    document.removeEventListener('mouseup', stopResize);
    document.body.style.userSelect = '';
    document.body.style.cursor = '';
    
    // Save the new windowed state
    if (!isMaximized) {
      saveWindowedState();
    }
    
    // Final terminal resize - fit all terminals after manual resize ends
    fitAllTerminals();
  }

  // Drag-to-move functionality
  function startDrag(event: MouseEvent) {
    if (isMaximized) return; // Don't drag in fullscreen
    
    // Check if click is on a button or other interactive element
    const target = event.target as HTMLElement;
    if (target.tagName === 'BUTTON' || target.closest('button')) {
      return;
    }
    
    isDragging = true;
    initialMouseX = event.clientX;
    initialMouseY = event.clientY;
    
    const rect = terminalContainer.getBoundingClientRect();
    initialLeft = rect.left;
    initialTop = rect.top;
    
    event.preventDefault();
    event.stopPropagation();
    document.addEventListener('mousemove', handleDrag);
    document.addEventListener('mouseup', stopDrag);
    document.body.style.userSelect = 'none';
    document.body.style.cursor = 'move';
    
    // Ensure terminal is positioned absolutely for dragging
    terminalContainer.style.position = 'absolute';
    terminalContainer.style.left = `${initialLeft}px`;
    terminalContainer.style.top = `${initialTop}px`;
    terminalContainer.style.margin = '0';
    terminalContainer.style.zIndex = '1000';
  }

  function handleDrag(event: MouseEvent) {
    if (!isDragging) return;
    
    const deltaX = event.clientX - initialMouseX;
    const deltaY = event.clientY - initialMouseY;
    
    // Calculate new position with viewport constraints
    const newLeft = Math.max(0, Math.min(window.innerWidth - terminalContainer.offsetWidth, initialLeft + deltaX));
    const newTop = Math.max(0, Math.min(window.innerHeight - terminalContainer.offsetHeight, initialTop + deltaY));
    
    terminalContainer.style.left = `${newLeft}px`;
    terminalContainer.style.top = `${newTop}px`;
  }

  function stopDrag() {
    isDragging = false;
    document.removeEventListener('mousemove', handleDrag);
    document.removeEventListener('mouseup', stopDrag);
    document.body.style.userSelect = '';
    document.body.style.cursor = '';
    
    // Save the new windowed state
    if (!isMaximized) {
      saveWindowedState();
    }
  }
</script>

<div class="shell-terminal-container" bind:this={terminalContainer}>
  <div class="shell-header" on:mousedown={startDrag} role="button" tabindex="0" title="Drag to move terminal">
    <div class="flex items-center justify-between">
      <span class="text-sm font-medium text-gray-600 dark:text-gray-300 pointer-events-none">
        Shell - {runnerName}
      </span>
      <div class="flex items-center space-x-2 pointer-events-auto">
        <button
          on:click={toggleMaximize}
          class="w-3 h-3 rounded-full bg-green-500 hover:bg-green-400 transition-colors cursor-pointer flex items-center justify-center pointer-events-auto"
          title={isMaximized ? "Restore" : "Maximize"}
          aria-label={isMaximized ? "Restore" : "Maximize"}
        >
          {#if isMaximized}
            <svg class="w-2 h-2 text-green-900" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h4a1 1 0 010 2H6.414l2.293 2.293a1 1 0 11-1.414 1.414L5 6.414V8a1 1 0 01-2 0V4zm9 1a1 1 0 010-2h4a1 1 0 011 1v4a1 1 0 01-2 0V6.414l-2.293 2.293a1 1 0 11-1.414-1.414L13.586 5H12zm-9 7a1 1 0 012 0v1.586l2.293-2.293a1 1 0 111.414 1.414L6.414 15H8a1 1 0 010 2H4a1 1 0 01-1-1v-4zm13-1a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 010-2h1.586l-2.293-2.293a1 1 0 111.414-1.414L15 13.586V12a1 1 0 011-1z" clip-rule="evenodd"></path>
            </svg>
          {:else}
            <svg class="w-2 h-2 text-green-900" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v12a1 1 0 01-1 1H4a1 1 0 01-1-1V4zm2 2v8h10V6H5z" clip-rule="evenodd"></path>
            </svg>
          {/if}
        </button>
        <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
        <button
          on:click={onClose}
          class="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors cursor-pointer flex items-center justify-center pointer-events-auto"
          title="Close Shell"
          aria-label="Close Shell"
        >
          <svg class="w-2 h-2 text-red-900" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
          </svg>
        </button>
      </div>
    </div>
  </div>

  <!-- Tab Bar -->
  <div class="tab-bar">
    <div class="tabs-container">
      {#each tabs as tab}
        <div 
          class="tab {tab.id === activeTabId ? 'active' : ''}"
          on:click={() => switchToTab(tab.id)}
          on:keydown={(event) => {
            if (event.key === 'Enter' || event.key === ' ') {
              event.preventDefault();
              switchToTab(tab.id);
            }
          }}
          role="button"
          tabindex="0"
        >
          <span class="tab-title">{tab.title}</span>
          <button 
            class="tab-close"
            on:click|stopPropagation={() => closeTab(tab.id)}
            title="Close tab"
            aria-label="Close tab"
          >
            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
            </svg>
          </button>
        </div>
      {/each}
      <button 
        class="new-tab-button"
        on:click={addNewTab}
        title="New tab"
        aria-label="New tab"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
        </svg>
      </button>
    </div>
  </div>

  <div class="shell-body">
    {#each tabs as tab}
      <div 
        class="terminal-tab {tab.id === activeTabId ? 'active' : ''}" 
        data-tab-id={tab.id}
      >
        {#if tab.isConnecting}
          <div class="shell-status">
            <div class="flex items-center justify-center space-x-3">
              <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
              <span>Connecting to shell...</span>
            </div>
          </div>
        {:else if tab.error}
          <div class="shell-status error">
            <div class="text-center">
              <div class="text-red-400 mb-2">
                <svg class="w-8 h-8 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <p class="text-red-300">Connection Error</p>
              <p class="text-sm text-red-200 mt-1">{tab.error}</p>
            </div>
          </div>
        {/if}
        <!-- Terminal content will be added here by xterm.js -->
      </div>
    {/each}
  </div>
  
  <!-- Resize handle - only show when not in fullscreen -->
  {#if !isMaximized}
    <div 
      class="resize-handle"
      on:mousedown={startResize}
      role="button"
      tabindex="0"
      title="Resize terminal"
      aria-label="Resize terminal by dragging"
    ></div>
  {/if}
</div>

<style>
  .shell-terminal-container {
    background-color: rgb(255 255 255 / 0.85);
    backdrop-filter: blur(8px);
    border: 1px solid rgb(229 231 235 / 0.5);
    border-radius: 0.5rem;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.25);
    width: 800px; /* Set a reasonable initial width */
    min-width: 300px;
    height: 500px;
    min-height: 200px;
    max-height: 90vh;
    max-width: none; /* Remove max-width constraint */
    display: flex;
    flex-direction: column;
    position: relative;
    resize: none; /* Disable default resize to use custom handle */
  }

  .tab-bar {
    background-color: rgb(229 231 235 / 0.6);
    border-bottom: 1px solid rgb(209 213 219 / 0.7);
    padding: 0;
    flex-shrink: 0;
  }

  :global(.dark) .tab-bar {
    background-color: rgb(55 65 81 / 0.6);
    border-bottom: 1px solid rgb(75 85 99 / 0.7);
  }

  .tabs-container {
    display: flex;
    align-items: center;
    overflow-x: auto;
    scrollbar-width: none;
  }

  .tabs-container::-webkit-scrollbar {
    display: none;
  }

  .tab {
    background-color: rgb(243 244 246 / 0.7);
    border-right: 1px solid rgb(209 213 219 / 0.5);
    padding: 0.5rem 0.75rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 120px;
    max-width: 200px;
    user-select: none;
    transition: background-color 0.15s ease;
  }

  .tab:hover {
    background-color: rgb(243 244 246);
  }

  .tab.active {
    background-color: rgb(255 255 255 / 0.9);
    border-bottom: 2px solid rgb(59 130 246);
  }

  :global(.dark) .tab {
    background-color: rgb(31 41 55 / 0.7);
    border-right: 1px solid rgb(55 65 81 / 0.5);
  }

  :global(.dark) .tab:hover {
    background-color: rgb(31 41 55);
  }

  :global(.dark) .tab.active {
    background-color: rgb(30 41 59 / 0.9);
    border-bottom: 2px solid rgb(59 130 246);
  }

  .tab-title {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.875rem;
    color: rgb(75 85 99);
  }

  :global(.dark) .tab-title {
    color: rgb(209 213 219);
  }

  .tab.active .tab-title {
    color: rgb(17 24 39);
    font-weight: 500;
  }

  :global(.dark) .tab.active .tab-title {
    color: rgb(243 244 246);
  }

  .tab-close {
    opacity: 0;
    transition: opacity 0.15s ease;
    padding: 0.125rem;
    border-radius: 0.25rem;
    color: rgb(107 114 128);
  }

  .tab:hover .tab-close,
  .tab.active .tab-close {
    opacity: 1;
  }

  .tab-close:hover {
    background-color: rgb(239 68 68 / 0.1);
    color: rgb(239 68 68);
  }

  :global(.dark) .tab-close {
    color: rgb(156 163 175);
  }

  :global(.dark) .tab-close:hover {
    background-color: rgb(239 68 68 / 0.2);
    color: rgb(248 113 113);
  }

  .new-tab-button {
    background-color: transparent;
    border: none;
    padding: 0.5rem 0.75rem;
    cursor: pointer;
    color: rgb(107 114 128);
    transition: all 0.15s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .new-tab-button:hover {
    background-color: rgb(243 244 246 / 0.7);
    color: rgb(59 130 246);
  }

  :global(.dark) .new-tab-button {
    color: rgb(156 163 175);
  }

  :global(.dark) .new-tab-button:hover {
    background-color: rgb(31 41 55 / 0.7);
    color: rgb(96 165 250);
  }

  :global(.dark) .shell-terminal-container {
    background-color: rgb(30 41 59 / 0.85);
    border: 1px solid rgb(71 85 105 / 0.3);
  }

  .shell-terminal-container:fullscreen {
    backdrop-filter: none;
    border-radius: 0;
    height: 100vh;
    max-height: none;
    background-color: rgb(255 255 255);
    border: none;
  }

  :global(.dark) .shell-terminal-container:fullscreen {
    background-color: rgb(17 24 39);
  }

  .shell-header {
    background-color: rgb(243 244 246 / 0.9);
    padding: 0.5rem 1rem;
    border-bottom: 1px solid rgb(229 231 235 / 0.7);
    flex-shrink: 0;
    cursor: move;
    user-select: none;
  }

  .shell-header:hover {
    background-color: rgb(243 244 246);
  }

  :global(.dark) .shell-header {
    background-color: rgb(31 41 55 / 0.9);
    border-bottom: 1px solid rgb(55 65 81 / 0.7);
  }

  :global(.dark) .shell-header:hover {
    background-color: rgb(31 41 55);
  }

  .shell-terminal-container:fullscreen .shell-header {
    background-color: rgb(243 244 246);
    border-bottom: 1px solid rgb(229 231 235);
    cursor: default; /* Don't show move cursor in fullscreen */
  }

  :global(.dark) .shell-terminal-container:fullscreen .shell-header {
    background-color: rgb(31 41 55);
    border-bottom: 1px solid rgb(55 65 81);
  }

  .shell-terminal-container:fullscreen .shell-header:hover {
    background-color: rgb(243 244 246); /* No hover effect in fullscreen */
  }

  :global(.dark) .shell-terminal-container:fullscreen .shell-header:hover {
    background-color: rgb(31 41 55); /* No hover effect in fullscreen */
  }

  .shell-body {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .shell-body::-webkit-scrollbar {
    display: none !important;
  }

  .shell-body {
    scrollbar-width: none !important; /* Firefox */
    -ms-overflow-style: none !important; /* IE and Edge */
  }

  .shell-status {
    color: rgb(209 213 219);
    padding: 2rem;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .shell-status.error {
    color: rgb(252 165 165);
  }

  .terminal-tab {
    flex: 1;
    height: 100%;
    padding: 0;
    overflow: hidden;
    border: none !important;
    outline: none !important;
    position: relative;
    width: 100%;
    display: none; /* Hidden by default */
  }

  .terminal-tab.active {
    display: block; /* Show only the active tab */
  }

  .terminal-tab::-webkit-scrollbar {
    display: none !important;
  }

  .terminal-tab {
    scrollbar-width: none !important; /* Firefox */
    -ms-overflow-style: none !important; /* IE and Edge */
  }

  .resize-handle {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 20px;
    height: 20px;
    cursor: nw-resize;
    background: linear-gradient(
      135deg,
      transparent 0%,
      transparent 30%,
      rgb(156 163 175 / 0.3) 30%,
      rgb(156 163 175 / 0.3) 35%,
      transparent 35%,
      transparent 45%,
      rgb(156 163 175 / 0.3) 45%,
      rgb(156 163 175 / 0.3) 50%,
      transparent 50%,
      transparent 60%,
      rgb(156 163 175 / 0.3) 60%,
      rgb(156 163 175 / 0.3) 65%,
      transparent 65%
    );
    z-index: 10;
    border-bottom-right-radius: 0.5rem;
  }

  .resize-handle:hover {
    background: linear-gradient(
      135deg,
      transparent 0%,
      transparent 30%,
      rgb(156 163 175 / 0.6) 30%,
      rgb(156 163 175 / 0.6) 35%,
      transparent 35%,
      transparent 45%,
      rgb(156 163 175 / 0.6) 45%,
      rgb(156 163 175 / 0.6) 50%,
      transparent 50%,
      transparent 60%,
      rgb(156 163 175 / 0.6) 60%,
      rgb(156 163 175 / 0.6) 65%,
      transparent 65%
    );
  }

  :global(.dark) .resize-handle {
    background: linear-gradient(
      135deg,
      transparent 0%,
      transparent 30%,
      rgb(209 213 219 / 0.3) 30%,
      rgb(209 213 219 / 0.3) 35%,
      transparent 35%,
      transparent 45%,
      rgb(209 213 219 / 0.3) 45%,
      rgb(209 213 219 / 0.3) 50%,
      transparent 50%,
      transparent 60%,
      rgb(209 213 219 / 0.3) 60%,
      rgb(209 213 219 / 0.3) 65%,
      transparent 65%
    );
  }

  :global(.dark) .resize-handle:hover {
    background: linear-gradient(
      135deg,
      transparent 0%,
      transparent 30%,
      rgb(209 213 219 / 0.6) 30%,
      rgb(209 213 219 / 0.6) 35%,
      transparent 35%,
      transparent 45%,
      rgb(209 213 219 / 0.6) 45%,
      rgb(209 213 219 / 0.6) 50%,
      transparent 50%,
      transparent 60%,
      rgb(209 213 219 / 0.6) 60%,
      rgb(209 213 219 / 0.6) 65%,
      transparent 65%
    );
  }

  /* Hide resize handle in fullscreen */
  .shell-terminal-container:fullscreen .resize-handle {
    display: none;
  }
</style>