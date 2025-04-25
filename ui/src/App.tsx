import { useState, useEffect } from "react";
import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import Templates from "./component/templates";
import Containers from "./component/containers";
import Navbar from "./component/navbar";
import Images from "./component/images";
import { type ReactElement } from "react";
import LoginPage from "./component/login";
import { storage_token, UI_VERSION } from "./config";
import { apiRequest } from "./api";

export interface User {
  user_id: string;
  username: string;
}



export function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const token = localStorage.getItem(storage_token);
    if (!token) {
      setLoading(false);
      return;
    }

    apiRequest(`/api/v1/auth/check`)
      .then((data) => {
        setUser({
          user_id: data.data.id,
          username: data.data.username,
        });
      })
      .catch(() => {
        localStorage.removeItem(storage_token);
        setUser(null);
      })
      .finally(() => setLoading(false));
  }, []);

  // PrivateRoute Component
  const PrivateRoute = ({ children }: { children: ReactElement }) => {
    if (!user) {
      // If no user is logged in, redirect to login page
      return <Navigate to="/login" replace />;
    }
    return children; // If the user is authenticated, allow access to the route
  };
  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-white bg-gray-900">
        Checking session...
      </div>
    );
  }
  return (
    <Router>
      <div className="min-h-screen bg-gray-900 text-gray-100 font-sans">
        {/* Navbar */}
        {user && <Navbar user={user} setUser={setUser} />}
        <Routes>
          <Route path="/templates" element={<PrivateRoute><Templates /></PrivateRoute>} />
          <Route path="/images" element={<PrivateRoute><Images /></PrivateRoute>} />
          <Route path="/containers" element={<PrivateRoute><Containers user={user} /></PrivateRoute>} />
          <Route path="/" element={<PrivateRoute><Containers user={user} /></PrivateRoute>} />
          <Route path="/login" element={<LoginPage setUser={setUser} />} />
        </Routes>
        <div className="text-sm text-gray-600 text-center p-4 fixed bottom-0 left-1/2 transform -translate-x-1/2">
          version: {UI_VERSION}
        </div>
      </div>
    </Router>

  );
}

export default App;
