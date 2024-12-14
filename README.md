# Кошелек

## Сборка и запуск
```sh
git clone https://github.com/HennOgyrchik/java_code.git
cd java_code
docker compose up
```

## Endpoints
1. `PUT /api/v1/wallet/{uuid}` - создание нового кошелька

   Пример: `/api/v1/wallets/111e4567-e89b-12d3-a456-426655440000`

2. `POST /api/v1/wallet` - пополнение или списание средств
 
   Пополнение:
```json
{
    "walletId": "111e4567-e89b-12d3-a456-426655440000",
    "operationType": "deposit",
    "amount": 2000
}
```
   Списание:
```json
{
    "walletId": "111e4567-e89b-12d3-a456-426655440000",
    "operationType": "withdraw",
    "amount": 2000
}
```

3. `GET /api/v1/wallets/{uuid}` - получение баланса
