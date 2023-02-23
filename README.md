# Golang VSS module
Windows API bindings for the `Volume Shadow Copy Service` in Golang for 32 and 64-bit systems. Enables the user to duplicate entire drives during runtime without any file access issues. The API bindings are accompanied by a simple CLI tool that creates and symlinks Shadow Copies of a given drive.
## Build
You can either import the vss api bindings into your project or use the CLI application. The CLI application can be called with the following command:
```shell
make build
```

## Usage
```sh
./vss -h
```
```
usage: 

  -D string
        Drive letter to copy (required)
  -S string
        Path of symlink folder that points to the snapshot
  -f    Creates snapshots if available shadow storage is low. Could delete old copies!
  -timeout int
        Snapshot creation timeout in seconds (default 30)
```