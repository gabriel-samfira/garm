<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { createShellConnection, type ShellConnection } from '$lib/utils/shell';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import '@xterm/xterm/css/xterm.css';

  export let runnerName: string;
  export let onClose: () => void;

  let terminalElement: HTMLDivElement;
  let connection: ShellConnection | null = null;
  let isConnecting = true;
  let isConnected = false;
  let error = '';
  let terminal: Terminal | null = null;
  let fitAddon: FitAddon | null = null;
  let isMaximized = false;
  let terminalContainer: HTMLDivElement;

  async function createConnection() {
    try {
      connection = await createShellConnection(
        runnerName,
        handleData,
        handleReady,
        handleExit,
        handleError
      );
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to connect';
      isConnecting = false;
    }
  }

  function initializeTerminal() {
    if (!terminal || !terminalElement) return;

    // Open terminal in the DOM element
    terminal.open(terminalElement);


    fitAddon?.fit();

    // Handle terminal input
    terminal.onData((data) => {
      if (connection && isConnected) {
        const encoder = new TextEncoder();
        connection.sendData(encoder.encode(data));
      }
    });

    // Handle terminal resize
    terminal.onResize(({ cols, rows }) => {
      if (connection && isConnected) {
        connection.resize(cols, rows);
      }
    });

    // Focus the terminal
    terminal.focus();
  }

  function fitTerminal() {
    if (!terminal || !fitAddon) return;

    fitAddon.fit();
    
    if (connection && isConnected) {
      connection.resize(terminal.cols, terminal.rows);
    }

    // Adjust container height to fit terminal content (only when not in fullscreen)
    if (!isMaximized && terminalContainer && terminalElement) {
      setTimeout(() => {
        const xtermScreen = terminalElement.querySelector('.xterm-screen') as HTMLElement;
        const headerElement = terminalContainer.querySelector('.shell-header') as HTMLElement;
        
        if (xtermScreen && headerElement) {
          const screenHeight = xtermScreen.offsetHeight;
          const headerHeight = headerElement.offsetHeight;
          const totalHeight = screenHeight + headerHeight;
          
          terminalContainer.style.height = `${totalHeight}px`;
        }
      }, 0);
    } else if (isMaximized && terminalContainer) {
      // Reset height in fullscreen mode
      terminalContainer.style.height = '';
    }
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
    if (!terminal) return;
    const isDarkMode = document.documentElement.classList.contains('dark');
    const theme = isDarkMode ? solarizedDark : solarizedLight;
    terminal.options.theme = theme;
    // Re-fit terminal after theme change
    if (fitAddon) {
      setTimeout(() => fitAddon?.fit(), 0);
    }
  }

  onMount(() => {
    // Detect initial theme from document class
    const isDarkMode = document.documentElement.classList.contains('dark');
    const theme = isDarkMode ? solarizedDark : solarizedLight;

    // Create xterm.js terminal first
    terminal = new Terminal({
      cursorBlink: true,
      theme,
      fontSize: 13,
      fontFamily: 'Monaco, "Menlo", "Ubuntu Mono", monospace',
      allowTransparency: true
    });

    fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);

    // Create the connection
    createConnection();

    // Handle window resize
    function onWindowResize() {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        if (terminal && fitAddon && connection && isConnected) {
          fitTerminal();
        }
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

    return () => {
      window.removeEventListener('resize', onWindowResize);
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
      observer.disconnect();
    };
  });

  onDestroy(() => {
    if (connection) {
      connection.close();
    }
    if (terminal) {
      terminal.dispose();
    }
  });

  function handleData(data: Uint8Array) {
    if (terminal) {
      const text = new TextDecoder().decode(data);
      terminal.write(text);
    }
  }

  function handleReady() {
    isConnecting = false;
    isConnected = true;

    // Send resize message immediately after receiving ShellReadyMessage
    if (connection && terminal && fitAddon && terminalElement) {
      fitTerminal();
    }
  }

  function handleExit() {
    isConnected = false;
    if (terminal) {
      terminal.write('\r\n[Shell session ended]');
    }
  }

  function handleError(errorMsg: string) {
    error = errorMsg;
    isConnecting = false;
    isConnected = false;
  }

  // Handle window resize
  let resizeTimeout: NodeJS.Timeout;

  $: if (!isConnecting && !error && terminalElement && terminal) {
    initializeTerminal();

    // Send resize if connection is ready but we missed it during handleReady
    if (connection && isConnected && fitAddon) {
      fitTerminal();
    }
  }

  function toggleMaximize() {
    isMaximized = !isMaximized;

    if (isMaximized) {
      terminalContainer.requestFullscreen?.();
    } else {
      document.exitFullscreen?.();
    }

    // Send resize after fullscreen change
    setTimeout(() => {
      if (connection && terminal && fitAddon && isConnected) {
        fitTerminal();
      }
    }, 100);
  }

  // Handle fullscreen change events
  function handleFullscreenChange() {
    isMaximized = !!document.fullscreenElement;

    // Send resize when exiting fullscreen
    if (!isMaximized && connection && terminal && fitAddon && isConnected) {
      setTimeout(() => {
        if (fitAddon && connection && terminal) {
          fitTerminal();
        }
      }, 100);
    }
  }
