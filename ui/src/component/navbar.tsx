import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useState, useRef, useEffect } from 'react';
import { type ReactElement } from 'react';
import { User } from '../App';
import { storage_token } from '../config';
import { apiRequest } from '@/api';

interface NavbarProps {
    user: User | null;
    setUser: React.Dispatch<React.SetStateAction<User | null>>;
}

export default function Navbar({ user, setUser }: NavbarProps): ReactElement {
  const [isDropdownOpen, setDropdownOpen] = useState(false);  // State to manage dropdown visibility
  const location = useLocation();
  const navigate = useNavigate();
  const menuRef = useRef<HTMLDivElement>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const isActive = (path: string) =>
    location.pathname === path ? 'border-b-2 border-green-500 text-white' : 'text-gray-400';

  const toggleDropdown = () => {
    setDropdownOpen((prevState) => !prevState);
  };

  const handleLogout = () => {
    apiRequest('/api/v1/auth/logout', { method: 'POST' })
      .then((data) => {
        // Clear user authentication, for example by removing the token or user info from local storage
        localStorage.removeItem(storage_token);
        setUser(null);  // Clear user state
        // Redirect to login paged
        navigate('/login');
      })
      .catch((err) => {
        setErrorMsg('Could not log out. Please try again later.');
      });

  };



  // Function to handle clicking outside the menu
  const handleClickOutside = (e: MouseEvent) => {
    if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
      setDropdownOpen(false);
    }
  };

  // Add event listener for clicks outside of the menu
  useEffect(() => {
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <nav className="flex justify-between items-center p-4 bg-gray-800 border-b border-gray-700">
      <div className="flex gap-6">
        <Link to="/containers" className={`hover:text-white pb-1 ${isActive('/containers')}`}>Containers</Link>
        <Link to="/templates" className={`hover:text-white pb-1 ${isActive('/templates')}`}>Templates</Link>
        <Link to="/images" className={`hover:text-white pb-1 ${isActive('/images')}`}>Images</Link>
      </div>
      <div className="relative">
        <button
          onClick={toggleDropdown}
          className="text-sm font-medium flex items-center gap-2 cursor-pointer"
        >
                    👤 {user?.username}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            className="w-4 h-4"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>

        {/* Dropdown Menu */}
        {isDropdownOpen && (
          <div ref={menuRef} className="absolute right-0 mt-2 w-48 bg-gray-800 text-gray-100 rounded-md shadow-lg border border-gray-600">
            <div className="py-2">
              <button
                onClick={handleLogout}
                className="w-full text-left px-4 py-2 text-sm hover:bg-gray-700"
              >
                                Logout
              </button>
            </div>
          </div>
        )}
      </div>
      {errorMsg && (
        <div className="mb-4 p-3 bg-red-600 text-white rounded-xl">
          {errorMsg}
        </div>
      )}
    </nav>

  );
}
