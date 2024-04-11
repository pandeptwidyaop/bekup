<br>
<img src="https://raw.githubusercontent.com/pandeptwidyaop/bekup/main/docs/bekup-dall-e.webp" width="64"> </img>

# Bekup #

> Under Development! 

Bekup is a tool for automating database backups and uploading them to cloud storage. The databases currently supported are: MySQL, MariaDB, PostgreSQL, and MongoDB. As for the cloud storage options supported, they include AWS S3, MinIO, FTP, and SFTP.

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

For convenience, you can use a ready-to-use Docker image, [See docker instruction](#for-docker-user).

## How to use

Simply, you just need to download Bekup from the release page, then run Bekup with the command `./bekup --config=path/to/config.json`

## Configuration


## For Docker User

## TODO 

- [x] MySQL & MariaDB Driver
- [x] PostgreSQL Driver
- [x] MongoDB Driver
- [x] AWS S3 & MinIO Driver
- [ ] FTP & SFTP Driver
- [ ] API Documentation
- [ ] Dockerfile
- [ ] Add password to archive zip file