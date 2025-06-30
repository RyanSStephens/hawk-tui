# Hawk TUI Installation Script for Windows
# PowerShell script for Windows installation

param(
    [string]$Version = "latest",
    [string]$InstallDir = "",
    [switch]$Source = $false,
    [switch]$Help = $false
)

# Colors for output
$ErrorColor = "Red"
$SuccessColor = "Green" 
$WarningColor = "Yellow"
$InfoColor = "Cyan"

function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $InfoColor
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $SuccessColor
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $WarningColor
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $ErrorColor
}

function Show-Help {
    Write-Host @"
Hawk TUI Installer for Windows

Usage: .\install.ps1 [OPTIONS]

Options:
  -Version VERSION     Install specific version (default: latest)
  -InstallDir DIR      Installation directory (default: %USERPROFILE%\.local\bin)
  -Source              Install from source instead of binary
  -Help                Show this help message

Examples:
  .\install.ps1                           # Install latest version
  .\install.ps1 -Version v1.0.0           # Install specific version
  .\install.ps1 -InstallDir C:\tools      # Install to custom directory
  .\install.ps1 -Source                   # Install from source

"@
}

function Test-IsAdmin {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Set-InstallDirectory {
    if ($InstallDir -eq "") {
        if (Test-IsAdmin) {
            $InstallDir = "$env:ProgramFiles\Hawk TUI"
            Write-Status "Installing system-wide as administrator"
        } else {
            $InstallDir = "$env:USERPROFILE\.local\bin"
            Write-Status "Installing to user directory"
        }
    }
    
    # Create directory if it doesn't exist
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-Status "Created directory: $InstallDir"
    }
    
    return $InstallDir
}

