
import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Clock, Lock, Server } from 'lucide-react';

const Header: React.FC = () => {
  const { isAuthenticated, tokenExpiryTime } = useAuth();

  // Calculate minutes remaining until token expiry
  const getMinutesRemaining = () => {
    if (!tokenExpiryTime) return 0;
    
    const now = new Date();
    const remaining = tokenExpiryTime.getTime() - now.getTime();
    
    return Math.max(0, Math.floor(remaining / 1000 / 60));
  };

  const minutesRemaining = getMinutesRemaining();

  return (
    <header className="border-b bg-white">
      <div className="container mx-auto px-4 py-3">
        <div className="flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2">
            <Server className="h-6 w-6 text-fileserver-blue" />
            <h1 className="text-xl font-semibold">Secure File Server</h1>
          </Link>
          
          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <>
                <div className="flex items-center gap-1">
                  <Clock size={16} className="text-muted-foreground" />
                  <Badge variant={minutesRemaining < 10 ? "destructive" : "secondary"}>
                    {minutesRemaining} min
                  </Badge>
                </div>
                <Button asChild size="sm" variant="outline">
                  <Link to="/auth">Manage Token</Link>
                </Button>
              </>
            ) : (
              <Button asChild size="sm">
                <Link to="/auth">
                  <Lock size={16} className="mr-1" />
                  Authenticate
                </Link>
              </Button>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
