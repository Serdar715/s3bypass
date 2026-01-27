# S3Bypass - Advanced S3 Scanner

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.19%2B-cyan.svg)

**S3Bypass** is a high-performance, concurrent penetration testing tool designed to identify sensitive files and configuration exposures in public AWS S3 buckets. 

Built with Go for speed and reliability, it uses a worker-pool architecture to efficiently scan thousands of potential paths across multiple buckets.

## 🚀 Key Features

*   **High Performance**: Tuned worker pool usage (default 100 threads) for rapid scanning.
*   **Smart Parsing**: Automatically handles various input formats (`http://...`, `Protected S3 Bucket: ...`).
*   **Comprehensive Payloads**: Checks for 60+ critical file types including:
    *   Secrets (`.env`, `credentials`, `id_rsa`)
    *   Backups (`.sql`, `.tar.gz`, `.bak`)
    *   Configs (`config.php`, `settings.py`, `terraform.tfstate`)
*   **Resource Efficient**: Streamed file reading to handle massive input lists without memory bloat.
*   **Modular Architecture**: Clean, maintainable Go codebase.

## 📦 Installation

Ensure you have [Go](https://go.dev/) installed.

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
./s3bypass -l targets.txt
```

### Single Target Scan
Quickly check a single bucket.

```bash
./s3bypass -u http://target-bucket.s3.amazonaws.com/
```

### Advanced Options

```bash
# Scan with 200 threads and save to custom output
./s3bypass -l targets.txt -t 200 -o secrets_found.txt

# Help menu
./s3bypass -h
```

## 📋 Input Formats

The tool is designed to be pipe-friendly and robust. It accepts:
*   Standard URLs: `http://bucket.s3.amazonaws.com`
*   Bucket Names: `my-production-bucket`
*   Tool Output: `Protected S3 Bucket: http://bucket.s3.amazonaws.com/`

## ⚠️ Disclaimer

This tool is for educational purposes and authorized security testing only. Scanning targets without permission is illegal. The author is not responsible for any misuse.

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.
