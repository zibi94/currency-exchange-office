# currency-exchange-office

# Description
```text
Simple currency exchange office app with two endpoints:
- /rates
    * required query param `currencies` - comma separated list of currencies (min 2)
- /exchange
    * required query params:
        ** from - string
        ** to - string
        ** amount - number
```

# Installation

### Required AppId - access for service openexchangerates.org !

## Build docker image
```bash
export APP_ID=<YOUR_APP_ID>

sudo docker build -t currency-exchange-office .

sudo docker run -dp 8080:8080 -e APP_ID=${APP_ID} currency-exchange-office
```

## Import docker image
```bash
export APP_ID=<YOUR_APP_ID>

sudo docker import --change 'ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin' \
--change 'CMD ["./main"]' --change 'WORKDIR /app' currency_exchange_office.tar currency-exchange-office

sudo docker container run -dp 8080:8080 -e APP_ID=${APP_ID} currency-exchange-office
```

##### Docker image exported wit command
```bash
docker export --output="currency_exchange_office.tar" <CONTAINER_ID>
```

