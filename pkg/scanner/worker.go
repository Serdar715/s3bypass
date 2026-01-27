package scanner

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
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
		
		// Random User-Agent
		req.Header.Set("User-Agent", utils.GetRandomUserAgent())

		resp, err := s.client.Do(req)
		if err == nil {
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
}
