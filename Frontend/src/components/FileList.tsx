
import React from 'react';
import FileItem from './FileItem';
import DirectoryPath from './DirectoryPath';
import { useAuth } from '@/context/AuthContext';
import { AlertTriangle, FolderOpen } from 'lucide-react';

interface File {
  name: string;
  path: string;
  isDirectory: boolean;
  size?: number;
  modTime?: string;
}

interface FileListProps {
  files: File[];
  currentPath: string;
  isLoading?: boolean;
}

const FileList: React.FC<FileListProps> = ({ files, currentPath, isLoading = false }) => {
  const { isAuthenticated } = useAuth();

  if (isLoading) {
    return (
      <div className="space-y-2">
        <div className="h-8 bg-muted animate-pulse rounded-md"></div>
        {[...Array(5)].map((_, i) => (
          <div key={i} className="h-16 bg-muted animate-pulse rounded-md"></div>
        ))}
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="text-center p-8 border rounded-md bg-muted">
        <AlertTriangle className="mx-auto h-12 w-12 text-amber-500 mb-4" />
        <h3 className="text-lg font-medium mb-2">Authentication Required</h3>
        <p className="text-muted-foreground mb-4">
          Please obtain and enter a valid token to view files
        </p>
      </div>
    );
  }

  if (files.length === 0) {
    return (
      <div>
        <DirectoryPath path={currentPath} />
        <div className="text-center p-8 border rounded-md">
          <FolderOpen className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
          <h3 className="text-lg font-medium mb-2">Empty Directory</h3>
          <p className="text-muted-foreground">
            This directory does not contain any files or folders
          </p>
        </div>
      </div>
    );
  }

  // Separate directories and files
  const directories = files.filter(file => file.isDirectory);
  const regularFiles = files.filter(file => !file.isDirectory);
  
  // Sort directories and files by name
  const sortedDirectories = [...directories].sort((a, b) => a.name.localeCompare(b.name));
  const sortedFiles = [...regularFiles].sort((a, b) => a.name.localeCompare(b.name));
  
  // Combine sorted lists with directories first
  const sortedItems = [...sortedDirectories, ...sortedFiles];

  return (
    <div>
      <DirectoryPath path={currentPath} />
      <div className="space-y-2">
        {sortedItems.map((file) => (
          <FileItem
            key={file.path}
            name={file.name}
            path={file.path}
            isDirectory={file.isDirectory}
            size={file.size}
            modTime={file.modTime}
          />
        ))}
      </div>
    </div>
  );
};

export default FileList;
