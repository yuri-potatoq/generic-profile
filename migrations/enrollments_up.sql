CREATE TABLE IF NOT EXISTS parents (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone_number TEXT NOT NULL,
    full_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS child_profile (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,    
    full_name TEXT NOT NULL,
    birthdate DATETIME NOT NULL,
    gender TEXT NOT NULL CHECK (gender IN ('Female', 'Male')),
    medical_info TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    zipcode TEXT NOT NULL,
    street TEXT NOT NULL,
    house_number INTEGER NOT NULL,
    state TEXT NOT NULL,
    city TEXT NOT NULL
);

CREATE TABLE modalities (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS enrollments_shift (
   ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
   name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS enrollments (
    ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	registered BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS enrollments_terms (
	enrollment_fk INTEGER NOT NULL,
	FOREIGN KEY(enrollment_fk) REFERENCES enrollments(ID),
	-- we can't allow enrollments without terms accepted
	terms_agreement BOOLEAN NOT NULL
		CHECK ( terms_agreement == true )
);


---------------------------------
----- Many to Many tables -------
---------------------------------

CREATE TABLE IF NOT EXISTS enrollments_shifts (
    enrollment_fk INTEGER NOT NULL,
    enrollments_shift_fk INTEGER NOT NULL,
    FOREIGN KEY(enrollments_shift_fk) REFERENCES enrollments_shift(ID),
    FOREIGN KEY(enrollment_fk) REFERENCES enrollments(ID)
);

CREATE TABLE IF NOT EXISTS child_parents (
    enrollment_fk INTEGER NOT NULL,
    parents_fk INTEGER NOT NULL,
    FOREIGN KEY(parents_fk) REFERENCES parents(ID),
    FOREIGN KEY(enrollment_fk) REFERENCES enrollments(ID)
);

CREATE TABLE IF NOT EXISTS student_modalities (
    enrollment_fk INTEGER NOT NULL,
    modalities_fk INTEGER NOT NULL,
    FOREIGN KEY(modalities_fk) REFERENCES modalities(ID),
    FOREIGN KEY(enrollment_fk) REFERENCES enrollments(ID)
);

CREATE TABLE IF NOT EXISTS enrollments_addresses (
    enrollment_fk INTEGER NOT NULL,
    addresses_fk INTEGER NOT NULL,
    FOREIGN KEY(addresses_fk) REFERENCES addresses(ID),
    FOREIGN KEY(enrollment_fk) REFERENCES enrollments(ID)
);

CREATE TABLE IF NOT EXISTS enrollments_child (
    child_profile_fk INTEGER NOT NULL,
    enrollment_fk INTEGER NOT NULL,
    FOREIGN KEY(enrollment_fk) REFERENCES enrollments(ID),
    FOREIGN KEY(child_profile_fk) REFERENCES child_profile(ID)
);


