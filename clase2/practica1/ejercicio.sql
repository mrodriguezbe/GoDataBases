USE movies_db;

-- Mostrar el título y el nombre del género de todas las series.
select s.title, g.name from series as s
inner join genres as g on g.id = s.genre_id;

-- Mostrar el título de los episodios, el nombre y apellido de los actores que trabajan en cada uno de ellos.
select aux.first_name, aux.last_name, e.title from
	(select a.first_name, a.last_name, ae.episode_id from actors a
	inner join actor_episode ae on ae.actor_id = a.id) as aux
inner join episodes e on e.id = aux.episode_id


-- Mostrar el título de todas las series y el total de temporadas que tiene cada una de ellas.
select s.title, SUM(1) as seasons_amount from series as s
inner join seasons ss on s.id = ss.serie_id
group by s.title;

-- Mostrar el nombre de todos los géneros y la cantidad total de películas por cada uno, siempre que sea mayor o igual a 3.
select g.name, SUM(1) as genre_amount from genres g
inner join movies m on g.id = m.genre_id 
group by g.name
having genre_amount > 2;

-- Mostrar sólo el nombre y apellido de los actores que trabajan en todas las películas de la guerra de las galaxias y que estos no se repitan.
select distinct CONCAT(first_name," ",last_name) from actors a
inner join (select am.actor_id, am.movie_id from actor_movie am 
	inner join (select title, id from movies m
		where title LIKE "%galaxia%") as swmovies on am.movie_id  = swmovies.id) as aux on a.id = aux.actor_id;







