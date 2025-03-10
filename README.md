
# Internal transfers System


## Run Locally

1. Clone the project

```bash
  git clone git@github.com:Sant470/accounts-svc.git
```

2. Go to the project directory

```bash
  cd accounts-svc
```

3. Change DB Credential in config.yaml


4. Install dependencies
```bash
  go mod tidy
```

5. Start the Server
```bash
  go run main.go
```

## Test it locally
Use the following curl to test different apis

Create an account
```bash
  curl --location 'http://localhost:8000/api/v1/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "account_id": "105",
    "initial_balance": 10000.7654
}'
```

Fetch Balance
```bash
curl --location 'http://localhost:8000/api/v1/accounts/101'
```

Make a transaction 
```bash
  curl --location 'http://localhost:8000/api/v1/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "source_account_id": "101",
    "destination_account_id": "105",
    "amount": 100.12345
}'
```
## Dependency && Installation

It requires go version 1.22.0, you can download it following the guide mentioned below 

```bash
   https://go.dev/dl/
```
    