#!/usr/bin/env bash
# Sets up a demo git repo with varied, vibey file changes for showcasing rizz.
# Usage:   ./demo/setup.sh [target-dir]
# Default: /tmp/rizz-demo

set -euo pipefail

TARGET="${1:-/tmp/rizz-demo}"

rm -rf "$TARGET"
mkdir -p "$TARGET"
cd "$TARGET"

git init -q -b main
git config user.email "demo@rizz.dev"
git config user.name "rizz demo"

# ---------- initial commit ----------

mkdir -p api styles internal
cat > README.md <<'EOF'
# sample app

A small demo service. Does some stuff, returns some bytes.

## development

Run `go run .` and send requests to `:8080`.
EOF

cat > main.go <<'EOF'
package main

import (
	"log"
	"net/http"

	"sample/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", api.HandleHello)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
EOF

cat > api/handlers.go <<'EOF'
package api

import (
	"fmt"
	"net/http"
)

// HandleHello greets the caller by name.
func HandleHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	fmt.Fprintf(w, "hello, %s", name)
}
EOF

cat > api/handlers_test.go <<'EOF'
package api

import (
	"net/http/httptest"
	"testing"
)

func TestHandleHello(t *testing.T) {
	req := httptest.NewRequest("GET", "/hello?name=ole", nil)
	rr := httptest.NewRecorder()
	HandleHello(rr, req)
	if got := rr.Body.String(); got != "hello, ole" {
		t.Fatalf("unexpected body: %q", got)
	}
}
EOF

cat > styles/app.css <<'EOF'
:root {
  --accent: #3fb950;
  --text: #c9d1d9;
  --bg: #0d1117;
}

body {
  background: var(--bg);
  color: var(--text);
  font-family: ui-sans-serif, system-ui;
}

.button {
  padding: 0.5rem 1rem;
  background: var(--accent);
  color: #000;
  border: none;
  border-radius: 6px;
}
EOF

cat > config.yaml <<'EOF'
server:
  port: 8080
  timeout: 30s
features:
  greeting: true
  metrics: false
log_level: info
EOF

cat > package.json <<'EOF'
{
  "name": "sample",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "start": "go run ."
  }
}
EOF

cat > internal/legacy.py <<'EOF'
# Old Python helper, scheduled for removal.
def old_helper(x):
    return x * 2
EOF

git add -A
git -c commit.gpgsign=false commit -q -m "initial commit"

# ---------- varied changes ----------

# 1. README: add a new section (pure addition)
cat >> README.md <<'EOF'

## deployment

Pushes to `main` auto-deploy via the `deploy.yaml` workflow.
EOF

# 2. main.go: register a new handler and tweak a comment (mixed edit)
cat > main.go <<'EOF'
package main

import (
	"log"
	"net/http"

	"sample/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", api.HandleHello)
	mux.HandleFunc("/health", api.HandleHealth)

	log.Println("rizz demo listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
EOF

# 3. api/handlers.go: change a function signature, add a health handler
cat > api/handlers.go <<'EOF'
package api

import (
	"encoding/json"
	"net/http"
)

// HandleHello greets the caller by name, returning JSON.
func HandleHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"greeting": "hello, " + name})
}

// HandleHealth reports service liveness.
func HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
EOF

# 4. styles/app.css: tweak values (intra-line word diff territory)
cat > styles/app.css <<'EOF'
:root {
  --accent: #ffd700;
  --text: #c9d1d9;
  --bg: #0d1117;
  --radius: 8px;
}

body {
  background: var(--bg);
  color: var(--text);
  font-family: ui-monospace, monospace;
}

.button {
  padding: 0.6rem 1.25rem;
  background: var(--accent);
  color: #000;
  border: none;
  border-radius: var(--radius);
  font-weight: 600;
}
EOF

# 5. config.yaml: flip a feature flag + bump timeout
cat > config.yaml <<'EOF'
server:
  port: 8080
  timeout: 60s
features:
  greeting: true
  metrics: true
log_level: debug
EOF

# 6. package.json: version bump
cat > package.json <<'EOF'
{
  "name": "sample",
  "version": "0.2.0",
  "private": true,
  "scripts": {
    "start": "go run .",
    "test": "go test ./..."
  }
}
EOF

# 7. new file: a bit of new Go
cat > api/metrics.go <<'EOF'
package api

import (
	"sync/atomic"
)

var requestCount atomic.Uint64

func RecordRequest() {
	requestCount.Add(1)
}

func RequestCount() uint64 {
	return requestCount.Load()
}
EOF

# 8. delete the legacy file
rm internal/legacy.py

git add -A -N  # intent-to-add for the new file so it appears in `git diff`

echo
echo "Demo repo ready at: $TARGET"
echo "Run:  cd $TARGET && rizz"
