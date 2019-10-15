# 🥤 `cupt`

`cupt`: the Cognito User Pool Tool.  Simplified AWS Cognito User Pool operations via a CLI, such as:

* Add a pre-verified user with a permanent password.
* List all users.
* Backup all users to a file.
* Restore all users from a backup (with new passwords).

## Installation

* Download the binary from the latest tag.
* Call `cupt` from your favorite terminal emulator.

## Source Installation

* Install Go and set it up.
* Install the latest version of Visual Studio Code.
* Install the Go VS Code plugin.
* Clone this repository.
* In the repository directory, run `go build`

## Usage

* `cupt --help`

*Note* that you must wrap values with special characters in `'single quotes'`.

### Notes on Attributes

The User Pool should be configured such that:

* `username` is a random GUID.
* `email` is a login alias.

Since the `sub` property cannot be managed out of Cognito, we use `username` as a unique identifier to Cognito User Pool users.

### Authentication

`cupt` calls require a path to an AWS JSON configuration file with credentials and region defined.  The syntax of this file is:

```json
{
    "accessKeyId": "xxx",
    "secretAccessKey": "xxx",
    "region": "us-east-1"
}
```

## Building

To build for Windows, Darwin, and Linux (64bit only):

* `./build.sh`

The binaries can be found in the `/builds` directory.

## Contributing

* Ensure Go code is auto-formatted on save.