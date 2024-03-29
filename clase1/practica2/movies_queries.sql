USE movies_db;

-- Mostrar todos los registros de la tabla de movies.
SELECT * FROM movies as m

-- Mostrar el nombre, apellido y rating de todos los actores.
SELECT first_name, last_name, rating FROM actors as a

-- Mostrar el título de todas las series y usar alias para que tanto el nombre de la tabla como el campo estén en español.
SELECT title as "Titulo" FROM series as s

-- Mostrar el nombre y apellido de los actores cuyo rating sea mayor a 7.5.
SELECT first_name, last_name, rating FROM actors as a
WHERE rating > 7.5;

-- Mostrar el título de las películas, el rating y los premios de las películas con un rating mayor a 7.5 y con más de dos premios.
SELECT title, rating, awards FROM movies as m
WHERE rating > 7.5 AND awards > 2;

-- Mostrar el título de las películas y el rating ordenadas por rating en forma ascendente.
SELECT title, rating FROM movies as m
ORDER BY rating ASC;

-- Mostrar los títulos de las primeras tres películas en la base de datos.
SELECT title FROM movies as m
LIMIT 3;

-- Mostrar el top 5 de las películas con mayor rating.
SELECT title, rating FROM movies as m
ORDER BY rating DESC
LIMIT 5;


-- Listar los primeros 10 actores.
SELECT first_name, last_name FROM actors as a
LIMIT 10;

-- Mostrar el título y rating de todas las películas cuyo título sea de Toy Story.
SELECT m.title, m.rating FROM movies as m
WHERE m.title LIKE "Toy%";

-- Mostrar a todos los actores cuyos nombres empiezan con Sam.
SELECT first_name, last_name FROM actors as a
WHERE a.first_name LIKE "Sam%";

-- Mostrar el título de las películas que salieron entre el 2004 y 2008.
SELECT title, release_date  FROM movies as m
WHERE YEAR(m.release_date) BETWEEN 2004 AND 2008;

-- Traer el título de las películas con el rating mayor a 3, con más de 1 premio y con fecha de lanzamiento entre el año 1988 al 2009. Ordenar los resultados por rating.
SELECT title, rating, awards, release_date FROM movies as m
WHERE awards > 1 AND rating > 3 AND YEAR(m.release_date) BETWEEN 1988 AND 2009;

