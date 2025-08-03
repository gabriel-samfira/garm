# GARM Web Interface

This is the HTMX-based web interface for GARM (GitHub Actions Runner Manager). It provides a user-friendly way to manage GitHub Actions runners, pools, repositories, and organizations.

## Features

The web interface provides CRUD operations for:

- **Repositories**: Manage GitHub repositories with runner pools
- **Organizations**: Manage GitHub organizations with runner pools  
- **Pools**: View and manage runner pools across all entities
- **Instances**: Monitor and manage individual runner instances
- **Credentials**: Manage GitHub/Gitea authentication credentials (future)

## Architecture

The web interface is built using:

- **HTMX**: For dynamic HTML updates without full page reloads
- **Tailwind CSS**: For responsive and modern styling
- **Alpine.js**: For additional client-side interactivity
- **Go Templates**: For server-side HTML rendering

## Directory Structure

```
web/
├── handlers/           # HTTP handlers for web routes
│   ├── web.go         # Main handlers and dashboard
│   ├── organizations.go # Organization CRUD handlers
│   ├── pools.go       # Pool management handlers
│   └── instances.go   # Instance management handlers
├── templates/          # HTML templates
│   ├── base.html      # Base template with navigation
│   ├── dashboard.html # Dashboard page
│   ├── repositories.html # Repository list page
│   ├── organizations.html # Organization list page
│   ├── pools.html     # Pool list page
│   ├── instances.html # Instance list page
│   └── *-form.html    # Form modals for creating/editing
├── static/            # Static assets
│   ├── css/           # Custom CSS files
│   └── js/            # Custom JavaScript files
└── routers/           # Web route definitions
    └── routers.go     # Route setup and middleware
```

## Routes

### Page Routes (serve full HTML pages)
- `GET /web/` - Dashboard
- `GET /web/repositories` - Repository list
- `GET /web/organizations` - Organization list  
- `GET /web/pools` - Pool list
- `GET /web/instances` - Instance list

### Form Routes (serve modal forms)
- `GET /web/repositories/new` - New repository form
- `GET /web/repositories/{id}/edit` - Edit repository form
- `GET /web/organizations/new` - New organization form
- `GET /web/organizations/{id}/edit` - Edit organization form

### API Routes (return HTML fragments for HTMX)
- `GET /web/api/repositories` - Repository table rows
- `POST /web/api/repositories` - Create repository
- `PUT /web/api/repositories/{id}` - Update repository
- `DELETE /web/api/repositories/{id}` - Delete repository
- Similar patterns for organizations, pools, and instances

## Usage

1. Start the GARM server with the web interface enabled
2. Navigate to `http://localhost:9997/web/` in your browser
3. Use the navigation menu to access different sections
4. Click "Add" buttons to create new resources
5. Use "Edit" and "Delete" actions in table rows

## HTMX Features

The interface uses HTMX for:

- **Dynamic table updates**: Tables refresh automatically every 5 seconds
- **Modal forms**: Create/edit forms appear as modals without page reloads
- **Inline actions**: Delete and other actions update content in-place
- **Loading indicators**: Visual feedback during requests
- **Error handling**: Toast notifications for errors

## Styling

The interface uses Tailwind CSS for:

- **Responsive design**: Works on desktop and mobile devices
- **Component styling**: Cards, tables, forms, and buttons
- **Color coding**: Status indicators and visual hierarchy
- **Dark/light themes**: Consistent color scheme

## Security

The web interface:

- **Inherits authentication**: Uses the same JWT middleware as the API
- **Admin access only**: Requires admin privileges for all operations
- **CSRF protection**: Forms use proper HTTP methods
- **Input validation**: Server-side validation for all inputs

## Future Enhancements

Planned improvements include:

- **Real-time updates**: WebSocket integration for live updates
- **Advanced filtering**: Search and filter capabilities
- **Bulk operations**: Multi-select actions for batch operations
- **Metrics dashboard**: Charts and graphs for monitoring
- **Mobile optimization**: Enhanced mobile experience
- **Credentials management**: Full CRUD for GitHub/Gitea credentials
- **Scale set management**: Enhanced Azure scale set support