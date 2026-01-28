# S3Bypass - Advanced S3 Scanner

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.20%2B-cyan.svg)

**S3Bypass** is a high-performance, concurrent penetration testing tool designed to identify sensitive files and configuration exposures in public AWS S3 buckets. 

Built with Go for speed and reliability, it uses a worker-pool architecture to efficiently scan thousands of potential paths across multiple buckets.

## 🚀 Key Features

*   **High Performance**: Tuned worker pool usage (default 10 threads) for rapid scanning.
*   **Smart Parsing**: Automatically handles various input formats (`http://...`, `Protected S3 Bucket: ...`).
*   **Comprehensive Coverage**: 
    *   **90+ Prefixes**: Media, DevOps, data, security, and cloud directories
    *   **200+ Payloads**: Critical file types including secrets, configs, backups
*   **🔍 Payload Categories**:
    *   Secrets (`.env`, `.aws/credentials`, `id_rsa`, `master.key`)
    *   Backups (`.sql`, `.tar.gz`, `mongodb.json`, `redis.rdb`)
    *   Configs (`terraform.tfstate`, `kubeconfig`, `docker-compose.yml`)
    *   CI/CD (`Jenkinsfile`, `.gitlab-ci.yml`, `azure-pipelines.yml`)
    *   Certificates (`fullchain.pem`, `keystore.jks`, `tls.key`)
    *   API Docs (`swagger.json`, `openapi.yaml`, `postman_collection.json`)
*   **🛡️ WAF Evasion**:
    *   **Rate Limiting**: Configurable delay with automatic jitter to mimic human behavior.
    *   **Random User-Agents**: Rotates modern browser signatures for every request.
*   **🎯 Smart Filtering**: Filter by status code, size, words, and lines (`-fc`, `-fs`, `-fw`, `-fl`)
*   **Resource Efficient**: Streamed file reading to handle massive input lists without memory bloat.
*   **Modular Architecture**: Clean, maintainable Go codebase.

## 📦 Installation

Ensure you have [Go 1.20+](https://go.dev/) installed.

```bash
git clone https://github.com/Serdar715/s3bypass
cd s3bypass
go build -o s3bypass ./cmd/s3bypass
sudo mv s3bypass /usr/local/bin/
```

## 🛠️ Usage

### Basic Scan (From File)
Scan a list of buckets from a text file. The tool automatically cleans prefixes.

```bash
s3bypass -l targets.txt
```

### Single Target Scan
Quickly check a single bucket.

```bash
s3bypass -u http://target-bucket.s3.amazonaws.com/
```

### WAF Evasion Mode
Scan with 500ms delay (plus random jitter) to avoid rate limits.

```bash
s3bypass -l targets.txt -delay 500 -t 20
```

### Advanced Options

```bash
# Scan with 50 threads and save to custom output
s3bypass -l targets.txt -t 50 -o secrets_found.txt

# Use custom wordlist
s3bypass -l targets.txt -w my_payloads.txt

# Filter out 403 and 404 responses
s3bypass -l targets.txt -fc 403,404

# Verbose mode for debugging
s3bypass -l targets.txt -v

# Help menu
s3bypass -h
```

## 📋 Input Formats

The tool is designed to be pipe-friendly and robust. It accepts:
*   Standard URLs: `http://bucket.s3.amazonaws.com`
*   Bucket Names: `my-production-bucket`
*   Tool Output: `Protected S3 Bucket: http://bucket.s3.amazonaws.com/`

## 🗂️ Directory Prefixes (90+)

| Category | Examples |
|----------|----------|
| Core | `backup/`, `config/`, `staging/`, `env/` |
| Media | `assets/`, `uploads/`, `media/`, `cdn/` |
| DevOps | `terraform/`, `kubernetes/`, `docker/`, `ansible/` |
| Data | `data/`, `analytics/`, `datalake/`, `exports/` |
| Security | `secrets/`, `credentials/`, `certs/`, `api-keys/` |
| Cloud | `lambda/`, `serverless/`, `cloudformation/` |

## ⚠️ Disclaimer

This tool is for educational purposes and authorized security testing only. Scanning targets without permission is illegal. The author is not responsible for any misuse.

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.
