#!/bin/bash

# Detect container runtime (Docker or Podman)
detect_runtime() {
    if command -v docker &> /dev/null; then
        echo "docker"
    elif command -v podman &> /dev/null; then
        echo "podman"
    else
        echo "none"
    fi
}

# Cleanup function
cleanup() {
    local runtime=$1

    echo "=== Cleaning up unused $runtime resources ==="

    # Remove dangling/untagged images
    echo "Removing dangling images..."
    $runtime image prune -f

    # Remove all unused images (not just dangling)
    echo "Removing unused images..."
    $runtime image prune -a -f

    # Remove unused volumes
    echo "Removing unused volumes..."
    $runtime volume prune -f

    # Remove unused networks (Docker only, Podman skips)
    if [ "$runtime" = "docker" ]; then
        echo "Removing unused networks..."
        $runtime network prune -f
    fi

    # Clean build cache (Docker only)
    if [ "$runtime" = "docker" ]; then
        echo "Clearing build cache..."
        $runtime builder prune -f
    fi

    echo "=== Cleanup complete! ==="
    echo "Disk space freed:"
    df -h / | awk 'NR==2 {print "Available: " $4}'
}

# Main script
runtime=$(detect_runtime)

case $runtime in
    "docker"|"podman")
        cleanup "$runtime"
        ;;
    *)
        echo "Error: Neither Docker nor Podman is installed."
        exit 1
        ;;
esac