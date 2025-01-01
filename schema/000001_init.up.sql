create table subjects
(
    id          serial primary key,
    title       varchar(30) not null,
    image       text        not null,
    description text        not null,
    isCertificated VARCHAR(30)
);

insert into subjects (title, image, description,isCertificated)
values ('Базы Данны. SQL', 'test-photo-db.png',
        'Курс по SQL предлагает практическое изучение SQL с нуля, начиная с основ и заканчивая профессиональным уровнем. В рамках курса студенты получат возможность освоить основные концепции SQL.', 'false'),
       ('Верстка сайтов', 'test-photo-web.png',
        'Предлагаем уникальную возможность освоить основы создания стильных и современных веб-страниц, начиная с HTML и заканчивая CSS, открывая перед вами новые перспективы в области веб-разработки.', 'false'),
       ('БАС', 'praxisbas.png',
        'Курс по Беспилотным авиационным средствам (БАС) предназначен для всех, кто хочет погрузиться в захватывающий мир дронов и их применения в различных сферах. На протяжении курса участники изучат ключевые аспекты проектирования, эксплуатации и управления беспилотными летательными аппаратами.', 'false');


create table themes
(
    id          serial primary key,
    title       varchar(125) not null,
    description text         not null,
    subject_id  integer      not null,
    foreign key (subject_id) references subjects (id)
);

INSERT INTO themes (title, description, subject_id)
VALUES ('Основы SQL', 'Изучение основ SQL, создание таблиц, выполнение запросов SELECT, INSERT, UPDATE, DELETE.', 1),
       ('Оптимизация баз данных', 'Методы оптимизации баз данных, индексы, проектирование эффективных структур данных.',1),
       ('Работа с различными СУБД', 'Изучение особенностей работы с различными системами управления базами данных, такими как MySQL, PostgreSQL, SQL Server.', 1),
       ('Основы HTML', 'Изучение основных тегов и структуры HTML для создания веб-страниц.', 2),
       ('Основы CSS', 'Изучение каскадных таблиц стилей (CSS) для стилизации и оформления веб-страниц.', 2),
       ('Адаптивная верстка', 'Принципы создания адаптивных и отзывчивых веб-страниц для различных устройств.', 2),
       ('Введение в беспилотные авиационные системы (БПЛА)', 'Основные понятия и история развития БПЛА. Классификация беспилотников по типам и назначениям.', 3),
       ('Технологии и компоненты БПЛА', 'Беспилотные авиационные системы (БПЛА) используют ряд технологий, которые обеспечивают их функциональность и эффективность.', 3);


create table lessons
(
    id       serial primary key,
    upkeep   text    not null,
    theme_id integer not null,
    foreign key (theme_id) references themes (id)
);



