
import React, { createContext, useState, useContext, useEffect } from 'react';
import { toast } from 'sonner';
import { generateToken, APIError } from '@/utils/api';

interface AuthContextType {
  token: string | null;
  setToken: (token: string | null) => void;
  tokenExpiryTime: Date | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  logOut: () => void;
  generateNewToken: () => Promise<void>;
  refreshToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);
  const [tokenExpiryTime, setTokenExpiryTime] = useState<Date | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  useEffect(() => {
    // Check if token exists in localStorage
    const storedToken = localStorage.getItem('gofs_tkn_4a7f');
    const storedExpiryTime = localStorage.getItem('tokenExpiryTime');
    
    if (storedToken && storedExpiryTime) {
      const expiryTime = new Date(storedExpiryTime);
      
      // Check if token is still valid (with 5 minute buffer)
      const bufferTime = new Date(expiryTime.getTime() - 5 * 60 * 1000);
      if (bufferTime > new Date()) {
        setToken(storedToken);
        setTokenExpiryTime(expiryTime);
        setIsAuthenticated(true);
      } else {
        // Token expired or will expire soon, clean up
        localStorage.removeItem('gofs_tkn_4a7f');
        localStorage.removeItem('tokenExpiryTime');
        toast.warning("Your authentication token has expired", {
          description: "Please generate a new token to continue",
        });
      }
    }
  }, []);

  useEffect(() => {
    // Set up expiry timer when token changes
    if (token && tokenExpiryTime) {
      const timeUntilExpiry = tokenExpiryTime.getTime() - new Date().getTime();
      
      // Show warning 5 minutes before expiry
      const warningTime = timeUntilExpiry - (5 * 60 * 1000);
      if (warningTime > 0) {
        setTimeout(() => {
          toast.warning("Your token will expire soon", {
            description: "Consider refreshing your token to avoid interruption",
          });
        }, warningTime);
      }

      // Set up expiry timer
      const expiryTimer = setTimeout(() => {
        logOut();
        toast.error("Your authentication token has expired", {
          description: "Please generate a new token to continue",
        });
      }, timeUntilExpiry);
      
      return () => clearTimeout(expiryTimer);
    }
  }, [token, tokenExpiryTime]);

  const logOut = () => {
    setToken(null);
    setTokenExpiryTime(null);
    setIsAuthenticated(false);
    localStorage.removeItem('gofs_tkn_4a7f');
    localStorage.removeItem('tokenExpiryTime');
  };

  const generateNewToken = async () => {
    setIsLoading(true);
    try {
      const response = await generateToken();
      const expiryTime = new Date(response.expiresAt);
      
      setToken(response.token);
      setTokenExpiryTime(expiryTime);
      setIsAuthenticated(true);
      
      localStorage.setItem('gofs_tkn_4a7f', response.token);
      localStorage.setItem('tokenExpiryTime', response.expiresAt);
      
      toast.success("Authentication token generated successfully", {
        description: `Token expires in ${Math.floor(response.expiresIn / 60)} minutes`,
      });
    } catch (error) {
      console.error('Error generating token:', error);
      if (error instanceof APIError) {
        toast.error("Failed to generate token", {
          description: error.message,
        });
      } else {
        toast.error("Failed to generate token", {
          description: "Please check your connection and try again",
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const refreshToken = async () => {
    await generateNewToken();
  };

  return (
    <AuthContext.Provider value={{ 
      token, 
      setToken, 
      tokenExpiryTime, 
      isAuthenticated, 
      isLoading,
      logOut, 
      generateNewToken,
      refreshToken
    }}>
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
