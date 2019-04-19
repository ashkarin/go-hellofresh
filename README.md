# hellofresh-api
This mini project demonstrates how the service to work with recipes can be organized. Since we don't need to do projections over the data, the MongoDB was used as a data storage.

The implementation was driven by the integration tests, which can be turned on by setting `TEST` in `docker-compose.yaml` to `'true'`.

Additionally, a set of Python tools were provided. These tools allow downloading the data from the website in JSON format, transform it according to the recipe schema and push to the database.

## Example
```
docker-compose up --build
python scripts/steal_recepies.py 0 500 --json-pretty --output data.json
python scripts/store_recepies.py --filepath data.json  --address http://localhost:8080
```
