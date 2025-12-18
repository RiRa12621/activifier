# Activifier

**Activifier** is a tiny cross-platform desktop utility that prevents your system from going idle by periodically
nudging the mouse **one pixel up and one pixel down**.

It is intentionally boring, predictable, and transparent.
---

## What it does

- Moves the mouse **+1px / −1px** at a configurable interval
- Runs as a **system tray application**
- Hides to tray when the window is closed
- Works on **Windows** and **macOS**
- Built with **Go**, **Fyne**, and **RobotGo**

Typical use cases:
- Prevent screen lock during presentations
- Keep VPNs or corporate presence tools active
- Avoid sleep during long-running tasks

---

## Usage

1. Launch Activifier
2. Use the slider to choose the interval (in seconds)
3. Click **Start**
4. Close the window — the app keeps running in the tray
5. Use the tray menu to **Start / Stop / Quit**

The mouse will move one pixel up and one pixel down at the selected interval.

---

## Platform notes

### Windows
- Works out of the box
- No additional permissions required

### macOS
- You must grant **Accessibility permissions** to Activifier  
  (`System Settings → Privacy & Security → Accessibility`)
- Without this, mouse movement will be blocked by the OS

---

## Installation

### Prebuilt binaries

Download the latest release from GitHub:

👉 **Releases:** https://github.com/rira12621/activifier/releases

- **macOS:** `.dmg`
- **Windows:** `.zip` containing a single `.exe`

### Build from source

Requirements:
- Go 1.25+
- A working C toolchain (required by RobotGo)
- Fyne dependencies for your platform

```bash
git clone https://github.com/rira12621/activifier.git
cd activifier
go build ./cmd/activifier
