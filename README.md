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
backup:
  path:
    local: /home/user/mysql/backup
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
        "zcat ${-f} | mysql -u user -ppass ${-db}"
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

###Backup
````
backup:
  path:
    local: /home/user/mysql/backup
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
        "zcat ${-f} | mysql -u user -ppass ${-db}"
    ]
  remove: false
````
- backup
    - path
        - local Путь на хост машине, куда будут загружены файлы
        - remote Путь на удаленном сервере с backup файлами
    - command Команда для запуска бекапа
        - main Основная команда доступная в `$PATH`
        - args Аргументы
    - remove удаление файла после выполнения команды