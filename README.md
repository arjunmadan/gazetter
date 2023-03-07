# gazetter
Get articles from the Gazette of India

## Usage

Download the correct binary from the [releases tab](https://github.com/arjunmadan/gazetter/releases/latest).

The CLI takes in 4 arguments:

- url: Base URL to make requests to. Defaults to https://egazette.nic.in/WriteReadData
- start: Starting range of article ID. Defaults to 160000.
- end: Ending range of article ID. Defaults to 180000.
- year: Year when the articles were published. Defaults to 2016.
- dir: Directory where the PDFs will be saved.


Example Usage:

On Linux/Mac:

```bash
./gazetter -url <url> -year 2020 -start 200000 -end 220000 -dir /Users/arjunmadan/Downloads
```


On Windows:

```shell
gazetter.exe -url <url> -year 2020 -start 200000 -end 220000 -dir C:\\Users\arjunmadan\Downloads
```
