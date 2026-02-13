import { Link, Outlet, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { Button } from './ui/button';
import { Plane } from 'lucide-react';

export function Layout() {
  const { isAuthenticated, user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navigation Bar */}
      <nav className="bg-[#3b5998] text-white shadow-md">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center space-x-2 text-xl font-bold hover:opacity-90">
              <Plane className="h-6 w-6" />
              <span>AMS</span>
            </Link>

            {/* Navigation Links */}
            <div className="flex items-center space-x-6">
              <Link to="/" className="hover:opacity-80 transition-opacity">
                Home
              </Link>
              
              {isAuthenticated ? (
                <>
                  <Link to="/bookings" className="hover:opacity-80 transition-opacity">
                    Booking
                  </Link>
                  
                  <Link to="/profile" className="hover:opacity-80 transition-opacity">
                    Profile
                  </Link>

                  {(user?.role === 'STAFF' || user?.role === 'ADMIN') && (
                    <Link to="/admin" className="hover:opacity-80 transition-opacity">
                      Admin
                    </Link>
                  )}

                  <Button 
                    variant="ghost" 
                    className="text-white hover:bg-white/20"
                    onClick={handleLogout}
                  >
                    Logout
                  </Button>
                </>
              ) : (
                <Link to="/login">
                  <Button variant="ghost" className="text-white hover:bg-white/20">
                    Login
                  </Button>
                </Link>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main>
        <Outlet />
      </main>
    </div>
  );
}