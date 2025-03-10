create table investigator (
	id integer primary key,
	first_name varchar(64),
	last_name varchar(64),
	email_address varchar(128),
	start_date date,
	end_date date,
	role_code varchar,
	nsf_id varchar
);

create table program_element (
	id integer primary key,
	code varchar(32),
	text text
);

create table program_reference(
	id integer primary key,
	code varchar(32),
	text text
);

create table foa_information(
	id integer primary key,
	code varchar(32),
	name varchar(64)
);

create table institution(
	id integer primary key,
	name varchar(128),
	city_name varchar(64),
	zip_code varchar(32),
	phone_number varchar(16),
	street_address varchar(128),
	country_name varchar(64),
	street_name varchar(64),
	state_code varchar(16)
);

create table directorate(
	id integer primary key,
	long_name text,
	abbreviation varchar(64)
);

create table division(
	id integer primary key,
	long_name text,
	abbreviation varchar(64)
);

create table organization(
	id integer primary key,
	code varchar(64),
	directorate_id integer references directorate(id),
	division_id integer references division(id)
);

create table program_officer(
	id integer primary key,
	sign_block_name varchar(128),
	email varchar(128),
	phone varchar(16)
);

create table award_instrument(
	id integer primary key,
	value varchar(128)
);


create table award (
	id integer primary key,
	title varchar(128),
	award_effective_date date,
	award_expiration_date date,
	award_amount integer,
	program_officer_id integer references program_officer(id),
	organization_id integer references organization(id),
	award_instrument_id integer references award_instrument(id)
);

create table award_foa_information(
	foa_information_id integer references foa_information(id),
	award_id integer references award(id)
);

create table award_institution(
	institution_id integer references institution(id),
	award_id integer references award(id)
);

create table award_investigator(
	investigator_id integer references investigator(id),
	award_id integer references award(id)
);

create table award_program_element(
	program_element_id integer references program_element(id),
	award_id integer references award(id)
);

create table award_program_reference(
	program_reference_id integer references program_reference(id),
	award_id integer references award(id)
);


