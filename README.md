# 🔐 pswd

**pswd** is a minimal and secure command-line password manager
inspired by [pass](https://www.passwordstore.org/), but built without GPG.

Instead of relying on external tools, `pswd` uses **Elliptic Curve Cryptography (ECC)**
internally to securely encrypt and decrypt your passwords.

Passwords are stored locally in `~/.pswd`.

## ✨ Features

* 🔐 Secure password encryption using ECC (Curve25519 + AES-GCM)
* 📁 Local, file-based password store
* 🚫 No GPG or external dependencies
* 🖥️ Simple CLI interface

## 📦 Installation

```bash
git clone https://github.com/yourusername/pswd.git
cd pswd
go build -o pswd ./cmd/cli.go
```

## 🚀 Usage

```bash
pswd [command]
```

### Available Commands

| Command      | Description                          |
| ------------ | ------------------------------------ |
| `init`       | Initialize a new password store      |
| `insert`     | Insert a new password by name        |
| `show`       | Show a stored password by name       |
| `help`       | Show help for commands               |

### Examples

Initialize the password store:

```bash
pswd init
```

Insert a password:

```bash
pswd insert example.com
```

Show a password:

```bash
pswd show example.com
```

## 🔒 How It Works

* Uses **Curve25519 (X25519)** for key exchange
* Derives a symmetric key via **HKDF-SHA256**
* Encrypts data with **AES-256-GCM**
* Encrypted passwords are stored as individual files in `~/.pswd/`

No external GPG or OpenSSL binaries are needed — all crypto is handled natively in Go.

## 🛣 Roadmap

Planned features and improvements:

- [x] 🔢 Built-in password generator
- [ ] 🔐 Password editing and deletion
- [ ] 🌳 Git integration
- [ ] 📥 Import from other password managers
- [ ] 🔍 Full-text search
- [ ] 🧪 Unit + integration test coverage
- [x] 🗂️ Hierarchical folder-like structure for entries
