CREATE DATABASE wallet;

CREATE TABLE accounts (
   id TEXT PRIMARY KEY,
   name TEXT NOT NULL
);

CREATE TABLE payments (
   id TEXT PRIMARY KEY,
   source TEXT NOT NULL REFERENCES accounts(id),
   destination TEXT NOT NULL REFERENCES accounts(id),
   amount INT NOT NULL
);
