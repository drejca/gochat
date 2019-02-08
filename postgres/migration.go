package postgres

func Migrate(conn *Conn) {
	MessageDown000(conn)
	ChannelUsersDown000(conn)
	ChannelDown000(conn)
	UserDown000(conn)

	UserUp000(conn)
	ChannelUp000(conn)
	ChannelUsersUp000(conn)
	MessageUp000(conn)
}

func UserUp000(conn *Conn) {
	mustExec(conn, `
CREATE TABLE app_user (
	id serial primary key,
	username varchar(20) unique not null,
	password_hash bytea not null
);`)
}

func UserDown000(conn *Conn) {
	dropTableIfExists(conn, "app_user")
}

func ChannelUp000(conn *Conn) {
	mustExec(conn, `
CREATE TABLE channel (
	id serial primary key,
	name varchar(25) unique not null,
	owner_id integer not null
);`)
}

func ChannelDown000(conn *Conn) {
	dropTableIfExists(conn, "channel")
}

func ChannelUsersUp000(conn *Conn) {
	mustExec(conn, `
CREATE TABLE channel_users (
	channel_id int references channel,
	app_user_id int references app_user,
	primary key(channel_id, app_user_id)
);`)
}

func ChannelUsersDown000(conn *Conn) {
	dropTableIfExists(conn, "channel_users")
}

func MessageUp000(conn *Conn) {
	mustExec(conn, `
CREATE TABLE message (
	id serial primary key,
	owner_id int references app_user,
	channel_id int references channel,
	message text
);`)
}

func MessageDown000(conn *Conn) {
	dropTableIfExists(conn, "message")
}

func dropTableIfExists(conn *Conn, tableName string) {
	mustExec(conn, "DROP TABLE IF EXISTS "+tableName+";")
}

func mustExec(conn *Conn, sql string) {
	_, err := conn.Exec(sql)
	if err != nil {
		panic(err)
	}
}
