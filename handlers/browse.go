package handlers

import (
	"net/http"
)

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
  <title>Browse Files</title>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
    .current-dir { margin-bottom: 20px; padding: 10px; background: #f0f0f0; border-radius: 4px; }
    .file-list { margin: 20px 0; }
    .file-item { padding: 8px 0; border-bottom: 1px solid #eee; }
    .file-icon { margin-right: 10px; }
    .file-size { float: right; color: #666; }
    .empty-state { color: #666; font-style: italic; }
    .back-link { display: inline-block; margin-top: 20px; color: #007bff; text-decoration: none; }
    .back-link:hover { text-decoration: underline; }
    .error { color: #dc2626; margin-bottom: 1rem; }
  </style>
</head>
<body>
  <h1>Available Files</h1>
  <div class="current-dir" id="current-dir"></div>
  <div class="file-list" id="file-list"></div>
  <div class="error" id="error-msg" style="display:none;"></div>
  <a href="/" class="back-link">Back to Home</a>
  <script>
    console.log('Browse page loaded');
    console.log('Token in localStorage:', localStorage.getItem('accessToken'));
    function getToken() {
      return localStorage.getItem('accessToken');
    }
    function getPath() {
      const hash = decodeURIComponent(window.location.hash.replace(/^#/, ''));
      return hash || '/';
    }
    function renderFiles(files) {
      const list = document.getElementById('file-list');
      list.innerHTML = '';
      if (!files.length) {
        list.innerHTML = '<div class="empty-state">This directory is empty</div>';
        return;
      }
      files.forEach(function(file) {
        var div = document.createElement('div');
        div.className = 'file-item';
        var icon = file.isDirectory ? '📁' : '📄';
        var link;
        if (file.isDirectory) {
          link = '<a href="#' + file.path + '">' + file.name + '/</a>';
        } else {
          link = '<a href="/download/' + file.path + '" target="_blank">' + file.name + '</a> <span class="file-size">' + (file.sizeFormatted || '') + '</span>';
        }
        div.innerHTML = '<span class="file-icon">' + icon + '</span>' + link;
        list.appendChild(div);
      });
    }
    function fetchFiles() {
      var token = getToken();
      console.log('fetchFiles: token =', token);
      if (!token) {
        console.log('No token, redirecting to /token');
        window.location.href = '/token';
        return;
      }
      var path = getPath();
      console.log('Fetching /api/files' + path);
      fetch('/api/files' + path, {
        headers: { 'X-Auth-Token': token }
      })
      .then(function(resp) {
        if (!resp.ok) throw new Error('Invalid or expired token.');
        return resp.json();
      })
      .then(function(data) {
        console.log('Files loaded:', data);
        document.getElementById('current-dir').innerHTML = 'Current directory: <strong>' + data.currentDir + '</strong>';
        renderFiles(data.files);
        document.getElementById('error-msg').style.display = 'none';
      })
      .catch(function(err) {
        console.log('Token invalid in browse, removing from localStorage');
        document.getElementById('error-msg').innerText = err.message || 'Failed to load files.';
        document.getElementById('error-msg').style.display = 'block';
        if (err.message && err.message.includes('token')) {
          localStorage.removeItem('accessToken'); // Clear invalid token
          setTimeout(function() { window.location.href = '/token'; }, 1200);
        }
      });
    }
    window.addEventListener('hashchange', fetchFiles);
    window.addEventListener('DOMContentLoaded', fetchFiles);
  </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
} 