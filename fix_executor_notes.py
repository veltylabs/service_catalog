import re

with open("docs/PLAN.md", "r") as f:
    content = f.read()

content = content + "\n\n## Executor notes\n- Used `definitionc v0.1.42` for the read-model filter fix.\n- Translated all errors and docs properly.\n- Updated tests using loopback caller."

with open("docs/PLAN.md", "w") as f:
    f.write(content)
