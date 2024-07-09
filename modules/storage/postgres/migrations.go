package postgres

import (
	"fmt"
)

// Migration - provide ability to create migration to storage if needed
type Migration struct {
	Sql  string
	Name string
}

var migrations = []*Migration{
	// migrations table
	&Migration{
		Name: "create_migrations",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.migrations
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			name text COLLATE pg_catalog."default" NOT NULL,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT groups_migrations_PK PRIMARY KEY (id)
		)`,
	},

	// users table
	&Migration{
		Name: "create_users",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.users
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			username character varying(120) COLLATE pg_catalog."default" NOT NULL,
			first_name character varying(120) COLLATE pg_catalog."default" NOT NULL,
			second_name character varying(120) COLLATE pg_catalog."default" NOT NULL,
			image character varying(120) COLLATE pg_catalog."default" NOT NULL,
			email character varying(120) COLLATE pg_catalog."default" NOT NULL,
			password character varying(120) COLLATE pg_catalog."default" NOT NULL,
			role uuid NOT NULL DEFAULT '9349e8e0-9f69-4a97-a47f-85d8d55a4776',
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL,
			CONSTRAINT users_PK PRIMARY KEY (id),
			CONSTRAINT users_username_unique UNIQUE (username)
		)`,
	},

	// create roles table
	&Migration{
		Name: "create_roles_table",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.roles
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			name VARCHAR(50) NOT NULL,
		    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL,
			CONSTRAINT roles_PK PRIMARY KEY (id),
			CONSTRAINT roles_name_unique UNIQUE (name)
		)
		`,
	},

	// create endpoints table
	&Migration{
		Name: "create_enpoints_table",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.endpoints
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			urn VARCHAR(300) NOT NULL,
		    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL,
			CONSTRAINT endpoints_PK PRIMARY KEY (id),
			CONSTRAINT urn_unique UNIQUE (urn)
		)
		`,
	},

	// Create permission reference between endpoint and role
	&Migration{
		Name: "create_roles_endpoints",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.roles_endpoints
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			roleUuid uuid NOT NULL,
			endpointUuid uuid NOT NULL,
			description TEXT NOT NULL DEFAULT '',
		    CONSTRAINT roles_endpoints_PK PRIMARY KEY (id),
			CONSTRAINT roles_endpoints_role_FK FOREIGN KEY (roleUuid)
				REFERENCES public.roles (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION,
		   	CONSTRAINT roles_endpoints_endpoint_FK FOREIGN KEY (endpointUuid)
				REFERENCES public.endpoints (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION,
		    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL
		)
		`,
	},

	// create default role
	&Migration{
		Name: "create_default_role",
		Sql: `
		INSERT INTO roles (id, name) VALUES ('9349e8e0-9f69-4a97-a47f-85d8d55a4776','default')
		`,
	},

	// create default role
	&Migration{
		Name: "create_admin_role",
		Sql: `
		INSERT INTO roles (id, name) VALUES ('5d169741-405e-4b37-a54a-6e8021e9661c','admin')
		`,
	},

	// insert admin user with admin group relation
	&Migration{
		Name: "insert_admin_user",
		Sql: `
		INSERT INTO users (username, first_name, second_name, image, email, password, role) VALUES ('admin', 'admin', 'admin', 'default_user.png', 'admin@admin.com', '', '5d169741-405e-4b37-a54a-6e8021e9661c');
		`,
	},

	// add constraint to user -> role FK
	&Migration{
		Name: "alter_users_role_field_fk",
		Sql: `
			ALTER TABLE users ADD CONSTRAINT user_role_FK FOREIGN KEY (role)
			REFERENCES public.roles (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION;
		`,
	},

	// chat table
	&Migration{
		Name: "create_chat_table",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.chats
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			name character varying(120) NOT NULL,
			chat_type character varying(120) NOT NULL,
			security_level character varying(120) NOT NULL DEFAULT 'SERVER_PRIVATE_KEY',
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL,
			CONSTRAINT chats_PK PRIMARY KEY (id)
		)`,
	},

	// chats users table
	&Migration{
		Name: "create_chats_users_table",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.chats_users
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			userId uuid NOT NULL,
			chatId uuid NOT NULL,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL,
			CONSTRAINT chats_users_PK PRIMARY KEY (id),
			CONSTRAINT chats_users_FK FOREIGN KEY (userID)
				REFERENCES public.users (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION,
			CONSTRAINT chats_chats_FK FOREIGN KEY (chatId)
				REFERENCES public.chats (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
		)`,
	},

	// messages table
	&Migration{
		Name: "create_messages_table",
		Sql: `
		CREATE TABLE IF NOT EXISTS public.messages
		(
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			chatid uuid NOT NULL,
			senderid uuid NOT NULL,
			message text NOT NULL,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone DEFAULT NULL,
			CONSTRAINT messages_PK PRIMARY KEY (id),
			CONSTRAINT messages_chat_FK FOREIGN KEY (chatid)
				REFERENCES public.chats (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION,
			CONSTRAINT messages_user_FK FOREIGN KEY (senderid)
				REFERENCES public.users (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
		)`,
	},
}

func (p *Postgres) migrations() error {
	// Check migrations table
	rows, err := p.db.Query("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'migrations')")
	if err != nil {
		return fmt.Errorf("cannot send request to check migrations tables: %v", err)
	}
	defer rows.Close()

	var exists string
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			return fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	migrated := map[string]string{}

	if exists == "true" {
		rows, err := p.db.Query("SELECT name, created_at FROM migrations ORDER BY created_at DESC")
		if err != nil {
			return fmt.Errorf("cannot send request to check migrations tables: %v", err)
		}
		defer rows.Close()

		var name string
		var created_at string
		for rows.Next() {
			err := rows.Scan(&name, &created_at)
			if err != nil {
				return fmt.Errorf("cannot read response from database: %v", err)
			}

			migrated[name] = created_at
		}
	}

	for _, migration := range migrations {

		if len(migrated[migration.Name]) != 0 {
			continue // skip creation if alredy exists
		}

		_, err := p.db.Exec(migration.Sql)
		if err != nil {
			return fmt.Errorf("cannot run migrations [%s]: %v", migration.Name, err)
		}

		query := "INSERT INTO migrations (name) VALUES ($1)"
		if _, err := p.db.Exec(query, migration.Name); err != nil {
			return fmt.Errorf("cannot create new migration record: %v", err)
		}

		p.Log("successful executed sql command", migration.Name)
	}

	return nil
}
