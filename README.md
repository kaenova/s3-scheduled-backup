# S3 Scheduled Docker Volumes Backup

Application to regularly backup your `docker volumes` to S3 storage. This application is till hardcoded to backup everyday on midnight. There's 2 modes to run this application. Please set up the environment variables using variables below:
```env
# Application Config
# 2 Modes: "docker" or "local"
MODE=

# Backup Service Config
# If in "docker" mode, PATH_BACKUP will be ignored
MAXIMUM_BACKUP_WINDOW=
PATH_BACKUP=

# S3 Config
S3_ENDPOINT=
S3_BUCKET_NAME=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_USE_SSL=
``` 

This application will bind to your `/var/lib/docker/volumes` and will zips, and upload every child directory from that path to S3 storage.

**NOTE**: This application only tested on Ubuntu >20.04 system.