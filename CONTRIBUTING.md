# Contributing Guidelines

## Build

To be able to build a project which uses private repositories:
```bash
# Run this command
echo 'GITLAB_USER=\nGITLAB_TOKEN=\n' > .env
```

Then fill the .env file with your Gitlab username and an acces token with read_repository scope (User > Settings > Access Tokens).
