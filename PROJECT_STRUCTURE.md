# 📁 Project Structure

Complete overview of the Multi-Strategy Routing Engine codebase.

```
shipt-route-optimizer/
│
├── backend/                          # Go REST API Server
│   ├── cmd/
│   │   └── main.go                  # Application entry point, server setup
│   │
│   ├── internal/                    # Private application code
│   │   ├── api/
│   │   │   ├── handlers/            # HTTP request handlers
│   │   │   └── router.go            # Route registration
│   │   │
│   │   ├── models/
│   │   │   └── models.go            # Data structures (Order, Shopper, etc.)
│   │   │
│   │   ├── optimizer/
│   │   │   ├── analytics.go         # Route optimization analytics
│   │   │   └── optimizer.go         # Route optimization algorithms
│   │   │
│   │   └── data/
│   │       └── generator.go         # Mock data generation
│   │
│   ├── go.mod                       # Go module definition
│   ├── go.sum                       # Go dependency checksums
│   ├── Makefile                     # Build automation
│   └── .gitignore                   # Git ignore rules
│
├── frontend/                         # React SPA
│   ├── public/
│   │   └── vite.svg                 # Favicon
│   │
│   ├── src/
│   │   ├── components/              # React components
│   │   │   ├── features/            # Feature-specific components
│   │   │   │   ├── AnalyticsDashboard.jsx
│   │   │   │   ├── MapView.jsx
│   │   │   │   ├── Sidebar.jsx
│   │   │   │   └── SummaryPanel.jsx
│   │   │   ├── layout/              # Layout components
│   │   │   │   ├── Navbar.jsx
│   │   │   │   └── SettingsDialog.jsx
│   │   │   └── ui/                  # UI primitives
│   │   │
│   │   ├── pages/                   # Application views
│   │   │   ├── AnalyticsPage.jsx
│   │   │   ├── DashboardPage.jsx
│   │   │   ├── HistoryPage.jsx
│   │   │   └── OptimizePage.jsx
│   │   │
│   │   ├── services/
│   │   │   └── optimizer.js         # Backend API client
│   │   │
│   │   ├── App.jsx                  # Main application component
│   │   ├── main.jsx                 # React entry point
│   │   └── index.css                # Global styles & Tailwind imports
│   │
│   ├── index.html                   # HTML template
│   ├── package.json                 # npm dependencies
│   ├── vite.config.js               # Vite build configuration
│   ├── tailwind.config.js           # Tailwind CSS configuration
│   ├── postcss.config.js            # PostCSS configuration
│   └── .gitignore                   # Git ignore rules
│
├── docs/                             # Documentation assets
│   └── assets/
│       ├── image.png
│       └── image2.png
│
├── README.md                         # Main documentation
├── QUICKSTART.md                     # Quick setup guide
├── PROJECT_STRUCTURE.md              # This file
└── .gitignore                        # Root git ignore

```

## 🔍 Key Files Explained

### Backend

**`cmd/main.go`**
- Server initialization
- CORS configuration
- Route registration
- Entry point for `go run`

**`internal/api/handlers/handlers.go`**
- `/api/health` - Health check
- `/api/sample-data` - Mock data generation
- `/api/optimize` - Route optimization endpoint

**`internal/optimizer/optimizer.go`**
- Haversine distance calculation
- Nearest-neighbor assignment algorithm
- Route sorting and optimization logic

**`internal/models/models.go`**
- `Order` - Delivery order structure
- `Shopper` - Shopper/driver structure
- `Assignment` - Optimized route assignment
- Request/response types

**`internal/data/generator.go`**
- Generates random orders around Birmingham, AL
- Creates realistic test data with delivery windows

### Frontend

**`src/App.jsx`**
- Main application state management
- API integration
- About modal
- Error handling

**`src/pages/DashboardPage.jsx`**
**`src/pages/OptimizePage.jsx`**
**`src/pages/AnalyticsPage.jsx`**
**`src/pages/HistoryPage.jsx`**

**`src/components/features/MapView.jsx`**
- Leaflet map integration
- Custom marker icons
- Route polylines
- Interactive popups

**`src/components/features/Sidebar.jsx`**
- Action buttons (Load Data, Optimize)
- Shopper list with assignments
- Order list with details
- Framer Motion animations

**`src/components/features/SummaryPanel.jsx`**
- Optimization statistics
- Animated number counters
- Improvement calculations
- Route breakdown

**`src/services/optimizer.js`**
- `getSampleData()` - Fetch mock data
- `optimizeRoutes()` - Trigger optimization
- `healthCheck()` - Backend health status

## 🎨 Styling Architecture

- **Tailwind CSS** - Utility-first CSS framework
- **Custom Colors** - Shipt green (#00C389) defined in `tailwind.config.js`
- **Framer Motion** - Smooth animations and transitions
- **Lucide Icons** - Modern icon library

## 📦 Dependencies

### Backend (Go)
- `gin-gonic/gin` - HTTP web framework
- `gin-contrib/cors` - CORS middleware

### Frontend (npm)
- `react` + `react-dom` - UI library
- `vite` - Build tool
- `tailwindcss` - CSS framework
- `leaflet` + `react-leaflet` - Mapping
- `framer-motion` - Animations
- `lucide-react` - Icons

## 🔄 Data Flow

1. **User** clicks "Load Sample Data"
2. **Frontend** calls `/api/sample-data`
3. **Backend** generates random orders/shoppers
4. **Frontend** displays markers on map
5. **User** clicks "Optimize Routes"
6. **Frontend** sends data to `/api/optimize`
7. **Backend** runs optimization algorithm
8. **Frontend** draws routes and shows stats

## 🚀 Build Outputs

### Backend
```bash
make build
# Creates: bin/shipt-route-optimizer
```

### Frontend
```bash
npm run build
# Creates: dist/ directory (static files)
```

---

**For setup instructions, see [QUICKSTART.md](QUICKSTART.md)**

