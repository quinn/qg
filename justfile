
copy:
    mkdir -p appdirs || true
    curl \
        -o appdirs/appdirs.go \
        https://raw.githubusercontent.com/hay-kot/scaffold/refs/heads/main/internal/appdirs/appdirs.go