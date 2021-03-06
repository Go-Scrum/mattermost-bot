// This file is automatically generated. Do not modify it manually.

const manifest = JSON.parse(`
{
    "id": "io.goscrum.mattermost-bot",
    "name": "GoScrum Bot",
    "description": "GoScrum Bot.",
    "version": "0.1.0",
    "min_server_version": "5.12.0",
    "server": {
        "executables": {
            "linux-amd64": "server/dist/plugin-linux-amd64",
            "darwin-amd64": "server/dist/plugin-darwin-amd64",
            "windows-amd64": "server/dist/plugin-windows-amd64.exe"
        },
        "executable": ""
    },
    "webapp": {
        "bundle_path": "webapp/dist/main.js"
    },
    "settings_schema": {
        "header": "",
        "footer": "",
        "settings": [
            {
                "key": "URL",
                "display_name": "GoScrum API URL",
                "type": "text",
                "help_text": "The base URL for using the plugin with a GoScrum installation. Examples: https://api.goscrum.io",
                "placeholder": "",
                "default": "https://api.goscrum.io"
            },
            {
                "key": "Token",
                "display_name": "API Token",
                "type": "text",
                "help_text": "Token to access GoScrum.io api",
                "placeholder": "",
                "default": null
            }
        ]
    }
}
`);

export default manifest;
export const id = manifest.id;
export const version = manifest.version;
