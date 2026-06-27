param(
    [switch]$Down,
    [switch]$Build,
    [switch]$Logs,
    [string]$Service
)

$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

function Start-Dev {
    Write-Host "🚀 Starting Spark development environment..." -ForegroundColor Green
    
    # Create .env from example if not exists
    if (-not (Test-Path ".env")) {
        Copy-Item ".env.example" ".env"
        Write-Host "📝 Created .env from .env.example" -ForegroundColor Yellow
    }
    
    if ($Build) {
        docker-compose build $Service
    }
    
    docker-compose up -d $Service
    Write-Host "✅ Environment started" -ForegroundColor Green
    Write-Host "   PostgreSQL: localhost:5432" -ForegroundColor Cyan
    Write-Host "   Redis:      localhost:6379" -ForegroundColor Cyan
    Write-Host "   Kafka:      localhost:9092" -ForegroundColor Cyan
    Write-Host "   MailHog:    localhost:8025" -ForegroundColor Cyan
    Write-Host "   Kafka UI:   localhost:8081" -ForegroundColor Cyan
}

function Stop-Dev {
    Write-Host "🛑 Stopping Spark development environment..." -ForegroundColor Yellow
    docker-compose down
    Write-Host "✅ Environment stopped" -ForegroundColor Green
}

function Show-Logs {
    if ($Service) {
        docker-compose logs -f $Service
    } else {
        docker-compose logs -f
    }
}

# Main
if ($Down) {
    Stop-Dev
} elseif ($Logs) {
    Show-Logs
} else {
    Start-Dev
}
