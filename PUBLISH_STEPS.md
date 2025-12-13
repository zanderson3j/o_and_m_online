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
```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

## Step 5: Upload DMG to Release
1. Go to https://github.com/zanderson3j/o_and_m_online/releases
2. Click "Create a new release"
3. Choose tag: v1.0.0
4. Release title: "O&M Game Room v1.0.0"
5. Copy contents from RELEASE_NOTES.md
6. Attach the DMG file from `build/darwin/`
7. Click "Publish release"

## Your Download Links

Once published, users can download from:
- Latest release: https://github.com/zanderson3j/o_and_m_online/releases/latest
- Direct DMG: Will be shown on the release page

## Testing Auto-Update

1. Install the current version
2. Edit `updater.go` to change version to "1.0.1"
3. Build and create a new release
4. Launch the v1.0.0 app - it should log "New version available"

## Future Releases

For next versions:
1. Update version in `updater.go`
2. Commit changes
3. Tag: `git tag v1.0.1 && git push origin v1.0.1`
4. Upload new DMG to the GitHub release