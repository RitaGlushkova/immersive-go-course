CREATE DATABASE "go-server-database"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE TABLE IF NOT EXISTS public.images
(
    id serial NOT NULL,
    title text NOT NULL,
    url text NOT NULL,
    alt_text text,
    pixels int,
    PRIMARY KEY (id)
);

INSERT INTO public.images (title, url, alt_text, pixels) VALUES ('Sunset', 'https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80', 'Clouds at sunset',9408494),('Mountain','https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80','A mountain at sunset',79404),('Cat','https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80','A cool cat',7899);