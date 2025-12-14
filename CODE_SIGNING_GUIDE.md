# Code Signing & Notarization Guide

This guide walks you through setting up code signing and notarization for seamless auto-updates.

## Overview

Code signing makes your app trusted by macOS:
- ✅ No security warnings
- ✅ No "damaged app" errors
- ✅ Seamless auto-updates work perfectly
- ✅ Professional user experience

## Step 1: Get Apple Developer Account

1. Go to https://developer.apple.com/programs/enroll/
2. Sign up for the Apple Developer Program
3. Pay $99/year (renews annually)
4. Wait for approval (usually 24-48 hours)

## Step 2: Create Developer ID Certificate

Once your account is approved:

1. Go to https://developer.apple.com/account/resources/certificates/list
2. Click **"+"** to create a new certificate
3. Select **"Developer ID Application"** (for apps distributed outside Mac App Store)
4. Follow the instructions to create a Certificate Signing Request (CSR):
   - Open **Keychain Access** on your Mac
   - Menu: Keychain Access → Certificate Assistant → Request a Certificate from a Certificate Authority
   - Enter your email
   - Select "Saved to disk"
   - Click Continue
5. Upload the CSR file to Apple's website
6. Download the certificate (**.cer** file)
7. Double-click to install it in Keychain Access

## Step 3: Export Certificate for GitHub Actions

To automate signing in GitHub Actions, you need to export your certificate:

1. Open **Keychain Access**
2. Find your "Developer ID Application" certificate
3. Right-click → Export
4. Save as **.p12** file
5. Set a password (you'll need this for GitHub secrets)
6. Keep this file secure!

## Step 4: Set Up GitHub Secrets

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

1. **APPLE_CERTIFICATE_BASE64**
   ```bash
   # Encode your .p12 file to base64
   base64 -i YourCertificate.p12 | pbcopy
   # Paste the output as the secret value
   ```

2. **APPLE_CERTIFICATE_PASSWORD**
   - The password you set when exporting the .p12 file

3. **APPLE_ID**
   - Your Apple ID email (e.g., your@email.com)

4. **APPLE_ID_PASSWORD**
   - Create an **app-specific password** at https://appleid.apple.com
   - Apple ID → Security → App-Specific Passwords → Generate
   - Use this, NOT your regular Apple ID password

5. **APPLE_TEAM_ID**
   - Find at https://developer.apple.com/account
   - Click "Membership" in the sidebar
   - Copy your Team ID (10 characters, e.g., "ABCD123456")

## Step 5: Update Build Script

The GitHub Actions workflow needs to:
1. Import the certificate
2. Sign the app
3. Notarize with Apple
4. Staple the notarization ticket

This will be done in the updated workflow.

## Step 6: Local Testing (Optional)

To test signing locally before pushing to GitHub:

```bash
# Find your certificate identity
security find-identity -v -p codesigning

# Sign the app
codesign --force --deep --sign "Developer ID Application: Your Name (TEAMID)" \
  "build/darwin/O&M Game Room.app"

# Verify the signature
codesign --verify --verbose "build/darwin/O&M Game Room.app"
spctl --assess --verbose "build/darwin/O&M Game Room.app"

# Notarize (requires credentials)
xcrun notarytool submit "build/darwin/OandM_Game_Room.dmg" \
  --apple-id "your@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password" \
  --wait

# Staple the ticket
xcrun stapler staple "build/darwin/O&M Game Room.app"
```

## Step 7: Enable Auto-Updates

Once code signing is working:
1. The app will be fully trusted by macOS
2. Auto-updates can download and install seamlessly
3. Users just click "Update" and it happens automatically
4. No security warnings, no manual steps

## Troubleshooting

### Certificate Not Found
- Make sure the certificate is in your **login** keychain
- Run `security find-identity -v -p codesigning` to verify

### Notarization Failed
- Check that you're using an **app-specific password**, not your regular password
- Verify your Team ID is correct
- Check notarization logs: `xcrun notarytool log <submission-id>`

### "Developer cannot be verified"
- This means the app isn't notarized yet
- Complete the notarization step
- Staple the ticket to the app

## Timeline

- **Day 1**: Sign up for Apple Developer Program
- **Day 2-3**: Wait for approval
- **Day 3**: Create certificate, set up GitHub secrets
- **Day 3**: Update workflow and test
- **Day 4**: Seamless auto-updates working! 🎉

## Cost

- **One-time**: $99/year for Apple Developer Program
- **Ongoing**: Auto-renews at $99/year

## Next Steps

1. ✅ Sign up for Apple Developer account
2. ⏳ Wait for approval
3. 📝 Create certificate
4. 🔐 Set up GitHub secrets
5. 🚀 Deploy signed app with auto-updates

---

**Ready to start?** The first step is signing up at https://developer.apple.com/programs/enroll/
