# Contributing Guidelines

## Build

To be able to build a project which uses private repositories:
```bash
# Run this command
export SSH_KEY=`cat ~/.ssh/id_rsa`

# Or add is to your .bashrc
echo 'export SSH_KEY=`cat ~/.ssh/id_rsa`' >> ~/.bashrc
```

*Note: if your private repository ssh key is not id_rsa, replace it in the above command.*