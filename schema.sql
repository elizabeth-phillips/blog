CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tools(
	ID uuid         DEFAULT uuid_generate_v4(),
	CREATOR         TEXT       		NOT NULL,
	TITLE           TEXT       		NOT NULL,
	CONTENT    	    TEXT       		NOT NULL,
	CREATED_AT     	DATE default current_date,
	UPDATE_AT		DATE default current_date
);