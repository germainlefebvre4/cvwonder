# LinkedIn Integration

CVWonder now supports converting your LinkedIn profile to YAML format using the `convert` command.

## Prerequisites

1. Create a LinkedIn Application to get your Client ID and Client Secret:
   - Go to [LinkedIn Developers](https://www.linkedin.com/developers/)
   - Create a new app
   - Add OAuth 2.0 redirect URLs (e.g., `http://localhost:8080/callback`)
   - Note your Client ID and Client Secret

2. Required OAuth scopes:
   - `openid`
   - `profile`
   - `email`
   - `w_member_social`

## Usage

### Basic Usage

To convert a LinkedIn profile, you must specify the profile username:

```bash
cvwonder convert \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --profile germainlefebvre4 \
  --output cv.yml
```

The profile username is the part after `/in/` in the LinkedIn URL. For example:
- URL: `https://www.linkedin.com/in/germainlefebvre4/`
- Username: `germainlefebvre4`

### Convert Your Own Profile

To convert your own profile, use your LinkedIn username:

```bash
cvwonder convert \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --profile YOUR_USERNAME \
  --output cv.yml
```

### Fetch Another User's Profile

You can fetch any public LinkedIn profile by specifying their username:

```bash
cvwonder convert \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --profile johndoe \
  --output cv.yml
```

### With Custom Redirect URI

```bash
cvwonder convert \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --redirect-uri http://localhost:9000/callback \
  --output my-cv.yml
```

### Complete Example with All Options

```bash
cvwonder convert \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --profile johndoe \
  --redirect-uri http://localhost:8080/callback \
  --output johndoe-cv.yml
```

### Using Aliases

The convert command supports aliases for convenience:

```bash
cvwonder c --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET --profile username

cvwonder conv --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET -p username
```

## Authentication Flow

1. The command starts a local HTTP server to handle the OAuth callback
2. The command generates an authorization URL
3. You open the URL in your browser and authorize the application
4. LinkedIn redirects you back to the local server with an authorization code
5. The server automatically captures the code and shuts down
## Example Session

## Example Session

### Converting a LinkedIn Profile

```
$ cvwonder convert --client-id abc123 --client-secret xyz789 --profile germainlefebvre4

CV Wonder - LinkedIn Converter
  Client ID: abc...123
  Redirect URI: http://localhost:8080/callback
  Output file: cv.yml

Step 1: Authorize the application
Please open the following URL in your browser:

https://www.linkedin.com/oauth/v2/authorization?client_id=abc123&...

Waiting for authorization callback...

Step 2: Exchanging authorization code for access token...
✓ Successfully obtained access token

Step 3: Fetching LinkedIn profile for user: germainlefebvre4
✓ Successfully fetched profile for: Germain Lefebvre

Step 4: Converting profile to YAML...
✓ Successfully converted and saved to: cv.yml

You can now use 'cvwonder generate' to create your CV from this file.
``` can now use 'cvwonder generate' to create your CV from this file.
```

## What Gets Converted

The LinkedIn converter extracts the following information from your profile:

- **Personal Information**: Name, headline, location, email, phone
- **Work Experience**: All positions with titles, companies, dates, and descriptions
- **Education**: Schools, degrees, fields of study, dates
- **Skills**: Professional skills
- **Languages**: Spoken languages with proficiency levels
- **Projects**: Side projects with descriptions and links
- **Certifications**: Professional certifications with issuers and dates

## Generated YAML Structure

The converter generates a YAML file compatible with CVWonder's CV model:

```yaml
person:
  name: John Doe
  profession: Software Engineer
  location: San Francisco, CA
  email: john.doe@example.com

socialNetworks:
  linkedin: https://www.linkedin.com/in/johndoe

career:
  - companyName: Tech Corp
    missions:
      - position: Senior Developer
        company: Tech Corp
        location: New York, NY
        dates: 2020-01 - 2023-12
        summary: Developed applications

education:
  - schoolName: MIT
    degree: Bachelor of Science in Computer Science
    dates: 2015 - 2019

technicalSkills:
  domains:
    - name: Technical Skills
      competencies:
        - name: Go
          level: 3
        - name: Python
          level: 3

### Profile data missing
- Check that you've granted all required OAuth scopes
- Some profile fields may be private or empty on LinkedIn
- When fetching another user's profile, only public information is accessible
- The target profile must have public visibility settings enabled
```

## Next Steps

After converting your LinkedIn profile to YAML:

1. Review and edit the generated `cv.yml` file
2. Generate your CV: `cvwonder generate -i cv.yml`
3. Serve it locally: `cvwonder serve -i cv.yml`
4. Export to PDF: `cvwonder generate -i cv.yml -f pdf`

## Troubleshooting

### Authorization fails
- Ensure your redirect URI matches what's configured in your LinkedIn app
- Verify your Client ID and Client Secret are correct
- Make sure the port (default 8080) is not already in use

### Callback server won't start
- Check if another application is using port 8080
- Use a different port with `--redirect-uri http://localhost:9000/callback`
- Make sure to update the redirect URI in your LinkedIn app settings

### Browser doesn't redirect back
- Check your LinkedIn app's redirect URI configuration
- Ensure the redirect URI exactly matches what you're using (including http vs https)
- Clear your browser cache and try again

### Profile data missing
- Check that you've granted all required OAuth scopes
- Some profile fields may be private or empty on LinkedIn

### API rate limits
- LinkedIn API has rate limits; wait before retrying if you hit them

## Security Notes

- Never commit your Client ID and Client Secret to version control
- Use environment variables for sensitive credentials:
  ```bash
  export LINKEDIN_CLIENT_ID=your_client_id
  export LINKEDIN_CLIENT_SECRET=your_client_secret
  cvwonder convert --client-id $LINKEDIN_CLIENT_ID --client-secret $LINKEDIN_CLIENT_SECRET
  ```
- The access token is only used during the conversion and is not stored
