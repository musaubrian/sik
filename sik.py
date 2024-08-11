#!/usr/bin/env python3

import os
import re
import json
import argparse

INDICES_DIR = f"{os.path.expanduser('~')}/.sik"
INDEX_LOCATION = f"{INDICES_DIR}/index.sik"
OUTPUT_STYLE = "ansi"


def read_markdown_files(directory):
    files_content = {}
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith(".md"):
                filepath = os.path.join(root, file)
                with open(filepath, "r", encoding="utf-8") as f:
                    content = f.read()
                files_content[filepath] = content
    return files_content


def create_index(content):
    index = {}
    lines = content.split("\n")
    for line_no, line in enumerate(lines, start=1):
        for match in re.finditer(r"\w+", line.lower()):
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
        indices[file_path] = {"index": index, "content": content}
    return indices


def store_indices(indices, file_path):
    with open(file_path, "w", encoding="utf-8") as file:
        json.dump(indices, file)


def load_indices(file_path):
    try:
        with open(file_path, "r", encoding="utf-8") as file:
            indices = json.load(file)
        return indices
    except FileNotFoundError:
        raise Exception(
            f"File not found: {file_path}, Try creating an index first")


def search_index(query, index, file_path):
    query_words = [query.lower()]
    results = []

    for word in index:
        if any(query_word in word for query_word in query_words):
            for line_no, col_no in index[word]:
                results.append((line_no, col_no, file_path))
    return results


def search_indices(query, indices):
    results = {}
    for file_path, file_data in indices.items():
        index = file_data["index"]
        file_results = search_index(query, index, file_path)
        if file_results:
            results[file_path] = file_results
    return results


def highlight_query_in_line(query, line):
    query_words = re.findall(r"\w+", query.lower())
    words = re.findall(r"\w+", line)
    highlighted_line = line
    for word in words:
        for query_word in query_words:
            if query_word in word.lower():
                highlighted_line = re.sub(
                    f"({word})", highlight(r"\1"), highlighted_line, flags=re.IGNORECASE
                )
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
            content = indices[file_path]["content"]
            line = content.split("\n")[line_no - 1]
            highlighted_line = highlight_query_in_line(query, line)
            location = f"{os.path.basename(file_path)}:{line_no}"
            grouped_results[parent_dir].append((location, highlighted_line))

    for parent_dir, file_results in grouped_results.items():
        print(print_bold(darken(remove_prefix(parent_dir))))
        for location, highlighted_line in file_results:
            print(
                f"   {darken(print_bold(location), level=2)}"
                f"  {darken(highlighted_line.strip(), level=4)}"
            )
        print()


def remove_prefix(line: str) -> str:
    if line.startswith("/"):
        return line.removeprefix("/")
    elif line.startswith("\\"):
        return line.removeprefix("\\")
    else:
        raise Exception(f"Unexpected prefix in {line}")


def print_bold(text):
    if OUTPUT_STYLE == "ansi":
        return f"\033[1m{text}\033[0m"
    else:
        return f"<b>{text}</b>"


def highlight(text):
    if OUTPUT_STYLE == "ansi":
        return f"\033[38;5;208m{text}\033[0m"
    else:
        return f"<span style='color: #ff8c00;'>{text}</span>"


def darken(text, level=1):
    # '#f4f4f4'
    base_r, base_g, base_b = 244, 244, 244
    darken_amount = 15
    # Calculate the new RGB values
    r = max(0, base_r - level * darken_amount)
    g = max(0, base_g - level * darken_amount)
    b = max(0, base_b - level * darken_amount)

    return f"\033[38;2;{r};{g};{b}m{text}\033[0m"


def create_index_dir():
    if not os.path.exists(INDICES_DIR):
        os.makedirs(INDICES_DIR)
        print(print_bold(darken(f"Created {INDICES_DIR}", level=5)))
        print(print_bold(darken(
            f"Run <sik.py --index --dir [path/to/dir]> to create an index", level=5)))


def json_query_result(query, index_path):
    global OUTPUT_STYLE
    OUTPUT_STYLE = "html"
    indices = load_indices(index_path)
    results = search_indices(query, indices)
    if not results:
        return []

    json_results = []
    for file_path, locations in results.items():
        occurrences = []
        for line_no, col_no, _ in locations[:3]:
            content = indices[file_path]["content"]
            line = content.split("\n")[line_no - 1]
            highlighted_line = highlight_query_in_line(query, line)
            occurrences.append(f"<b>{line_no}:</b> {highlighted_line.strip()}")

        json_results.append({
            "full_path": file_path,
            "filename": os.path.basename(file_path),
            "hits": occurrences
        })

    return json_results


def main():
    create_index_dir()

    parser = argparse.ArgumentParser(description="Query your markdown files")
    parser.add_argument("-q", "--query", help="Word(s) to search for")
    parser.add_argument("-d", "--dir", help="The directory to index")
    parser.add_argument(
        "--json",
        action="store_true",
        help="Get query results in json for easy parsing.")
    parser.add_argument(
        "--index",
        action="store_true",
        help="Create indices for the markdown files."
    )

    args = parser.parse_args()
    directory = args.dir

    if args.index:
        if directory is None:
            print("Please specify a directory to index.")
            return
        json_index_path = INDEX_LOCATION
        print("Indexing files...")
        indices = create_indices_for_directory(directory)
        store_indices(indices, json_index_path)
        print(f"Created index: {print_bold(json_index_path)}")
    elif args.query:
        json_index_path = INDEX_LOCATION
        if json_index_path is None:
            return
        query = args.query
        if args.json:
            json_results = json_query_result(query, json_index_path)
            print(json.dumps(json_results, indent=4))
        else:
            search_kb(query, json_index_path)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()
