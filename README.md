# Task-runner

Вспомогательная консольная утилита включает в себя следующие компоненты:

- Автоматизированные ресторы баз данных
- Frontend сборку на основе npm или yaml

Установка
---------

- загрузить bin файл https://gitlab.com/airollc/task-runner/-/tree/master/bin
- поместить бинарный файл в корень проекта
- создать файл конфигурации tr.config.yaml

Файл конфигурации
---------

Пример полного файла конфигурации:

````
connections:
  ssh:
    username: user
    host: ip
    port: 22
    password: ''
    private_key: 'id_rsa'
restore:
  db:
    path:
      local: /home/user/backup
      remote: /var/mysql
    command:
      main: "docker"
      args: [
          "exec",
          "-w",
          "/var/mysql/backup",
          "-i",
          "mysql",
          "/bin/bash",
          "-c",
          "zcat <-f> | mysql -u root -papolon13 <-db>"
      ]
    remove: false
grpc:
  client:
    root: '${CLIENT_PROTO_ROOT}'
    plugin: '/usr/bin/grpc_php_plugin'
    out: '${CLIENT_GENERATED}'
    common: '${CLIENT_PROTO_ROOT}/common'
    clear: [
        '${CLIENT_GENERATED}/Airo/Proto/CleaningPlan',
        '${CLIENT_GENERATED}/GPBMetadata/CleaningPlan',
        '${CLIENT_GENERATED}/Airo/Proto/Common'
    ]
  server:
    root: '${SERVER_PROTO_ROOT}'
    plugin: '/usr/bin/protoc-gen-php-grpc'
    out: '${SERVER_GENERATED}'
    common: '${SERVER_PROTO_ROOT}/common'
    clear: [
        '${SERVER_GENERATED}/Airo/Proto/CleaningPlan',
        '${SERVER_GENERATED}/GPBMetadata/CleaningPlan',
        '${SERVER_GENERATED}/Airo/Common'
    ]    
build:
  frontend:
    check-file: 'vue/app.js'
    root: '/home/user/project/modules'
    cut-exec-path: '/home/user/project/modules/'
    parallel: 5
    command:
      main: 'yarn'
      args: [
          '--cwd',
          '/home/user/project',
          '<-mode>',
          '-- --env.root'
      ]
    recursive: [
        "modules"
    ]
````

Конфигурация отдельных компонентов
---------

### Restore

````
restore:
  db:
    path:
      local: /home/user/backup
      remote: /var/mysql
    command:
      main: "docker"
      args: [
          "exec",
          "-w",
          "/var/mysql/backup",
          "-i",
          "mysql",
          "/bin/bash",
          "-c",
          "zcat <-f> | mysql -u root -papolon13 <-db>"
      ]
    remove: false
````

- restore
    - db
        - path
            - local Путь на хост машине, куда будут загружены файлы
            - remote Путь на удаленном сервере с dump файлами
        - command Команда для запуска рестора
            - main Основная команда доступная в `$PATH` или путь к исп. файлу
            - args Аргументы
        - remove Удаление файла после выполнения команды

### Grpc

````
grpc:
  client:
    root: '${CLIENT_PROTO_ROOT}'
    plugin: '/usr/bin/grpc_php_plugin'
    out: '${CLIENT_GENERATED}'
    common: '${CLIENT_PROTO_ROOT}/common'
    clear: [
        '${CLIENT_GENERATED}/Airo/Proto/CleaningPlan',
        '${CLIENT_GENERATED}/GPBMetadata/CleaningPlan',
        '${CLIENT_GENERATED}/Airo/Proto/Common'
    ]
  server:
    root: '${SERVER_PROTO_ROOT}'
    plugin: '/usr/bin/protoc-gen-php-grpc'
    out: '${SERVER_GENERATED}'
    common: '${SERVER_PROTO_ROOT}/common'
    clear: [
        '${SERVER_GENERATED}/Airo/Proto/CleaningPlan',
        '${SERVER_GENERATED}/GPBMetadata/CleaningPlan',
        '${SERVER_GENERATED}/Airo/Common'
    ]    
````

- grpc
    - client | server
        - root Корневой каталог proto интерфейсов
        - plugin Плагин для генерации
        - out Каталог куда будет записан результат 
        - common Каталог с общими proto интерфейсами
        - clear Список каталогов которые необходимо почистить перед запуском компилятора
    
### Build

````
build:
  frontend:
    check-file: 'vue/app.js'
    root: '/home/user/project/modules'
    cut-exec-path: '/home/user/project/modules/'
    parallel: 5
    command:
      main: 'yarn'
      args: [
          '--cwd',
          '/home/user/project',
          '<-mode>',
          '-- --env.root'
      ]
    recursive: [
        "modules"
    ]
````

- build
    - frontend
        - check-file Запускать сборку если в директории находится файл
        - root Корень обхода сборщиком
        - cut-exec-path Обрезать аргумент с названием директории (См. *)
        - parallel Количество параллельных сборок
        - command Команда сборки
            - main yarn vs npm
            - args Аргументы
        - recursive

_*Команда при передаче на исполнение будет выглядеть yarn production /home/user/project/modules/lsystem, если нам
необходимо передать только часть пути, мы можем указать cut-exec-path /home/user/project/modules/, таким образом команда
будет выглядеть yarn production lsystem._

Команды и аргументы task-runner -h
---------

````
Usage: task-runner restore-db
  -cnf string
        config file path (default "/home/user/project/config.yaml")
  -db string
        database
  -f string
        dump file name
Usage: task-runner build-frontend
  -cnf string
        config file path (default "/home/user/project/config.yaml")
  -mode string
        production or development (default "production")
Usage: task-runner grpc
  -cnf string
        config file path (default "/home/user/project/config.yaml)
  -pattern string
        <client or server>[:<service_name>]

````