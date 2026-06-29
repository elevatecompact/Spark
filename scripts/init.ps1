# =============================================================================
# Spark Monorepo - Initialization Script
# =============================================================================
# Run this script to bootstrap the monorepo for local development.
# Usage: .\scripts\init.ps1
# =============================================================================

[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$InformationPreference = "Continue"

function Write-Step {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "    $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "    WARNING: $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "    ERROR: $Message" -ForegroundColor Red
}

# =============================================================================
# 1. Check prerequisites
# =============================================================================
Write-Step "Checking prerequisites..."

# Go
$goVersion = "0.0"
try {
    $goVersion = (go version 2>$null) -replace '.*go(\d+\.\d+).*', '$1'
} catch {}
if ([version]$goVersion -ge [version]"1.25") {
    Write-Success "Go $goVersion detected"
} else {
    Write-Error "Go >= 1.25 is required (found $goVersion)"
    Write-Error "Download from: https://go.dev/dl/"
    exit 1
}

# Rust / Cargo
try {
    $rustVersion = (cargo --version 2>$null)
    Write-Success "Rust detected: $rustVersion"
} catch {
    Write-Warning "Rust not found. Install from: https://rustup.rs/"
}

# Node.js
try {
    $nodeVersion = (node --version 2>$null)
    Write-Success "Node.js detected: $nodeVersion"
} catch {
    Write-Warning "Node.js not found. Install from: https://nodejs.org/"
}

# npm
try {
    $npmVersion = (npm --version 2>$null)
    Write-Success "npm detected: $npmVersion"
} catch {
    Write-Warning "npm not found."
}

# Docker
try {
    $dockerVersion = (docker --version 2>$null)
    Write-Success "Docker detected: $dockerVersion"
} catch {
    Write-Warning "Docker not found. Install from: https://docs.docker.com/get-docker/"
}

# =============================================================================
# 2. Install golangci-lint
# =============================================================================
Write-Step "Installing golangci-lint..."
try {
    $null = (golangci-lint version 2>$null)
    Write-Success "golangci-lint already installed"
} catch {
    Write-Step "Installing golangci-lint v1.59.1..."
    # binary installation for Windows
    $installDir = "$env:USERPROFILE\.local\bin"
    $null = New-Item -ItemType Directory -Path $installDir -Force -ErrorAction SilentlyContinue

    $url = "https://github.com/golangci/golangci-lint/releases/download/v1.59.1/golangci-lint-1.59.1-windows-amd64.zip"
    $zipPath = "$env:TEMP\golangci-lint.zip"
    try {
        Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
        Expand-Archive -Path $zipPath -DestinationPath "$env:TEMP\golangci-lint" -Force
        Copy-Item "$env:TEMP\golangci-lint\golangci-lint-1.59.1-windows-amd64\golangci-lint.exe" -Destination "$installDir\golangci-lint.exe" -Force
        Remove-Item -Path $zipPath -Force -ErrorAction SilentlyContinue
        Remove-Item -Path "$env:TEMP\golangci-lint" -Recurse -Force -ErrorAction SilentlyContinue

        # Add to PATH for current session
        $env:Path = "$installDir;$env:Path"
        Write-Success "golangci-lint installed to $installDir"
    } catch {
        Write-Warning "Failed to install golangci-lint: $_"
        Write-Warning "Install manually: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1"
    }
}

# =============================================================================
# 3. Install protoc plugins
# =============================================================================
Write-Step "Installing protoc plugins..."

try {
    $null = (protoc --version 2>$null)
    Write-Success "protoc already installed"
} catch {
    Write-Warning "protoc not found. Install from: https://github.com/protocolbuffers/protobuf/releases"
}

$protocPlugins = @(
    "google.golang.org/protobuf/cmd/protoc-gen-go@latest",
    "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
    "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest",
    "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest"
)

foreach ($plugin in $protocPlugins) {
    try {
        Write-Step "Installing $plugin..."
        go install $plugin
        Write-Success "Installed $plugin"
    } catch {
        Write-Warning "Failed to install $plugin`: $_"
    }
}

# =============================================================================
# 4. Install pre-commit hooks
# =============================================================================
Write-Step "Installing pre-commit hooks..."

try {
    $null = (pre-commit --version 2>$null)
    Write-Success "pre-commit already installed"
} catch {
    Write-Step "Installing pre-commit via pip..."
    try {
        pip install pre-commit --quiet
        Write-Success "pre-commit installed"
    } catch {
        Write-Warning "Failed to install pre-commit. Install manually: pip install pre-commit"
    }
}

if (Get-Command pre-commit -ErrorAction SilentlyContinue) {
    try {
        pre-commit install --hook-type pre-commit --hook-type pre-push --hook-type commit-msg
        Write-Success "Git hooks installed"
    } catch {
        Write-Warning "Failed to install git hooks: $_"
    }
} else {
    Write-Warning "pre-commit not available. Run 'pre-commit install' manually after installing it."
}

# =============================================================================
# 5. Initialize Go workspaces and tidy services
# =============================================================================
Write-Step "Initializing Go workspaces..."

$goServices = @(
    "services/identity-service",
    "services/creator-service",
    "services/viewer-service",
    "services/stream-service",
    "services/wallet-service",
    "services/chat-service",
    "services/messaging-service",
    "services/subscription-service",
    "services/gift-service",
    "services/payment-service",
    "services/analytics-service",
    "services/notification-service",
    "services/recommendation-service",
    "services/search-service",
    "services/translation-service",
    "services/moderation-service",
    "services/community-service",
    "services/event-service",
    "services/competition-service",
    "services/advertising-service",
    "services/commerce-service",
    "services/media-service",
    "services/licensing-service",
    "services/discovery-service",
    "services/trust-service"
)

$root = $PSScriptRoot | Split-Path -Parent

foreach ($svc in $goServices) {
    $svcPath = Join-Path $root $svc
    $goModPath = Join-Path $svcPath "go.mod"

    if (Test-Path $goModPath) {
        Write-Step "Tidying $svc..."
        try {
            Push-Location $svcPath
            go mod tidy
            Write-Success "go mod tidy completed for $svc"
            Pop-Location
        } catch {
            Write-Warning "go mod tidy failed for $svc`: $_"
            if ($PWD.Path -eq $svcPath) { Pop-Location }
        }
    } else {
        Write-Warning "No go.mod found in $svc (skip)"
    }
}

# =============================================================================
# 6. Initialize git hooks (additional)
# =============================================================================
Write-Step "Setting up git hooks..."

$hooksDir = Join-Path $root ".git\hooks"
if (Test-Path $hooksDir) {
    # Ensure hooks are executable (Windows handles this differently)
    Write-Success "Git hooks directory exists at $hooksDir"
} else {
    Write-Warning "No .git directory found. Initialize with 'git init' first."
}

# =============================================================================
# 7. Verify tooling installation
# =============================================================================
Write-Step "Verifying tooling installation..."

$checks = @()
$checks += if (Get-Command go -ErrorAction SilentlyContinue) { "go ok" } else { "go MISSING" }
$checks += if (Get-Command golangci-lint -ErrorAction SilentlyContinue) { "golangci-lint ok" } else { "golangci-lint MISSING" }
$checks += if (Get-Command protoc -ErrorAction SilentlyContinue) { "protoc ok" } else { "protoc MISSING" }
$checks += if (Get-Command protoc-gen-go -ErrorAction SilentlyContinue) { "protoc-gen-go ok" } else { "protoc-gen-go MISSING" }
$checks += if (Get-Command pre-commit -ErrorAction SilentlyContinue) { "pre-commit ok" } else { "pre-commit MISSING" }
$checks += if (Get-Command docker -ErrorAction SilentlyContinue) { "docker ok" } else { "docker MISSING" }
$checks += if (Get-Command node -ErrorAction SilentlyContinue) { "node ok" } else { "node MISSING" }
$checks += if (Get-Command npm -ErrorAction SilentlyContinue) { "npm ok" } else { "npm MISSING" }

foreach ($check in $checks) {
    if ($check -match "ok$") {
        Write-Success "  $check"
    } else {
        Write-Warning "  $check"
    }
}

# =============================================================================
# Summary
# =============================================================================
Write-Host ""
Write-Host "=============================================================================" -ForegroundColor Green
Write-Host "  Spark Monorepo Initialization Complete" -ForegroundColor Green
Write-Host "=============================================================================" -ForegroundColor Green
Write-Host ""
Write-Host "  Next steps:" -ForegroundColor Yellow
Write-Host "    make help          - List all available commands" -ForegroundColor White
Write-Host "    make dev-up        - Start local development environment" -ForegroundColor White
Write-Host "    make build         - Build all services" -ForegroundColor White
Write-Host "    make test          - Run all tests" -ForegroundColor White
Write-Host "    make ci            - Run full CI pipeline" -ForegroundColor White
Write-Host ""

