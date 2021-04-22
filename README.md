# Task-runner

Вспомогательная консольная утилита включает в себя следующие компоненты:

- Автоматизированные бекапы
- Frontend сборку на основе npm или yaml
- Release и Deploy скрипты с поддержкой промежуточных веток

Установка
---------

- загрузить bin файл https://gitlab.com/airollc/task-runner/-/tree/master/bin
- поместить бинарный файл в корень проекта
- создать файл конфигурации config.yaml

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
          "zcat ${-f} | mysql -u root -papolon13 ${-db}"
      ]
    remove: false
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
          '${-mode}',
          '-- --env.root'
      ]
    recursive: [
        "modules"
    ]
git:
  release:
    intermediate:
        - name: "test"
          amend: true
          command:
              main: '/home/user/project/task-runner'
              args: [
                "build-frontend"
              ]
````

Конфигурация отдельных компонентов
---------

### Backup

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
          "zcat ${-f} | mysql -u root -papolon13 ${-db}"
      ]
    remove: false
````
- restore
    - db
        - path
            - local Путь на хост машине, куда будут загружены файлы
            - remote Путь на удаленном сервере с backup файлами
        - command Команда для запуска бекапа
            - main Основная команда доступная в `$PATH` или путь к исп. файлу
            - args Аргументы
        - remove Удаление файла после выполнения команды
    
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
          '${-mode}',
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
    
_*Команда при передаче на исполнение будет выглядеть yarn production /home/user/project/modules/lsystem,
если нам необходимо передать только часть пути, мы можем указать cut-exec-path /home/user/project/modules/,
таким образом команда будет выглядеть yarn production lsystem._

### Git

````
git:
  release:
    intermediate:
        - name: "test"
          amend: true
          command:
              main: '/home/user/project/task-runner'
              args: [
                "build-frontend"
              ]
````
- git
    - release
        - intermediate Если требуется держать несколько веток в состоянии master
            - name Имя ветки
            - amend Требуется ли делать commit amend после исполнения команды
            - command Команда, которая будет выполнена на ветке name
                - main Основная команда доступная в `$PATH` или путь к исп. файлу
                - args Аргументы


Команды и аргументы task-runner -h
---------
````
Usage: task-runner restore-db
  -cnf string
        config file path (default "/home/user/project/config.yaml")
  -db string
        database
  -f string
        backup file name
Usage: task-runner build-frontend
  -cnf string
        config file path (default "/home/user/project/config.yaml")
  -mode string
        production or development (default "production")
Usage: task-runner release
  -branch string
        release branch (default "current")
  -cnf string
        config file path (default "/home/user/project/config.yaml")
Usage: task-runner deploy
  -branch string
        deploy branch (default "current")
  -stand string
        test stand

````