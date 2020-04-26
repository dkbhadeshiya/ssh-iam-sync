#! /usr/bin/env node

const AWS = require("aws-sdk");
const yargs = require('yargs');
const log4js = require("log4js");
const path = require("path")
const logger = log4js.getLogger("SyncIAMKeys");
const fs = require("fs");

const args = yargs
    .option("g", {
        alias: "groups",
        demandOption: true,
        describe: "Pass the IAM Groups that are allowed to access this server",
        type: "array"
    })
    .option("S", {
        alias: "ssh-path",
        default: "~/.ssh/authorized_keys",
        describe: "Path to authorized_keys file.",
        type: "string"
    })
    .option("f", {
        alias: "force",
        default: false,
        describe: "Force overwrite existing key files. Warning: this will wipe your current ssh key files",
        type: "boolean"
    })
    .option("L", {
        alias: "log-level",
        default: "info",
        describe: "Log level to print. Supports log4js log levels"
    })
    .argv

logger.level = args.logLevel;

process(args).then().catch(e => console.error(e));

async function getUsers(groups, IAM) {
    return new Promise((resolve, reject) => {
        let Users = [];
        let promises = [];
        groups.forEach(group => {
            promises.push(IAM.getGroup({GroupName: group}).promise().then(data => {
                Users.push(...data.Users)
            }).catch(e => console.error(e)));
        });
        Promise.all(promises).then(d => resolve(Users)).catch(e => reject(e));
    })

}

async function getKeyIds(users, IAM) {
    return new Promise((resolve, reject) => {
        const promises = [];
        const KeyIds = [];
        users.forEach(user => {
            promises.push(IAM.listSSHPublicKeys({
                UserName: user.UserName,
            }).promise().then(res => {
                KeyIds.push(...res.SSHPublicKeys.map(e => {
                    return {SSHPublicKeyId: e.SSHPublicKeyId, UserName: e.UserName, Encoding: 'SSH'}
                }))
            }).catch(e => console.error(e)));
        });
        Promise.all(promises).then(d => resolve(KeyIds)).catch(e => reject(e))
    })
}

async function getKeyBody(KeyIds, IAM) {
    return new Promise((resolve, reject) => {
        const keyBodies = [];
        const promises = [];
        KeyIds.forEach(keyid => {
            promises.push(IAM.getSSHPublicKey(keyid).promise().then(res => {
                if(res.SSHPublicKey.Status === 'Active')
                    keyBodies.push(res.SSHPublicKey.SSHPublicKeyBody)
            }).catch(e => console.error(e)));
        });
        Promise.all(promises).then(d => resolve(keyBodies)).catch(e => reject(e));
    })
}

async function process(args) {
    // const credentials = new AWS.EC2MetadataCredentials();
    const credentials = new AWS.SharedIniFileCredentials({profile: "default"});
    AWS.config.credentials = credentials;

    logger.info("Starting syncing keys...");
    logger.debug("Using AWS access key: "+ credentials.accessKeyId);
    const IAM = new AWS.IAM();

    logger.debug("Fetching Users");
    const users = await getUsers(args.groups, IAM);
    logger.trace("Got following users", users);

    logger.debug("Fetching key ids for Users: " + users.length);
    const keyIds = await getKeyIds(users, IAM);
    logger.trace("Got following key ids: ", keyIds);

    logger.debug("Fetching key bodies for Keys: " + keyIds.length);
    const keys = await getKeyBody(keyIds, IAM);
    logger.trace("Got following key bodies: ", keys);

    logger.info("Got " + keys.length + " key(s) from "+ users.length + " user(s)!")

    logger.info("Writing authorized_key file");

    const filePath = path.normalize(args.sshPath)
    if(args.force) {
        logger.info("Replacing file: "+ filePath)
        fs.writeFileSync(filePath, keys.join("\n") + "\n", {flag: "w", mode: 644});
    } else {
        logger.info("Appending file: "+ filePath)
        fs.writeFileSync(filePath, keys.join("\n") + "\n", {flag: "a", mode: 644});
    }

    logger.info("SSH Keys have been synced")
}
