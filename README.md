# sik

Sik allows you to search through Markdown files within a directory and quickly find the information you need.


## Features
- **Flexible Querying:** You can search for any query within the Markdown files. The search is not case-sensitive and supports partial word matches.
- **Highlighting:** Search results are intelligently highlighted, making it easy to spot the relevant information at a glance.

## Usage

### Indexing
Before searching, you need to index the Markdown files. Use the `--index` flag along with the directory path to initiate indexing.
You can index as many directories as you wish

```bash
python3 sik.py --index -d /path/to/directory

or ./sik.py --index -d /path/to/directory
```

### Searching
Once indexed, you can search for queries within the Markdown files. Provide your query using the `-q` or `--query` flag.
If you have more than one index, you'll be prompted to pick one
```bash
python3 sik.py -q <random_query>
```

## GUI

Sik now has a gui version.
It does still depend on the sik cli version to handle the searching and indexing functionality

### Usage
Currently you have to build it from source.

```sh
cd sik
npm tauri dev
# This will spin up the development server
#or build application using
npm tauri build
```


## Contributing
This project is open for contributions!
Feel free to fork the repository, make improvements, and submit pull requests.

## License
This project is licensed under the MIT License, allowing you to use, modify, and distribute the code for both commercial and non-commercial purposes. See the [LICENSE](./LICENSE) file for more details.


> [!NOTE]
>
> Inspired by [seroost](https://github.com/tsoding/seroost)
