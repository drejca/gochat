package postgres

import (
	"errors"
	"github.com/drejca/gochat"
	"github.com/jackc/pgx"
)

var EntryNotFound = errors.New("Db entry not found")

type Conn struct {
	*pgx.ConnPool
}

func New() (*Conn, error) {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		return nil, err
	}

	config.Host = "localhost"
	config.User = "postgres"
	config.Password = "postgres"
	config.Database = "chatapp"

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: config})
	if err != nil {
		return nil, err
	}
	return &Conn{ConnPool: pool}, nil
}

type UserRepository struct {
	conn *Conn
}

func NewUserRepository(conn *Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u *UserRepository) Store(username string, passwordHash []byte) (gochat.User, error) {
	sql := `INSERT INTO app_user(username, password_hash) VALUES($1, $2) 
		RETURNING id, username`
	row := u.conn.QueryRow(sql, username, passwordHash)
	return scanUser(row)
}

func (u *UserRepository) Find(id int) (gochat.User, error) {
	sql := `SELECT id, username FROM app_user WHERE id = $1`
	row := u.conn.QueryRow(sql, id)

	return scanUser(row)
}

func (u *UserRepository) GetUserPasswordHash(username string) (passwordHash []byte, err error) {
	sql := `SELECT password_hash FROM app_user WHERE username = $1;`
	row := u.conn.QueryRow(sql, username)
	if row == nil {
		return passwordHash, EntryNotFound
	}

	err = row.Scan(&passwordHash)
	return passwordHash, err
}

func scanUser(row *pgx.Row) (gochat.User, error) {
	if row == nil {
		return gochat.User{}, EntryNotFound
	}

	user := gochat.User{}

	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

type ChannelRepository struct {
	conn *Conn
}

func NewChannelRepository(conn *Conn) *ChannelRepository {
	return &ChannelRepository{conn: conn}
}

func (c *ChannelRepository) Store(channelName string, owner gochat.User) (gochat.Channel, error) {
	sql := `INSERT INTO channel(name, owner_id) VALUES($1, $2) 
		RETURNING id, name, owner_id`
	row := c.conn.QueryRow(sql, channelName, owner.Id)

	return scanChannel(row)
}

func (c *ChannelRepository) Find(id int) (gochat.Channel, error) {
	sql := `SELECT id, name, owner_id FROM channel WHERE id = $1`
	row := c.conn.QueryRow(sql, id)

	return scanChannel(row)
}

func scanChannel(row *pgx.Row) (gochat.Channel, error) {
	if row == nil {
		return gochat.Channel{}, EntryNotFound
	}

	channel := gochat.Channel{}
	err := row.Scan(&channel.Id, &channel.Name, &channel.OwnerId)
	if err != nil {
		return channel, err
	}
	return channel, nil
}

type ChannelUsersRepository struct {
	conn *Conn
}

func NewChannelUsersRepository(conn *Conn) *ChannelUsersRepository {
	return &ChannelUsersRepository{conn: conn}
}

func (c *ChannelUsersRepository) Store(channel gochat.Channel, user gochat.User) error {
	sql := `INSERT INTO channel_users(channel_id, app_user_id) VALUES($1, $2)`
	_, err := c.conn.Exec(sql, channel.Id, user.Id)

	return err
}

func (c *ChannelUsersRepository) GetChannelUsers(channelId int) ([]gochat.User, error) {
	var users []gochat.User

	sql := `SELECT id, username FROM channel_users 
		INNER JOIN app_user ON (channel_users.app_user_id = app_user.id) 
		WHERE channel_id = $1`
	rows, err := c.conn.Query(sql, channelId)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := gochat.User{}
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

type MessageRepository struct {
	conn *Conn
}

func NewMessageRepository(conn *Conn) *MessageRepository {
	return &MessageRepository{conn: conn}
}

func (m *MessageRepository) Store(ownerId int, channelId int, message string) error {
	sql := `INSERT INTO message(owner_id, channel_id, message) VALUES($1, $2, $3)`
	_, err := m.conn.Exec(sql, ownerId, channelId, message)
	return err
}