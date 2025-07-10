package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/pavandhadge/goFileShare/auth"
)

// BasePath and IsDirMode are used by all handlers. Set them at server startup.
var (
	BasePath  string
	IsDirMode bool
)

var sessionDuration time.Duration

// Init sets the base path and directory mode for handlers.
func Init(basePath string, isDirMode bool, duration ...time.Duration) {
	BasePath = basePath
	IsDirMode = isDirMode
	if len(duration) > 0 {
		sessionDuration = duration[0]
	}
}

func TokenEntryPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("token-entry").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Enter Access Token - goFileServer</title>
  <link rel="icon" href="https://fav.farm/🔑" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap" rel="stylesheet">
  <style>
    body { font-family: 'Inter', sans-serif; background: #f8fafc; color: #1e293b; margin: 0; }
    .container { max-width: 420px; margin: 3.5rem auto; background: #fff; border-radius: 14px; box-shadow: 0 4px 24px #0002; padding: 2.5rem 2rem; }
    .title { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.7rem; }
    .instructions { color: #475569; font-size: 1.08rem; margin-bottom: 1.5rem; line-height: 1.6; }
    .token-input { width: 100%; padding: 0.8rem; font-size: 1.08rem; border: 1px solid #cbd5e1; border-radius: 6px; margin-bottom: 1.2rem; background: #f1f5f9; transition: border 0.15s; }
    .token-input:focus { border: 1.5px solid #2563eb; outline: none; }
    .submit-btn { background: #2563eb; color: #fff; border: none; border-radius: 6px; padding: 0.8rem 1.2rem; font-size: 1.08rem; cursor: pointer; width: 100%; font-weight: 500; transition: background 0.15s; }
    .submit-btn:hover, .submit-btn:focus { background: #1d4ed8; }
    .error { color: #dc2626; margin-bottom: 1rem; font-size: 1.02rem; }
    @media (max-width: 600px) { .container { padding: 1.1rem 0.5rem; } }
  </style>
</head>
<body>
  <div class="container">
    <div class="title">Enter Access Token</div>
    <div class="instructions">
      <ol style="margin:0 0 1.2rem 1.2rem;padding:0;">
        <li>Paste the access token you received from the file owner below.</li>
        <li>Click <b>Continue</b> to unlock the shared files or folders.</li>
        <li>If you do not have a token, please contact the owner.</li>
      </ol>
      <span style="color:#64748b;font-size:0.98rem;">Your token is only used to verify your access and is never shared.</span>
    </div>
    <form id="token-form" autocomplete="off">
      <input type="text" id="token-input" class="token-input" placeholder="Paste your access token here" required autofocus />
      <div id="error-msg" class="error" style="display:none;"></div>
      <button type="submit" class="submit-btn">Continue</button>
    </form>
  </div>
  <script>
    function checkAndRedirect(token) {
      fetch('/api/files', { headers: { 'X-Auth-Token': token } })
        .then(resp => {
          if (!resp.ok) throw new Error('Invalid or expired token.');
          return resp.json();
        })
        .then(data => {
          localStorage.setItem('gofs_tkn_4a7f', token);
          if (data.baseInfo && data.baseInfo.isDirectory) {
            window.location.href = '/browse?token=' + encodeURIComponent(token);
          } else if (data.baseInfo && !data.baseInfo.isDirectory) {
            window.location.href = '/file?file=' + encodeURIComponent(data.baseInfo.name) + '&token=' + encodeURIComponent(token);
          } else {
            document.getElementById('error-msg').innerText = 'Unexpected server response.';
            document.getElementById('error-msg').style.display = 'block';
          }
        })
        .catch(() => {
          localStorage.removeItem('gofs_tkn_4a7f');
          document.getElementById('error-msg').innerText = 'Invalid or expired token.';
          document.getElementById('error-msg').style.display = 'block';
          document.getElementById('token-input').focus();
        });
    }
    document.getElementById('token-form').onsubmit = function(e) {
      e.preventDefault();
      const token = document.getElementById('token-input').value.trim();
      if (!token) return;
      checkAndRedirect(token);
    };
  </script>
</body>
</html>
`))
	tmpl.Execute(w, nil)
}

// responseWriterCapture is used to capture HTTP responses for internal handler calls
type responseWriterCapture struct {
	headersWritten bool
	header        http.Header
	body          strings.Builder
	status        int
}
func (rw *responseWriterCapture) Header() http.Header { return rw.header }
func (rw *responseWriterCapture) Write(b []byte) (int, error) { return rw.body.WriteString(string(b)) }
func (rw *responseWriterCapture) WriteHeader(statusCode int) { rw.status = statusCode; rw.headersWritten = true }

// /file handler: show image/preview and download button
func FilePageHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	token := r.URL.Query().Get("token")
	if file == "" || token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing file or token"))
		return
	}
	file = strings.TrimPrefix(file, "/")
	ext := strings.ToLower(filepath.Ext(file))
	isImage := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".svg" || ext == ".webp"
	isText := ext == ".txt" || ext == ".md" || ext == ".go" || ext == ".js" || ext == ".ts" || ext == ".json" || ext == ".css" || ext == ".html" || ext == ".py" || ext == ".sh" || ext == ".c" || ext == ".cpp" || ext == ".h" || ext == ".java" || ext == ".rs" || ext == ".toml" || ext == ".yaml" || ext == ".yml"
	var previewHTML string
	if isImage {
		previewHTML = `<img src="/api/download/` + file + `?token=` + token + `" alt="Preview" style="max-width:100%;border-radius:10px;margin-bottom:1.5rem;box-shadow:0 2px 8px #0001;" />`
	} else if isText {
		path := filepath.Join(BasePath, file)
		b, err := os.ReadFile(path)
		if err != nil {
			previewHTML = `<div style='color:#dc2626;'>Could not read file.</div>`
		} else {
			max := 100 * 1024
			if len(b) > max {
				b = b[:max]
				previewHTML = `<pre style="background:#f1f5f9;padding:1.2rem 1rem 1.2rem 1.2rem;border-radius:8px;max-height:500px;overflow:auto;text-align:left;font-size:1.04rem;line-height:1.6;font-family:monospace;box-shadow:0 1px 4px #0001;">` + template.HTMLEscapeString(string(b)) + `\n... (truncated)</pre>`
			} else {
				previewHTML = `<pre style="background:#f1f5f9;padding:1.2rem 1rem 1.2rem 1.2rem;border-radius:8px;max-height:500px;overflow:auto;text-align:left;font-size:1.04rem;line-height:1.6;font-family:monospace;box-shadow:0 1px 4px #0001;">` + template.HTMLEscapeString(string(b)) + `</pre>`
			}
		}
	} else {
		previewHTML = `<div style='font-size:3rem;margin-bottom:1.5rem;'>&#128196;</div><div style='color:#64748b;'>No preview available</div>`
	}
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Preview: ` + file + ` - goFileServer</title>
  <link rel="icon" href="https://fav.farm/📁" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap" rel="stylesheet">
  <style>
    body { font-family: 'Inter', sans-serif; background: #f8fafc; color: #1e293b; margin: 0; }
    .container { max-width: 800px; margin: 2.5rem auto; background: #fff; border-radius: 14px; box-shadow: 0 4px 24px #0002; padding: 2.5rem 2rem; text-align: left; }
    .header { display: flex; align-items: center; margin-bottom: 2rem; }
    .back-link { color: #2563eb; text-decoration: none; font-size: 1.08rem; font-weight: 600; margin-right: 1.2rem; transition: color 0.15s; }
    .back-link:hover { color: #1d4ed8; }
    .file-title { font-size: 1.5rem; font-weight: 600; flex: 1; }
    .download-btn { background: #2563eb; color: #fff; border: none; border-radius: 6px; padding: 0.7rem 1.4rem; font-size: 1.08rem; cursor: pointer; margin-top: 1.5rem; font-weight: 500; transition: background 0.15s; }
    .download-btn:hover, .download-btn:focus { background: #1d4ed8; }
    @media (max-width: 600px) {
      .container { padding: 1.1rem 0.5rem; }
      .file-title { font-size: 1.1rem; }
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <a href="javascript:history.back()" class="back-link" title="Back to directory">&larr; Back</a>
      <span class="file-title">Preview: ` + file + `</span>
    </div>
    ` + previewHTML + `
    <a href="/api/download/` + file + `?token=` + token + `" download class="download-btn" title="Download this file">Download File</a>
  </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// Improved BrowsePageHandler UI
func BrowsePageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing token"))
		return
	}
	if !auth.ValidateToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid or expired token"))
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}
	req, _ := http.NewRequest("GET", "/api/files"+path, nil)
	req.Header.Set("X-Auth-Token", token)
	rw := &responseWriterCapture{header: http.Header{}}
	FilesAPIHandler(rw, req)
	if rw.status != 0 && rw.status != 200 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid or expired token"))
		return
	}
	var resp struct {
		Files []struct {
			Name string `json:"name"`
			Path string `json:"path"`
			IsDirectory bool `json:"isDirectory"`
		} `json:"files"`
	}
	json.Unmarshal([]byte(rw.body.String()), &resp)
	parentPath := "/"
	if path != "/" {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) > 1 {
			parentPath = "/" + strings.Join(parts[:len(parts)-1], "/")
		} else {
			parentPath = "/"
		}
	}
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Browse Files - goFileServer</title>
  <link rel="icon" href="https://fav.farm/📁" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap" rel="stylesheet">
  <style>
    body { font-family: 'Inter', sans-serif; background: #f8fafc; color: #1e293b; margin: 0; }
    .container { max-width: 950px; margin: 2.5rem auto; background: #fff; border-radius: 14px; box-shadow: 0 4px 24px #0002; padding: 2.5rem 2rem; }
    .header { position: sticky; top: 0; background: #fff; z-index: 2; padding-bottom: 1.2rem; margin-bottom: 1.2rem; border-bottom: 1px solid #e2e8f0; display: flex; align-items: center; }
    .path { font-size: 1.08rem; color: #64748b; margin-left: 0.7rem; flex: 1; white-space: nowrap; overflow-x: auto; }
    .parent-link { color: #2563eb; text-decoration: none; font-size: 1.05rem; font-weight: 600; margin-right: 1.2rem; transition: color 0.15s; }
    .parent-link:hover { color: #1d4ed8; }
    .file-table { width: 100%; border-collapse: collapse; margin-top: 1.5rem; }
    .file-table th, .file-table td { padding: 1rem 0.7rem; text-align: left; }
    .file-table th { background: #f1f5f9; font-weight: 600; color: #475569; border-bottom: 2px solid #e2e8f0; }
    .file-table tr { transition: background 0.13s; }
    .file-table tr:hover { background: #f1f5f9; }
    .icon { font-size: 1.3rem; margin-right: 0.7rem; vertical-align: middle; }
    .action-btn { background: #2563eb; color: #fff; border: none; border-radius: 5px; padding: 0.4rem 1rem; font-size: 1rem; cursor: pointer; margin-right: 0.5rem; font-weight: 500; transition: background 0.15s; text-decoration: none; }
    .action-btn.preview { background: #10b981; }
    .action-btn.preview:hover, .action-btn.preview:focus { background: #059669; }
    .action-btn.download:hover, .action-btn.download:focus { background: #1d4ed8; }
    @media (max-width: 700px) {
      .container { padding: 1.1rem 0.5rem; }
      .file-table th, .file-table td { padding: 0.7rem 0.3rem; font-size: 0.98rem; }
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      `
	if path != "/" {
		html += `<a href="/browse?token=` + token + `&path=` + parentPath + `" class="parent-link" title="Go to parent directory">&uarr; Parent Directory</a>`
	}
	html += `<span class="path">Current: ` + template.HTMLEscapeString(path) + `</span>`
	html += `</div>`
	html += `<table class="file-table">
      <tr><th></th><th>Name</th><th>Type</th><th>Actions</th></tr>`
	for _, f := range resp.Files {
		if f.IsDirectory {
			dirPath := strings.TrimPrefix(f.Path, "/")
			html += `<tr><td><span class="icon">&#128193;</span></td><td><a href="/browse?token=` + token + `&path=/` + dirPath + `" class="dir-link">` + f.Name + `</a></td><td>Directory</td><td></td></tr>`
		} else {
			filePath := strings.TrimPrefix(f.Path, "/")
			icon := "&#128196;"
			ftype := "File"
			if strings.HasSuffix(f.Name, ".png") || strings.HasSuffix(f.Name, ".jpg") || strings.HasSuffix(f.Name, ".jpeg") || strings.HasSuffix(f.Name, ".gif") || strings.HasSuffix(f.Name, ".svg") || strings.HasSuffix(f.Name, ".webp") {
				icon = "&#128247;"; ftype = "Image"
			} else if strings.HasSuffix(f.Name, ".txt") || strings.HasSuffix(f.Name, ".md") || strings.HasSuffix(f.Name, ".go") || strings.HasSuffix(f.Name, ".js") || strings.HasSuffix(f.Name, ".ts") || strings.HasSuffix(f.Name, ".json") || strings.HasSuffix(f.Name, ".css") || strings.HasSuffix(f.Name, ".html") || strings.HasSuffix(f.Name, ".py") || strings.HasSuffix(f.Name, ".sh") || strings.HasSuffix(f.Name, ".c") || strings.HasSuffix(f.Name, ".cpp") || strings.HasSuffix(f.Name, ".h") || strings.HasSuffix(f.Name, ".java") || strings.HasSuffix(f.Name, ".rs") || strings.HasSuffix(f.Name, ".toml") || strings.HasSuffix(f.Name, ".yaml") || strings.HasSuffix(f.Name, ".yml") {
				icon = "&#128441;"; ftype = "Text"
			}
			html += `<tr><td><span class="icon">` + icon + `</span></td><td>` + f.Name + `</td><td>` + ftype + `</td><td><a href="/file?file=` + filePath + `&token=` + token + `" class="action-btn preview" title="Preview this file">Preview</a><a href="/api/download/` + filePath + `?token=` + token + `" download class="action-btn download" title="Download this file">Download</a></td></tr>`
		}
	}
	html += `</table></div></body></html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// SessionEndedHandler displays a session expired message.
func SessionEndedHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.WriteHeader(http.StatusUnauthorized)
  w.Write([]byte(`
    <!DOCTYPE html>
    <html>
    <head><title>Session Expired</title></head>
    <body>
      <h2>Session Expired</h2>
      <p>Your session has ended. Please <a href="/token">enter a new token</a> to continue.</p>
    </body>
    </html>
  `))
}

// TestPageHandler is a minimal test page for token entry and download
func TestPageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
  <title>Token Download Test</title>
</head>
<body>
  <input type="text" id="token" placeholder="Enter token" />
  <button onclick="checkToken()">Check Token</button>
  <button onclick="downloadFile()">Download File</button>
  <div id="msg"></div>
  <script>
    let validToken = null;
    function checkToken() {
      const token = document.getElementById('token').value.trim();
      fetch('/api/files', { headers: { 'X-Auth-Token': token } })
        .then(resp => {
          if (!resp.ok) throw new Error('Invalid token');
          return resp.json();
        })
        .then(data => {
          validToken = token;
          document.getElementById('msg').innerText = 'Token valid! Ready to download.';
        })
        .catch(() => {
          validToken = null;
          document.getElementById('msg').innerText = 'Invalid token!';
        });
    }
    function downloadFile() {
      const token = localStorage.getItem('gofs_tkn_4a7f');
      if (!token) {
        window.location.href = '/token';
        return;
      }
      // Direct download with token as query param
      window.location.href = '/api/download/%s?token=' + encodeURIComponent(token);
    }
    // Check for token in localStorage
    if (!localStorage.getItem('gofs_tkn_4a7f')) {
      window.location.href = '/token';
    }
  </script>
</body>
</html>
  `))
}

// HomeHandler (owner page)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Owner Panel - goFileServer</title>
  <link rel="icon" href="https://fav.farm/🗂️" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap" rel="stylesheet">
  <style>
    body { font-family: 'Inter', sans-serif; background: #f8fafc; color: #1e293b; margin: 0; }
    .container { max-width: 520px; margin: 3.5rem auto; background: #fff; border-radius: 14px; box-shadow: 0 4px 24px #0002; padding: 2.7rem 2.2rem; }
    .title { font-size: 1.6rem; font-weight: 600; margin-bottom: 0.7rem; }
    .instructions { color: #475569; font-size: 1.09rem; margin-bottom: 1.7rem; line-height: 1.7; }
    .token-box { background: #f1f5f9; border-radius: 7px; padding: 1.1rem 1rem; font-size: 1.08rem; font-family: monospace; color: #2563eb; margin-bottom: 1.2rem; word-break: break-all; display: flex; align-items: center; gap: 0.7rem; flex-wrap: wrap; }
    .label { color: #64748b; font-size: 0.98rem; margin-bottom: 0.3rem; display: block; }
    .reset-btn { background: #dc2626; color: #fff; border: none; border-radius: 6px; padding: 0.7rem 1.2rem; font-size: 1.08rem; cursor: pointer; font-weight: 500; transition: background 0.15s; margin-top: 0.7rem; }
    .reset-btn:hover, .reset-btn:focus { background: #b91c1c; }
    .copy-btn { background: #10b981; color: #fff; border: none; border-radius: 6px; padding: 0.5rem 1.1rem; font-size: 1.02rem; cursor: pointer; font-weight: 500; transition: background 0.15s; outline: none; box-shadow: 0 1px 4px #0001; }
    .copy-btn:hover, .copy-btn:focus { background: #059669; }
    .copy-feedback { color: #10b981; font-size: 0.98rem; margin-left: 0.7rem; transition: opacity 0.2s; }
    .user-link { display: inline-block; margin-top: 1.5rem; color: #2563eb; text-decoration: underline; font-size: 1.08rem; font-weight: 500; transition: color 0.15s; }
    .user-link:hover { color: #1d4ed8; }
    @media (max-width: 600px) {
      .container { padding: 1.1rem 0.5rem; }
      .token-box { flex-direction: column; align-items: stretch; gap: 0.5rem; }
      .copy-btn, .copy-feedback { margin-left: 0; }
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="title">Owner Panel</div>
    <div class="instructions">
      <ol style="margin:0 0 1.2rem 1.2rem;padding:0;">
        <li>Share the access token below with people you want to give access to your files or folders.</li>
        <li>Keep this token secret. Anyone with the token can access the shared content until the session expires.</li>
        <li>To revoke access, reset the token. This will immediately invalidate the old token.</li>
      </ol>
      <span style="color:#64748b;font-size:0.98rem;">Session expires automatically after the configured duration.</span>
    </div>
    <span class="label">Current Access Token:</span>
    <div class="token-box" id="token-box">{{.Token}}
      <button class="copy-btn" id="copy-btn" title="Copy token to clipboard">Copy</button>
      <span class="copy-feedback" id="copy-feedback" style="opacity:0;">Copied!</span>
    </div>
    <form method="POST" action="/reset-token">
      <button type="submit" class="reset-btn">Reset Token</button>
    </form>
    <a href="/token" class="user-link" title="Open user access page">Go to User Access Page &rarr;</a>
  </div>
  <script>
    document.getElementById('copy-btn').onclick = function() {
      const token = document.getElementById('token-box').childNodes[0].nodeValue.trim();
      navigator.clipboard.writeText(token).then(function() {
        const feedback = document.getElementById('copy-feedback');
        feedback.style.opacity = 1;
        setTimeout(() => { feedback.style.opacity = 0; }, 1200);
      });
      return false;
    };
  </script>
</body>
</html>
  `))
  token, _ := auth.GetOwnerToken()
  tmpl.Execute(w, struct{ Token string }{Token: token})
}