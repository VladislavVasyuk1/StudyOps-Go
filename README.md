# StudyOps — Go / GitHub / Linux

Доработка учебного демо https://github.com/studybsuir/Dev-course/tree/main/go.
Исходные файлы сохранены в go/. Рабочее приложение и Dockerfile находятся
в корне этого репозитория.

## Фича

Русскоязычная адаптивная страница: таймер фокуса 5/15/25 минут,
пауза, сброс, счётчик завершённых сессий в localStorage.
Панель сервера показывает окружение, версию (SHA коммита) и uptime.
Веб-страница встроена в Go-бинарник, внешние зависимости не нужны.

## GitLab Flow на GitHub

feature/* → Pull Request → main → Pull Request → staging → Pull Request → production.
Используйте merge commits при продвижении между ветками окружений.

- feature/* и main: CI на GitHub-hosted Ubuntu.
- staging: тесты, Docker-сборка и деплой на своём Linux-раннере, localhost:8891.
- production: тесты, Docker-сборка и деплой на своём Linux-раннере, localhost:8892.

## Linux-раннер

Settings → Actions → Runners → New self-hosted runner → Linux.
Выберите архитектуру VM, выполните команды скачивания и проверки checksum
с этой страницы. При ./config.sh добавьте --labels studyops.
Не публикуйте токен регистрации. Запускайте runner обычным пользователем.
Раннеру нужны Docker, bash, curl и доступ к интернету. Go на хосте не нужен.
Метки workflow: self-hosted, linux, studyops.

В папке раннера после config:

```bash
sudo ./svc.sh install "$USER"
sudo ./svc.sh start
```

Раннер должен появиться в Settings со статусом Idle.
Workflow Deploy запускается только на push в staging и production,
не на pull_request. Запускайте на выделенной учебной VM.

## Локальный запуск

```bash
docker build --build-arg VERSION=local -t studyops:local .
docker run --rm -p 127.0.0.1:8890:8888 studyops:local
```

Откройте http://127.0.0.1:8890 на машине с Docker.
Для удалённой VM используйте SSH-туннель:

```bash
ssh -L 8892:127.0.0.1:8892 -L 8891:127.0.0.1:8891 USER@LINUX_HOST
```

## Проверка и защита

```bash
curl -f http://127.0.0.1:8892/healthz
curl -f http://127.0.0.1:8892/api/info
docker ps --filter name=studyops
docker logs studyops-production
```

Покажите собственную фичу, историю feature-ветки и PR, успешные CI/Deploy,
свой Linux-раннер Idle, контейнер healthy и версию на странице, совпадающую
с SHA Deploy run. Новый push через цепочку PR должен обновлять приложение.

Деплой проверяет health нового контейнера и при ошибке пытается вернуть
предыдущий контейнер. При первом запуске предыдущего контейнера нет.
Замена вызывает короткий перерыв. Staging и production собираются из
собственных SHA, которые могут отличаться из-за merge-коммитов.

Подготовленный workflow сам по себе не подтверждает успешный деплой:
результат проверяется в Actions и по ответам работающего сервера.
