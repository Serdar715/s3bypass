package scanner

import (
	"fmt"
	"log/slog"
	"os"
	"s3bypass/pkg/config"
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

// worker processes jobs from the channel using RequestStrategy
func (scan *Scanner) worker(jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		scan.limiter.Wait()

		url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s%s", job.Bucket, job.Prefix, job.Payload)
		
		// RequestStrategy kullanarak request yap
		respData, err := scan.strategy.Execute(scan.client, url)
		if err != nil {
			slog.Debug("Request failed", "url", url, "error", err)
			continue
		}

		// Filter Check - ResponseData üzerinden
		if scan.shouldSkipResponse(respData) {
			continue
		}

		// Verbose logging
		if scan.cfg.Verbose {
			fmt.Fprintf(os.Stderr, "%s [Status: %d, Size: %d, Words: %d, Lines: %d]   %s\n", 
				job.Payload, respData.StatusCode, respData.ContentLength, 
				respData.WordCount, respData.LineCount, url)
		}

		if respData.StatusCode == config.SuccessStatusCode {
			results <- Result{
				URL:  url,
				Size: respData.ContentLength,
			}
		}
	}
}

// shouldSkipResponse ResponseData üzerinden filtreleme yapar
func (scan *Scanner) shouldSkipResponse(resp *ResponseData) bool {
	// Filter by Status Code
	if len(scan.filter.Codes) > 0 {
		if _, ok := scan.filter.Codes[resp.StatusCode]; ok {
			return true
		}
	}

	// Filter by Content Size
	if len(scan.filter.Sizes) > 0 {
		if _, ok := scan.filter.Sizes[int(resp.ContentLength)]; ok {
			return true
		}
	}

	// Filter by Words
	if len(scan.filter.Words) > 0 {
		if _, ok := scan.filter.Words[resp.WordCount]; ok {
			return true
		}
	}
	
	// Filter by Lines
	if len(scan.filter.Lines) > 0 {
		if _, ok := scan.filter.Lines[resp.LineCount]; ok {
			return true
		}
	}

	return false
}

// Global lists moved here to avoid clutter in main
var Prefixes = []string{
	// Core & Legacy
	"", "v1/", "v2/", "v3/", "backup/", "config/", "staging/", "env/", "old/", "builds/", "test/", "deploy/", "aws/", "conf/", "db/", "tmp/",
	// Expanded v2.2
	"exports/", "db_dumps/", "financial/", "private/", "ssl/", "keys/", "users/", "customers/", "secure/", "archive/", "logs/",
	// Media & Static Files
	"assets/", "uploads/", "media/", "images/", "static/", "cdn/", "files/", "documents/", "attachments/", "downloads/",
	"photos/",
	// Development & DevOps
	"dev/", "development/", "prod/", "production/", "release/", "releases/", "ci/", "cicd/", "pipeline/", "artifacts/",
	"packages/", "dist/", "build/", "output/",
	// Data & Analytics
	"data/", "raw/", "processed/", "analytics/", "reports/", "etl/", "warehouse/", "lake/", "datalake/",
	"metrics/", "stats/", "insights/",
	// Backup & Archive
	"backups/", "snapshots/", "restore/", "dump/", "dumps/", "recovery/", "daily/", "weekly/", "monthly/", "yearly/",
	// Security & Identity
	"secrets/", "credentials/", "certs/", "certificates/", "pki/", "auth/", "tokens/", "api-keys/", "ssh/", ".ssh/",
	// System & Infrastructure
	"system/", "internal/", "admin/", "management/", "ops/", "infra/", "terraform/", "ansible/", "docker/",
	"kubernetes/", "k8s/", "helm/", "charts/",
	// Application Specific
	"api/", "web/", "app/", "mobile/", "backend/", "frontend/", "public/", "www/", "html/", "htdocs/",
	// Cloud Provider
	"lambda/", "functions/", "serverless/", "cloudformation/", "cdk/", "sam/",
}

