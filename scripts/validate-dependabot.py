import yaml, sys, os

with open(os.path.join(os.path.dirname(__file__), '..', '.github', 'dependabot.yml')) as f:
    data = yaml.safe_load(f)

print(f"YAML syntax: OK")
print(f"Updates count: {len(data.get('updates', []))}")
for u in data.get('updates', []):
    directory = u['directory']
    full_path = os.path.join(os.path.dirname(__file__), '..', directory.lstrip('/'))
    has_manifest = os.path.exists(os.path.join(full_path, 'go.mod')) if u['package-ecosystem'] == 'gomod' else True
    status = "OK" if has_manifest else "MISSING go.mod"
    print(f"  {u['package-ecosystem']:20s} {directory:30s} [{status}]")
