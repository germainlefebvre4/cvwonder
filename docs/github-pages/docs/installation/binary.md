---
sidebar_position: 1
---
# Binary

---

The binary is available for Windows, MacOS, and Linux. You can download it from the [releases page](https://github.com/germainlefebvre4/cvwonder/releases).

You can install it using the following command:

```bash
apt install curl jq

DISTRIBUTION=linux   # linux, darwin, windows
CPU_ARCH=amd64       # amd64, arm64, i386

VERSION=$(curl -s "https://api.github.com/repos/germainlefebvre4/cvwonder/releases/latest" | jq -r '.tag_name')
curl -L -o cvwonder "https://github.com/germainlefebvre4/cvwonder/releases/download/${VERSION}/cvwonder_${DISTRIBUTION}_${CPU_ARCH}"
chmod +x cvwonder
sudo mv cvwonder /usr/local/bin/
```

This command will download the latest version of CV Wonder and install it in your `PATH`.
