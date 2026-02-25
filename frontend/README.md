# Online Leadership Frontend

Production-ready React frontend application for Online Leadership platform.

## Tech Stack

- **React 18** - UI library
- **Vite** - Build tool and dev server
- **TypeScript** - Type safety
- **React Router** - Client-side routing
- **Fetch API** - HTTP client
- **Nginx** - Static hosting and API proxy
- **Docker** - Containerization

## Features

- JWT authentication with localStorage token storage
- Protected routes with automatic redirect to login
- Responsive design
- Pagination support for list views
- API proxy via Nginx (no CORS needed)

## Project Structure

```
frontend/
├── src/
│   ├── components/      # Reusable components (Layout, ProtectedRoute)
│   ├── contexts/        # React contexts (AuthContext)
│   ├── pages/           # Page components (Login, Games, Leaderboard)
│   ├── services/        # API service layer
│   ├── types/           # TypeScript type definitions
│   ├── App.tsx          # Main app component with routing
│   ├── main.tsx         # Entry point
│   └── index.css        # Global styles
├── nginx.conf           # Nginx configuration
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Docker Compose configuration
└── package.json         # Dependencies and scripts
```

## Development

### Prerequisites

- Node.js 20+ and npm
- Docker and Docker Compose (for production deployment)

### Install Dependencies

```bash
npm install
```

### Run Development Server

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
```

Output will be in the `dist/` directory.

## Docker Deployment

### Build Docker Image

```bash
docker build -t online-leadership-frontend .
```

### Run with Docker Compose

```bash
docker-compose up -d
```

The frontend will be available at `http://localhost` and will proxy API requests to the backend service.

### Environment Configuration

The backend service name in `docker-compose.yml` should match your actual backend service. Update the `backend` service configuration with your backend image and environment variables.

## API Integration

All API requests use relative paths:
- `/api/*` - Proxied to backend:8080
- `/auth/*` - Proxied to backend:8080
- `/admin/*` - Proxied to backend:8080

The Nginx configuration handles all proxying, so no CORS configuration is needed.

## Pages

- `/login` - Login page (public)
- `/games` - Games list (protected)
- `/leaderboard/global` - Global leaderboard (protected)
- `/leaderboard/my` - User's leaderboard (protected)

## Authentication

- Access tokens are stored in `localStorage`
- Authorization header is automatically attached to API requests
- Protected routes redirect to `/login` if user is not authenticated
- Logout clears tokens and redirects to login

## Production Considerations

- Multi-stage Docker build for optimized image size
- Nginx serves static files with caching headers
- SPA routing handled with `try_files` directive
- Gzip compression enabled
- No hardcoded URLs - all API calls use relative paths
