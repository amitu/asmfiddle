import glob

DOIT_CONFIG = {
    "verbosity": 2,
    "default_tasks": ["pip", "gtests", "js", "css"]
}


def all_files():
    return (
        glob.glob("src/asmfiddle/*/*.go")
        + glob.glob("src/asmfiddle/*.go")
        + ["dodo.py"]
    )

ALL_FILES = all_files()


def task_gtests():
    return {
        "actions": [
            "go test -v -race -timeout 1s asmfiddle/...",
            "go vet asmfiddle/...",
            "./bin/errcheck -asserts asmfiddle/...",
        ],
        "file_dep": ALL_FILES,
    }


def task_pip():
    return {
        "actions": ["pip install -r requirements.txt"],
        "file_dep": ["requirements.txt"]
    }


def task_css():
    return {
        "actions": [
            "mkdir -p build",
            (
                "sassc -I scss/foundation-sites/scss/ --sourcemap "
                "scss/main.scss build/style.css"
            )
        ],
        "targets": ["build/style.css"],
        "file_dep": glob.glob("scss/*.scss"),
    }


def task_js():
    return {
        "actions": [
            "gopherjs build -m asmfiddle/cmd/script",
            "mkdir -p build",
            "mv script.js script.js.map build",
        ],
        "file_dep": ALL_FILES,
    }
