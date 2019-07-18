# Comics-DL
[![license](https://img.shields.io/github/license/The-Eye-Team/Comics-DL.svg)](https://github.com/The-Eye-Team/Comics-DL/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/302796547656253441.svg)](https://discord.gg/py3kX3Z)

A comic and hentai archiver that saves as `.cbz` with support for the following sites:

- https://readcomicsonline.ru/
- https://www.tsumino.com/
- https://e-hentai.org/
- https://myreadingmanga.info/

## Download
```
go get -u github.com/The-Eye-Team/Comics-DL
```

## Usage
```
./Comics-DL --url {URL}
```

### Flags
|      | Name | Default Value | Description |
|------|------|---------------|-------------|
| `-u` | `--url` | Required. | URL of comic to save. |
| `-c` | `--concurrency` | `4` | The number of simultaneous downloads to run. |
| `-o` | `--output-dir` | `./results/` | Path to directory to save files to. |
| `-k` | `--keep-jpg` | `false` | Flag to keep/destroy the `.jpg` page data. |

## License
MIT.
