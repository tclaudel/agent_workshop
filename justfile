# --------------------
# @BMad

# Installs BMad in the current project with predefined settings and updates .gitignore.
install-bmad:
    #!/usr/bin/env bash
    set -euo pipefail
    user_name="Thomas"
    project_name=$(basename "$PWD")
    npx bmad-method install --yes \
        --directory=. \
        --modules bmm,bmb \
        --tools claude-code \
        --set core.user_name="${user_name}" \
        --set core.project_name="${project_name}" \
        --set core.communication_language="English" \
        --set bmm.user_skill_level="intermediate" \
        --all-stable
    entries=(
        "# bmad generated"
        "_bmad/"
        "skills/bmad-*/"
        ".agents/skills/bmad-*/"
        ".claude/skills/bmad-*/"
    )
    for entry in "${entries[@]}"; do
        grep -qxF "${entry}" .gitignore 2>/dev/null || echo "${entry}" >> .gitignore
    done