insert into lessons (upkeep, theme_id)
values ('<h1 class="lh1">Адаптивная верстка</h1>
<p class="lps">У нас уже получилось что-то прикольное, но если отркрыть сайт на телефоне нас ждет не приятный сюрприз</p>
<p class="lps">Адаптивная верстка - проще говоря это когда сайт без всяких проблем отображается на устройствах с разными пропорциями экранов</p>
', 6),
       ('<h1 class="lh1">Основы CSS</h1>
<p class="lps">HTML - это круто, но хочеться что бы глазам было приятно смотреть на страницу</p>
<p class="lps">И так <b>CSS</b>(каскадные таблицы стилей) - создан для того что бы стилизировать сайты</p>
', 5),
       ('<h1 class="lh1">Основы HTML</h1>
<p class="lps">Основой для всех веб страниц служит html</p>
<p class="lps"><b>HTML</b> - это язык скриптовой разметки</p>
', 4),
       ('<h1 class="lh1">Работа с различными СУБД</h1>
<p class="lps">Есть различные системами управления базами данных, такие как:</p>
<ul>
<li>MySQL</li>
<li>PostgreSQL</li>
<li>SQL Server</li>
</ul>
', 3),
       ('<h1 class="lh1">Оптимизация баз данных</h1>
<p class="lps">В этой части мы поговорим про:</p>
<ol>
<li>Методы оптимизации БД</li>
<li>Индексы</li>
<li>Проектирование эффективных структур данных</li>
</ol>
', 2),
       ('<h1 class="lh1">Основы SQL</h1>
<p class="lps">Для начала создадим базу данных:</p>
<code class="lcmd">
	<p>CREATE DATABASE db_name</p>
</code>
<p class="lps">Всместо  <b>db_name</b> можно написать любое название для базы данных</p>', 1),
       (' <h1 class="lh1">Основные понятия и история развития БПЛА.</h1>
<p class="lps">
Беспилотные авиационные системы (БПЛА) представляют собой летательные аппараты, которые могут управляться без пилота на борту. Эти системы начали развиваться в начале 20 века, но получили широкое распространение только в последние десятилетия благодаря достижениям в области технологий и миниатюризации. Важными этапами в истории БПЛА стали создание первых радиоуправляемых моделей, а затем и более сложных систем, использующих GPS и автоматизированные технологии.
</p>
<h1 class="lh1">Классификация беспилотников</h1>
<ul>
<li><b>Типы:</b> фиксированные крылья, вертолеты, мультикоптеры и т.д.</li>
<li><b>Назначение:</b> гражданские (например, для сельского хозяйства, мониторинга) и военные (разведка, ударные операции).</li>
<li><b>Размеры:</b> от микро-БПЛА до больших беспилотников, способных нести значительные грузы.</li>
</ul>
<h1 class="lh1">Значение БПЛА</h1>
<p class="lps">БПЛА открывают новые возможности для сбора данных, мониторинга и выполнения задач, которые ранее были труднодоступны или небезопасны для человека. Они становятся важным инструментом в борьбе с изменением климата, обеспечении безопасности и повышении эффективности в различных отраслях.<p>', 7),
       ('<h1 class="lh1">Технологии и компоненты БПЛА</h1>
<ul>
<li><b>Сенсоры:</b>  БПЛА оснащаются различными типами сенсоров, такими как камеры, инфракрасные и ультразвуковые датчики, которые позволяют собирать данные о окружающей среде. Это важно для выполнения задач, таких как мониторинг, картографирование и поиск объектов.</li>
<li><b>Навигация:</b> БПЛА используют глобальные навигационные спутниковые системы (GNSS), такие как GPS, для определения своего местоположения. Кроме того, применяются инерциальные навигационные системы (INS) для повышения точности.</li>
<li><b>Связь:</b> БПЛА используют глобальные навигационные спутниковые системы (GNSS), такие как GPS, для определения своего местоположения. Кроме того, применяются инерциальные навигационные системы (INS) для повышения точности.</li>
</ul>', 8);


create table users
(
    id          text primary key,
    email       varchar(255) not null,
    password    text         not null,
    role        VARCHAR(100) NOT NULL,
    create_date text         not null
);
create table done_lessons
(
    id       serial primary key,
    theme_id integer  not null,
    user_id  text not null,
    foreign key (theme_id) references themes (id),
    foreign key (user_id) references users (id)
);


create table profiles
(
    id          serial primary key,
    user_id     TEXT         not null,
    description varchar(125) not null,
    phone       varchar(100) not null,
    full_name   varchar(125) not null,
    image       text         not null,
    score       integer      not null default 0,
    foreign key (user_id) references users (id)
);


create table last_subjects
(
    id        serial primary key,
    user_id   text not null,
    subjects_id integer  not null
);

INSERT INTO last_subjects(user_id,subjects_id) values ('b43a1721-2bc3-4421-8e70-b7bd932ad802',1);


INSERT INTO users (id, email, password, role, create_date)
VALUES ('b43a1721-2bc3-4421-8e70-b7bd932ad802', 'themountok@gmail.com',
        '$2a$10$yvIcyoUWkoBmT8PDMQ.L9.9zDCvel76DOexLupCD4m/CJB1jToEAy', 'admin', '15.08.2024-22:18');

INSERT INTO profiles (user_id, description, phone, full_name, image, score)
VALUES ('b43a1721-2bc3-4421-8e70-b7bd932ad802', 'Привет меня зовут Ислам.', '+79280229349', 'Ислам Дашуев',
        'avatar_for_user_b43a1721-2bc3-4421-8e70-b7bd932ad802IMG_2907.jpg', 120);

INSERT INTO users (id, email, password, role, create_date)
VALUES ('84203e90-3c01-4531-8479-7501e3b92882', 'admin@gmail.com',
        '$2a$10$yvIcyoUWkoBmT8PDMQ.L9.9zDCvel76DOexLupCD4m/CJB1jToEAy', 'admin', '01.09.2024-11:09');

INSERT INTO profiles (user_id, description, phone, full_name, image, score)
VALUES ('84203e90-3c01-4531-8479-7501e3b92882', 'Администратор', '+79280229349', 'Admin', 'admin.png', 666);



-- ############################# ДЛЯ ТЕСТОВ

CREATE TABLE subject_test (
       id serial PRIMARY KEY,
       title VARCHAR(125) NOT NULL,
       subject_id INTEGER NOT NULL,
       FOREIGN KEY (subject_id) REFERENCES subjects(id)
);

INSERT INTO subject_test(title, subject_id) VALUES ('Основные операторы SQL',1);
INSERT INTO subject_test(title, subject_id) VALUES ('Построение запросов SQL',1);
INSERT INTO subject_test(title, subject_id) VALUES ('Основные теги HTML',2);

CREATE TABLE test_questions (
       id serial PRIMARY KEY,
       question TEXT NOT NULL,
       options TEXT NOT NULL,
       answer TEXT NOT NULL,
       test_id INTEGER NOT NULL,
       FOREIGN KEY (test_id) REFERENCES subject_test(id)
);

INSERT INTO test_questions(question,options,answer,test_id) VALUES ('Какой оператор позволяет делать выборку из таблицы?','INSERT INTO;SELECT;UPDATE;DELETE','SELECT',1);
INSERT INTO test_questions(question,options,answer,test_id) VALUES ('Какой оператор позволяет делать удалять из таблицы данные?','INSERT INTO;SELECT;UPDATE;DELETE','DELETE',1);
INSERT INTO test_questions(question,options,answer,test_id) VALUES ('Тест вопрос?','INSERT INTO;SELECT;UPDATE;DELETE','SELECT',2);



INSERT INTO test_questions(question,options,answer,test_id) VALUES 
('Что такое база данных?','Хранилище данных;Программа для обработки данных;Система управления данными;Интерфейс пользователя','Хранилище данных',4),
('Какой язык используется для работы с базами данных?','HTML;SQL;JavaScript;Python','SQL',4),
('Что такое SQL?','Язык программирования;Язык запросов к базе данных;Язык разметки;Операционная система','Язык запросов к базе данных',4),
('Какой из следующих типов баз данных является реляционной?','MongoDB;MySQL;Redis;Cassandra','Что такое SQL?',4),
('Что такое первичный ключ?','Уникальный идентификатор записи;Индекс для ускорения поиска;Ссылка на другую таблицу;Необязательное поле','Уникальный идентификатор записи',4),
('Какой оператор используется для извлечения данных из базы данных?','INSERT;SELECT;UPDATE;DELETE','SELECT',4),
('Что делает оператор INSERT?','Удаляет данные;Обновляет данные;Вставляет новые данные;Извлекает данные','Вставляет новые данные',4),
('Что такое JOIN в SQL?','Оператор для объединения таблиц;Функция для агрегации данных;Команда для удаления записей;Тип базы данных','Оператор для объединения таблиц',4),
('Какой оператор используется для изменения существующих данных в таблице?','INSERT;SELECT;UPDATE;DELETE','UPDATE',4),
('Что такое индекс в базе данных?','Структура для ускорения поиска;Тип базы данных;Функция для резервного копирования;Метод шифрования','Структура для ускорения поиска',4),
('Какой тип связи между таблицами обозначает один ко многим?','1:1;1:М;n:m;М:1','1:М',4),
('Что делает оператор DROP DATABASE?','Удаляет таблицу;Удаляет базу данных;Удаляет данные из таблицы;Изменяет структуру таблицы','Удаляет базу данных',4),
('Какой из следующих операторов используется для фильтрации результатов запроса?','ORDER BY;GROUP BY;WHERE;HAVING','WHERE',4),
('Что такое подзапрос в SQL?','Запрос внутри другого запроса;Запрос на создание новой таблицы;Запрос на удаление данных;Запрос на изменение структуры базы данных','Запрос внутри другого запроса',4),
('Какой оператор используется для сортировки результатов запроса?','ORDER BY;SORT BY;GROUP BY;FILTER','ORDER BY',4),
('Что делает функция COUNT() в SQL?','Считает количество строк;Считает сумму значений;Находит максимальное значение;Находит минимальное значение','Считает количество строк',4),
('Что делает оператор GROUP BY?','Объединяет строки с одинаковыми значениями;Сортирует строки по определенному столбцу;Фильтрует строки по условию;Создает новую таблицу','Объединяет строки с одинаковыми значениями',4),
('Какой из этих операторов используется для объединения нескольких условий в запросе?','AND;OR;NOT;IN','AND',4),
('Что такое нормализация базы данных?','Процесс оптимизации структуры базы данных;Процесс резервного копирования данных;Процесс шифрования данных;Процесс создания индексов','Процесс оптимизации структуры базы данных',4),
('Какой из следующих типов баз данных не является реляционным?','PostgreSQL;SQLite;MongoDB;Oracle','MongoDB',4),
('Что делает оператор DISTINCT?','Удаляет дубликаты из результата запроса;Сортирует результаты;Объединяет таблицы;Создает новый столбец','Удаляет дубликаты из результата запроса',4),
('Какой оператор используется для добавления условий в группированные запросы?','HAVING;WHERE;ORDER BY;GROUP BY','HAVING',4),
('Что такое транзакция в контексте баз данных?','Набор операций, которые выполняются как единое целое;Процесс резервного копирования данных;Метод шифрования данных;Тип базы данных','Набор операций, которые выполняются как единое целое',4),
('Что такое внешние ключи?','Поля, которые связывают таблицы;Индексы для ускорения поиска;Запросы для извлечения данных;Типы данных','Поля, которые связывают таблицы',4),
('Какой из следующих операторов используется для изменения структуры таблицы?','ALTER TABLE;CREATE TABLE;DROP TABLE;INSERT INTO','ALTER TABLE',4),
('Какой оператор используется для создания новой таблицы?','CREATE TABLE;ADD TABLE;NEW TABLE;INSERT TABLE','CREATE TABLE',4),
('Какой из следующих операторов используется для объединения двух или более таблиц?','JOIN;UNION;MERGE;LINK','JOIN',4),
('Какой оператор используется для ограничения количества возвращаемых строк?','LIMIT;TOP;FETCH;COUNT','LIMIT',4),
('Что делает функция AVG() в SQL?','Находит среднее значение;Находит максимальное значение;Находит минимальное значение;Считает количество строк','Находит среднее значение',4),
('Что такое SQL-инъекция?','Уязвимость безопасности в веб-приложениях;Метод резервного копирования данных;Тип базы данных;Процесс нормализации','Уязвимость безопасности в веб-приложениях',4),
('Какой оператор используется для проверки наличия значения в подмножестве?','IN;EXISTS;ANY;ALL','IN',4),
('Что такое временная таблица?','Таблица, которая существует только во время сессии пользователя;Таблица для хранения резервных копий данных;Тип базы данных;Индекс для ускорения поиска','Таблица, которая существует только во время сессии пользователя',4);

-- SELECT 
--     dt.id, 
--     dt.test_id, 
--     st.subject_id, 
--     dt.user_id, 
--     dt.points,
--     COUNT(tq.id) AS question_count
-- FROM 
--     done_test dt 
-- JOIN 
--     subject_test st ON st.id = dt.test_id  
-- LEFT JOIN 
--     test_questions tq ON tq.test_id = dt.test_id
-- WHERE 
--     dt.user_id = 'b43a1721-2bc3-4421-8e70-b7bd932ad802'
--     AND st.subject_id = 1
-- GROUP BY 
--     dt.id, dt.test_id, st.subject_id, dt.user_id, dt.points
-- ORDER BY 
--     dt.test_id ASC;

CREATE TABLE done_test (
       id serial PRIMARY KEY,
       test_id INTEGER NOT NULL,
       user_id TEXT NOT NULL,
       points INTEGER NOT NULL,
       FOREIGN KEY (test_id) REFERENCES subject_test(id),
       FOREIGN KEY (user_id) REFERENCES users(id)
);

-- INSERT INTO done_test(test_id,user_id,points) VALUES (1,'b43a1721-2bc3-4421-8e70-b7bd932ad802',1);
-- INSERT INTO done_test(test_id,user_id,points) VALUES (2,'b43a1721-2bc3-4421-8e70-b7bd932ad802',1);





-- Надо добавить в хостинг БД 
CREATE TABLE certificates (
       id serial PRIMARY KEY,
       user_id TEXT NOT NULL,
       subject_id INTEGER NOT NULL,
       get_date DATE NOT NULL DEFAULT CURRENT_DATE,
       FOREIGN KEY (user_id) REFERENCES users(id),
       FOREIGN KEY (subject_id) REFERENCES subjects(id)
);