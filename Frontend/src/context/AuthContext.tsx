
import React, { createContext, useState, useContext, useEffect } from 'react';
import { toast } from 'sonner';

interface AuthContextType {
  token: string | null;
  setToken: (token: string | null) => void;
  tokenExpiryTime: Date | null;
  isAuthenticated: boolean;
  logOut: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);
  const [tokenExpiryTime, setTokenExpiryTime] = useState<Date | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  useEffect(() => {
    // Check if token exists in localStorage
    const storedToken = localStorage.getItem('authToken');
    const storedExpiryTime = localStorage.getItem('tokenExpiryTime');
    
    if (storedToken && storedExpiryTime) {
      const expiryTime = new Date(storedExpiryTime);
      
      // Check if token is still valid
      if (expiryTime > new Date()) {
        setToken(storedToken);
        setTokenExpiryTime(expiryTime);
        setIsAuthenticated(true);
      } else {
        // Token expired, clean up
        localStorage.removeItem('authToken');
        localStorage.removeItem('tokenExpiryTime');
      }
    }
  }, []);

  useEffect(() => {
    // Set up expiry timer when token changes
    if (token && tokenExpiryTime) {
      const expiryTimer = setTimeout(() => {
        logOut();
        toast.warning("Your authentication token has expired", {
          description: "Please refresh the page to obtain a new token",
        });
      }, tokenExpiryTime.getTime() - new Date().getTime());
      
      return () => clearTimeout(expiryTimer);
    }
  }, [token, tokenExpiryTime]);

  const logOut = () => {
    setToken(null);
    setTokenExpiryTime(null);
    setIsAuthenticated(false);
    localStorage.removeItem('authToken');
    localStorage.removeItem('tokenExpiryTime');
  };

  // Update localStorage when token changes
  useEffect(() => {
    if (token && tokenExpiryTime) {
      localStorage.setItem('authToken', token);
      localStorage.setItem('tokenExpiryTime', tokenExpiryTime.toISOString());
      setIsAuthenticated(true);
    }
  }, [token, tokenExpiryTime]);

  return (
    <AuthContext.Provider value={{ token, setToken, tokenExpiryTime, isAuthenticated, logOut }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
