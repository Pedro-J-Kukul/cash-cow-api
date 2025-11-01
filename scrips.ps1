# Cattle Marketplace API - Windows Management Script

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("start", "stop", "restart", "migrate-up", "migrate-down", "migrate-create", "logs", "clean")]
    [string]$Action,
    
    [string]$MigrationName = ""
)

function Start-Services {
    Write-Host "Starting Cattle Marketplace services..." -ForegroundColor Green
    docker-compose up -d postgres
    Start-Sleep -Seconds 10
    docker-compose up -d api
    Write-Host "Services started! API available at http://localhost:8080" -ForegroundColor Green
}

function Stop-Services {
    Write-Host "Stopping services..." -ForegroundColor Yellow
    docker-compose down
    Write-Host "Services stopped." -ForegroundColor Green
}

function Restart-Services {
    Stop-Services
    Start-Services
}

function MigrateUp {
    Write-Host "Running database migrations UP..." -ForegroundColor Blue
    docker-compose --profile migrate run --rm migrate -path=/migrations -database="postgres://cash-cow-user:cash-cow-password!@cash-cow:5432/cattle_marketplace?sslmode=disable" up
    Write-Host "Migrations completed." -ForegroundColor Green
}

function MigrateDown {
    Write-Host "Running database migrations DOWN..." -ForegroundColor Red
    docker-compose --profile migrate run --rm migrate -path=/migrations -database="postgres://cash-cow-user:cash-cow-password!@cash-cow:5432/cattle_marketplace?sslmode=disable" down 1
    Write-Host "Migration rollback completed." -ForegroundColor Green
}

function CreateMigration {
    if ([string]::IsNullOrEmpty($MigrationName)) {
        Write-Host "Migration name is required. Use: .\scripts.ps1 migrate-create -MigrationName 'your_migration_name'" -ForegroundColor Red
        return
    }
    
    Write-Host "Creating migration: $MigrationName" -ForegroundColor Blue
    
    if (!(Test-Path "migrations")) {
        New-Item -ItemType Directory -Path "migrations"
    }
    
    $timestamp = Get-Date -Format "yyyyMMddHHmmss"
    $upFile = "migrations\${timestamp}_${MigrationName}.up.sql"
    $downFile = "migrations\${timestamp}_${MigrationName}.down.sql"
    
    New-Item -ItemType File -Path $upFile
    New-Item -ItemType File -Path $downFile
    
    Write-Host "Created: $upFile" -ForegroundColor Green
    Write-Host "Created: $downFile" -ForegroundColor Green
}

function Show-Logs {
    Write-Host "Showing logs..." -ForegroundColor Blue
    docker-compose logs -f
}

function CleanEnvironment {
    Write-Host "Cleaning Docker environment..." -ForegroundColor Yellow
    docker-compose down -v
    docker system prune -f
    Write-Host "Environment cleaned." -ForegroundColor Green
}

# Execute action
switch ($Action) {
    "start" { Start-Services }
    "stop" { Stop-Services }
    "restart" { Restart-Services }
    "migrate-up" { MigrateUp }
    "migrate-down" { MigrateDown }
    "migrate-create" { CreateMigration }
    "logs" { Show-Logs }
    "clean" { CleanEnvironment }
}