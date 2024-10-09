CREATE TABLE IF NOT EXISTS parents
(
	ID           SERIAL PRIMARY KEY,
	email        TEXT UNIQUE NOT NULL,
	phone_number TEXT        NOT NULL,
	full_name    TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS child_profile
(
	ID           SERIAL PRIMARY KEY,
	full_name    TEXT NOT NULL,
	birthdate    DATE NOT NULL,
	gender       TEXT NOT NULL CHECK (gender IN ('Female', 'Male')),
	medical_info TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses
(
	ID           SERIAL PRIMARY KEY,
	zipcode      TEXT    NOT NULL,
	street       TEXT    NOT NULL,
	house_number INTEGER NOT NULL,
	state        TEXT    NOT NULL,
	city         TEXT    NOT NULL
);

CREATE TABLE modalities
(
	ID   SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS enrollments_shift
(
	ID   SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS enrollments
(
	ID         SERIAL PRIMARY KEY,
	registered BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS enrollments_terms
(
	enrollment_fk   INTEGER NOT NULL,
	terms_agreement BOOLEAN NOT NULL CHECK (terms_agreement = TRUE),
	FOREIGN KEY (enrollment_fk) REFERENCES enrollments (ID)
);


---------------------------------
----- Many to Many tables -------
---------------------------------

CREATE TABLE IF NOT EXISTS enrollments_shifts
(
	enrollment_fk        INTEGER NOT NULL,
	enrollments_shift_fk INTEGER NOT NULL,
	FOREIGN KEY (enrollments_shift_fk) REFERENCES enrollments_shift (ID),
	FOREIGN KEY (enrollment_fk) REFERENCES enrollments (ID)
);

CREATE TABLE IF NOT EXISTS child_parents
(
	enrollment_fk INTEGER NOT NULL,
	parents_fk    INTEGER NOT NULL,
	FOREIGN KEY (parents_fk) REFERENCES parents (ID),
	FOREIGN KEY (enrollment_fk) REFERENCES enrollments (ID)
);

CREATE TABLE IF NOT EXISTS student_modalities
(
	enrollment_fk INTEGER NOT NULL,
	modalities_fk INTEGER NOT NULL,
	FOREIGN KEY (modalities_fk) REFERENCES modalities (ID),
	FOREIGN KEY (enrollment_fk) REFERENCES enrollments (ID)
);

CREATE TABLE IF NOT EXISTS enrollments_addresses
(
	enrollment_fk INTEGER NOT NULL,
	addresses_fk  INTEGER NOT NULL,
	FOREIGN KEY (addresses_fk) REFERENCES addresses (ID),
	FOREIGN KEY (enrollment_fk) REFERENCES enrollments (ID)
);

CREATE TABLE IF NOT EXISTS enrollments_child
(
	child_profile_fk INTEGER NOT NULL,
	enrollment_fk    INTEGER NOT NULL,
	FOREIGN KEY (enrollment_fk) REFERENCES enrollments (ID),
	FOREIGN KEY (child_profile_fk) REFERENCES child_profile (ID)
);

--------------------
---- Procedures ----
--------------------

CREATE OR REPLACE PROCEDURE insert_child_profile(
	enrollmentId INTEGER,
	fullName TEXT,
	birthdate_arg DATE,
	gender_arg TEXT,
	medicalInfo TEXT
) AS
$BODY$
DECLARE
	insertedId INTEGER;
BEGIN
	INSERT INTO child_profile ( full_name
							  , birthdate
							  , gender
							  , medical_info)
	VALUES (fullName, birthdate_arg, gender_arg, medicalInfo)
	RETURNING ID INTO insertedId;

	INSERT INTO enrollments_child (enrollment_fk, child_profile_fk) VALUES (enrollmentId, insertedId);
END
$BODY$
	LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE insert_child_parent(
	enrollmentId INTEGER,
	fullName TEXT,
	phoneNumber TEXT,
	email_arg TEXT
) AS
$BODY$
DECLARE
	insertedId INTEGER;
BEGIN
	INSERT INTO parents( email
					   , phone_number
					   , full_name)
	VALUES (email_arg, phoneNumber, fullName)
	RETURNING ID INTO insertedId;

	INSERT INTO child_parents (enrollment_fk, parents_fk)
	VALUES (enrollmentId, insertedId);
END
$BODY$
	LANGUAGE plpgsql;


CREATE OR REPLACE PROCEDURE insert_address(
	enrollmentId INTEGER,
	zipcode_arg TEXT,
	street_arg TEXT,
	state_arg TEXT,
	city_arg TEXT,
	houseNumber INTEGER
) AS
$BODY$
DECLARE
	insertedId INTEGER;
BEGIN
	INSERT INTO addresses( zipcode
						 , street
						 , house_number
						 , state
						 , city)
	VALUES (zipcode_arg, street_arg, houseNumber, state_arg, city_arg);

	INSERT INTO enrollments_addresses (enrollment_fk, addresses_fk)
	VALUES (enrollmentId, insertedId);
END
$BODY$
	LANGUAGE plpgsql;



CREATE OR REPLACE PROCEDURE insert_modality(
	enrollmentId INTEGER,
	modalities TEXT[]
) AS $BODY$
BEGIN
	DELETE FROM student_modalities where enrollment_fk = enrollmentId;

	INSERT INTO student_modalities(enrollment_fk, modalities_fk)
	SELECT enrollmentId, m.ID FROM modalities m where m.name = ANY(modalities);
END
$BODY$
	LANGUAGE plpgsql;


CREATE OR REPLACE PROCEDURE insert_shift(
	enrollmentId INTEGER,
	shift TEXT
) AS $BODY$
BEGIN
	DELETE FROM enrollments_shifts where enrollment_fk = enrollmentId;

	INSERT INTO enrollments_shifts(enrollment_fk, enrollments_shift_fk)
	VALUES (enrollmentId, (SELECT ID FROM enrollments_shift where name = shift));
END
$BODY$
	LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE insert_term(
	enrollmentId INTEGER,
	term BOOLEAN
) AS $BODY$
BEGIN
	DELETE FROM enrollments_terms where enrollment_fk = enrollmentId;

	INSERT INTO enrollments_terms(enrollment_fk, terms_agreement)
	VALUES (enrollmentId, term);
END
$BODY$
	LANGUAGE plpgsql;
