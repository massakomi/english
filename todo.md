

- getBooksSelector нужно сделать selected = $_GET['book'] ?: $_POST['book'] ?: $_COOKIE['book']
- вынести posgres connect в .env

- реализовать функции
  - getWords
  плюс какие-то методы из init.php нужно будет

- реализовать блоки index.php
  - Внесение новых слов в базу чтения
  - readpage Регистрация прочитанных страниц
  - get_word Получение истории переводов слов
  - get_page Последняя страница книги
  - save_word Редактирование слов
  - edit, remind, log
  - import

- перенос файлов
  - exercise.php
  - category.php
  - memory.php
  - friends.php
