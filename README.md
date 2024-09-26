# sik

> [!IMPORTANT]
>
> The project was rewritten to use Go and a **proper-ish* indexing method
> If you would still like to use the old python version, it will always be available at the [py](https://github.com/musaubrian/sik/tree/py) branch.

Sik allows you to search through Markdown files within a directory and quickly find the information you need.

## TODO:
- [x] **Flexible Querying(WIP):** You can search for any query within the Markdown files. The search is not case-sensitive and supports partial word matches.
- [ ] **Highlighting/previews**


## Usage
### Installation

Recommended:
```sh
go install github.com/musaubrian/sik/cmd@latest
```

Manual:
```sh
git clone https://github.com/musaubrian/sik

cd sik
go build -o sik ./cmd
```

or, get it from the [releases](https://github.com/musaubrian/sik/releases/latest)

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

## Contributing
This project is open for contributions!
Feel free to fork the repository, make improvements, and submit pull requests.

## License
See the [LICENSE](./LICENSE) file for more details.


> [!NOTE]
>
> Inspired by [seroost](https://github.com/tsoding/seroost)
