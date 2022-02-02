# `make-doordash-jwt`: A command-line tool to generate a DoorDash Developer JWT

This simple CLI tool takes a DoorDash Access Key (a JSON object comprised of a `developerId`, `key_id`, and `signing_secret` and creates a JWT.

## How to install

If you have Go installed:

```sh
go install github.com/infin8x/make-doordash-jwt
```

If you don't have Go installed:

1. Select the latest release from the sidebar on the right (look for "v1.0.[something]")
1. Download the version appropriate for your system
    1. Tip: `darwin-amd64` is probably what you want if you have a Mac; `windows-amd64.exe`, if you have a Windows machine.

## How to use

1. Create an [Access Key](https://developer.doordash.com/portal/integration/drive/credentials) in the DoorDash Developer Portal
1. Click "Copy" to copy your Access Key details
1. Save them to a file (e.g. `key.json`) **not in source control!**
1. Open a shell, navigate to the directory where you saved `key.json` and run `make-doordash-jwt -f key.json`
    1. If you downloaded the tool instead of running `go install`, you may need to run `.\make-doordash-jwt-darwin-amd64 -f key.json` on Mac or `.\make-doordash-jwt-windows-amd64.exe -f key.json` on PowerShell on Windows
    1. Alternatively, you can provide the Access Key directly on the command line. Make sure to surround it with single quotes ("`").
  
```sh
make-doordash-jwt -o `{
    "developer_id": "14e84291-d900-4c20-8528-ed6ca8de660f",
    "key_id": "d4e87d9c-432b-4ab4-b06f-254ce7a1ef30",
    "signing_secret": "xVu_RIEHqVw0ISBOqklCrKTIrlxX47TiexoJIY8_naw"
  }`
```
