# Wordle-Vuetify Project 
### Authors: Merle Epping, Luca Völker and Matteo Nollenberg
### Note: This README is also included in the project itself
#
#
## To start the game WITH a Database do following things:

Setup a local Database with mysql:
- install mysql and set the password to: 1234
- Run: CREATE DATABASE wordleDB;
- Run: CREATE TABLE data (Name varchar(15), Versuche int, Wort varchar(5));

Setup the game:
- open a terminal in this directory: wordle-proket-ibt/Wordle-Vuetify
- Run: npm intsall
- Run: npm install chart.js vue-chartjs
- Run: npm run build
Run the game:
- Run: go run .\src\main.go


## To start the game WITHOUT a Database do following things:

Setup the game:
- open a terminal in this directory: wordle-proket-ibt/Wordle-Vuetify
- Run: npm intsall
- Run: npm install chart.js vue-chartjs
Run the game:
- Run: npm run dev
