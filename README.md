﻿# Profile APP

Minimal API doing CRUD with partial updates to profile system.


### TODO
- [x] Submit feature. When the profile is fulled with all information needed, 
we submit the profile to check data and save as registered profile.
- [ ] Handle multiple relations at once, create guards to MAX relations number. For now, just order and pick latest one.
- [ ] Use transactions to do atomic updates for all fields.
- [ ] Create migrations flow (pick a migrations lib or DIY).

### API regular flow diagram
```mermaid
sequenceDiagram
    %% primeiro patch para criar nova matricula
    Client->>+API: PATCH /profile ["childInfos" sem ID]

    API-->>-Client: Response com ID temp criado [erro caso dados inválidos]

    API->>+DB: Cria matricula TEMP [registered=false]
    DB-->>-API: ok

    %% patch para atualizar os dados
    Note over Client,API: Agora seguem quantos PATCHS subsequentes forem necessários.
    loop para cada campo: childInfos, parent, address, modalities, enrollmentShift, terms
        Client->>+API: PATCH /profile ["payload" com ID]
        API-->>-Client: response [erro caso dados inválidos]

        API->>+DB: Atualiza matricula TEMP com ID recebido
        DB-->>-API: ok
    end

    %% Submit da matricula
    Client->>+API: POST /submit ["ID" da matricula TEMP]
    API-->>-Client: response [erro caso dados inválidos]

    API->>+DB: Atualiza matricula TEMP com ID recebido [registered=true]
    DB-->>-API: ok
```

### Database Modeling
```mermaid
---
title: API Diagram
config:
  theme: dark
---
erDiagram
    parents {
        INTEGER ID PK
        TEXT email
        TEXT phone_number
        TEXT full_name
    }

    child_profile {
        INTEGER ID PK
        TEXT full_name
        DATE birthdate
        TEXT gender "CHECK (gender = 'Female' OR gender = 'Male')"
        TEXT medical_info
    }

    addresses {
        INTEGER ID PK
        TEXT zipcode
        TEXT street
        INTEGER house_number
        TEXT state
        TEXT city
    }

    modalities {
        INTEGER ID PK
        TEXT name
    }

    enrollments_shift {
        INTEGER ID PK
        TEXT name
    }

    enrollments_terms {
        INTEGER enrollment_fk FK
        BOOLEAN terms_agreement
    }

    enrollments_shifts {
        INTEGER enrollment_fk FK
        INTEGER enrollments_shift_fk FK
    }

    child_parents {
        INTEGER enrollment_fk FK
        INTEGER parents_fk FK
    }

    student_modalities {
        INTEGER enrollment_fk FK
        INTEGER modalities_fk FK
    }

    enrollments_addresses {
        INTEGER enrollment_fk FK
        INTEGER addresses_fk FK
    }

    enrollments_child {
        INTEGER child_profile_fk FK
        INTEGER enrollment_fk FK
    }

    enrollments {
        INTEGER ID PK
    }

    %% Relationships
    parents ||--o{ child_parents : "has"
    child_profile ||--o{ enrollments_child : "enrolled in"    
    enrollments_shift ||--o{ enrollments_shifts : "has shift"
    addresses ||--o{ enrollments_addresses : "located at"
    enrollments ||--o{ enrollments_terms : "agrees to"
    enrollments ||--o{ enrollments_shifts : "has shift"
    enrollments ||--o{ student_modalities : "in"
    enrollments ||--o{ enrollments_addresses : "related to"
    enrollments ||--o{ child_parents : "related to"
    enrollments ||--o{ enrollments_child : "related to"
    modalities ||--o{ student_modalities : "offered in"
```

### Refs:
- https://github.com/go-jet/jet
- https://github.com/nofeaturesonlybugs/sqlh
- https://github.com/jackskj/carta
- https://github.com/harness/harness/blob/main/cmd/gitness/wire.go#L141
- https://github.com/grafana/grafana/blob/main/pkg/storage/unified/sql/db/dbimpl/db.go#L32
- https://eradman.com/posts/pg-admin-queries.html

