@echo off
setlocal enabledelayedexpansion

echo Running tests...
go test ./... -v
if errorlevel 1 (
    echo Tests failed!
    exit /b 1
)

echo Building binary...
go build -o golab.exe ./cmd/golab
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Building Docker image...
docker build -t golab:latest .
if errorlevel 1 (
    echo Docker build failed!
    exit /b 1
)

echo Done.
