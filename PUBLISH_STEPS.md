# Publishing O&M Game Room to GitHub

## Step 1: Initialize Git (if needed)
```bash
git init
git add .
git commit -m "Initial release of O&M Game Room v1.0.0"
```

## Step 2: Create GitHub Repository
1. Go to https://github.com/new
2. Name: `o_and_m_online`
3. Description: "Olive & Millie's Game Room - Online multiplayer games"
4. Keep it PUBLIC (so people can download)
5. Don't initialize with README (we already have one)
6. Click "Create repository"

## Step 3: Push Your Code
```bash
git remote add origin https://github.com/zanderson3j/o_and_m_online.git
git branch -M main
git push -u origin main
```

## Step 4: Create Your First Release

GitHub Actions will automatically build and release when you push a tag:

```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions will:
- Build DMG files for Intel and Apple Silicon
- Create a GitHub release
- Upload both DMG files automatically
- No manual steps required!

## Your Download Links

Once published, users can download from:
- Latest release: https://github.com/zanderson3j/o_and_m_online/releases/latest
- Direct DMG: Will be shown on the release page

## Testing Auto-Update

1. Install the current version (e.g., v1.0.0)
2. Update version in `updater.go` and `Info.plist` to "1.0.1"
3. Commit, tag, and push: `git tag v1.0.1 && git push --tags`
4. Launch the v1.0.0 app
5. A green notification will appear in the lobby (bottom-left)
6. Click it to open browser to download page

## Future Releases

For next versions:
1. Update version in `updater.go` and `resources/Info.plist`
2. Commit changes: `git commit -m "Bump version to X.X.X"`
3. Tag and push: `git tag vX.X.X && git push --tags`
4. GitHub Actions automatically builds and publishes the release
5. Users get notified of the update in the lobby