function Add-ToPath {
    param([string]$Directory)
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$Directory*") {
        Write-Warning "Adding $Directory to user PATH"
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Warning "Please restart your terminal or run: `$env:PATH += ';$Directory'"
    }
}

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { 
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

function Test-Dependencies {
    Write-Status "Checking dependencies..."
    
    # Check if we can download files
    try {
        $null = Invoke-WebRequest -Uri "https://www.google.com" -UseBasicParsing -TimeoutSec 5
        Write-Success "Internet connection verified"
    } catch {
        Write-Error "Internet connection required for installation"
        exit 1
    }
    
    # Check if we have extraction capabilities
    if (!(Get-Command Expand-Archive -ErrorAction SilentlyContinue)) {
        Write-Error "PowerShell 5.0+ required for archive extraction"
        exit 1
    }
    
    Write-Success "All dependencies found"
}

function Get-LatestVersion {
    if ($Version -eq "latest") {
        Write-Status "Fetching latest release version..."
        try {
            $response = Invoke-RestMethod -Uri "https://api.github.com/repos/hawk-tui/hawk-tui/releases/latest"
            $Version = $response.tag_name
            Write-Status "Latest version: $Version"
        } catch {
            Write-Warning "Could not fetch latest version, using v1.0.0"
            $Version = "v1.0.0"
        }
    }
    return $Version
}

function Install-Binary {
    param([string]$InstallDir)
    
    $arch = Get-Architecture
    $downloadUrl = "https://github.com/hawk-tui/hawk-tui/releases/download/$Version/hawk-tui-$Version-windows-$arch.zip"
    $tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    $tempFile = "$tempDir\hawk-tui.zip"
    
    try {
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        
        Write-Status "Downloading Hawk TUI from $downloadUrl"
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
        
        Write-Status "Extracting archive..."
        Expand-Archive -Path $tempFile -DestinationPath $tempDir -Force
        
        Write-Status "Installing binary to $InstallDir"
        $binaryPath = Get-ChildItem -Path $tempDir -Name "hawk.exe" -Recurse | Select-Object -First 1
        if ($binaryPath) {
            $sourcePath = Join-Path $tempDir $binaryPath
            Copy-Item -Path $sourcePath -Destination "$InstallDir\hawk.exe" -Force
        } else {
            # Try alternative path
            Copy-Item -Path "$tempDir\hawk.exe" -Destination "$InstallDir\hawk.exe" -Force
        }
        
        Write-Success "Binary installed successfully"
        return $true
        
    } catch {
        Write-Error "Failed to install binary: $($_.Exception.Message)"
        return $false
    } finally {
        if (Test-Path $tempDir) {
            Remove-Item -Path $tempDir -Recurse -Force
        }
    }
}

function Install-FromSource {
    param([string]$InstallDir)
    
    Write-Status "Installing from source..."
    
    # Check for Go
    if (!(Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Error "Go is required for source installation"
        Write-Error "Please install Go 1.21+ from https://golang.org/dl/"
        exit 1
    }
    
    # Check for Git
    if (!(Get-Command git -ErrorAction SilentlyContinue)) {
        Write-Error "Git is required for source installation"
        Write-Error "Please install Git from https://git-scm.com/download/win"
        exit 1
    }
    
    $tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    
    try {
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        Set-Location $tempDir
        
        Write-Status "Cloning repository..."
        git clone "https://github.com/hawk-tui/hawk-tui.git" .
        
        Write-Status "Building binary..."
        go build -o hawk.exe .\cmd\hawk
        
        Write-Status "Installing binary to $InstallDir"
        Copy-Item -Path "hawk.exe" -Destination "$InstallDir\hawk.exe" -Force
        
        Write-Success "Source installation completed"
        return $true
        
    } catch {
        Write-Error "Source installation failed: $($_.Exception.Message)"
        return $false
    } finally {
        Set-Location $PSScriptRoot
        if (Test-Path $tempDir) {
            Remove-Item -Path $tempDir -Recurse -Force
        }
    }
}

function Test-Installation {
    param([string]$InstallDir)
    
    $hawkPath = "$InstallDir\hawk.exe"
    if (Test-Path $hawkPath) {
        try {
            $versionOutput = & $hawkPath --version 2>$null
            $installedVersion = if ($versionOutput -match 'v\d+\.\d+\.\d+') { $matches[0] } else { "unknown" }
            Write-Success "Hawk TUI installed successfully!"
            Write-Success "Version: $installedVersion"
            Write-Success "Location: $hawkPath"
            return $true
        } catch {
            Write-Error "Installation verification failed: hawk.exe found but not working"
            return $false
        }
    } else {
        Write-Error "Installation verification failed: hawk.exe not found"
        return $false
    }
}

function Show-Examples {
    Write-Host ""
    Write-Success "Installation complete! Here are some quick examples:"
    Write-Host ""
    Write-Host "  # Basic usage with any application:"
    Write-Host "  your-app.exe | hawk"
    Write-Host ""
    Write-Host "  # With Node.js application:"
    Write-Host "  node app.js | hawk"
    Write-Host ""
    Write-Host "  # With Python application:"
    Write-Host "  python app.py | hawk"
    Write-Host ""
    Write-Host "  # With PowerShell:"
    Write-Host "  Get-Process | ConvertTo-Json | hawk"
    Write-Host ""
    Write-Host "  # Show help:"
    Write-Host "  hawk --help"
    Write-Host ""
    Write-Host "  # Show version:"
    Write-Host "  hawk --version"
    Write-Host ""
    Write-Status "Visit https://hawktui.dev for documentation and examples"
}

# Main installation function
function Main {
    Write-Host "╔══════════════════════════════════════════════════════════════════════════════╗"
    Write-Host "║                              🦅 Hawk TUI Installer                          ║"
    Write-Host "║                     Universal TUI Framework for Any Language                ║"
    Write-Host "╚══════════════════════════════════════════════════════════════════════════════╝"
    Write-Host ""
    
    if ($Help) {
        Show-Help
        exit 0
    }
    
    $InstallDir = Set-InstallDirectory
    Test-Dependencies
    $Version = Get-LatestVersion
    
    $success = $false
    if ($Source) {
        $success = Install-FromSource -InstallDir $InstallDir
    } else {
        $success = Install-Binary -InstallDir $InstallDir
        if (!$success) {
            Write-Warning "Binary installation failed, trying source installation..."
            $success = Install-FromSource -InstallDir $InstallDir
        }
    }
    
    if (!$success) {
        Write-Error "Installation failed"
        exit 1
    }
    
    if (Test-Installation -InstallDir $InstallDir) {
        Add-ToPath -Directory $InstallDir
        Show-Examples
        Write-Host ""
        Write-Success "🎉 Happy monitoring with Hawk TUI!"
    } else {
        exit 1
    }
}

# Run main function
Main