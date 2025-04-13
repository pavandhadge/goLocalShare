
import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { FileQuestion, Home } from "lucide-react";

const NotFound: React.FC = () => {
  const location = useLocation();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 p-4">
      <div className="max-w-md w-full text-center">
        <div className="mb-8">
          <FileQuestion className="h-16 w-16 text-fileserver-blue mx-auto" />
        </div>
        <h1 className="text-4xl font-bold mb-4">404</h1>
        <p className="text-xl text-muted-foreground mb-6">
          The resource at <span className="font-mono text-sm">{location.pathname}</span> was not found
        </p>
        <Button asChild className="gap-2">
          <Link to="/">
            <Home size={16} />
            Return to Home
          </Link>
        </Button>
      </div>
    </div>
  );
};

export default NotFound;
