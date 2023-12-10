# Gopher mart
It's yandex Go-programmer graduation project

## Setup
You should setup variables in `.env` file.
Example of this file you can find in `.env.sample`.

## TODO
- [x] POST /api/user/register — регистрация пользователя;
- [x] POST /api/user/login — аутентификация пользователя;
- [x] POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
- [x] GET  /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
- [x] GET  /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
- [x] POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
- [x] GET  /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
- [ ] взаимодействие с системой начисления баллов

## TODO features
### Tests
- [ ] POST /api/user/register
- [ ] POST /api/user/login
- [ ] POST /api/user/orders
- [ ] GET  /api/user/orders
- [ ] GET  /api/user/balance
- [ ] POST /api/user/balance/withdraw
- [ ] GET  /api/user/withdrawals
- [ ] gomock for DB