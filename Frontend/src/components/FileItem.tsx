
import React from 'react';
import { useAuth } from '@/context/AuthContext';
import { getFileType, buildDownloadUrl, buildBrowseUrl } from '@/utils/fileUtils';
import { FileIcon, FolderIcon, FileTextIcon, FileImageIcon, FileAudioIcon, 
         FileVideoIcon, ArchiveIcon, FileCodeIcon, FileType2Icon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { DownloadIcon, ChevronRightIcon } from 'lucide-react';
import { toast } from 'sonner';

interface FileItemProps {
  name: string;
  path: string;
  isDirectory: boolean;
  size?: number;
  modTime?: string;
}

const FileItem: React.FC<FileItemProps> = ({
  name,
  path,
  isDirectory,
  size,
  modTime,
}) => {
  const { token, isAuthenticated } = useAuth();

  const handleDownload = (e: React.MouseEvent) => {
    if (!isAuthenticated) {
      e.preventDefault();
      toast.error('Authentication required', {
        description: 'Please obtain a token to download files',
      });
    }
  };

  const handleNavigate = (e: React.MouseEvent) => {
    if (!isAuthenticated) {
      e.preventDefault();
      toast.error('Authentication required', {
        description: 'Please obtain a token to browse directories',
      });
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const formatSize = (bytes?: number) => {
    if (bytes === undefined) return '-';
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const renderFileIcon = () => {
    if (isDirectory) {
      return <FolderIcon className="file-icon file-icon-folder" />;
    }
    
    const fileType = getFileType(name);
    
    switch (fileType) {
      case 'pdf':
        return <FileIcon className="file-icon file-icon-pdf" />;
      case 'image':
        return <FileImageIcon className="file-icon file-icon-image" />;
      case 'audio':
        return <FileAudioIcon className="file-icon file-icon-audio" />;
      case 'video':
        return <FileVideoIcon className="file-icon file-icon-video" />;
      case 'archive':
        return <ArchiveIcon className="file-icon file-icon-archive" />;
      case 'document':
        return <FileTextIcon className="file-icon file-icon-document" />;
      case 'code':
        return <FileCodeIcon className="file-icon file-icon-document" />;
      default:
        return <FileType2Icon className="file-icon file-icon-other" />;
    }
  };

  return (
    <div className="flex items-center p-3 border rounded-md hover:bg-accent group">
      <div className="mr-3">
        {renderFileIcon()}
      </div>
      <div className="flex-1 min-w-0">
        <div className="font-medium truncate">{name}</div>
        <div className="text-xs text-muted-foreground flex space-x-2">
          {!isDirectory && <span>{formatSize(size)}</span>}
          {modTime && <span>{formatDate(modTime)}</span>}
        </div>
      </div>
      <div className="ml-2">
        {isDirectory ? (
          <Button 
            variant="ghost" 
            size="sm" 
            asChild
            onClick={handleNavigate}
            className="opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <a href={buildBrowseUrl(path, token)}>
              <ChevronRightIcon size={18} />
            </a>
          </Button>
        ) : (
          <Button 
            variant="ghost" 
            size="sm" 
            asChild
            onClick={handleDownload}
            className="opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <a href={buildDownloadUrl(path, token)} download={name}>
              <DownloadIcon size={18} />
            </a>
          </Button>
        )}
      </div>
    </div>
  );
};

export default FileItem;
