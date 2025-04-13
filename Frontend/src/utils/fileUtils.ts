
// Define file types and their corresponding icons
export const FILE_TYPES = {
  PDF: ['pdf'],
  IMAGE: ['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp', 'bmp', 'ico'],
  AUDIO: ['mp3', 'wav', 'ogg', 'flac', 'aac', 'm4a'],
  VIDEO: ['mp4', 'webm', 'avi', 'mov', 'mkv', 'flv', 'wmv'],
  ARCHIVE: ['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz'],
  DOCUMENT: ['doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'rtf', 'csv', 'json', 'md', 'html', 'css', 'js', 'ts', 'jsx', 'tsx'],
  CODE: ['go', 'java', 'py', 'rb', 'php', 'c', 'cpp', 'h', 'hpp', 'cs', 'swift', 'kt', 'sh', 'rs'],
};

// Get the type of a file based on its extension
export function getFileType(fileName: string): string {
  const extension = fileName.split('.').pop()?.toLowerCase() || '';
  
  if (FILE_TYPES.PDF.includes(extension)) return 'pdf';
  if (FILE_TYPES.IMAGE.includes(extension)) return 'image';
  if (FILE_TYPES.AUDIO.includes(extension)) return 'audio';
  if (FILE_TYPES.VIDEO.includes(extension)) return 'video';
  if (FILE_TYPES.ARCHIVE.includes(extension)) return 'archive';
  if (FILE_TYPES.DOCUMENT.includes(extension)) return 'document';
  if (FILE_TYPES.CODE.includes(extension)) return 'code';
  
  return 'other';
}

// Format file size to human-readable format
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format date to a readable format
export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleString();
}

// Create a URL with token
export function createTokenizedUrl(baseUrl: string, token: string): string {
  const url = new URL(baseUrl, window.location.origin);
  url.searchParams.append('token', token);
  return url.toString();
}

// Function to build a download URL
export function buildDownloadUrl(path: string, token: string | null): string {
  if (!token) return '';
  
  // Ensure path is properly encoded
  const encodedPath = encodeURIComponent(path).replace(/%2F/g, '/');
  return `/download/${encodedPath}?token=${token}`;
}

// Function to build a browse URL for directories
export function buildBrowseUrl(path: string, token: string | null): string {
  if (!token) return '';
  
  // Ensure path is properly encoded and ends with a slash for directories
  let encodedPath = encodeURIComponent(path).replace(/%2F/g, '/');
  if (!encodedPath.endsWith('/')) {
    encodedPath += '/';
  }
  return `/browse/${encodedPath}?token=${token}`;
}
