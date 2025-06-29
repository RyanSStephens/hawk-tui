#!/usr/bin/env python3
"""
Setup script for Hawk TUI Python Client Library

This setup script enables easy installation and distribution of the Hawk TUI
Python client library with all its dependencies and optional features.

Installation methods:
1. Development install: pip install -e .
2. Regular install: pip install .
3. From PyPI (future): pip install hawk-tui
4. With extras: pip install .[flask,advanced]

The package supports Python 3.7+ and has minimal required dependencies.
"""

from setuptools import setup, find_packages
import os
import sys

# Ensure we're running on Python 3.7+
if sys.version_info < (3, 7):
    print("Hawk TUI requires Python 3.7 or higher")
    sys.exit(1)

# Read version from package
def get_version():
    """Get version from the package without importing it."""
    version = {}
    with open("hawk.py", "r") as f:
        content = f.read()
        # Look for __version__ = "x.x.x" pattern
        for line in content.split('\n'):
            if line.startswith('__version__'):
                exec(line, version)
                break
    
    # If no version found in code, use a default
    return version.get('__version__', '1.0.0')

# Read long description from README
def get_long_description():
    """Get long description from README file."""
    readme_path = os.path.join(os.path.dirname(__file__), 'README.md')
    if os.path.exists(readme_path):
        with open(readme_path, 'r', encoding='utf-8') as f:
            return f.read()
    else:
        return """
# Hawk TUI Python Client

Dead simple TUI integration for Python applications.

## Quick Start

```python
import hawk
hawk.auto()  # That's it!

hawk.log("Hello, world!")
hawk.metric("requests_per_second", 145)
hawk.config("port", default=8080)
```

## Features

- Zero configuration required
- Works with or without Hawk TUI running
- Thread-safe and high-performance
- Extensive monitoring capabilities
- Enterprise-grade features available

See the full documentation and examples in the repository.
"""

# Define package requirements
INSTALL_REQUIRES = [
    # No required dependencies - the library is designed to be standalone
]

# Optional dependencies for different use cases
EXTRAS_REQUIRE = {
    'flask': [
        'flask>=1.0.0',
    ],
    'advanced': [
        # For advanced features like compression, encryption, etc.
        # These are optional and the library gracefully degrades without them
    ],
    'dev': [
        'pytest>=6.0.0',
        'pytest-cov>=2.10.0',
        'black>=21.0.0',
        'flake8>=3.8.0',
        'mypy>=0.800',
        'pre-commit>=2.10.0',
    ],
    'examples': [
        'flask>=1.0.0',
        'requests>=2.20.0',
    ]
}

# Add 'all' extra that includes everything except dev
EXTRAS_REQUIRE['all'] = [
    dep for extra in ['flask', 'advanced', 'examples'] 
    for dep in EXTRAS_REQUIRE[extra]
]

setup(
    # Basic package information
    name="hawk-tui-client",
    version=get_version(),
    description="Dead simple TUI integration for Python applications",
    long_description=get_long_description(),
    long_description_content_type="text/markdown",
    
    # Author and contact information
    author="Hawk TUI Project",
    author_email="contact@hawk-tui.dev",
    url="https://github.com/hawk-tui/hawk-tui",
    project_urls={
        "Documentation": "https://docs.hawk-tui.dev",
        "Source": "https://github.com/hawk-tui/hawk-tui",
        "Tracker": "https://github.com/hawk-tui/hawk-tui/issues",
    },
    
    # Package discovery and contents
    py_modules=["hawk", "hawk_advanced"],
    packages=find_packages(exclude=["tests", "tests.*"]),
    
    # Include additional files
    include_package_data=True,
    package_data={
        '': ['*.md', '*.txt', '*.yml', '*.yaml'],
    },
    
    # Dependencies
    install_requires=INSTALL_REQUIRES,
    extras_require=EXTRAS_REQUIRE,
    python_requires=">=3.7",
    
    # Package classification
    classifiers=[
        # Development status
        "Development Status :: 4 - Beta",
        
        # Intended audience
        "Intended Audience :: Developers",
        "Intended Audience :: System Administrators",
        
        # License
        "License :: OSI Approved :: MIT License",
        
        # Programming language
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        
        # Operating systems
        "Operating System :: OS Independent",
        "Operating System :: POSIX",
        "Operating System :: Microsoft :: Windows",
        "Operating System :: MacOS",
        
        # Topics
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: System :: Monitoring",
        "Topic :: System :: Logging",
        "Topic :: Software Development :: User Interfaces",
        "Topic :: Terminals",
        
        # Framework
        "Framework :: Flask",
    ],
    
    # Keywords for discovery
    keywords=[
        "tui", "terminal", "monitoring", "logging", "metrics", 
        "dashboard", "cli", "observability", "flask", "web"
    ],
    
    # Entry points for command-line tools
    entry_points={
        'console_scripts': [
            'hawk-demo=flask_demo:main',
        ],
    },
    
    # Minimum setuptools version
    setup_requires=["setuptools>=40.0.0"],
    
    # Additional metadata
    zip_safe=True,
    license="MIT",
    
    # Test configuration
    test_suite="tests",
    tests_require=EXTRAS_REQUIRE['dev'],
)