#!/usr/bin/env python3
import os
import re
import sys
from collections import defaultdict, deque

def extract_dependencies(content):
    """Extract dependencies from markdown content."""
    lines = content.split('\n')
    deps = []
    in_deps = False
    for line in lines:
        stripped = line.strip()
        if stripped.startswith('## Dependencies'):
            in_deps = True
            continue
        if in_deps and stripped.startswith('##'):
            break
        if in_deps and stripped.startswith('-'):
            # Check for "All previous workstreams (01-10)"
            if 'All previous workstreams' in stripped:
                # Add all 01-10
                for i in range(1, 11):
                    deps.append(f'00-001-{i:02d}')
            else:
                match = re.search(r'00-\d{3}-\d{2}', stripped)
                if match:
                    deps.append(match.group(0))
    return deps

def parse_workstream(filepath):
    """Parse workstream file and return (ws_id, dependencies)."""
    with open(filepath, 'r') as f:
        content = f.read()
    # Extract ws_id from frontmatter
    ws_id_match = re.search(r'ws_id:\s*(00-\d{3}-\d{2})', content)
    if not ws_id_match:
        # fallback from filename
        ws_id = os.path.basename(filepath).replace('.md', '')
    else:
        ws_id = ws_id_match.group(1)
    
    deps = extract_dependencies(content)
    return ws_id, deps

def topological_sort(graph):
    """Kahn's algorithm for topological sort."""
    in_degree = defaultdict(int)
    for node in graph:
        for neighbor in graph[node]:
            in_degree[neighbor] += 1
    for node in graph:
        if node not in in_degree:
            in_degree[node] = 0
    
    queue = deque([node for node in graph if in_degree[node] == 0])
    result = []
    while queue:
        node = queue.popleft()
        result.append(node)
        for neighbor in graph.get(node, []):
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)
    
    if len(result) != len(graph):
        raise ValueError("Cycle detected")
    return result

def main():
    ws_dir = 'docs/workstreams/backlog'
    files = sorted([f for f in os.listdir(ws_dir) if f.startswith('00-001-') and f.endswith('.md')])
    
    graph = defaultdict(list)
    nodes = {}
    
    for f in files:
        ws_id, deps = parse_workstream(os.path.join(ws_dir, f))
        nodes[ws_id] = deps
        # Reverse edges: dependency -> ws_id
        for dep in deps:
            graph[dep].append(ws_id)
        # Ensure node exists in graph
        if ws_id not in graph:
            graph[ws_id] = []
    
    # Add missing nodes (dependencies that aren't workstreams themselves)
    all_deps = set()
    for deps in nodes.values():
        all_deps.update(deps)
    for dep in all_deps:
        if dep not in graph:
            graph[dep] = []
    
    print("Dependency graph:")
    for ws in sorted(nodes.keys()):
        print(f"{ws}: depends on {nodes[ws] if nodes[ws] else 'none'}")
    
    try:
        order = topological_sort(graph)
        print("\nTopological execution order:")
        for i, ws in enumerate(order, 1):
            print(f"{i:2d}. {ws}")
    except ValueError as e:
        print(f"\nError: {e}")
        sys.exit(1)
    
    return order

if __name__ == '__main__':
    main()