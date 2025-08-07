# GARM SPA (SvelteKit)

This is a Single Page Application (SPA) implementation of the GARM web interface using SvelteKit.

## Features

- **Lightweight**: Uses SvelteKit for minimal bundle size and fast performance
- **Modern**: TypeScript-first development with full type safety
- **Responsive**: Mobile-first design using Tailwind CSS
- **Real-time**: WebSocket integration for live updates
- **API-driven**: Uses the existing GARM REST API endpoints

## Development

```bash
cd webapp
npm install
npm run dev
```

## Build

```bash
npm run build
```

This generates static files in the `build/` directory that can be served by the Go backend.

## Architecture

- **SvelteKit**: Frontend framework with SSG/SPA capabilities
- **TypeScript**: Full type safety with auto-generated API client
- **Tailwind CSS**: Utility-first CSS framework matching the HTMX version
- **REST API**: Uses existing `/api/v1/*` endpoints
- **JWT Auth**: Same authentication system as the main app

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