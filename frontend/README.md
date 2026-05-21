# KubeWatcher Frontend

A React + TypeScript frontend for the KubeWatcher application that monitors Kubernetes clusters and manages employee reviews.

## Features

- **Dashboard**: Overview of cluster resources (nodes, pods, deployments, services)
- **Nodes**: Detailed view of Kubernetes nodes with status and specifications
- **Pods**: List of pods with search functionality and status indicators
- **Employees**: Employee management with review creation
- **Reviews**: View all employee reviews with ratings and feedback
- **Authentication**: Login and registration system

## Tech Stack

- React 18
- TypeScript
- Material-UI (MUI)
- Vite
- Axios for API calls

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- Backend server running on http://localhost:7979

### Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build for Production

```bash
npm run build
```

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   └── Sidebar.tsx          # Navigation sidebar
│   ├── pages/
│   │   ├── Login.tsx            # Authentication page
│   │   ├── Dashboard.tsx        # Cluster overview
│   │   ├── Nodes.tsx            # Nodes table
│   │   ├── Pods.tsx             # Pods table with search
│   │   ├── Employees.tsx        # Employees management
│   │   └── Reviews.tsx          # Reviews display
│   ├── services/
│   │   ├── api.ts               # Axios configuration
│   │   ├── monitorService.ts    # Kubernetes monitoring API
│   │   └── reviewService.ts     # Employee reviews API
│   ├── hooks/
│   │   └── useAuth.tsx          # Authentication hook
│   ├── App.tsx                  # Main app component
│   └── main.tsx                 # Entry point
├── package.json
├── tsconfig.json
├── vite.config.ts
└── index.html
```

## API Integration

The frontend communicates with the backend API at `http://localhost:7979`. The API endpoints include:

- Authentication: `/api/v1/auth/*`
- Monitoring: `/api/v1/monitor/*`
- Reviews: `/api/v1/reviews/*`

## Styling

The application uses a dark theme inspired by Kubernetes Dashboard:

- Dark background (#121212)
- Blue accents (#2196f3)
- Material-UI components with custom styling
- Responsive design

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run lint` - Run ESLint
- `npm run preview` - Preview production build

### Code Style

- TypeScript for type safety
- ESLint for code quality
- Prettier for code formatting (recommended)

## Contributing

1. Follow the existing code style
2. Use TypeScript interfaces for API responses
3. Test components and API calls
4. Update documentation as needed

## License

MIT License