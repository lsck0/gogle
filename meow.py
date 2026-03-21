import os
import re


def replace_words_in_file(file_path):
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()

        # Replace all words (sequences of letters/numbers) with "meow"
        new_content = re.sub(r"\b\w+\b", "meow", content)

        with open(file_path, "w", encoding="utf-8") as f:
            f.write(new_content)

    except Exception as e:
        print(f"Skipping {file_path}: {e}")


def process_directory(directory):
    for root, _, files in os.walk(directory):
        for name in files:
            file_path = os.path.join(root, name)
            replace_words_in_file(file_path)


if __name__ == "__main__":
    target_dir = "."  # change to your directory path
    process_directory(target_dir)
