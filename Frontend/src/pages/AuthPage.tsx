
import React from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '@/components/Header';
import Authentication from '@/components/Authentication';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';

const AuthPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Header />
      <main className="flex-1 container mx-auto px-4 py-8">
        <div className="mb-6">
          <Button 
            variant="ghost" 
            className="mb-4"
            onClick={() => navigate(-1)}
          >
            <ArrowLeft size={16} className="mr-2" />
            Back
          </Button>
          <h1 className="text-2xl font-bold">Authentication</h1>
          <p className="text-muted-foreground">
            Manage your authentication token
          </p>
        </div>

        <div className="max-w-lg mx-auto">
          <Authentication />
          
          <div className="mt-8 p-4 bg-accent rounded-md text-sm space-y-2">
            <h3 className="font-medium mb-2">About Authentication</h3>
            <p>Tokens are valid for 1 hour from the time they are issued.</p>
            <p>You can get a token by visiting the server's homepage in your browser.</p>
            <p>For API access, include the token either as a URL parameter (?token=...) or in the X-Auth-Token header.</p>
          </div>
        </div>
      </main>
    </div>
  );
};

export default AuthPage;
