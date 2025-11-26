


- $content главной flow.php
- getBooksSelector нужно сделать selected = $_GET['book'] ?: $_POST['book'] ?: $_COOKIE['book']


- реализовать функции
  - baseForm
  - fs
  - addBookPage
  + getBooks
  + getBooksSelector
  - getBookName
  - averagePageTime
  - autoPagination
  - getLastBookPage
  - getWords
  плюс какие-то методы из init.php нужно будет

- реализовать блоки
  - Внесение новых слов в базу чтения
  - readpage Регистрация прочитанных страниц
  - get_word Получение истории переводов слов
  - get_page Последняя страница книги
  - save_word Редактирование слов
  - edit, remind, log
  - import

- перенос файлов
  - category.php
  - exercise.php
  - memory.php
  - friends.php
  - flow.php
