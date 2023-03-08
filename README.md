# gazetter
Get articles from the Gazette of India

## Usage

Download the correct binary from the [releases tab](https://github.com/arjunmadan/gazetter/releases/latest).

The CLI takes in 5 arguments:

- url: Base URL to make requests to. Defaults to https://egazette.nic.in/WriteReadData
- start: Starting range of article ID. Defaults to 160000.
- end: Ending range of article ID. Defaults to 180000.
- year: Year when the articles were published. Defaults to 2016.
- dir: Directory where the PDFs will be saved.
- gaz: Which gazette to download. Defaults to an empty string "" which will use the Gazette of India. The other option is "TN" which uses the Gazette of Tamil Nadu. When specifying `-gaz TN` `start`, `end`, and `url` are not required. 


### Example Usage for Gazette of India:

#### On Linux/Mac:

```bash
./gazetter -url <url> -year 2020 -start 200000 -end 220000 -dir /Users/arjunmadan/Downloads
```


#### On Windows:

```shell
gazetter.exe -url <url> -year 2020 -start 200000 -end 220000 -dir C:\\Users\arjunmadan\Downloads
```


### Example Usage for Gazette of Tamil Nadu:

#### On Linux/Mac:

```bash
./gazetter -gaz TN -year 2020 -dir /Users/arjunmadan/Downloads
```


#### On Windows:

```shell
gazetter.exe -gaz TN -year 2020 -dir C:\\Users\arjunmadan\Downloads
```
