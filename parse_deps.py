#!/usr/bin/env python3
import os
import re
import sys
from collections import defaultdict, deque

def parse_dependencies(ws_file):
    """Parse workstream file and return list of dependencies."""
    deps = []
    ws_id = None
    in_frontmatter = False
    with open(ws_file, 'r') as f:
        lines = f.readlines()
    
    i = 0
    while i < len(lines):
        line = lines[i].strip()
        if line == '---':
            in_frontmatter = not in_frontmatter
            i += 1
            continue
        if in_frontmatter:
            if line.startswith('ws_id:'):
                ws_id = line.split(':', 1)[1].strip()
            i += 1
            continue
        # Outside frontmatter
        if line.startswith('## Dependencies'):
            i += 1
            # Read next lines until next section
            while i < len(lines):
                dep_line = lines[i].strip()
                if dep_line.startswith('##'):
                    break
                if dep_line.startswith('- '):
                    match = re.search(r'00-\d{3}-\d{2}', dep_line)
                    if match:
                        deps.append(match.group(0))
                i += 1
            continue
        i += 1
    if ws_id is None:
        # fallback: extract from filename
        ws_id = os.path.basename(ws_file).replace('.md', '')
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
    pattern = '00-001-*.md'
    files = sorted([f for f in os.listdir(ws_dir) if f.startswith('00-001-') and f.endswith('.md')])
    
    graph = defaultdict(list)
    nodes = {}
    
    for f in files:
        ws_id, deps = parse_dependencies(os.path.join(ws_dir, f))
        nodes[ws_id] = deps
        # Reverse edges: dependency -> ws_id
        for dep in deps:
            graph[dep].append(ws_id)
        # Ensure node exists in graph
        if ws_id not in graph:
            graph[ws_id] = []
    
    # Compute topological order
    try:
        order = topological_sort(graph)
        print("Topological order:", ' -> '.join(order))
        # Verify all nodes present
        missing = set(nodes.keys()) - set(order)
        if missing:
            print("Warning: missing nodes:", missing)
    except ValueError as e:
        print("Error:", e)
        sys.exit(1)
    
    # Print dependency tree
    print("\nDependency graph:")
    for ws in sorted(nodes.keys()):
        print(f"{ws}: depends on {nodes[ws] if nodes[ws] else 'none'}")

if __name__ == '__main__':
    main()