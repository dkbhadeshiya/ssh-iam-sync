SSH Public Key Sync for IAM Users
---------------------------------

## Purpose
- The purpose of the project is to sync SSH public files from AWS IAM users to a server.
- Use AWS IAM Groups to manage access to your servers. 
- The project is built using Go v1.20.

## Introduction
- All the sources are under `cmd/ssh-iam-sync`
- Whole code compiles in a binary `ssh-iam-sync`
- Config file needs to be defined in order to run it, see config file below
- Config files will be read from following folders according to priority:
    - `./config.yaml` current folder from where the binary is running
    - `/etc/ssh-iam-sync/config.yaml` from ETC

## Configurations
Here is the config reference file:
```yaml
aws:
  method: accessKey # Either accessKey or profile or instanceProfile
  profile: default
  region: ap-south-1
  accessKey: <your-access-key>
  secretKey: <your-secret-key>
  groups:
    - projec1
    - project2

authorizedKeys: ~/.ssh/authorized_keys  # Path to authorized key file
overwrite: true # Overwrite existing key file, false appends the keys to file
```

    
## Libraries used: 
| Library Name                             | Version  |
|------------------------------------------|----------|
| github.com/aws/aws-sdk-go-v2             | v1.32.3  |
| github.com/aws/aws-sdk-go-v2/config       | v1.28.1  |
| github.com/aws/aws-sdk-go-v2/credentials | v1.17.42 |
| github.com/aws/aws-sdk-go-v2/service/iam | v1.37.3  |
| github.com/kkyr/fig                       | v0.4.0   |


## Ideal server setup
- Use cron to run the binary every 10 minutes or so. 
- Assign IAM role to your EC2 instance and use `instanceProfile` method instead.

## Scope of improvements
- Binary can be compiled for windows servers as well 
- Distribute using package manager
