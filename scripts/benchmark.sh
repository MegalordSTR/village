#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

OUTPUT_DIR="${OUTPUT_DIR:-./benchmark_results}"
mkdir -p "$OUTPUT_DIR"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
OUTPUT_FILE="$OUTPUT_DIR/benchmark_$TIMESTAMP.txt"

echo "Running simulation benchmarks..."
echo "Results will be saved to $OUTPUT_FILE"
echo ""

# Run benchmarks with benchtime 10 seconds, include memory stats
go test ./internal/simulation -bench=. -benchtime=10s -benchmem 2>&1 | tee "$OUTPUT_FILE"

# Extract ns/op for each benchmark
echo "" >> "$OUTPUT_FILE"
echo "=== Performance Summary ===" >> "$OUTPUT_FILE"

# Parse benchmark results and check against targets
awk '/^BenchmarkProcessWeek[0-9]+Residents[-0-9]*/ {
    # Extract resident count
    match($1, /[0-9]+/)
    residents = substr($1, RSTART, RLENGTH)
    # ns/op is the third field
    ns = $3
    # Convert to seconds
    seconds = ns / 1000000000
    # Target thresholds
    target = (residents == 10 ? 1.0 : (residents == 50 ? 3.0 : (residents == 100 ? 5.0 : 10.0)))
    status = (seconds < target) ? "✅ PASS" : "❌ FAIL"
    printf "  %s residents: %.6fs (target <%.1fs) %s\n", residents, seconds, target, status
}' "$OUTPUT_FILE" >> "$OUTPUT_FILE"

echo "Benchmark complete."