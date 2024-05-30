#!/usr/bin/env python3

import os
import re
import json
import argparse


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

    grouped_results = {}
    for filename, locations in results.items():
        parent_dir = os.path.dirname(filename)
        if parent_dir not in grouped_results:
            grouped_results[parent_dir] = []

        for line_no, col_no, file_path in locations:
            content = indices[file_path]['content']
            line = content.split('\n')[line_no - 1]
            highlighted_line = highlight_query_in_line(query, line)
            location = f"{os.path.basename(file_path)}:{line_no},{col_no}"
            grouped_results[parent_dir].append((location, highlighted_line))

    for parent_dir, file_results in grouped_results.items():
        print(print_bold(parent_dir))
        for location, highlighted_line in file_results:
            print(f"   {print_bold(location)}")
            print(f"       {highlighted_line}")
        print()


def print_bold(text):
    return f"\033[1m{text}\033[0m"


def highlight(text):
    return f"\033[38;5;208m{text}\033[0m"


def main():
    parser = argparse.ArgumentParser(
        description="Search markdown files for a query.")
    parser.add_argument(
        "-q", "--query", help="The query to search for in the markdown files.")
    parser.add_argument(
        "-d", "--dir", help="The directory to start indexing and searching")
    parser.add_argument("--index", action='store_true',
                        help="Create indices for the markdown files.")

    args = parser.parse_args()
    directory = args.dir
    json_index_path = 'index.sik'

    if args.index:
        print("Indexing files...")
        indices = create_indices_for_directory(directory)
        store_indices(indices, json_index_path)
        print("Indexing complete.")
    elif args.query:
        query = args.query
        search_kb(query, json_index_path)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()
