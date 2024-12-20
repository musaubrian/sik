# sik

> [!IMPORTANT]
>
> If you still wish to use the old python version, it's available at the [py](https://github.com/musaubrian/sik/tree/py) branch.

Sik allows you to search through Markdown files within a directory and quickly find the information you need.

## TODO:
- [x] **Flexible Querying:** You can search for any query within the Markdown files. The search is not case-sensitive and supports partial word matches.
- [x] **Highlighting/previews**
- [x] Loads up the file directly in the browser, making it easier to view what you were searching for
- [x] Bundle up [marked](https://marked.js.org/) and make it completely local
- [x] LiveSearch: Search on type


## Usage
### Installation

Recommended:
```sh
go install github.com/musaubrian/sik/cmd/sik@latest
```

Manual:
```sh
git clone https://github.com/musaubrian/sik

cd sik
go build ./cmd/sik
```

### Indexing
Before searching, you need to index the Markdown files. Use the `-index` flag along with the directory path to initiate indexing.

```bash
sik -index </path/to/directory/to/index>
```

### Searching
Once indexed, you can search your notes by providing the `-b` flag, this will start a webserver at port `8990`.
You can search through your data at this page
```bash
sik -b

# You can run it in the background with `sik -b &`
```
- You don't have to restart the webserver to re-index your information


> You can always check the help `sik -h`

### Dev mode
If you wish to make changes and already have an instance running,
prefix run command with `SIK_DEV=<anything_you_wish> air|go run cmd/sik/sik.go`

## Contributing
This project is open for contributions!
Feel free to fork the repository, make improvements, and submit pull requests.

## License
See the [LICENSE](./LICENSE) file for more details.


> Inspired by [seroost](https://github.com/tsoding/seroost)
