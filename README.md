# Comics-DL
[![license](https://img.shields.io/github/license/The-Eye-Team/Comics-DL.svg)](https://github.com/The-Eye-Team/Comics-DL/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/302796547656253441.svg)](https://discord.gg/the-eye)

A comic and doujinshi archiver that saves as `.cbz` with support for the following sites:

- https://readcomicsonline.ru/
- https://e-hentai.org/
- https://myreadingmanga.info/
- https://doujins.com/
- https://nhentai.net/

## Download
```
go get -u github.com/The-Eye-Team/Comics-DL
```

## Usage
```
./Comics-DL --url {URL}
./Comics-Dl --file {path}
```

### Flags
|      | Name | Default Value | Description |
|------|------|---------------|-------------|
| `-u` | `--url` | None. | URL of comic to save. |
| `-c` | `--concurrency` | `20` | The number of simultaneous downloads to run. |
| `-o` | `--output-dir` | `./results/` | Path to directory to save files to. |
| `-k` | `--keep-jpg` | `false` | Flag to keep/destroy the `.jpg` page data. |
| `-f` | `--file` | None. | Path to file with list of URLs to download. |

## License
MIT.
