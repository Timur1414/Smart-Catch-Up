param(
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$ternArgs
)

$envPath = Join-Path $PSScriptRoot "../../.env"
if (Test-Path $envPath) {
    Get-Content $envPath | ForEach-Object {
        if ($_ -match "^(.*?)=(.*)$") {
            [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
        }
    }
    Write-Host "Loaded environment from $envPath"
}

tern @ternArgs

$seed = Read-Host "Do you want to seed the database? (y/n)"
if ($seed -eq "y") {
    Write-Host "Seeding database..." -ForegroundColor Green
    go run seed/main.go
}
