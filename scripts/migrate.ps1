param(
    [ValidateSet("up", "down", "status")]
    [string]$Command = "status",
    [string]$Service
)

$root = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$envFile = Join-Path $root ".env"

if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match "^\s*([^#=]+?)\s*=\s*(.+?)\s*$") {
            [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
        }
    }
}

$dbUrl = $env:SPARK_DATABASE_URL
if (-not $dbUrl) {
    $dbUrl = "postgres://spark:spark_dev@localhost:5432/spark?sslmode=disable"
}

function Get-MigrationDirs {
    $base = Join-Path $root "services"
    if ($Service) {
        $dir = Join-Path $base $Service "migrations"
        if (Test-Path $dir) { return @($dir) }
        return @()
    }
    return Get-ChildItem -Path "$base/*/migrations" -Directory | ForEach-Object { $_.FullName }
}

function Ensure-MigrationsTable {
    $connString = $dbUrl -replace "^postgres://", "postgresql://"
    try {
        & "psql" $connString -c @"
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    filename VARCHAR(512) NOT NULL,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    checksum VARCHAR(64)
);
"@ 2>$null
    } catch {
        Write-Warning "Could not connect to database. Make sure PostgreSQL is running."
        exit 1
    }
}

function Get-AppliedVersions {
    $connString = $dbUrl -replace "^postgres://", "postgresql://"
    $result = & "psql" $connString -t -A -c "SELECT version FROM schema_migrations ORDER BY version"
    return @($result -split "`n" | Where-Object { $_ -ne "" })
}

function Get-MigrationFiles {
    param([string]$Dir)
    $upFiles = Get-ChildItem -Path $Dir -Filter "*.up.sql" | Sort-Object Name
    $downFiles = Get-ChildItem -Path $Dir -Filter "*.down.sql" | Sort-Object Name
    return @{ Up = $upFiles; Down = $downFiles }
}

function Get-FileChecksum {
    param([string]$Path)
    if (Get-Command "Get-FileHash" -ErrorAction SilentlyContinue) {
        return (Get-FileHash -Path $Path -Algorithm SHA256).Hash
    }
    return ""
}

function Invoke-Up {
    $dirs = Get-MigrationDirs
    $applied = Get-AppliedVersions
    
    foreach ($dir in $dirs) {
        $files = Get-MigrationFiles $dir
        foreach ($file in $files.Up) {
            $version = $file.BaseName -replace "\.up$", ""
            if ($version -in $applied) {
                Write-Host "  [SKIP] $($file.Name) (already applied)" -ForegroundColor DarkGray
                continue
            }
            
            $fullPath = $file.FullName
            $connString = $dbUrl -replace "^postgres://", "postgresql://"
            
            Write-Host "  [UP]   $($file.Name)" -ForegroundColor Green
            $sql = Get-Content $fullPath -Raw
            & "psql" $connString -c $sql 2>&1 | Out-Null
            
            if ($LASTEXITCODE -eq 0) {
                $checksum = Get-FileChecksum $fullPath
                & "psql" $connString -c "INSERT INTO schema_migrations (version, filename, checksum) VALUES ('$version', '$($file.Name)', '$checksum')" 2>&1 | Out-Null
                Write-Host "         Applied successfully" -ForegroundColor Green
            } else {
                Write-Host "         FAILED" -ForegroundColor Red
                exit 1
            }
        }
    }
}

function Invoke-Down {
    $dirs = Get-MigrationDirs
    $applied = Get-AppliedVersions
    
    foreach ($dir in $dirs) {
        $files = Get-MigrationFiles $dir
        foreach ($file in $files.Down) {
            $version = $file.BaseName -replace "\.down$", ""
            if ($version -notin $applied) {
                continue
            }
            
            $fullPath = $file.FullName
            $connString = $dbUrl -replace "^postgres://", "postgresql://"
            
            Write-Host "  [DOWN] $($file.Name)" -ForegroundColor Yellow
            $sql = Get-Content $fullPath -Raw
            & "psql" $connString -c $sql 2>&1 | Out-Null
            
            if ($LASTEXITCODE -eq 0) {
                & "psql" $connString -c "DELETE FROM schema_migrations WHERE version = '$version'" 2>&1 | Out-Null
                Write-Host "         Rolled back successfully" -ForegroundColor Yellow
            } else {
                Write-Host "         FAILED" -ForegroundColor Red
                exit 1
            }
        }
    }
}

function Show-Status {
    $dirs = Get-MigrationDirs
    $applied = Get-AppliedVersions
    
    Write-Host "Migration Status:" -ForegroundColor Cyan
    Write-Host "=================" -ForegroundColor Cyan
    
    foreach ($dir in $dirs) {
        $relPath = [IO.Path]::GetRelativePath($root, $dir)
        Write-Host "`n$relPath" -ForegroundColor White
        
        $files = Get-MigrationFiles $dir
        foreach ($file in $files.Up) {
            $version = $file.BaseName -replace "\.up$", ""
            if ($version -in $applied) {
                Write-Host "  [✅] $($file.Name)" -ForegroundColor Green
            } else {
                Write-Host "  [⬜] $($file.Name)" -ForegroundColor DarkGray
            }
        }
    }
}

# Main
Ensure-MigrationsTable

switch ($Command) {
    "up"     { Invoke-Up }
    "down"   { Invoke-Down }
    "status" { Show-Status }
}