var Payloads = []string{
	// === SECRETS & CREDENTIALS ===
	".env", ".env.local", ".env.prod", ".env.production", ".env.staging", ".env.development", ".env.test",
	".env.bak", ".env.backup", ".env.old", ".env.example", ".env.sample",
	"secrets.json", "secrets.yaml", "secrets.yml", "secrets.xml", "secrets.txt",
	"credentials.json", "credentials.yaml", "credentials.xml", "credentials", ".credentials",
	"password.txt", "passwords.txt", "passwd", "passwords.json",

	// === AWS SPECIFIC ===
	".aws/credentials", ".aws/config", "aws-credentials", "aws-config.json",
	"aws-exports.js", "aws-exports.json", "s3-config.json", "aws-creds.json",
	"cloudwatch.json", "lambda.json", "sam-template.yaml", "cdk.json", "cdk.context.json",

	// === SSH & KEYS ===
	"id_rsa", "id_rsa.pub", "id_dsa", "id_dsa.pub", "id_ed25519", "id_ed25519.pub", "id_ecdsa",
	"authorized_keys", "known_hosts", ".ssh/id_rsa", ".ssh/config", ".ssh/authorized_keys",
	"private.key", "private.pem", "server.key", "server.pem", "ssl.key", "ssl.pem",
	"deployment-key", "deploy_key", "github_rsa", "gitlab_rsa", "id_rsa_deploy",
	"master.key", "rails_master_key", "secret_key_base",

	// === API KEYS & TOKENS ===
	"api_key", "api_keys.json", "apikey.txt", "api-key.json", "api_keys.txt",
	"token.txt", "tokens.json", "access_token", "refresh_token", "auth_token.txt",
	"bearer_token", "jwt_secret", "jwt.json", "auth.json",
	"stripe_key", "stripe.json", "twilio.json", "sendgrid.json", "mailgun.json",
	"github_token", "gitlab_token", "npm_token", ".npmrc", ".pypirc",
	"slack_token", "discord_token", "telegram_token",

	// === DATABASE BACKUPS ===
	"backup.sql", "db.sql", "dump.sql", "database.sql", "db_backup.sql",
	"mysql.sql", "postgres.sql", "postgresql.sql", "mariadb.sql",
	"mongodb.json", "mongo_dump.json", "mongo.json", "redis.rdb",
	"data.sql", "migrate.sql", "schema.sql", "seed.sql", "init.sql",
	"db.sqlite", "db.sqlite3", "database.db", "app.db", "production.db",
	"production.sql", "prod.sql", "staging.sql", "dev.sql", "test.sql",
	"full_backup.sql", "incremental.sql", "dump.tar.gz", "db_export.csv",
	"users.sql", "customers.sql", "orders.sql", "transactions.sql",

	// === CONFIGURATION FILES ===
	"config.json", "config.yaml", "config.yml", "config.xml", "config.php", "config.js", "config.ts",
	"settings.json", "settings.yaml", "settings.py", "local_settings.py", "settings.ini",
	"application.yml", "application.yaml", "application.properties", "application.json",
	"bootstrap.yml", "parameters.yml", "parameters.yaml", "parameters.xml",
	"web.config", "app.config", "appsettings.json", "appsettings.Development.json", "appsettings.Production.json",
	"connections.xml", "database.yml", "db.conf.php", "wp-config.php", "wp-config.php.bak",
	"configuration.php", "configuration.yml", "env.json", "runtime.json",
	"firebase.json", ".firebaserc", "firebaseConfig.js",

	// === INFRASTRUCTURE & DevOps ===
	"terraform.tfstate", "terraform.tfvars", "terraform.tfstate.backup", "terraform.tfvars.json",
	"main.tf", "variables.tf", "outputs.tf", "backend.tf", "providers.tf",
	"ansible.cfg", "inventory.yml", "inventory.ini", "hosts.yml", "playbook.yml", "vars.yml", "vault.yml",
	"docker-compose.yml", "docker-compose.yaml", "docker-compose.prod.yml", "docker-compose.override.yml",
	"Dockerfile", ".dockerenv", "docker-stack.yml",
	"kubeconfig", ".kube/config", "k8s/secrets.yaml", "kubernetes.yaml", "k8s/kubeconfig",
	"helm/values.yaml", "values.prod.yaml", "Chart.yaml", "kustomization.yaml",
	"Vagrantfile", "serverless.yml", "serverless.yaml", "cloudformation.yaml", "template.yaml",

	// === CI/CD ===
	".travis.yml", ".gitlab-ci.yml", ".github/workflows/main.yml", "Jenkinsfile", "jenkins.xml",
	"circle.yml", ".circleci/config.yml", "azure-pipelines.yml", "bitbucket-pipelines.yml",
	".drone.yml", "cloudbuild.yaml", "buildspec.yml", "codebuild.yaml", "appveyor.yml",
	"jenkins/secrets/master.key",

	// === PACKAGE MANAGERS ===
	"package.json", "package-lock.json", ".npmrc", ".yarnrc", "yarn.lock", ".nvmrc",
	"composer.json", "composer.lock", "Gemfile", "Gemfile.lock",
	"requirements.txt", "Pipfile", "Pipfile.lock", "poetry.lock", "pyproject.toml",
	"go.mod", "go.sum", "Cargo.toml", "Cargo.lock", "Gopkg.lock",

	// === LOGS & DEBUG ===
	"debug.log", "error.log", "access.log", "app.log", "server.log", "application.log",
	"laravel.log", "storage/logs/laravel.log", "var/log/app.log", "logs/error.log",
	"phpinfo.php", "info.php", "test.php", "debug.php", "adminer.php",
	".htaccess", ".htpasswd", "nginx.conf", "httpd.conf", "apache.conf",

	// === VERSION CONTROL ===
	".git/config", ".git/HEAD", ".gitignore", ".gitconfig", ".git-credentials",
	".svn/entries", ".svn/wc.db", ".hg/hgrc",
	"CHANGELOG.md", "VERSION", ".version", "version.json",

	// === CLOUD PROVIDER SPECIFIC ===
	"service-account.json", "service-account-key.json", "gcloud-key.json", "gcp-key.json",
	"azure.json", "serviceprincipal.json", "azuredeploy.json", "azureDeploy.parameters.json",
	"digitalocean.json", "linode.json", "vultr.json",
	"deployment-key.json",

	// === API DOCUMENTATION ===
	"swagger.json", "swagger.yaml", "openapi.json", "openapi.yaml", "openapi3.yaml",
	"api-docs.json", "postman_collection.json", "graphql/schema.json", "schema.graphql",
	"insomnia.json", "insomnia_collection.json",

	// === CERTIFICATES & PKI ===
	"certificate.pem", "certificate.crt", "ca-bundle.crt", "fullchain.pem", "chain.pem",
	"privkey.pem", "keystore.jks", "truststore.jks", "client.p12", "server.pfx",
	"client.key", "client.pem", "server.crt", "ssl_certificate.pem", "ca.pem",
	"cert.pem", "key.pem", "tls.crt", "tls.key",

	// === ARCHIVE & BACKUP ===
	"backup.tar.gz", "backup.tar", "backup.zip", "backup.rar", "backup.7z",
	"site_backup.zip", "www_backup.tar.gz", "full_backup.tar.gz", "db_backup.tar.gz",
	"archive.zip", "export.zip", "data_export.zip", "prod.bak",
	"customers.zip", "users.zip", "data.tar.gz",

	// === DATA FILES ===
	"customer_data.csv", "users.csv", "customers.csv", "emails.csv", "contacts.csv",
	"employees.csv", "salaries.csv", "transactions.csv", "orders.csv",
	"data.json", "export.json", "dump.json", "backup.json",

	// === MISC SENSITIVE ===
	"shadow", "group", "gshadow", "license.key", "activation.key", "product.key",
	"storage.json", "client_secrets.json", "oauth.json", "oauth2.json",
	"admin/config.php", "admin/.env", "backup/database.sql", "db/prod.sqlite",
}
