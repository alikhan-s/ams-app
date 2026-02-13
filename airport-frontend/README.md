# Airport Management System - Frontend

A production-ready React frontend application for the Airport Management System, built with modern technologies and best practices.

## 🚀 Tech Stack

- **Framework**: React 18 + Vite + TypeScript
- **Styling**: Tailwind CSS + Shadcn/UI
- **State Management**: Zustand (Auth & Global State)
- **Data Fetching**: TanStack Query (React Query) v5
- **Routing**: React Router v6
- **HTTP Client**: Axios with JWT interceptors
- **Forms**: React Hook Form + Zod validation
- **Notifications**: React Hot Toast

## 📁 Project Structure

```
src/
├── api/              # API client and service functions
│   ├── axios.ts      # Axios instance with interceptors
│   └── index.ts      # API endpoints
├── components/       # Reusable UI components
│   ├── ui/          # Shadcn UI components
│   └── Layout.tsx   # Main layout with navbar
├── pages/           # Application pages
│   ├── admin/       # Admin/Staff pages
│   ├── Home.tsx     # Flight search & results
│   ├── Login.tsx    # Login page
│   ├── Register.tsx # Registration page
│   ├── Bookings.tsx # My bookings with baggage
│   ├── BookFlight.tsx # Booking flow
│   └── Profile.tsx  # User profile
├── store/           # Zustand stores
│   └── authStore.ts # Authentication state
├── types/           # TypeScript interfaces
│   └── index.ts     # All type definitions
├── lib/             # Utility functions
│   └── utils.ts     # Helper functions
├── App.tsx          # Main app component with routing
└── main.tsx         # Application entry point
```

## 🔌 Backend Integration

The frontend connects to a Go backend running on `http://localhost:8080/api/v1`

### API Endpoints Used:

**Auth:**
- `POST /auth/register` - User registration
- `POST /auth/login` - User authentication (returns JWT)

**Flights:**
- `GET /flights?origin=...&destination=...&date=...` - Search flights

**Bookings:**
- `POST /bookings` - Book a flight (with passport_no & phone for JIT passenger creation)
- `GET /bookings/my` - Get user's bookings
- `POST /bookings/:id/cancel` - Cancel a ticket
- `GET /bookings/baggage` - Get user's baggage

**Operations (Staff/Admin only):**
- `GET /ops/gates` - List gates
- `POST /ops/gates` - Create a gate
- `POST /ops/baggage` - Check-in baggage
- `GET /ops/baggage` - List all baggage with passenger info
- `PATCH /ops/baggage/:id` - Update baggage status

## 🎨 Features

### Public Pages:
- **Home**: Flight search form with real-time results
  - Search by origin, destination, and date
  - Visual flight status indicators (Scheduled=Green, Delayed=Yellow)
  - Direct booking flow

### Authenticated User Pages:
- **My Bookings**: View active and cancelled tickets
  - Cancel active tickets
  - Check baggage status for all bookings
- **Booking Flow**: 
  - Step-by-step booking process
  - Passport and phone information capture (JIT passenger profile creation)
  - Flight details confirmation

### Admin/Staff Pages:
- **Dashboard**:
  - **Gates Management**: Create and view gates with status indicators
  - **Baggage Operations**: 
    - View all baggage with passenger names
    - Update baggage status (RECEIVED → LOADED → IN_TRANSIT → DELIVERED)

## 🛠️ Setup Instructions

### Prerequisites
- Node.js 18+ and npm/yarn
- Backend server running on port 8080

### Installation

1. **Clone and navigate:**
   ```bash
   cd airport-frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Configure environment:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` if your backend is on a different URL:
   ```
   VITE_API_BASE_URL=http://localhost:8080/api/v1
   ```

4. **Start development server:**
   ```bash
   npm run dev
   ```
   
   The app will open at `http://localhost:3000`

5. **Build for production:**
   ```bash
   npm run build
   npm run preview
   ```

## 🔐 Authentication Flow

1. JWT token is stored in `localStorage` after login
2. Axios interceptor automatically attaches token to all requests
3. 401 responses trigger automatic logout and redirect to login
4. Protected routes check authentication state before rendering

## 🎯 Key Implementation Details

### Type Safety
All TypeScript interfaces mirror Go backend structs exactly:
```typescript
interface Flight {
  id: number;
  flight_no: string;
  origin: string;
  destination: string;
  // ... matches Go Flight struct
}
```

### JWT Token Handling
```typescript
// Stored in localStorage
localStorage.setItem('authToken', token);

// Auto-attached via interceptor
config.headers.Authorization = `Bearer ${token}`;
```

### Passenger Profile Logic
The booking flow handles JIT (Just-In-Time) passenger creation:
- If user has no passenger profile, `passport_no` and `phone` are required
- Backend creates passenger profile automatically during booking
- Subsequent bookings don't require this info

### Error Handling
- Global axios interceptor catches all errors
- User-friendly toast notifications
- 401 errors trigger automatic logout
- Network errors display appropriate messages

## 🎨 UI/UX Design

Based on the provided screenshots, the application features:
- **Color Scheme**: Blue (#3b5998) primary with orange (#ff6b35) accents
- **Clean Layout**: Card-based design with clear hierarchy
- **Responsive**: Mobile-first approach with grid layouts
- **Status Indicators**: Color-coded flight and baggage statuses
- **Smooth Interactions**: Loading states and transition animations

## 📝 Available Scripts

```bash
npm run dev      # Start development server
npm run build    # Build for production
npm run preview  # Preview production build
npm run lint     # Run ESLint
```

## 🔄 State Management

### Zustand (Auth Store)
```typescript
const { isAuthenticated, user, setAuth, logout } = useAuthStore();
```

### TanStack Query (Data Fetching)
```typescript
const { data: flights } = useQuery({
  queryKey: ['flights', params],
  queryFn: () => flightsApi.search(params),
});
```

## 🚦 Route Protection

```typescript
// Protected routes require authentication
<ProtectedRoute>
  <Bookings />
</ProtectedRoute>

// Admin routes require STAFF or ADMIN role
<AdminRoute>
  <AdminDashboard />
</AdminRoute>
```

## 📱 Responsive Design

All pages are fully responsive:
- Mobile: Single column layouts
- Tablet: 2-column grids
- Desktop: 3-4 column grids for optimal content display

## 🐛 Troubleshooting

**API Connection Issues:**
- Ensure backend is running on port 8080
- Check CORS settings in backend
- Verify `.env` file configuration

**Authentication Issues:**
- Clear localStorage and cookies
- Check JWT token expiration
- Verify backend JWT_SECRET matches

**Build Issues:**
- Clear `node_modules` and reinstall
- Check Node.js version (18+)
- Update dependencies if needed

## 📄 License

MIT License - see LICENSE file for details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

Built with using React + Vite + TypeScript