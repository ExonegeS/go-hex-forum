# API для получения вакансий 

Название источника
1. Название: HH.ru API

2. Тип: REST 

3. Ссылка на API в [гитхабе](https://github.com/hhru/api)

4. Как работает [API](https://api.hh.ru/openapi/redoc#section/Obshaya-informaciya)

Ответ возвращается только в формате JSON.  
Перед отправкой запроса необходимо представиться и показать апишке, кто делает запрос используя хедер User-Agent.  
Пример User-Agent: findVaccancies (some-app-using-hhapi@gmail.com)

## Географическая информация
- **Страна происхождения**: Россия
- **Целевые страны**: Россия, Беларусь, Казахстан, другие страны СНГ
- **Поддерживаемые домены**:
    - rabota.by
    - hh1.az
    - hh.uz
    - hh.kz
    - headhunter.ge
    - headhunter.kg

### Как воспользоваться им?
Нужно отправить запрос на https://api.hh.ru/  
Для получения полного доступа ко всему функционалу приложения необходим зарегестрировать свое приложение, существуют и платные методы, все они описаны в гитхабе.  
Для простого поиска вакансии можно не регистрировать свое приложение, но User-Agent предоставлять нужно, иначе прилетит Bad Request.


### Поиск вакансий
Поиск вакансий	Поддержка параметров: text, area, page, per_page, host

### Как получить данные о всех возможных вакансиях?
Отправить GET запрос на https://api.hh.ru/ c User-Agent headerом
 
Пример User-Agent: findVaccancies (some-app-using-hhapi@gmail.com)

Ответ от приложенения по запросу GET 
```json
{
    "items": [
        {
            "id": "119883635",
            "premium": false,
            "name": "Junior Go Developer",
            "department": null,
            "has_test": false,
            "response_letter_required": false,
            "area": {
                "id": "1",
                "name": "Москва",
                "url": "https://api.hh.ru/areas/1"
            },
            "salary": {
                "from": 85000,
                "to": null,
                "currency": "RUR",
                "gross": false
            },
            "salary_range": {
                "from": 85000,
                "to": null,
                "currency": "RUR",
                "gross": false,
                "mode": {
                    "id": "MONTH",
                    "name": "За месяц"
                },
                "frequency": {
                    "id": "TWICE_PER_MONTH",
                    "name": "Два раза в месяц"
                }
            },
            "type": {
                "id": "open",
                "name": "Открытая"
            },
            "address": null,
            "response_url": null,
            "sort_point_distance": null,
            "published_at": "2025-04-23T18:55:44+0300",
            "created_at": "2025-04-23T18:55:44+0300",
            "archived": false,
            "apply_alternate_url": "https://hh.ru/applicant/vacancy_response?vacancyId=119883635",
            "show_logo_in_search": null,
            "insider_interview": null,
            "url": "https://api.hh.ru/vacancies/119883635?host=hh.ru",
            "alternate_url": "https://hh.ru/vacancy/119883635",
            "relations": [],
            "employer": {
                "id": "11496690",
                "name": "Бетель филиал в г. Владимир",
                "url": "https://api.hh.ru/employers/11496690",
                "alternate_url": "https://hh.ru/employer/11496690",
                "logo_urls": {
                    "original": "https://img.hhcdn.ru/employer-logo-original/1341307.png",
                    "90": "https://img.hhcdn.ru/employer-logo/6985170.png",
                    "240": "https://img.hhcdn.ru/employer-logo/6985171.png"
                },
                "vacancies_url": "https://api.hh.ru/vacancies?employer_id=11496690",
                "accredited_it_employer": false,
                "employer_rating": {
                    "total_rating": "0",
                    "reviews_count": null
                },
                "trusted": true
            },
            "snippet": {
                "requirement": "Понимание ООП, базовых алгоритмов и структур данных. Опыт работы с SQL или NoSQL базами данных. Знание принципов RESTful API (опыт...",
                "responsibility": "Поддержка команды и наставничество. Готов к новым вызовам и развитию в backend? Присоединяйся к нашей команде!"
            },
            "show_contacts": true,
            "contacts": null,
            "schedule": {
                "id": "remote",
                "name": "Удаленная работа"
            },
            "working_days": [],
            "working_time_intervals": [],
            "working_time_modes": [],
            "accept_temporary": false,
            "fly_in_fly_out_duration": [],
            "work_format": [
                {
                    "id": "REMOTE",
                    "name": "Удалённо"
                }
            ],
            "working_hours": [
                {
                    "id": "HOURS_8",
                    "name": "8 часов"
                }
            ],
            "work_schedule_by_days": [
                {
                    "id": "FIVE_ON_TWO_OFF",
                    "name": "5/2"
                }
            ],
            "night_shifts": false,
            "professional_roles": [
                {
                    "id": "40",
                    "name": "Другое"
                }
            ],
            "accept_incomplete_resumes": true,
            "experience": {
                "id": "between1And3",
                "name": "От 1 года до 3 лет"
            },
            "employment": {
                "id": "full",
                "name": "Полная занятость"
            },
            "employment_form": {
                "id": "FULL",
                "name": "Полная"
            },
            "internship": false,
            "adv_response_url": null,
            "is_adv_vacancy": false,
            "adv_context": null
        }
    ],
    "found": 2907,
    "pages": 2000,
    "page": 0,
    "per_page": 1,
    "clusters": null,
    "arguments": null,
    "fixes": null,
    "suggests": null,
    "alternate_url": "https://hh.ru/search/vacancy?enable_snippets=true&items_on_page=1&text=go"
}
```

### Локализация вакансий по сайтам группы HH
Для получения локализаций, доступных на h необходимо сделать GET запрос на один из поддерживаемых сайтов, например https://api.hh.ru/locales?host=hh.kz, чтобы получить вакансии имеющиеся на hh.kz  
Запрос без ?host приведет к тому что вы получите вакансии со всех сайтов группы компании HeadHunter. 

### Получить определенную вакансию
Определенную вакансию можно найти написав айди вакансии в формате GET api.hh.ru/vacancies/{vacancy_id}  
Фильтрацию можно выполнить с помощью параметра text для того чтобы получить интересующую вакансию  
Полный список параметров при поиске вакансий можно найти [здесь](https://api.hh.ru/openapi/redoc#tag/Poisk-vakansij-dlya-soiskatelya) 



### Пагинация результатов
Пагинация
К любому запросу, подразумевающему выдачу списка объектов, можно в параметрах указать page=N&per_page=M. Нумерация идёт с нуля, по умолчанию выдаётся первая (нулевая) страница с 20 объектами на странице. Во всех ответах, где доступна пагинация, единообразный корневой объект:
```
{
  "found": 1,
  "per_page": 1,
  "pages": 1,
  "page": 0,
  "items": [{}]
}
```
При указании параметров пагинации (page, per_page) работает ограничение: глубина возвращаемых результатов не может быть больше 2000. Например, возможен запрос per_page=10&page=199 (выдача с 1991 по 2000 вакансию), но запрос с per_page=10&page=200 вернёт ошибку (выдача с 2001 по 2010 вакансию)



### Ошибки и коды ответов
API широко использует информирование при помощи кодов ответов. Приложение должно корректно их обрабатывать.


При каждой ошибке, помимо кода ответа, в теле ответа может быть выдана дополнительная информация, позволяющая разработчику понять причину соответствующего ответа.
