CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  name text NOT NULL ,
  description text ,
  author text,
  pages INT
);

insert into books (name, description, author, pages) values ('name1', 'description1', 'author1', 100);
insert into books (name, description, author, pages) values ('name2', 'description2', 'author2', 200);
insert into books (name, description, author, pages) values ('name3', 'description3', 'author3', 300);
insert into books (name, description, author, pages) values ('name4', 'description4', 'author4', 400);
insert into books (name, description, author, pages) values ('name5', 'description5', 'author5', 500);