</script>

<div class="shell-terminal-container" bind:this={terminalContainer}>
  <div class="shell-header">
    <div class="flex items-center justify-between">
      <span class="text-sm font-medium text-gray-600 dark:text-gray-300">
        Shell - {runnerName}
      </span>
      <div class="flex items-center space-x-2">
        <button
          on:click={toggleMaximize}
          class="w-3 h-3 rounded-full bg-green-500 hover:bg-green-400 transition-colors cursor-pointer flex items-center justify-center"
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
          class="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors cursor-pointer flex items-center justify-center"
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

  <div class="shell-body">
    {#if isConnecting}
      <div class="shell-status">
        <div class="flex items-center justify-center space-x-3">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
          <span>Connecting to shell...</span>
        </div>
      </div>
    {:else if error}
      <div class="shell-status error">
        <div class="text-center">
          <div class="text-red-400 mb-2">
            <svg class="w-8 h-8 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <p class="text-red-300">Connection Error</p>
          <p class="text-sm text-red-200 mt-1">{error}</p>
        </div>
      </div>
    {:else}
      <div
        bind:this={terminalElement}
        class="terminal"
      ></div>
    {/if}
  </div>
</div>

<style>
  .shell-terminal-container {
    background-color: rgb(255 255 255 / 0.85);
    backdrop-filter: blur(8px);
    border: 1px solid rgb(229 231 235 / 0.5);
    border-radius: 0.5rem;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.25);
    width: 100%;
    height: 500px;
    max-height: 70vh;
    display: flex;
    flex-direction: column;
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
  }

  :global(.dark) .shell-header {
    background-color: rgb(31 41 55 / 0.9);
    border-bottom: 1px solid rgb(55 65 81 / 0.7);
  }

  .shell-terminal-container:fullscreen .shell-header {
    background-color: rgb(243 244 246);
    border-bottom: 1px solid rgb(229 231 235);
  }

  :global(.dark) .shell-terminal-container:fullscreen .shell-header {
    background-color: rgb(31 41 55);
    border-bottom: 1px solid rgb(55 65 81);
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

  .terminal {
    flex: 1;
    height: 100%;
    padding: 0;
    overflow: hidden;
    border: none !important;
    outline: none !important;
    position: relative;
    width: calc(100% + 1px);
    margin-right: -1px;
  }

  .terminal::-webkit-scrollbar {
    display: none !important;
  }

  .terminal {
    scrollbar-width: none !important; /* Firefox */
    -ms-overflow-style: none !important; /* IE and Edge */
  }
</style>