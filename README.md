## API

Exhaustive list of needed APIs.

| method | route                     | status |
|--------|---------------------------|--------|
| GET    | /transactions             | DONE   |
| POST   | /transactions             | DONE   |
| PUT    | /transactions             | DONE   |
| DELETE | /transactions             | DONE   |
| GET    | /accounts                 | DONE   |
| POST   | /accounts                 | DONE   |
| PUT    | /accounts                 | DONE   |
| DELETE | /accounts                 | DONE   |
| GET    | /banks                    | DONE   |
| POST   | /banks                    | DONE   |
| PUT    | /banks                    | DONE   |
| DELETE | /banks                    | DONE   |
| GET    | /categories               | DONE   |
| POST   | /categories               | DONE   |
| PUT    | /categories               | DONE   |
| DELETE | /categories               | DONE   |
| GET    | /summary                  | TODO   |
| GET    | /report                   | TODO   |
| GET    | /report/debit             | TODO   |
| GET    | /report/credit            | TODO   |
| GET    | /report/categories        | TODO   |
| GET    | /report/debit/categories  | TODO   |
| GET    | /report/credit/categories | TODO   |
| GET    | /report/accounts          | TODO   |
| GET    | /report/debit/accounts    | TODO   |
| GET    | /report/credit/accounts   | TODO   |
| GET    | /report/banks             | TODO   |
| GET    | /report/debit/banks       | TODO   |
| GET    | /report/credit/banks      | TODO   |
| GET    | /users                    | TODO?  |
| POST   | /users                    | TODO?  |
| PUT    | /users                    | TODO?  |
| DELETE | /users                    | TODO?  |

``/report`` with the right filters ->
```json
{
  "debit": 9242.30,
  "credit": 4300
}
```

``/report/categories`` ->
```json
[
  {
    "categoryId": 11,
    "debit": 242.30,
    "credit": 400.00
  },
  {
    "categoryId": 21,
    "debit": 242.30,
    "credit": 0.00
  }
]
```

``/report/debit/categories`` ->
```json
[
  {
    "categoryId": 11,
    "debit": 242.30
  },
  {
    "categoryId": 21,
    "debit": 242.30
  }
]
```

``/report/credit/categories`` ->
```json
[
  {
    "categoryId": 11,
    "credit": 400.00
  },
  {
    "categoryId": 21,
    "credit": 0.00
  }
]
```
``/summary``
```json
{
  "total": 21324,
  "lastTransactionDate": 1224325543,
  "banks": [
    {
      "id": 2,
      "name": "My Bank x",
      "total": 21314,
      "accounts": [
        {
          "id": 2,
          "name": "account a",
          "total": 21300
        },
        {
          "id": 7,
          "name": "account b",
          "total": 14
        }
      ]
    },
    {
      "id": 3,
      "name": "My Bank y",
      "total": 10,
      "accounts": [
        {
          "id": 9,
          "name": "account g",
          "total": 10
        }
      ]
    }
  ],

  "categories": [
    {
      "id": 2,
      "name": "restaurant"
    },
    {
      "id": 9,
      "name": "technology"
    }
  ]
}
```

## TODO
  - simplify code for now (database definitions updated, good when it compile!)
  - most sql_ methods take the userId (banks, accounts, categories and transactions) (denormalize bdd to include id? or cache?)
  - add rest of the routes. First summary then report/categories then test then rest while frontin'
  - possibility to activate/deactive https from config (will surely have a reverse proxy anyway)
  - begin to program the front and webserver (!!)
  - add way to calculate the real total for each account (diff with transaction sum? / date + state at this date + transaction sum from date?)
  - add personal parsers for csv etc.
  - create web scrapper and cronjob to update?
  - add autofilters to automatically set categories (front or back?)
  - nested categories
  - crypt transactions/banks/accounts/categories?
  - switch all code to elixir or something
  - simplify that monstruosity
  - README To explain every API
  - implement cache somewhere (over the rainbow?)
  - become richer than uncle scrooge
  - sleep
