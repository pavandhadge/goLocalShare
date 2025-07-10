package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if !IsDirMode {
		fileInfo, err := os.Stat(BasePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		fileName := filepath.Base(BasePath)

		tmpl := `<!DOCTYPE html>
<html>
<head>
  <title>Download File</title>
  <style>
    body { font-family: Arial, sans-serif; background: #f8fafc; color: #1e293b; }
    .container { max-width: 500px; margin: 3rem auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 2rem; }
    .file-info { margin-bottom: 1.5rem; }
    .download-btn { background: #2563eb; color: #fff; border: none; border-radius: 4px; padding: 0.7rem 1.2rem; font-size: 1rem; cursor: pointer; }
    .download-btn:hover { background: #1d4ed8; }
    .error { color: #dc2626; margin-bottom: 1rem; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Download File</h2>
    <div class="file-info">
      <strong>Name:</strong> %s<br>
      <strong>Size:</strong> %.2f MB<br>
      <strong>Last Modified:</strong> %s
    </div>
    <div id="error-msg" class="error" style="display:none;"></div>
    <button class="download-btn" onclick="downloadFile()">Download</button>
    <a href="/" style="margin-left:1rem;">Back to Home</a>
  </div>
  <script>
    function downloadFile() {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        window.location.href = '/token';
        return;
      }
      fetch('/api/download/%s', {
        headers: { 'X-Auth-Token': token }
      }).then(resp => {
        if (!resp.ok) {
          document.getElementById('error-msg').innerText = 'Invalid or expired token.';
          document.getElementById('error-msg').style.display = 'block';
          if (resp.status === 401) setTimeout(() => window.location.href = '/token', 1200);
          return;
        }
        return resp.blob();
      }).then(blob => {
        if (!blob) return;
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = '%s';
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
      }).catch(() => {
        document.getElementById('error-msg').innerText = 'Network error.';
        document.getElementById('error-msg').style.display = 'block';
      });
    }
    // Check for token in localStorage
    if (!localStorage.getItem('accessToken')) {
      window.location.href = '/token';
    }
  </script>
</body>
</html>`

		page := fmt.Sprintf(tmpl, fileName, float64(fileInfo.Size())/1024/1024, fileInfo.ModTime().Format("2006-01-02 15:04:05"), fileName, fileName)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
		return
	}

	// Directory mode
	requestedPath := strings.TrimPrefix(r.URL.Path, "/download/")
	cleanPath := filepath.Join(BasePath, requestedPath)
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		http.Redirect(w, r, "/browse/"+requestedPath, http.StatusSeeOther)
		return
	}

	fileName := filepath.Base(cleanPath)
	tmpl := `<!DOCTYPE html>
<html>
<head>
  <title>Download File</title>
  <style>
    body { font-family: Arial, sans-serif; background: #f8fafc; color: #1e293b; }
    .container { max-width: 500px; margin: 3rem auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 2rem; }
    .file-info { margin-bottom: 1.5rem; }
    .download-btn { background: #2563eb; color: #fff; border: none; border-radius: 4px; padding: 0.7rem 1.2rem; font-size: 1rem; cursor: pointer; }
    .download-btn:hover { background: #1d4ed8; }
    .error { color: #dc2626; margin-bottom: 1rem; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Download File</h2>
    <div class="file-info">
      <strong>Name:</strong> %s<br>
      <strong>Size:</strong> %.2f MB<br>
      <strong>Last Modified:</strong> %s
    </div>
    <div id="error-msg" class="error" style="display:none;"></div>
    <button class="download-btn" onclick="downloadFile()">Download</button>
    <a href="/browse" style="margin-left:1rem;">Back to Browse</a>
  </div>
  <script>
    function downloadFile() {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        window.location.href = '/token';
        return;
      }
      fetch('/api/download/%s', {
        headers: { 'X-Auth-Token': token }
      }).then(resp => {
        if (!resp.ok) {
          document.getElementById('error-msg').innerText = 'Invalid or expired token.';
          document.getElementById('error-msg').style.display = 'block';
          if (resp.status === 401) setTimeout(() => window.location.href = '/token', 1200);
          return;
        }
        return resp.blob();
      }).then(blob => {
        if (!blob) return;
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = '%s';
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
      }).catch(() => {
        document.getElementById('error-msg').innerText = 'Network error.';
        document.getElementById('error-msg').style.display = 'block';
      });
    }
    // Check for token in localStorage
    if (!localStorage.getItem('accessToken')) {
      window.location.href = '/token';
    }
  </script>
</body>
</html>`

	page := fmt.Sprintf(tmpl, fileName, float64(fileInfo.Size())/1024/1024, fileInfo.ModTime().Format("2006-01-02 15:04:05"), requestedPath, fileName)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(page))
} 