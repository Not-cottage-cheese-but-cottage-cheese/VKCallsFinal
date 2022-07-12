# VK Видео и VK Звонки

Решение задач трека VK Видео и VK Звонки команды "Не творог, а творог"

- [10 и 20](#10-и-20)
- [30](#30)
- [40](#40)
- [50](#50)


# Общее руководство по запуску
Для корректной работы с VK API потребуется два токена: групповой токен и access токен. 
Гайды по их получению: 
- [групповой токен](https://dev.vk.com/api/bots/getting-started#%D0%9F%D0%BE%D0%BB%D1%83%D1%87%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%BA%D0%BB%D1%8E%D1%87%D0%B0%20%D0%B4%D0%BE%D1%81%D1%82%D1%83%D0%BF%D0%B0)
- [access токен](https://www.pandoge.com/socialnye-seti-i-messendzhery/poluchenie-klyucha-dostupa-access_token-dlya-api-vkontakte)

В корне проекта создать файл `.env` (задачи на 10 и 20) следующего вида:
```bash
# .env
GROUP_TOKEN=TOP_SECRET_TOKEN
SECRET=TOP_SECRET_TOKEN
```

Для запуска потребуется `go1.18`

### Запуск без сборки:
```bash
go run path/to/main.go arg1 arg2 ... argN
```

### Запуск со сброкой:
```bash
go build .
```
C последующим запуском бинарного файла

### Запуск в докере
```bash
docker-compose up --build -d
```

main.go файлы располагаются по пути 
```bash
/cmd/{TASK}/main.go
```
например 
```bash
/cmd/30/main.go
```

# 10 и 20
В рамках данных задач написан чат бот на Go, который знает следующие команды:

- [Звонок](#звонок)
- [Звонок оператору](#звонок-оператору)
- [Хочу быть оператором](#хочу-быть-оператором)
- [Я свободен](#я-свободен)

**_ВНИМАНИЕ:_**
Так как для работы бота требуются токены, которые имеют свойство истекать, а бессрочный получить нынче не удается, поэтому если вдруг что-то пойдет не так, то напишите нам, мы обновим ключи

## Звонок
Отправляет написавшему человеку уникальную ссылку на звонок

## Звонок оператору
*Соединяет* написавшего со свободным оператором или же ставит его в очередь, сообщая при появлении свободного оператора

## Хочу быть оператором
Добавляет написавшего человека в группу операторов. Изначально оператор считается **свободным**. При *соединении* с клиентом автоматически переходит в состояние **занят**

## Я свободен
Данное сообщение доступно для операторов. Переводит оператора в состояние **свободен**

# 30
Реализован скрипт на языке Go. Для запуска необходимы два аргумента - access токен и ID сообщества (c ведущим минусом):
```bash
go run ./cmd/30/main.go TOKEN -123456
```
Ссылки на все новые трансляции будут появлятся в консоли

# 40
Реализован скрипт на языке Go. Для запуска необходимы как минимум 4 аргумента - access токен, ID сообщества (c ведущим минусом), ID видео, набор регулярных выражений:
```bash
go run ./cmd/30/main.go TOKEN -123456 7891011 ".*" "\+" 
```
В результате в консоли будет выведен результат в следующем виде:
```bash
1) .* => 19
2) \+ => 1
```

# 50
Реализован скрипт на языке python. Для запуска необходимы 3 аргумента - access токен, ID сообщества (c ведущим минусом), название выходного видео:
```bash
python ./50/main.py TOKEN -123456 BEST_VIDEO_EVER 
```
Перед запуском необходимо поставить все нужные зависимости:
```bash
pip install -r ./50/requirements.txt
```

Пример сгенерированного видео можно увидеть [тут](./resources/video.avi)

# Описание задач

## 10
Напишите бота ВКонтакте, которому можно будет отправить сообщение «Звонок», а он ответит ссылкой-приглашением для присоединения к звонку. Каждый раз должен создаваться отдельный звонок.

## 20
Колл-центр. У вас есть 4 оператора. Напишите бота, который будет следить за тем, кто из операторов сейчас свободен (напишет об этом боту), и в соответствии с этим равномерно распределять посетителей бота между операторами в разных звонках. Оставляйте ссылку для каждого оператора постоянной, но не отправляйте посетителям ссылку до тех пор, пока очередной оператор не освободится.

## 30
Хотите не пропускать ни одного прямого эфира? Напишите скрипт, который будет раз в минуту проверять, появились ли в сообществе новые прямые трансляции. Выводите в консоль информацию о новых эфирах! Реализуйте это задание в качестве локального скрипта, который можно запустить на нашем компьютере.

## 40
Помогите администраторам сообществ считать голосования в комментариях. С помощью метода video.getComments посчитайте количество сообщений заданного вида в комментариях к видео.

Реализуйте это задание в качестве локального скрипта, который можно запустить на нашем компьютере. В параметрах для запуска передавайте идентификатор видео и варианты ответов. Чтобы получить полный балл, поддержите регулярные выражения в вариантах ответов.

## 50
Автоматически сгенерируйте видеотрейлер заданного сообщества. Экспортируйте обложки и названия публично доступных видеозаписей. Сформируйте из них видеотрейлер, в котором каждое видео канала будет представлено пятисекундным сегментом с обложкой видео и его названием. Упорядочивайте видео внутри трейлера в порядке убывания количества просмотров.
