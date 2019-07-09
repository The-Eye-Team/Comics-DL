# Comics-DL
[![license](https://img.shields.io/github/license/The-Eye-Team/Comics-DL.svg)](https://github.com/The-Eye-Team/Comics-DL/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/302796547656253441.svg)](https://discord.gg/py3kX3Z)

A comic and hentai archiver that saves as `.cbz` with support for the following sites:

- https://readcomicsonline.ru/
- https://www.tsumino.com/
- https://e-hentai.org/

## Download
```
go get -u github.com/The-Eye-Team/Comics-DL
```

## Usage
```
./Comics-DL --url {URL}
```

### Flags
| Name | Default | Description |
|------|---------|-------------|
| `--url` | Required. | URL of comic to save. |
| `--concurrency` | `4` | The number of simultaneous downloads to run. |
| `--output-dir` | `./results/` | Path to directory to save files to. |
| `--keep-jpg` | `false` | Flag to keep/destroy the `.jpg` page data. |

## License
MIT.
