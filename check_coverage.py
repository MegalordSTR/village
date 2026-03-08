#!/usr/bin/env python3
import subprocess
import re
import sys

# Run go tool cover -func
result = subprocess.run(['go', 'tool', 'cover', '-func=coverage.out'], 
                       capture_output=True, text=True, cwd='.')
if result.returncode != 0:
    print("Failed to run go tool cover")
    sys.exit(1)

lines = result.stdout.strip().split('\n')
# Parse each line
file_coverage = {}
for line in lines:
    if not line.strip() or line.startswith('total:'):
        continue
    # Format: filename:line:function\tcoverage%
    # Some lines have tabs, some spaces
    if '\t' in line:
        parts = line.split('\t')
        filename_part = parts[0]
        coverage_part = parts[-1]
    else:
        # Split by spaces
        parts = line.split()
        if len(parts) < 2:
            continue
        filename_part = parts[0]
        coverage_part = parts[-1]
    # Extract filename (before first colon)
    filename = filename_part.split(':')[0]
    # Extract coverage percentage (remove '%')
    coverage = float(coverage_part.replace('%', ''))
    if filename not in file_coverage:
        file_coverage[filename] = []
    file_coverage[filename].append(coverage)

# Compute average per file
print("Coverage per file:")
below_threshold = []
for filename, coverages in sorted(file_coverage.items()):
    avg = sum(coverages) / len(coverages)
    print(f"{filename}: {avg:.1f}% (based on {len(coverages)} functions)")
    if avg < 80.0:
        below_threshold.append((filename, avg))

# Check system files
system_files = {
    'environment': 'environment.go',
    'production': 'production.go',
    'social': 'social.go',
    'economy': 'economy.go',
    'events': 'events.go',
}

print("\nSystem coverage:")
for system, file in system_files.items():
    full_path = f'github.com/vano44/village/internal/simulation/{file}'
    if full_path in file_coverage:
        avg = sum(file_coverage[full_path]) / len(file_coverage[full_path])
        status = "✅" if avg >= 80.0 else "❌"
        print(f"{system}: {avg:.1f}% {status}")
    else:
        print(f"{system}: no coverage data")

if below_threshold:
    print("\nFiles below 80%:")
    for f, cov in below_threshold:
        print(f"  {f}: {cov:.1f}%")
    sys.exit(1)
else:
    print("\nAll files ≥80% coverage.")