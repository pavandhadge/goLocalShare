
import React from 'react';
import { useAuth } from '@/context/AuthContext';
import { ChevronRight, Home } from 'lucide-react';
import { buildBrowseUrl } from '@/utils/fileUtils';

interface DirectoryPathProps {
  path: string;
}

const DirectoryPath: React.FC<DirectoryPathProps> = ({ path }) => {
  const { token } = useAuth();

  // Parse the path into components
  const getPathSegments = (path: string) => {
    if (!path || path === '/') return [{ name: 'Home', path: '' }];
    
    // Remove leading and trailing slashes
    const cleanPath = path.replace(/^\/|\/$/g, '');
    
    // Split by slash
    const segments = cleanPath.split('/');
    
    // Build incrementally longer paths
    return [
      { name: 'Home', path: '' },
      ...segments.map((segment, index) => {
        const segmentPath = segments.slice(0, index + 1).join('/');
        return {
          name: segment,
          path: segmentPath,
        };
      }),
    ];
  };

  const segments = getPathSegments(path);

  return (
    <nav className="flex items-center space-x-1 text-sm py-2 mb-4 overflow-x-auto">
      {segments.map((segment, index) => (
        <React.Fragment key={segment.path}>
          <a
            href={buildBrowseUrl(segment.path, token)}
            className={`flex items-center hover:text-primary transition-colors ${
              index === segments.length - 1 ? 'font-semibold text-primary' : 'text-muted-foreground'
            }`}
          >
            {index === 0 ? (
              <Home size={16} className="inline mr-1" />
            ) : null}
            <span>{segment.name}</span>
          </a>
          {index < segments.length - 1 && (
            <ChevronRight size={14} className="text-muted-foreground" />
          )}
        </React.Fragment>
      ))}
    </nav>
  );
};

export default DirectoryPath;
