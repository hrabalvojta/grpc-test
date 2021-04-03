# DVD rental database 
## types
AGGREGATE - new func like avg(), min(), max(), ...
DOMAIN - variable (integer) with constraint (check) ex.:mail,p. number, ...
TYPE - new data type (combobox)
CONSTRAINT - primary foreign keys
SEQUENCE - seq fir something_id
INDEX - indexing column, btree or gist
TRIGGER - execute function, procedure
TEMPORARY TABLE - for session only

## Schema of database
```
DATABASE dvdrental
TABLE public.customer
    CONSTRAINT customer_pkey PRIMARY KEY (customer_id);
    CONSTRAINT customer_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE ON DELETE
    INDEX idx_last_name ON public.customer USING btree (last_name);
    INDEX idx_fk_store_id ON public.customer USING btree (store_id);
    INDEX idx_fk_address_id ON public.customer USING btree (address_id);
    SEQUENCE public.customer_customer_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.customer
    FUNCTION public.rewards_report
    TEMPORARY TABLE tmpCustomer
TABLE public.actor
    CONSTRAINT actor_pkey PRIMARY KEY (actor_id);
    INDEX idx_actor_last_name ON public.actor USING btree (last_name);
    SEQUENCE public.actor_actor_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.actor
TABLE public.category
    CONSTRAINT category_pkey PRIMARY KEY (category_id);
    SEQUENCE public.category_category_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.category
TABLE public.film
    CONSTRAINT film_pkey PRIMARY KEY (film_id);
    CONSTRAINT film_language_id_fkey FOREIGN KEY (language_id) REFERENCES public.language(language_id) ON UPDATE ON DELETE
    INDEX idx_title ON public.film USING btree (title);
    INDEX idx_fk_language_id ON public.film USING btree (language_id);
    INDEX film_fulltext_idx ON public.film USING gist (fulltext);
    SEQUENCE public.film_film_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.film
    TRIGGER film_fulltext_trigger BEFORE INSERT OR UPDATE ON public.film
    FUNCTION public.film_in_stock
    FUNCTION public.last_day # new year shenanigans 
    FUNCTION public.film_not_in_stock
    FUNCTION public.get_customer_balance # INTO, nice math
    VIEW public.film_list
    VIEW public.nicer_but_slower_film_list
TABLE public.film_actor
    CONSTRAINT film_actor_pkey PRIMARY KEY (actor_id, film_id);
    CONSTRAINT film_actor_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES public.actor(actor_id) ON UPDATE ON DELETE
    CONSTRAINT film_actor_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film(film_id) ON UPDATE ON DELETE
    INDEX idx_fk_film_id ON public.film_actor USING btree (film_id);
    TRIGGER last_updated BEFORE UPDATE ON public.film_actor
TABLE public.film_category
    CONSTRAINT film_category_pkey PRIMARY KEY (film_id, category_id);
    CONSTRAINT film_category_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(category_id) ON UPDATE ON DELETE
    CONSTRAINT film_category_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film(film_id) ON UPDATE ON DELETE
    TRIGGER last_updated BEFORE UPDATE ON public.film_category
TABLE public.address
    CONSTRAINT address_pkey PRIMARY KEY (address_id);
    CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES public.city(city_id);
    INDEX idx_fk_city_id ON public.address USING btree (city_id);
    SEQUENCE public.address_address_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.address
    VIEW public.actor_info
TABLE public.city
    CONSTRAINT city_pkey PRIMARY KEY (city_id);
    CONSTRAINT fk_city FOREIGN KEY (country_id) REFERENCES public.country(country_id);
    INDEX idx_fk_country_id ON public.city USING btree (country_id);
    SEQUENCE public.city_city_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.city
TABLE public.country
    CONSTRAINT country_pkey PRIMARY KEY (country_id);
    SEQUENCE public.country_country_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.country
    VIEW public.customer_list
TABLE public.inventory
    CONSTRAINT inventory_pkey PRIMARY KEY (inventory_id);
    CONSTRAINT inventory_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film(film_id) ON UPDATE ON DELETE;
    INDEX idx_store_id_film_id ON public.inventory USING btree (store_id, film_id);
    SEQUENCE public.inventory_inventory_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.inventory
    FUNCTION public.inventory_held_by_customer
    FUNCTION public.inventory_in_stock
TABLE public.language
    CONSTRAINT language_pkey PRIMARY KEY (language_id);
    SEQUENCE public.language_language_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.language
TABLE public.payment
    CONSTRAINT payment_pkey PRIMARY KEY (payment_id);
    CONSTRAINT payment_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer(customer_id) ON UPDATE ON DELETE;
    CONSTRAINT payment_rental_id_fkey FOREIGN KEY (rental_id) REFERENCES public.rental(rental_id) ON UPDATE ON DELETE;
    CONSTRAINT payment_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES public.staff(staff_id) ON UPDATE ON DELETE;
    INDEX idx_fk_staff_id ON public.payment USING btree (staff_id);
    INDEX idx_fk_rental_id ON public.payment USING btree (rental_id);
    INDEX idx_fk_customer_id ON public.payment USING btree (customer_id);
    SEQUENCE public.payment_payment_id_seq
TABLE public.rental
    CONSTRAINT rental_pkey PRIMARY KEY (rental_id);
    CONSTRAINT rental_staff_id_key FOREIGN KEY (staff_id) REFERENCES public.staff(staff_id);
    CONSTRAINT rental_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer(customer_id) ON UPDATE ON DELETE
    CONSTRAINT rental_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventory(inventory_id) ON UPDATE ON DELETE;
    UNIQUE INDEX idx_unq_rental_rental_date_inventory_id_customer_id ON public.rental USING btree (rental_date, inventory_id, customer_id);
    INDEX idx_fk_inventory_id ON public.rental USING btree (inventory_id);
    SEQUENCE public.rental_rental_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.rental
    VIEW public.sales_by_film_category
TABLE public.staff
    CONSTRAINT staff_pkey PRIMARY KEY (staff_id);
    CONSTRAINT staff_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE ON DELETE;
    SEQUENCE public.staff_staff_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.staff
TABLE public.store
    CONSTRAINT store_pkey PRIMARY KEY (store_id);
    CONSTRAINT store_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE ON DELETE;
    CONSTRAINT store_manager_staff_id_fkey FOREIGN KEY (manager_staff_id) REFERENCES public.staff(staff_id) ON UPDATE ON DELETE;
    UNIQUE INDEX idx_unq_manager_staff_id ON public.store USING btree (manager_staff_id);
    SEQUENCE public.store_store_id_seq
    TRIGGER last_updated BEFORE UPDATE ON public.store
    VIEW public.sales_by_store
    VIEW public.staff_list

AGGREGATE public.group_concat # is using FUNCTION public._group_concat
DOMAIN public.year #year
FUNCTION public._group_concat # merge 2 strings to 1
FUNCTION public.last_updated
TYPE public.mpaa_rating #PG rating
```