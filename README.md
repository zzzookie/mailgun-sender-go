# mailgun-sender-go
# mailgun-sender-node
## How to install
Copy ```.env.example``` file and rename it to ```.env```. Add environment variables according to the example.

## How to use
```go run main.go -camp <mailgun_template_name> -ml <maillist_file_path>```

### Parameters
```-camp``` or ```-campaign``` **(required)** — Mailgun template name in the database.

```-ml``` or ```-maillist``` (optional) — Path to a ```.csv``` file in the following format:

    name,email,lang,ext_id
    "Client One",example@mail.com,en,d7DfCaF
    "Client Two",example2@gmail.com,es,e8Afaba
    
    (The first line is optional)

You can pass a file path or just a file name, either with or without the ```.csv``` extension.


### Using example:
``go run main.go -ml test.csv -camp template-name``
