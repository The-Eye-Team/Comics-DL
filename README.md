# Comics-DL
[![license](https://img.shields.io/github/license/The-Eye-Team/Comics-DL.svg)](https://github.com/The-Eye-Team/Comics-DL/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/302796547656253441.svg)](https://discord.gg/the-eye)
[![CircleCI](https://circleci.com/gh/The-Eye-Team/Comics-DL.svg?style=svg)](https://circleci.com/gh/The-Eye-Team/Comics-DL)
[![goreportcard](https://goreportcard.com/badge/github.com/The-Eye-Team/Comics-DL)](https://goreportcard.com/report/github.com/The-Eye-Team/Comics-DL)

A comic and doujinshi archiver that saves as `.cbz` with support for the following sites:

- https://readcomicsonline.ru/
- https://e-hentai.org/
- https://myreadingmanga.info/
- https://doujins.com/
- https://nhentai.net/
- https://pururin.io/

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

Files will be save to the location `{--output-dir}/{domain}/{title}.cbz

## Builds
[![circleci](https://circleci.com/gh/The-Eye-Team/Comics-DL.svg?style=svg)](https://circleci.com/gh/The-Eye-Team/Comics-DL)

Pre-compiled binaries are published on Circle CI at https://circleci.com/gh/The-Eye-Team/Comics-DL. To download a binary, navigate to the most recent build and click on 'Artifacts'. Here there will be a list of files. Click on the one appropriate for your system.

Once downloaded, run the following with the values applicable to you.
```
$ ./comics-dl-{build_date}-{git_tag}-{os}-{arch}
```

## License
MIT.
