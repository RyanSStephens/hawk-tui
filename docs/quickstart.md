# Quick Start Guide

Get up and running with Hawk TUI in under 5 minutes.

## Step 1: Install Hawk TUI

Choose your preferred installation method:

```bash
# Linux/macOS - One liner
curl -sSL https://raw.githubusercontent.com/hawk-tui/hawk-tui/main/scripts/install.sh | bash

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/hawk-tui/hawk-tui/main/scripts/install.ps1 | iex

# Or using package managers
brew install hawk-tui/tap/hawk-tui  # Homebrew
npm install -g hawk-tui            # NPM
pip install hawk-tui               # Python
```

## Step 2: Choose Your Language

### Python

Install the client library:
```bash
pip install hawk-tui
```

Add to your application:
```python
import hawk
hawk.auto()  # One line - that's it!

# Your existing code works unchanged
print("Server starting...")
import logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
logger.info("Database connected")
```

### Node.js

Install the client library:
```bash
npm install hawk-tui
```

Add to your application:
```javascript
const hawk = require('hawk-tui');
hawk.auto();  // One line - that's it!

// Your existing code works unchanged
console.log('Server starting...');
const logger = require('winston').createLogger({
    level: 'info',
    format: winston.format.json(),
    transports: [new winston.transports.Console()]
});
logger.info('Database connected');
```

### Go

Add the client library:
```bash
go get github.com/hawk-tui/hawk-tui/clients/go
```

Add to your application:
```go
import _ "github.com/hawk-tui/hawk-tui/clients/go/auto"

func main() {
    // Your existing code works unchanged
    fmt.Println("Server starting...")
    log.Println("Database connected")
}
```

## Step 3: Run with TUI

Run your application and pipe it through Hawk TUI:

```bash
# Python
python your_app.py | hawk

# Node.js
node your_app.js | hawk

# Go
./your_app | hawk

# Any other application
your_command | hawk
```

## Step 4: Explore the Interface

Once running, you'll see a beautiful TUI with:

- **Logs Panel**: All your application output with syntax highlighting
- **Metrics Panel**: Automatically detected metrics and counters
- **Status Panel**: Application health and performance info

### Keyboard Shortcuts

- `Tab` / `Shift+Tab`: Navigate between panels
- `j/k` or `↑/↓`: Scroll within panels
- `Ctrl+C` or `q`: Quit
- `/`: Search logs
- `r`: Refresh/reload
- `h`: Show help

## Examples

### Web Server Monitoring

```python
# app.py
import hawk
from flask import Flask

hawk.auto()
app = Flask(__name__)

@app.route('/api/users')
def get_users():
    hawk.metric('api_requests', 1)
    users = fetch_users()
    hawk.log(f'Returned {len(users)} users')
    return users

if __name__ == '__main__':
    hawk.log('Starting web server', level='INFO')
    app.run(host='0.0.0.0', port=5000)
```

Run it:
```bash
python app.py | hawk
```

### Database Processing

```python
# processor.py
import hawk

hawk.auto()

def process_records():
    total = 10000
    progress = hawk.progress('Processing records', total=total)
    
    for i in range(total):
        process_record(i)
        progress.update(i + 1)
        
        if i % 100 == 0:
            hawk.metric('records_processed', i)
            hawk.log(f'Processed {i}/{total} records')
    
    hawk.log('Processing complete!', level='SUCCESS')

if __name__ == '__main__':
    process_records()
```

### Real-time Data Pipeline

```javascript
// pipeline.js
const hawk = require('hawk-tui');
hawk.auto();

async function processPipeline() {
    hawk.log('Starting data pipeline', { level: 'INFO' });
    
    while (true) {
        const batch = await fetchDataBatch();
        hawk.metric('batch_size', batch.length);
        
        const processed = await processData(batch);
        hawk.metric('processed_records', processed.length);
        
        hawk.log(`Processed batch: ${processed.length} records`);
        
        await new Promise(resolve => setTimeout(resolve, 1000));
    }
}

processPipeline().catch(console.error);
```

## Configuration

Customize Hawk TUI with command-line options:

```bash
# Dark theme (default)
python app.py | hawk --theme dark

# Light theme  
python app.py | hawk --theme light

# Custom refresh rate
python app.py | hawk --refresh-rate 500ms

# Remote monitoring
python app.py | hawk --remote --port 9090
```

## What's Next?

- **[API Reference](api-reference.md)**: Learn about all available methods
- **[Configuration](configuration.md)**: Customize the interface and behavior
- **[Examples](examples.md)**: See real-world usage examples
- **[Architecture](architecture.md)**: Understand how it works under the hood

## Troubleshooting

### Common Issues

**Hawk TUI doesn't start**
- Make sure Hawk TUI is installed: `hawk --version`
- Check that your application is sending output
- Try running without piping first to verify your app works

**No logs appear**
- Make sure you're using `print()` or proper logging
- Check that `hawk.auto()` is called before other output
- Verify your application isn't buffering output

**Performance issues**
- Reduce refresh rate: `hawk --refresh-rate 1000ms`
- Limit log lines: `hawk --max-logs 1000`
- Check system resources with `top` or Task Manager

**Installation problems**
- See the detailed [Installation Guide](installation.md)
- Try manual installation if scripts fail
- Check permissions and PATH settings

Need more help? Check our [Troubleshooting Guide](troubleshooting.md) or open an issue on [GitHub](https://github.com/hawk-tui/hawk-tui/issues).