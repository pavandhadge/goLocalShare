
import React, { useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { toast } from 'sonner';
import { Copy, CheckCircle2, AlertTriangle, Clock } from 'lucide-react';

const Authentication = () => {
  const { token, setToken, tokenExpiryTime, isAuthenticated, logOut } = useAuth();
  const [newToken, setNewToken] = useState('');
  const [copied, setCopied] = useState(false);

  // Calculate time remaining until token expiry
  const getTimeRemaining = () => {
    if (!tokenExpiryTime) return '';
    
    const now = new Date();
    const remaining = tokenExpiryTime.getTime() - now.getTime();
    
    if (remaining <= 0) return 'Expired';
    
    const minutes = Math.floor((remaining / 1000 / 60) % 60);
    const hours = Math.floor((remaining / 1000 / 60 / 60) % 24);
    
    return `${hours}h ${minutes}m remaining`;
  };

  const handleCopyToken = () => {
    if (!token) return;
    
    navigator.clipboard.writeText(token)
      .then(() => {
        setCopied(true);
        toast.success('Token copied to clipboard');
        setTimeout(() => setCopied(false), 2000);
      })
      .catch(() => {
        toast.error('Failed to copy token');
      });
  };

  const handleSetNewToken = () => {
    if (!newToken.trim()) {
      toast.error('Please enter a token');
      return;
    }

    // Set token with 1 hour expiry (as per the Go server behavior)
    const expiryTime = new Date();
    expiryTime.setHours(expiryTime.getHours() + 1);
    
    setToken(newToken.trim());
    setNewToken('');
    toast.success('Authentication token set', {
      description: 'Valid for 1 hour'
    });
  };

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Authentication</CardTitle>
        <CardDescription>
          Access tokens are valid for 1 hour from the time they are issued
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {isAuthenticated ? (
          <div className="space-y-4">
            <div className="flex flex-col space-y-2">
              <Label htmlFor="token">Current Token</Label>
              <div className="flex items-center">
                <Input 
                  id="token"
                  value={token || ''}
                  readOnly
                  className="font-mono text-sm pr-10"
                />
                <Button 
                  variant="ghost" 
                  size="sm" 
                  className="ml-[-2.5rem]" 
                  onClick={handleCopyToken}
                >
                  {copied ? <CheckCircle2 size={16} /> : <Copy size={16} />}
                </Button>
              </div>
            </div>
            
            <div className="flex items-center gap-2 text-sm">
              <Clock size={16} className="text-blue-500" />
              <span>{getTimeRemaining()}</span>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center gap-2 p-3 border rounded-md bg-yellow-50 text-yellow-800">
              <AlertTriangle size={18} />
              <p className="text-sm">You need a valid token to access protected resources</p>
            </div>
            
            <div className="flex flex-col space-y-2">
              <Label htmlFor="newToken">Enter Token</Label>
              <div className="flex gap-2">
                <Input
                  id="newToken"
                  value={newToken}
                  onChange={(e) => setNewToken(e.target.value)}
                  placeholder="Paste your token here"
                  className="font-mono"
                />
                <Button onClick={handleSetNewToken}>Set</Button>
              </div>
            </div>
          </div>
        )}
      </CardContent>
      {isAuthenticated && (
        <CardFooter>
          <Button variant="outline" className="w-full" onClick={logOut}>
            Clear Token
          </Button>
        </CardFooter>
      )}
    </Card>
  );
};

export default Authentication;
