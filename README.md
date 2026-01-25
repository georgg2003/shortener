# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Оптимизация 
Теперь почти вся память потребляется gzip сжатием и логгером 
```
File: shortener
Type: alloc_space
Time: 2026-01-25 18:41:31 MSK
Showing nodes accounting for -1910.86MB, 41.35% of 4620.71MB total
Dropped 345 nodes (cum <= 23.10MB)
      flat  flat%   sum%        cum   cum%
-1076.28MB 23.29% 23.29% -1163.78MB 25.19%  github.com/georgg2003/shortener/internal/repository/db/postgres.(*repository).GetUserURLs
 -656.21MB 14.20% 37.49%  -656.21MB 14.20%  bytes.growSlice
  333.18MB  7.21% 30.28%   397.54MB  8.60%  compress/flate.NewWriter (inline)
 -222.67MB  4.82% 35.10% -1607.46MB 34.79%  github.com/georgg2003/shortener/internal/usecase.(*useCase).GetUserURLs
 -113.50MB  2.46% 37.56%  -221.51MB  4.79%  github.com/georgg2003/shortener/internal/usecase.(*useCase).shortURLFromID (inline)
 -108.77MB  2.35% 39.91% -3060.44MB 66.23%  github.com/georgg2003/shortener/internal/delivery.(*delivery).GetUserURLs
 -104.50MB  2.26% 42.17%  -105.01MB  2.27%  fmt.Sprintf
     -88MB  1.90% 44.08%      -88MB  1.90%  github.com/jackc/pgx/v5/pgtype.scanPlanString.Scan
   60.36MB  1.31% 42.77%    60.36MB  1.31%  compress/flate.(*compressor).initDeflate (inline)
   30.92MB  0.67% 42.10%    30.92MB  0.67%  regexp.(*bitState).reset
   11.59MB  0.25% 41.85%    11.59MB  0.25%  github.com/jackc/pgx/v5/internal/iobufpool.init.0.func1
    8.52MB  0.18% 41.67%    49.12MB  1.06%  github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func3
    4.50MB 0.097% 41.57%     9.01MB  0.19%  net/http.(*conn).readRequest
    3.50MB 0.076% 41.49%    38.60MB  0.84%  github.com/jackc/pgx/v5.connect
    2.50MB 0.054% 41.44%    14.59MB  0.32%  github.com/jackc/pgx/v5/pgproto3.NewFrontend
    2.50MB 0.054% 41.39%    31.10MB  0.67%  github.com/jackc/pgx/v5/pgconn.connectOne
    1.50MB 0.032% 41.35%    10.50MB  0.23%  net.(*sysDialer).dialSingle
      -1MB 0.022% 41.38%    12.50MB  0.27%  github.com/golang-jwt/jwt/v4.(*Parser).ParseWithClaims
       1MB 0.022% 41.35%    59.43MB  1.29%  github.com/georgg2003/shortener/internal/repository/audit/service.(*auditServiceRepository).WriteLog
      -1MB 0.022% 41.38% -1930.87MB 41.79%  github.com/georgg2003/shortener/internal/delivery.(*delivery).GetNewRouter.NewAccessLogMiddleware.func1.1
    0.50MB 0.011% 41.36%       13MB  0.28%  github.com/georgg2003/shortener/pkg/jwthelper.(*JWTHelper).ReadAccessToken
    0.50MB 0.011% 41.35%   498.55MB 10.79%  github.com/georgg2003/shortener/internal/delivery.(*delivery).APIShortenURLBatch
         0     0% 41.35%  -452.28MB  9.79%  bytes.(*Buffer).Write
         0     0% 41.35%  -202.77MB  4.39%  bytes.(*Buffer).WriteString
         0     0% 41.35%  -656.21MB 14.20%  bytes.(*Buffer).grow
         0     0% 41.35%    13.53MB  0.29%  compress/flate.(*Writer).Close (inline)
         0     0% 41.35%    13.53MB  0.29%  compress/flate.(*compressor).close
         0     0% 41.35%    64.36MB  1.39%  compress/flate.(*compressor).init
         0     0% 41.35%    13.53MB  0.29%  compress/gzip.(*Writer).Close
         0     0% 41.35%   389.02MB  8.42%  compress/gzip.(*Writer).Write
         0     0% 41.35%  -270.19MB  5.85%  encoding/json.(*Encoder).Encode
         0     0% 41.35%  -658.71MB 14.26%  encoding/json.(*encodeState).marshal
         0     0% 41.35%  -658.71MB 14.26%  encoding/json.(*encodeState).reflectValue
         0     0% 41.35%  -656.71MB 14.21%  encoding/json.arrayEncoder.encode
         0     0% 41.35%  -656.71MB 14.21%  encoding/json.sliceEncoder.encode
         0     0% 41.35%  -454.28MB  9.83%  encoding/json.stringEncoder
         0     0% 41.35%  -656.71MB 14.21%  encoding/json.structEncoder.encode
         0     0% 41.35%   592.99MB 12.83%  github.com/georgg2003/shortener/internal/delivery.(*delivery).APIShortenURL
         0     0% 41.35% -1933.37MB 41.84%  github.com/georgg2003/shortener/internal/delivery.(*delivery).GetNewRouter.NewGzipCompressionMiddleware.func2.1
         0     0% 41.35% -1946.40MB 42.12%  github.com/georgg2003/shortener/internal/delivery.(*delivery).GetNewRouter.NewSimpleAuthMiddleware.func3.1
         0     0% 41.35%    13.53MB  0.29%  github.com/georgg2003/shortener/pkg/middlewares.gzipWriter.Close
         0     0% 41.35%   389.02MB  8.42%  github.com/georgg2003/shortener/pkg/middlewares.gzipWriter.Write
         0     0% 41.35% -1928.37MB 41.73%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 41.35% -1961.40MB 42.45%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 41.35%    51.42MB  1.11%  github.com/go-resty/resty/v2.(*Client).execute
         0     0% 41.35%    42.42MB  0.92%  github.com/go-resty/resty/v2.(*Client).executeBefore
         0     0% 41.35%    59.43MB  1.29%  github.com/go-resty/resty/v2.(*Request).Execute
         0     0% 41.35%    55.93MB  1.21%  github.com/go-resty/resty/v2.(*Request).Execute.func2
         0     0% 41.35%    59.43MB  1.29%  github.com/go-resty/resty/v2.(*Request).Post (inline)
         0     0% 41.35%    56.93MB  1.23%  github.com/go-resty/resty/v2.Backoff
         0     0% 41.35%    32.42MB   0.7%  github.com/go-resty/resty/v2.IsJSONType (inline)
         0     0% 41.35%    23.68MB  0.51%  github.com/go-resty/resty/v2.handleRequestBody
         0     0% 41.35%    23.68MB  0.51%  github.com/go-resty/resty/v2.parseRequestBody
         0     0% 41.35%     9.24MB   0.2%  github.com/go-resty/resty/v2.parseRequestHeader
         0     0% 41.35%      -88MB  1.90%  github.com/jackc/pgx/v5.(*baseRows).Scan
         0     0% 41.35%    39.10MB  0.85%  github.com/jackc/pgx/v5.ConnectConfig
         0     0% 41.35%    12.09MB  0.26%  github.com/jackc/pgx/v5/internal/iobufpool.Get
         0     0% 41.35%    33.10MB  0.72%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
         0     0% 41.35%    14.59MB  0.32%  github.com/jackc/pgx/v5/pgconn.ParseConfigWithOptions.func1
         0     0% 41.35%    31.60MB  0.68%  github.com/jackc/pgx/v5/pgconn.connectPreferred
         0     0% 41.35%    12.09MB  0.26%  github.com/jackc/pgx/v5/pgproto3.newChunkReader (inline)
         0     0% 41.35%    49.12MB  1.06%  github.com/jackc/puddle/v2.(*Pool[go.shape.*uint8]).initResourceValue.func1
         0     0% 41.35%     9.50MB  0.21%  github.com/sirupsen/logrus.(*Entry).Log
         0     0% 41.35%       10MB  0.22%  net.(*sysDialer).dialParallel.func1
         0     0% 41.35%    11.50MB  0.25%  net.(*sysDialer).dialSerial
         0     0% 41.35%        9MB  0.19%  net.(*sysDialer).dialTCP
         0     0% 41.35%        9MB  0.19%  net.(*sysDialer).doDialTCP (inline)
         0     0% 41.35%        9MB  0.19%  net.(*sysDialer).doDialTCPProto
         0     0% 41.35%        9MB  0.19%  net.internetSocket
         0     0% 41.35%        9MB  0.19%  net.socket
         0     0% 41.35%        9MB  0.19%  net/http.(*Client).Do (inline)
         0     0% 41.35%        9MB  0.19%  net/http.(*Client).do
         0     0% 41.35% -1917.36MB 41.49%  net/http.(*conn).serve
         0     0% 41.35% -1930.87MB 41.79%  net/http.HandlerFunc.ServeHTTP
         0     0% 41.35% -1928.37MB 41.73%  net/http.serverHandler.ServeHTTP
         0     0% 41.35%    32.42MB   0.7%  regexp.(*Regexp).MatchString (inline)
         0     0% 41.35%    32.42MB   0.7%  regexp.(*Regexp).backtrack
         0     0% 41.35%    32.42MB   0.7%  regexp.(*Regexp).doExecute
         0     0% 41.35%    32.42MB   0.7%  regexp.(*Regexp).doMatch (inline)
         0     0% 41.35%    13.59MB  0.29%  sync.(*Pool).Get
```