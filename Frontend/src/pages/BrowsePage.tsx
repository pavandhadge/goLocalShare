
import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import Header from '@/components/Header';
import FileList from '@/components/FileList';
import { useAuth } from '@/context/AuthContext';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AlertTriangle, RefreshCw, Shield } from 'lucide-react';
import { toast } from 'sonner';

interface File {
  name: string;
  path: string;
  isDirectory: boolean;
  size?: number;
  modTime?: string;
}

const BrowsePage: React.FC = () => {
  const { token, isAuthenticated } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [files, setFiles] = useState<File[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Extract path from location
  const getPathFromLocation = () => {
    const path = location.pathname.replace(/^\/browse\/?/, '');
    return path || '/';
  };

  const currentPath = getPathFromLocation();

  // Simulate fetch files (in a real app, this would call the Go server API)
  const fetchFiles = async () => {
    if (!isAuthenticated) {
      setError('Authentication required');
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      // In a real app, this would be a fetch call to the server
      // For now, we'll simulate a response with sample data
      await new Promise(resolve => setTimeout(resolve, 500)); // Simulate network delay

      // Sample data - in a real implementation, this would come from the server
      const sampleFiles: File[] = [
        { name: 'Documents', path: `${currentPath === '/' ? '' : currentPath}/Documents`, isDirectory: true, modTime: new Date().toISOString() },
        { name: 'Images', path: `${currentPath === '/' ? '' : currentPath}/Images`, isDirectory: true, modTime: new Date().toISOString() },
        { name: 'README.md', path: `${currentPath === '/' ? '' : currentPath}/README.md`, isDirectory: false, size: 2048, modTime: new Date().toISOString() },
        { name: 'main.go', path: `${currentPath === '/' ? '' : currentPath}/main.go`, isDirectory: false, size: 4096, modTime: new Date().toISOString() },
        { name: 'screenshot.png', path: `${currentPath === '/' ? '' : currentPath}/screenshot.png`, isDirectory: false, size: 102400, modTime: new Date().toISOString() },
        { name: 'data.zip', path: `${currentPath === '/' ? '' : currentPath}/data.zip`, isDirectory: false, size: 1048576, modTime: new Date().toISOString() },
      ];

      setFiles(sampleFiles);
    } catch (err) {
      console.error('Error fetching files:', err);
      setError('Failed to fetch files. Please try again.');
      toast.error('Error loading files', {
        description: 'Could not retrieve file list from the server',
      });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, [currentPath, isAuthenticated]);

  const handleRefresh = () => {
    fetchFiles();
    toast.success('Refreshed file list');
  };

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex flex-col bg-gray-50">
        <Header />
        <main className="flex-1 container mx-auto px-4 py-8">
          <Card className="max-w-md mx-auto">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-amber-500" />
                Authentication Required
              </CardTitle>
              <CardDescription>
                You need a valid token to browse files
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-4">
                Please obtain an authentication token from the homepage or authenticate with an existing token.
              </p>
            </CardContent>
            <CardFooter className="flex justify-between">
              <Button variant="outline" onClick={() => navigate('/')}>
                Back to Home
              </Button>
              <Button onClick={() => navigate('/auth')}>
                <Shield className="mr-2 h-4 w-4" />
                Authenticate
              </Button>
            </CardFooter>
          </Card>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Header />
      <main className="flex-1 container mx-auto px-4 py-8">
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-2xl font-bold">File Browser</h1>
          <Button variant="outline" size="sm" onClick={handleRefresh}>
            <RefreshCw size={16} className="mr-2" />
            Refresh
          </Button>
        </div>

        {error ? (
          <Card>
            <CardHeader>
              <CardTitle className="text-red-500 flex items-center gap-2">
                <AlertTriangle size={18} />
                Error
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p>{error}</p>
            </CardContent>
            <CardFooter>
              <Button onClick={fetchFiles}>Try Again</Button>
            </CardFooter>
          </Card>
        ) : (
          <FileList 
            files={files} 
            currentPath={currentPath} 
            isLoading={isLoading} 
          />
        )}
      </main>
    </div>
  );
};

export default BrowsePage;
