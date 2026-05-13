#!/usr/bin/bash

cd sql/schema
goose postgres "postgres://postgres@localhost:5432/chirpy" up
cd -
