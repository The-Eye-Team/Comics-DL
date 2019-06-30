# Comics-DL
[![license](https://img.shields.io/github/license/The-Eye-Team/Comics-DL.svg)](https://github.com/The-Eye-Team/Comics-DL/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/302796547656253441.svg)](https://discord.gg/py3kX3Z)

A comics scraper with support for readcomicsonline.ru.

## Download
```
go get -u github.com/The-Eye-Team/Comics-DL
go build
```

## Usage
```
./Comics-DL --comic-id {ID}
```

### Flags
| Name | Default | Description |
|------|---------|-------------|
| `--comic-id` | Required. | The slug for the comic in the site URL. |
| `--concurrency` | `8` | The number of simultaneous downloads to run. |

## License
MIT.
