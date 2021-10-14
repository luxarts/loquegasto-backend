# loquegasto-backend
Backend del bot de Telegram LoQueGasto

# Examples
### POST Transaction
```curl
curl --location --request POST 'localhost:8080/transaction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "msg_id": 123,
    "amount": 1000,
    "description": "Some item",
    "source": "debit"
}'
```