# Gd - command line tool for Google Drive

## Overview
gd is a command line tool for Google Drive. It allows you to view and list file on Google Drive.

## Installation
```bash
make build;
./gd --help;
```

## Usage
1. List file in a folder which id is abc123
```bash
./gd ls abc123;
```

2. ls can recursively list all files in a folder by using -r option
```bash
./gd ls -r abc123;
```

