#!/usr/bin/env python3

import os
import re
import json
import argparse

INDICES_DIR = f"{os.path.expanduser('~')}/.sik_indices"


def read_markdown_files(directory):
    files_content = {}
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith(".md"):
                filepath = os.path.join(root, file)
                with open(filepath, 'r', encoding='utf-8') as f:
                    content = f.read()
                files_content[filepath] = content
    return files_content


def create_index(content):
    index = {}
    lines = content.split('\n')
    for line_no, line in enumerate(lines, start=1):
        for match in re.finditer(r'\w+', line.lower()):
            word = match.group()
            col_no = match.start()
            if word not in index:
                index[word] = []
            index[word].append((line_no, col_no))
    return index


def create_indices_for_directory(directory):
    indices = {}
    files_content = read_markdown_files(directory)
    for file_path, content in files_content.items():
        index = create_index(content)
        indices[file_path] = {'index': index, 'content': content}
    return indices


def store_indices(indices, file_path):
    with open(file_path, 'w', encoding='utf-8') as file:
        json.dump(indices, file)


def load_indices(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        indices = json.load(file)
    return indices


def search_index(query, index, file_path):
    query_words = [query.lower()]
    results = []

    for word in index:
        if any(query_word in word for query_word in query_words):
            for (line_no, col_no) in index[word]:
                results.append((line_no, col_no, file_path))
    return results


def search_indices(query, indices):
    results = {}
    for file_path, file_data in indices.items():
        index = file_data['index']
        file_results = search_index(query, index, file_path)
        if file_results:
            results[file_path] = file_results
    return results


def highlight_query_in_line(query, line):
    query_words = re.findall(r'\w+', query.lower())
    words = re.findall(r'\w+', line)
    highlighted_line = line
    for word in words:
        for query_word in query_words:
            if query_word in word.lower():
                highlighted_line = re.sub(f"({word})", highlight(
                    r'\1'), highlighted_line, flags=re.IGNORECASE)
    return highlighted_line


def search_kb(query, index_path):
    indices = load_indices(index_path)
    results = search_indices(query, indices)
    if not results:
        print(f"Found nothing matching {print_bold(query)}")

    grouped_results = {}
    for filename, locations in results.items():
        parent_dir = os.path.dirname(filename)
        if parent_dir not in grouped_results:
            grouped_results[parent_dir] = []

        for line_no, col_no, file_path in locations:
            content = indices[file_path]['content']
            line = content.split('\n')[line_no - 1]
            highlighted_line = highlight_query_in_line(query, line)
            location = f"{os.path.basename(file_path)}:{line_no}"
            grouped_results[parent_dir].append((location, highlighted_line))

    for parent_dir, file_results in grouped_results.items():
        print(print_bold(darken(remove_prefix(parent_dir))))
        for location, highlighted_line in file_results:
            print(
                f"   {darken(print_bold(location), level=2)}  "
                f"  {darken(highlighted_line.strip(), level=4)}")
        print()


def remove_prefix(line: str) -> str:
    if line.startswith("/"):
        return line.removeprefix("/")
    elif line.startswith("\\"):
        return line.removeprefix("\\")
    else:
        raise Exception(f"Unexpected prefix in {line}")


def print_bold(text):
    return f"\033[1m{text}\033[0m"


def highlight(text):
    return f"\033[38;5;208m{text}\033[0m"


def darken(text, level=1):
    # '#f4f4f4'
    base_r, base_g, base_b = 244, 244, 244
    darken_amount = 15
    # Calculate the new RGB values
    r = max(0, base_r - level * darken_amount)
    g = max(0, base_g - level * darken_amount)
    b = max(0, base_b - level * darken_amount)

    return f"\033[38;2;{r};{g};{b}m{text}\033[0m"


def get_index_file(directory):
    dir_name = os.path.basename(os.path.normpath(directory))
    return f"{dir_name}_index.sik"


def list_available_indices():
    return [file for file in
            os.listdir(INDICES_DIR) if file.endswith('_index.sik')]


def create_indices_dir():
    if not os.path.exists(INDICES_DIR):
        os.makedirs(INDICES_DIR)
        print(print_bold(darken(f"Created {INDICES_DIR}", level=5)))


def select_index_file():
    available_indices = list_available_indices()
    if not available_indices:
        print(darken("No index files found. Please create an index first."))
        return None

    if len(available_indices) == 1:
        return available_indices[0]

    print(print_bold("Available indices:"))
    print(darken("=" * 30))
    for idx, file in enumerate(available_indices, 1):
        print(f"{darken(f'{idx:2d}:')} {print_bold(file)}")
    print(darken("=" * 30))

    choice = int(input(darken("Select the index file by number: "))) - 1
    if 0 <= choice < len(available_indices):
        return available_indices[choice]
    else:
        print(darken("Invalid selection."))
        return None


def main():
    create_indices_dir()
    parser = argparse.ArgumentParser(
        description="Query your markdown files")
    parser.add_argument(
        "-q", "--query", help="Word(s) to search for")
    parser.add_argument(
        "-d", "--dir", help="The directory to start index")
    parser.add_argument("--index", action='store_true',
                        help="Create indices for the markdown files.")

    args = parser.parse_args()
    directory = args.dir

    if args.index:
        if directory is None:
            print("Please specify a directory to index.")
            return
        json_index_path = f"{INDICES_DIR}/{get_index_file(directory)}"
        print("Indexing files...")
        indices = create_indices_for_directory(directory)
        store_indices(indices, json_index_path)
        print(f"Created index: {print_bold(json_index_path)}")
    elif args.query:
        json_index_path = f"{INDICES_DIR}/{select_index_file()}"
        if json_index_path is None:
            return
        query = args.query
        search_kb(query, json_index_path)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()
