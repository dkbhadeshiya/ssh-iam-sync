SSH Public Key Sync for IAM Users
---------------------------------

## Purpose
- The purpose of the project is to sync SSH public files from AWS IAM users to a server.
- Use AWS IAM Groups to manage access to your servers. 
- The project is built using NodeJS.

## Introduction
- `index.js` file is the only file in the project which acts as a command line tool. 
- The script accepts following arguments:
    - `-g` or `--groups`: List of IAM groups to fetch user list from. 
    - `-S` or `--ssh-path`: Absolute path of `authorized_keys` file to write to. Default is `~/.ssh/authorized_keys`
    - `-f` or `--force` : Overwrite current values in `authorized_keys` and replace it every time. 
    - `-L` or `--log-level` : Log4JS log level. Increasing log level might print sensitive information on console. 
    Handle with care.
    
## Libraries used: 
| Library Name | Version|
|--------------|--------|
| yargs        | 15.3.1 |
| aws-sdk      | 2.633.0|
| log4js       | 6.2.1  |


## Ideal server setup
- Use cron to run the script every 10 minutes or so. 
- Assign IAM role to your EC2 instance and use `EC2MetadataCredentials` instead.

## Scope of improvements
- This only syncs keys for a single user on server.
- Multi user syncs can be done manually by setting up multiple CRON jobs to sync different groups.
- Tested only on Linux servers. Windows servers might work but not tested. 
