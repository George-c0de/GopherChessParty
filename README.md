# GopherChessParty

**GopherChessParty** – это весёлое веб-приложение для игры в шахматы, реализованное на Go. Проект сочетает
собственноручно разработанную серверную часть с интеграцией продвинутой шахматной
библиотеки [corentings/chess](https://github.com/CorentinGS/chess). Особое внимание уделено обновлению состояния игры
в реальном времени, эффективной обработке запросов с использованием горутин и каналов, а также масштабированию с
использованием возможностей многопроцессорной обработки Go.

---

## Содержание

- [Обзор проекта](#Обзор-проекта)
- [Цели и задачи](#Цели-и-задачи)
- [Архитектура](#Архитектура)
- [Технологический стек](#Технологический-стек)
- [Параллелизм и каналы](#Параллелизм-и-каналы)
- [Структура репозитория](#Структура-репозитория)
- [Планы на будущее]()
- [Ресурсы и ссылки]()

---

## Обзор проекта

**GopherChessParty** — это платформа для онлайн-игры в шахматы, где:

- Серверная часть написана на Go и использует весёлые переименования компонентов для придания индивидуальности.
- Для шахматной логики интегрирована библиотека [corentings/chess](https://github.com/CorentinGS/chess), которая
  обеспечивает генерацию ходов, работу с PGN/FEN, проверку состояний партии и многое другое.
- Реализована система обновления состояния игры в реальном времени с использованием WebSocket.
- Фронтенд создаётся с помощью GopherJS, а доска снабжается забавными анимациями («танцующая доска»).

---

## Цели и задачи

- **Цель:** Создать весёлое, масштабируемое и высокопроизводительное веб-приложение для игры в шахматы, демонстрирующее
  возможности Go.
- **Задачи:**
    - Реализовать REST API с использованием **Giggle Gin** для управления игровыми сессиями.
    - Организовать двустороннюю связь в реальном времени через **Gorilla Groove** (WebSocket).
    - Интегрировать библиотеку [corentings/chess](https://github.com/CorentinGS/chess) для продвинутой шахматной
      логики.
    - Создать интерактивный фронтенд с помощью GopherJS и юмористических CSS-анимаций.
    - Обеспечить эффективную обработку запросов с использованием горутин и каналов, а также масштабирование на
      многоядерных системах.
    - Упаковать и развернуть проект с использованием **Docker Disco** и **K8s Kooky**.

---

## Архитектура

### 1. Серверная часть – **Giggle Gin**

- **Описание:** HTTP-сервер, построенный на Go, использующий Giggle Gin для обработки REST API запросов.
- **Функции:**
    - Аутентификация пользователей.
    - Управление запросами на совершение ходов и получение состояния шахматной доски.
    - Интеграция с библиотекой corentings/chess для обновления и проверки состояния партии.

### 2. Реальное время и WebSocket – **Gorilla Groove**

- **Описание:** Модуль, обеспечивающий двустороннюю связь с клиентами через WebSocket.
- **Функции:**
    - Обработка WebSocket-соединений.
    - Рассылка уведомлений о совершённых ходах и изменениях в игровом состоянии в реальном времени.

### 3. Фронтенд – **GopherJS & Танцующая доска**

- **Описание:** Одностраничное приложение, скомпилированное с помощью GopherJS, с забавными анимациями, чтобы доска
  «танцевала» при каждом ходе.
- **Функции:**
    - Отображение шахматной доски, фигур и анимаций.
    - Интерактивное управление игрой и отправка ходов на сервер.

### 4. Контейнеризация и деплой – **Docker Disco & K8s Kooky**

- **Docker Disco:** Контейнеризация приложения для упрощённого развертывания.
- **K8s Kooky:** Оркестрация контейнеров для масштабирования и управления производительностью.

### 5. Логирование и мониторинг – **Zap Zap! Logger**

- **Описание:** Использование библиотеки zap для сбора и анализа логов работы приложения.
- **Функции:**
    - Запись логов в реальном времени.
    - Мониторинг работы сервера и выявление ошибок.

---

## Технологический стек

| Компонент            | Технология/Библиотека                                   | Описание                                                                                             |
|----------------------|---------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| **Язык разработки**  | Go                                                      | Основной язык реализации серверной логики и бизнес-процессов.                                        |
| **REST API**         | Giggle Gin (форк Gin)                                   | Фреймворк для маршрутизации и обработки HTTP-запросов.                                               |
| **Шахматная логика** | [corentings/chess](https://github.com/CorentinGS/chess) | Библиотека для генерации ходов, работы с PGN/FEN, проверки состояний партии и UCI-интерфейса.        |
| **WebSocket**        | Gorilla WebSocket (переименованный в Gorilla Groove)    | Обеспечивает двустороннюю связь для обновления состояния игры в реальном времени.                    |
| **Фронтенд**         | GopherJS, HTML/CSS/JS                                   | Компиляция Go в JavaScript для создания интерактивного одностраничного приложения с анимацией доски. |
| **Контейнеризация**  | Docker                                                  | Упаковка приложения в контейнеры для упрощённого развертывания.                                      |
| **Оркестрация**      | Kubernetes                                              | Масштабирование и управление контейнерами в продакшене.                                              |
| **Логирование**      | zap                                                     | Эффективное логирование для мониторинга работы приложения.                                           |

---

## Параллелизм и каналы

Проект **GopherChessParty** активно использует возможности Go по работе с горутинами и каналами для многопроцессорной
обработки:

- **Обработка ходов:** При получении запроса на совершение хода сервер запускает горутину для валидации и обновления
  состояния партии с использованием библиотеки corentings/chess.
- **Каналы уведомлений:** Использование каналов позволяет передавать сообщения между горутинами, уведомляя всех
  подписанных клиентов (через Gorilla Groove) о новых ходах и изменениях в игре.
- **Масштабирование:** Благодаря встроенному планировщику Go, обработка запросов распределяется по доступным ядрам, что
  обеспечивает высокую производительность и масштабируемость.
- **Асинхронное логирование:** Логирование с использованием Zap Zap! Logger также может работать в отдельной горутине,
  чтобы не блокировать основной процесс обработки ходов.

Эта архитектура гарантирует, что система будет способна обрабатывать множество одновременных игровых сессий, эффективно
используя ресурсы многоядерных систем.

---

## Структура репозитория

```
gopherchessparty/
├── cmd/
│   └── server/
│       └── main.go          // Точка входа – сервер с Giggle Gin и интеграцией corentings/chess
├── chesslogic/             // Обёртка для работы с библиотекой corentings/chess
│   └── game.go
├── internal/
└── README.md
```