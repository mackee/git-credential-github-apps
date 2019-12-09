# git-credential-github-apps

A git credential helper with GitHub Apps

## Overview

`git-credential-github-apps` provides authentication behavior in GitHub Apps on git commands.

This command returns credentials that GitHub Token. Also, that response contains a cached token while during in not expire.

`git-credential-github-apps` is work as git-credential-helper. If you want to know more details, see the `Install` section in this document and [this document](https://git-scm.com/docs/api-credentials).

## Install

Download latest version from [Releases](https://github.com/mackee/git-credential-github-apps/releases).

Extract into a directory that written in your PATH environment variable.

## Usage

### Prepare

Using this tool requires a private key, App ID and Installation ID or organization name.

You will get a private key and App ID on the Config page at GitHub Apps.

More details for private key: [Generating a private key](https://developer.github.com/apps/building-github-apps/authenticating-with-github-apps/#generating-a-private-key)

Installation ID is the identifier of installation organization on GitHub Apps

Organization name can be alternate to installation ID. `git-credential-github-apps` detect installation ID from organization name.

### Install to git

Type following this. This is set credential helper to git configuration in global.

```console
$ git config --global credential.helper 'github-apps -privatekey <path to private key> -appid <App ID of GitHub Apps> -login <installation organization>'
```

If you want to set to repository local, you will type following this on directory of the repository.

```console
git config --global credential.helper 'github-apps -privatekey <path to private key> -appid <App ID of GitHub Apps> -login <installation organization>'
```

### More Options

If you want to know more options, execution `git-credential-github-apps` with `-h`.

## Author

mackee, KAYAC Inc.

## License

MIT.
