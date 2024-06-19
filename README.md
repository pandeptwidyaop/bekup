<br>
<img src="https://raw.githubusercontent.com/pandeptwidyaop/bekup/main/docs/bekup-dall-e.webp" width="64"> </img>

# Bekup #

> Under Development! 

Bekup is a tool for automating database backups and uploading them to cloud storage. The databases currently supported are: MySQL, MariaDB, PostgreSQL, MongoDB and Redis. As for the cloud storage options supported, they include AWS S3, MinIO, FTP, and SFTP.

## Features

Bekup is fully developed using GO with goroutines to carry out the dump process through to upload. Bekup can directly perform multi-database and multi-DBMS backup executions at once, which can be dynamically configured in the `config.json` file, and it also supports uploading to multiple cloud servers simultaneously. 

Bekup has four main processes, which are:

1. Dump 
2. Zip
3. Upload
4. Cleanup

## Requirements

To run Bekup, there are several tools that currently need to be installed first, namely:

- For backing up MySQL or MariaDB, `mysqldump` is required.
- For backing up PostgreSQL, `pg_dump` is needed.
- For backing up MongoDB, `mongodump` is necessary.
- For backing up Redis, `redis-cli` is exist.

For convenience, you can use a ready-to-use Docker image, [See docker instruction](#for-docker-user).

## How to use

Simply, you just need to download Bekup from the release page, then run Bekup with the command : 

`./bekup --config=path/to/config.json`

To run backup regulary using crontab you can add this command:

`crontab -e`

`* * 1 * * /path/to/bekup --config=/path/to/config.json`

## Configuration

To run Bekup, a configuration file with a `json` extension is required. You can download a sample of it [here](/configs/example.config.json).

| Json Key | Type|Description |
|----------|-----|-------|
| `sources.*.driver` | `string` |Database driver, available options : `mysql`,`postgres`,`mongodb`,`redis`|
| `sources.*.host`|`string` |Database host |
| `sources.*.port`|`string` |Database port |
| `sources.*.username` |`string`| Database username |
| `sources.*.password` |`string`| Database password |
| `sources.*.mongodb_uri`|`string`| For mongodb driver only, if `mongodb_uri` defined will ignore other host,port,username and password. Example `mongodb://username:password@host:port` (without trailing database). By default, the authSource is automatically defined to `admin` database.|
|`source.*.databases.*`| `string`| Database name want to backup |
| `destinations.*.driver`|`string`|Driver options for backup destination, now options are `s3`,`ftp`,`sftp`|
|`destinations.*.aws_access_key`|`string`|Your aws access key, required if using `s3` driver|
|`destinations.*.aws_secret_key`|`string`|Yout aws access secret key, required if using `s3` driver|



## For Docker User

## TODO 

- [x] MySQL & MariaDB Driver
- [x] PostgreSQL Driver
- [x] MongoDB Driver
- [x] Redis Driver
- [x] AWS S3 & MinIO Driver
- [ ] FTP & SFTP Driver
- [ ] API Documentation
- [ ] Dockerfile
- [ ] Add password to archive zip file