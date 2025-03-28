-- create table investigator (
-- 	id serial primary key,
-- 	investigator_id varchar(32),
-- 	first_name varchar(64),
-- 	last_name varchar(64),
-- 	email_address varchar(128),
-- 	start_date date,
-- 	end_date date,
-- 	role_code varchar,
-- 	nsf_id varchar
-- );
--
-- create table program_element (
-- 	id serial primary key,
-- 	code varchar(32),
-- 	text text
-- );
--
-- create table program_reference(
-- 	id serial primary key,
-- 	code varchar(32),
-- 	text text
-- );
--
-- create table institution(
-- 	id serial primary key,
-- 	name varchar(128),
-- 	city_name varchar(64),
-- 	zip_code varchar(32),
-- 	phone_number varchar(16),
-- 	street_address varchar(128),
-- 	country_name varchar(64),
-- 	street_name varchar(64),
-- 	state_code varchar(16)
-- );
--
-- create table directorate(
-- 	id serial primary key,
-- 	long_name text,
-- 	abbreviation varchar(64)
-- );
--
-- create table division(
-- 	id serial primary key,
-- 	long_name text,
-- 	abbreviation varchar(64)
-- );
--
-- create table organization(
-- 	id serial primary key,
-- 	code varchar(64),
-- 	directorate_id integer references directorate(id),
-- 	division_id integer references division(id)
-- );
--
-- create table program_officer(
-- 	id serial primary key,
-- 	sign_block_name varchar(128),
-- 	email varchar(128),
-- 	phone varchar(16)
-- );
--
-- create table award (
-- 	id serial primary key,
-- 	title text,
-- 	award_effective_date date,
-- 	award_expiration_date date,
-- 	award_amount numeric,
-- 	abstract_text text,
-- 	program_officer_id integer references program_officer(id),
-- 	organization_id integer references organization(id),
-- 	institution_id integer references institution(id)
-- );
--
-- create table award_investigator(
-- 	investigator_id integer references investigator(id),
-- 	award_id integer references award(id)
-- );
--
-- create table award_program_element(
-- 	program_element_id integer references program_element(id),
-- 	award_id integer references award(id)
-- );
--
-- create table award_program_reference(
-- 	program_reference_id integer references program_reference(id),
-- 	award_id integer references award(id)
-- );
--
