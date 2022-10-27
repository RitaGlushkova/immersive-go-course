CREATE DATABASE "go-server-database-test"
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


INSERT INTO public.images (title, url, alt_text) VALUES ('White Cat','https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg','White cat sitting and looking to the left'),('Catch a Ball','https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg','A dog jumping up catching a red ball');