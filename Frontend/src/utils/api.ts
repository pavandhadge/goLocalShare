const API_BASE_URL = 'http://localhost:8090/api';

export interface TokenResponse {
  token: string;
  expiresAt: string;
  expiresIn: number;
  serverTime: string;
}

export interface FileInfo {
  name: string;
  path: string;
  isDirectory: boolean;
  size?: number;
  modTime: string;
  sizeFormatted?: string;
}

export interface FilesResponse {
  files: FileInfo[];
  currentDir: string;
  basePath: string;
}

export interface ErrorResponse {
  error: string;
  message: string;
  code: number;
}

class APIError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string
  ) {
    super(message);
    this.name = 'APIError';
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    let errorData: ErrorResponse | null = null;
    try {
      errorData = await response.json();
    } catch {
      // If we can't parse the error response, use a generic message
    }

    throw new APIError(
      errorData?.message || `HTTP ${response.status}: ${response.statusText}`,
      response.status,
      errorData?.error
    );
  }

  return response.json();
}

export async function generateToken(): Promise<TokenResponse> {
  return apiRequest<TokenResponse>('/auth/token');
}

export async function getFiles(path: string = '/'): Promise<FilesResponse> {
  const token = localStorage.getItem('gofs_tkn_4a7f');
  if (!token) {
    throw new APIError('No authentication token found', 401);
  }

  return apiRequest<FilesResponse>(`/files${path}`, {
    headers: {
      'X-Auth-Token': token,
    },
  });
}

export function getDownloadUrl(path: string): string {
  const token = localStorage.getItem('gofs_tkn_4a7f');
  if (!token) {
    throw new APIError('No authentication token found', 401);
  }

  return `${API_BASE_URL}/download/${path}?token=${token}`;
}

export function validateToken(token: string): boolean {
  // Basic validation - in a real app, you might want to verify with the server
  return token.length === 64; // 32 bytes = 64 hex characters
}

export { APIError }; 