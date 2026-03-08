#!/bin/bash
set -e

REPO_URL="git@github.com:ax958888/skill.git"
BRANCH="main"

echo "=== Deploying to GitHub ==="
echo "Repository: $REPO_URL"
echo "Branch: $BRANCH"
echo ""

# Check if git is initialized
if [ ! -d ".git" ]; then
    echo "Initializing git repository..."
    git init
    git branch -M $BRANCH
else
    echo "Git repository already initialized"
fi

# Check if remote exists
if git remote | grep -q "origin"; then
    echo "Remote 'origin' already exists"
    git remote set-url origin $REPO_URL
else
    echo "Adding remote 'origin'..."
    git remote add origin $REPO_URL
fi

# Add all files
echo "Adding files..."
git add .

# Commit
echo "Committing changes..."
COMMIT_MSG="feat: implement skill-analyzer and skill-builder tools

- skill-analyzer: GitHub repository analysis tool
- skill-builder: Skill building and deployment tool
- Complete documentation and build scripts
- Static binary compilation support"

git commit -m "$COMMIT_MSG" || echo "No changes to commit"

# Push
echo "Pushing to GitHub..."
git push -u origin $BRANCH

echo ""
echo "=== Deployment Complete ==="
echo "Repository: https://github.com/ax958888/skill"
