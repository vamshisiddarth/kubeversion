# Kubeversion

A simple version manager for kubectl that allows you to switch between different kubectl versions easily.

## Features

- Interactive version selection
- Easy installation and switching between versions
- Supports all kubectl versions from official Kubernetes releases
- Search functionality to quickly find versions
- Progress bar for downloads
- No interference with system-installed kubectl

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Install from source 
Clone the repository

```bash
git clone https://github.com/yourusername/kubeversion.git
cd kubeversion
```

Build and install
```bash
go install ./cmd/kubeversion
export PATH="$HOME/.kubeversion/bin:$PATH"

# Then reload your shell:
source ~/.bashrc 
#    (or) 
source ~/.zshrc
```

## Usage

### List and Install Versions

Show interactive version selector

```bash
kubeversion list
```

### Install Specific Version

```bash
kubeversion install v1.25.0
```


### Switch to a Specific Version

```bash
kubeversion use v1.25.0
```

## Directory Structure

Kubeversion manages kubectl versions in your home directory:
```bash
~/.kubeversion/
├── bin/
│ └── kubectl -> ../versions/kubectl-v1.27.3
└── versions/
├── kubectl-v1.27.3
├── kubectl-v1.26.5
└── ...
```

## Commands

- `kubeversion list` - Interactive version selector
- `kubeversion install [version]` - Install specific version
- `kubeversion use [version]` - Switch to specific version

## Tips

1. You can type to search for specific versions in the interactive list
2. Version numbers can be specified with or without the 'v' prefix
3. The tool automatically downloads and switches to selected versions
4. Previously downloaded versions are reused without re-downloading

## Troubleshooting

### Common Issues

1. **kubectl not found**
   - Ensure `$HOME/.kubeversion/bin` is in your PATH
   - Verify shell configuration is sourced

2. **Permission denied**
   - Check file permissions in ~/.kubeversion
   - Ensure write permissions in home directory

3. **Version not switching**
   - Verify PATH order (kubeversion bin should be before system paths)
   - Check symlink in ~/.kubeversion/bin

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Built With

This project was developed using [Cursor](https://cursor.sh/), an AI-first code editor.

## License

MIT License

Copyright (c) 2024

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

