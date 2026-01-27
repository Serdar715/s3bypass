package scanner

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"s3bypass/pkg/config"
	"s3bypass/pkg/utils"
	"time"
)

// Job represents a single scan task
type Job struct {
	Bucket  string
	Prefix  string
	Payload string
}

// Result represents a successful finding
type Result struct {
	URL  string
	Size int64
}

// worker processes jobs from the channel
func (s *Scanner) worker(jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		// Rate Limiting / Delay
		if s.cfg.Delay > 0 {
			// Add jitter (+/- 10%)
			jitter := int(float64(s.cfg.Delay) * config.JitterPercentage)
			actualDelay := s.cfg.Delay + rand.Intn(jitter*config.JitterMultiplier+1) - jitter
			if actualDelay < 0 {
				actualDelay = 0
			}
			time.Sleep(time.Duration(actualDelay) * time.Millisecond)
		}

		url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s%s", job.Bucket, job.Prefix, job.Payload)
		
		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			slog.Error("Failed to create request", "url", url, "error", err)
			continue
		}
		
		req.Header.Set("User-Agent", utils.GetRandomUserAgent())

		resp, err := s.client.Do(req)
		if err == nil {
			// FILTER CHECK
			if len(s.filterCodes) > 0 {
				if _, ok := s.filterCodes[resp.StatusCode]; ok {
					resp.Body.Close()
					continue // Filtered by code
				}
			}
			if len(s.filterSizes) > 0 {
				if _, ok := s.filterSizes[int(resp.ContentLength)]; ok {
					resp.Body.Close()
					continue // Filtered by size
				}
			}
			// Note: Words/Lines filtering requires body parsing (GET). 
			// Since we use HEAD, these will be 0. 
			// If user filters -fw 0, it effectively filters everything out unless we change logic.
			// For now, we check against 0.
			if len(s.filterWords) > 0 {
				if _, ok := s.filterWords[0]; ok { // Words=0
					resp.Body.Close()
					continue 
				}
			}
			if len(s.filterLines) > 0 {
				if _, ok := s.filterLines[0]; ok { // Lines=0
					resp.Body.Close()
					continue 
				}
			}

			// FFUF-style verbose output
			if s.cfg.Verbose {
				// Format: filename [Status: CODE, Size: SIZE, Words: 0, Lines: 0]   URL
				fmt.Fprintf(os.Stderr, "%s [Status: %d, Size: %d, Words: 0, Lines: 0]   %s\n", 
					job.Payload, resp.StatusCode, resp.ContentLength, url)
			}

			if resp.StatusCode == 200 {
				results <- Result{
					URL:  url,
					Size: resp.ContentLength,
				}
			}
			resp.Body.Close()
		} else {
			slog.Debug("Request failed", "url", url, "error", err)
		}
	}
}

// Global lists moved here to avoid clutter in main
var Prefixes = []string{
	"", "v1/", "v2/", "backup/", "config/", "staging/", "env/", "old/", "builds/", "test/", "deploy/", "aws/", "conf/", "db/", "tmp/",
	// Expanded v2.2
	"exports/", "db_dumps/", "financial/", "private/", "ssl/", "keys/", "users/", "customers/", "secure/", "archive/", "logs/",
}

var Payloads = []string{
	// Secrets
	".env", ".env.local", ".env.prod", ".env.staging", ".env.bak",
	"secrets.json", "secrets.yaml", "credentials.json", "credentials",
	".aws/credentials", ".passwd", "id_rsa", "id_rsa.pub", "master.key",
	"token.txt", "access_token", "auth.json", "service-account.json",
	// Configs
	"config.json", "config.php", "config.js", "web.config", "settings.py",
	"local_settings.py", "application.yml", "bootstrap.yml", "firebase.json",
	"parameters.yml", "connections.xml", "db.conf.php", "docker-compose.yml",
	// Backups
	"backup.sql", "db.sql", "dump.sql", "database.sql", "db_backup.sql",
	"backup.tar.gz", "backup.zip", "full_backup.sql", "mysql.sql",
	"data.sql", "migrate.sql", "dump.gz", "prod.bak", "db.sqlite",
	// Dev
	"package-lock.json", ".npmrc", ".yarnrc", "composer.json", "Gemfile.lock",
	".gitignore", ".git/config", "terraform.tfstate", "terraform.tfvars",
	".travis.yml", ".gitlab-ci.yml", "jenkins.xml", "circle.yml",
	// Logs
	"phpinfo.php", "info.php", "debug.log", "error.log", "access.log",
	// Expanded v2.2
	"server.key", "api_keys.json", "customer_data.csv", "database.sqlite", 
	"auth_token.txt", "client_secrets.json", "keystore.jks", "backup.rar", 
	"shadow", "passwd", "id_dsa",
	// Deep Research v3.0
	"swagger.json", "swagger.yaml", "openapi.json", "graphql/schema.json",
	"admin/config.php", "admin/.env", "backup/database.sql", "db/prod.sqlite",
	"jenkins/secrets/master.key", "k8s/kubeconfig", "kubeconfig", ".kube/config",
	"id_rsa_deploy", "deployment-key.json", "service-account-key.json",
	"storage.json", "aws-creds.json", "s3-config.json",
}
