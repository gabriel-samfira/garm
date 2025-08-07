# GARM SPA (SvelteKit)

This is a Single Page Application (SPA) implementation of the GARM web interface using SvelteKit.

## Features

- **Lightweight**: Uses SvelteKit for minimal bundle size and fast performance
- **Modern**: TypeScript-first development with full type safety
- **Responsive**: Mobile-first design using Tailwind CSS
- **Real-time**: WebSocket integration for live updates
- **API-driven**: Uses the existing GARM REST API endpoints

## Development Environment Setup

### Prerequisites

- **Node.js 18+** and **npm** (or yarn/pnpm)
- **Go 1.21+** (for building the GARM backend)
- **openapi-generator-cli** in your PATH (for API client generation)
- A running GARM backend server

### Installing openapi-generator-cli

**Option 1: NPM Global Install (Recommended)**
```bash
npm install -g @openapitools/openapi-generator-cli
```

**Option 2: Manual Install**
Download from [OpenAPI Generator releases](https://github.com/OpenAPITools/openapi-generator/releases) and add to your PATH.

**Verify Installation:**
```bash
openapi-generator-cli version
```

### Quick Start

1. **Clone the repository** (if not already done)
   ```bash
   git clone https://github.com/cloudbase/garm.git
   cd garm
   ```

2. **Build the webapp** (from project root)
   ```bash
   go generate ./... && cd webapp && npm install && npm run build
   ```

   This will:
   - Generate OpenAPI spec and TypeScript client (`go generate ./...`)
   - Install webapp dependencies (`npm install`)
   - Build the SPA (`npm run build`)

3. **Build and test GARM with embedded webapp**
   ```bash
   # From project root - build, deploy, and restart GARM
   /bin/bash -c ./build-webapp.sh >/dev/null 2>&1 && rm -rf webapp/assets/_app && cp -r webapp/build/* webapp/assets/ && go install ./... && sudo systemctl stop garm && sudo cp /home/ubuntu/go/bin/garm /usr/local/bin/garm && sudo cp ~/go/bin/garm-cli /usr/local/bin/garm-cli && sudo systemctl start garm
   ```

4. **Access the webapp**
   - Navigate to `http://localhost:9997/webapp/` (or your configured port)

### Development Workflow

- **Make changes**: Edit source files in `webapp/src/`
- **Rebuild**: Run step 2-3 above to rebuild and deploy changes
- **No dev server**: The webapp is always served by the GARM binary at `/webapp/` endpoint

### Building for Production

The webapp is **embedded directly into the GARM binary** using Go's `embed` package for zero-dependency deployment.

**Full build process:**
```bash
# From project root
go generate ./... && cd webapp && npm install && npm run build && cd .. && go build -o garm ./cmd/garm-server/
```

The built webapp is served at `http://localhost:9997/webapp/` by the GARM binary with no external dependencies.

### Git Workflow

**DO NOT commit** the following directories:
- `webapp/node_modules/` - Dependencies (managed by package-lock.json)  
- `webapp/.svelte-kit/` - Build cache and generated files
- `webapp/build/` - Production build output

These are already included in `.gitignore`. Only commit source files in `webapp/src/` and configuration files.

### API Client Generation

The webapp uses auto-generated TypeScript clients from the GARM OpenAPI spec using `go generate`.

**When the backend API changes, regenerate everything:**

```bash
# From project root - regenerates OpenAPI spec and TypeScript client
go generate ./... && cd webapp && npm install && npm run build
```

This will:
- Generate fresh OpenAPI spec from Go annotations (`swagger.yaml`)
- Validate the OpenAPI spec  
- Remove old generated client
- Generate new TypeScript client using `openapi-generator-cli`
- Rebuild the webapp with updated types

The `go generate` process requires `openapi-generator-cli` to be installed and in your PATH.

## Architecture

- **SvelteKit**: Frontend framework with SSG/SPA capabilities
- **Go Embed**: Webapp assets embedded directly in GARM binary (`//go:embed all:*`)
- **TypeScript**: Full type safety with auto-generated API client
- **Tailwind CSS**: Utility-first CSS framework matching the HTMX version  
- **REST API**: Uses existing `/api/v1/*` endpoints
- **JWT Auth**: Same authentication system as the main app

### Asset Serving

The webapp is embedded using Go's `embed` package in `webapp/assets/assets.go`:

```go
//go:embed all:*
var EmbeddedSPA embed.FS
```

This allows GARM to serve the entire webapp with zero external dependencies. The webapp assets are compiled into the Go binary at build time.

## Comparison with HTMX Version

| Feature | HTMX Version | SvelteKit SPA |
|---------|-------------|---------------|
| Bundle Size | ~50KB | ~150KB |
| Initial Load | Faster (server-rendered) | Slower (client-rendered) |
| Navigation | Full page reload | Instant navigation |
| Offline Support | None | Possible with service workers |
| Development Experience | Template-based | Component-based |
| Type Safety | None | Full TypeScript |
| Real-time Updates | WebSocket + DOM updates | WebSocket + reactive stores |

## API Integration

The SPA uses a generated TypeScript client that mirrors the Go client structure:

```typescript
import { garmApi } from '$lib/api/client.js';

// List repositories
const repos = await garmApi.listRepositories();

// Create a repository
const newRepo = await garmApi.createRepository({
  name: 'my-repo',
  owner: 'my-org',
  credentials_name: 'github-creds'
});
```