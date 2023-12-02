# gopher-mart
It's yandex Go-programmer graduation project

## Setup
You should setup variables in `.env` file.
Example of this file you can find in `.env.sample`.

# TODO
- [x] POST /api/user/register — регистрация пользователя;
- [x] POST /api/user/login — аутентификация пользователя;
- [x] POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
- [x] GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
- [x] GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
- [x] POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
- [ ] GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.