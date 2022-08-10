package types

import (
	"database/sql"
	"fmt"
	"time"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Access string

const (
	Owner   Access = "Owner"
	Manager Access = "Manager"
	Super   Access = "Super"
	Member  Access = "Member"
	Guest   Access = "Guest"
)

var MinManager []Access = []Access{Owner, Manager}
var MinSuper []Access = []Access{Owner, Manager, Super}
var MinMember []Access = []Access{Owner, Manager, Super, Member}
var MinGuest []Access = []Access{Owner, Manager, Super, Member, Guest}

type User struct {
	Id        string `json:"id"`
	Submitter string `json:"submitter"`
	IsGroup   bool   `json:"isgroup"`
	Coins     int    `json:"coins"`
}

type Hit struct {
	Id      string `json:"id"`
	Appname string `json:"appname"`
	Command string `json:"command"`
	Hits    int    `json:"Hits"`
}

type Permission struct {
	Id      string    `json:"id"`
	Access  Access    `json:"access"`
	Issued  time.Time `json:"issued"`
	Expires time.Time `json:"expires"`
}

type baseAuth struct {
	database *sql.DB
}

type IAuth interface {
	UserRemove(id string) (bool, error)
	UserAdd(user User, perm Permission) (bool, error)
	UserExists(id string) (User, error)
	UserList() ([]User, error)

	PermissionRemove(id string, access Access) (bool, error)
	PermissionAdd(perm Permission) (bool, error)

	ValidManyPermission(id string, access []Access) ([]Permission, bool)
	ValidPermission(id string, access Access) (Permission, bool)

	HitCommand(id, appname, command string) (Hit, error)
}

func NewAuthenticator(db *sql.DB) IAuth {
	baseauth := baseAuth{}

	baseauth.database = db
	defer db.Close()

	baseauth.init()
	return &baseauth

}

func (a *baseAuth) init() {
	sql_tables := []string{
		`CREATE TABLE IF NOT EXISTS public.user(
			created_at timestamp default current_timestamp,
			id varchar(255) unique not null primary key,
			submitter varchar(255) not null,
			isgroup bool default false,
			coins int default 0 
			)`,
		`CREATE TABLE IF NOT EXISTS public.permission(
			created_at timestamp default current_timestamp,
			id varchar(255) not null,
			access varchar(255) not null,
			issued timestamp default current_timestamp,
			expires timestamp 
			)`,
		`CREATE TABLE IF NOT EXISTS public.hit(
			created_at timestamp default current_timestamp,
			id varchar(255) not null,
			appname varchar(255) not null,
			command varchar(255) not null,
			hits int default 1 
			)`,
	}

	if ctx, err := a.database.Begin(); err == nil {

		for _, sql_string := range sql_tables {
			_, err := ctx.Exec(sql_string)
			if err != nil {
				panic(err)
			}
		}
		ctx.Commit()
	}
}

func (a *baseAuth) PermissionRemove(id string, access Access) (bool, error) {
	for _, table := range []string{"permission"} {
		if ctx, err := a.database.Begin(); err != nil {
			return false, err
		} else {
			sql_string := fmt.Sprintf(`delete from public.%s where id = $1 and access = $2`, table)
			if _, err := ctx.Exec(sql_string, id, access); err != nil {
				return false, err
			} else {
				return true, err
			}
		}
	}
	return false, nil
}

func (a *baseAuth) PermissionAdd(perm Permission) (bool, error) {
	if ctx, err := a.database.Begin(); err != nil {
		return false, err
	} else {
		sql_string := `insert into public.permission(id, access, expires) values($1, $2, $3)`
		if prepare, err := ctx.Prepare(sql_string); err != nil {
			return false, err
		} else {
			if _, err := prepare.Exec(perm.Id, perm.Access, perm.Expires); err != nil {
				return false, err
			} else {
				if err := ctx.Commit(); err != nil {
					return false, err
				} else {
					return true, err
				}
			}
		}
	}
}

func (a *baseAuth) UserExists(id string) (User, error) {
	var user User
	if ctx, err := a.database.Begin(); err != nil {
		return user, err
	} else {
		sql_string := `select id, submitter, isgroup, coins from public.user where id = "%s"`
		if err := ctx.QueryRow(fmt.Sprintf(sql_string, id)).Scan(&user.Id, &user.Submitter, &user.IsGroup, &user.Coins); err != nil {
			return User{}, err
		} else {
			return user, err
		}
	}
}

func (a *baseAuth) UserRemove(id string) (bool, error) {

	for _, table := range []string{"user", "permission"} {
		if ctx, err := a.database.Begin(); err != nil {
			return false, err
		} else {
			sql_string := fmt.Sprintf(`delete from public.%s where id = $1 ;`, table)
			if _, err := ctx.Exec(sql_string, id); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (a *baseAuth) UserAdd(user User, perm Permission) (bool, error) {

	if ctx, err := a.database.Begin(); err != nil {
		return false, err
	} else {

		if prepare, err := ctx.Prepare(`insert into public.user(id, submitter, isgroup, coins) values($1, $2, $3, $4)`); err != nil {
			return false, err
		} else {

			if _, err := prepare.Exec(user.Id, user.Submitter, user.IsGroup, user.Coins); err != nil {
				return false, err
			} else {
				if ok, err := a.PermissionAdd(perm); err != nil {
					return ok, err
				}
			}

			if err := ctx.Commit(); err != nil {
				return false, err
			} else {
				return true, nil
			}

		}
	}

}

func (a *baseAuth) UserList() ([]User, error) {
	var users []User
	if ctx, err := a.database.Begin(); err != nil {
		return users, err
	} else {
		if rows, err := ctx.Query(`select id, submitter, isgroup, coins from public.user`); err != nil {
			return users, err
		} else {
			for rows.Next() {
				var user User
				if err := rows.Scan(&user.Id, &user.Submitter, &user.IsGroup, &user.Coins); err != nil {
					return users, err
				} else {
					users = append(users, user)
				}
			}
		}
	}

	return users, nil
}

func (a *baseAuth) ValidManyPermission(id string, access []Access) ([]Permission, bool) {
	var newperm []Permission
	var isvalid bool

	for _, acc := range access {

		perm, validity := a.ValidPermission(id, acc)
		newperm = append(newperm, perm)
		isvalid = isvalid || validity
	}

	return newperm, isvalid
}

func (a *baseAuth) ValidPermission(id string, access Access) (Permission, bool) {
	var newperm Permission
	var isvalid bool

	sql_string := `select id, access, issued, expires from public.permission where id = $1 and access = $2 limit 1;`

	if err := a.database.QueryRow(sql_string, id, access).Scan(&newperm.Id, &newperm.Access, &newperm.Issued, &newperm.Expires); err != nil {
		// LogCode(err)
		isvalid = false
	} else {
		isvalid = time.Now().Before(newperm.Expires) || AccessContain(access, MinManager)
	}

	return newperm, isvalid
}

func (a *baseAuth) HitCommand(id, appname, command string) (Hit, error) {
	var hit Hit = Hit{Id: id, Appname: appname, Command: command, Hits: 1}
	sql_string := `select id, appname, command, hits from public.hit where id = "%s" and appname = "%s" and command = "%s"`
	if ctx, err := a.database.Begin(); err != nil {
		return hit, err
	} else {
		if err := ctx.QueryRow(fmt.Sprintf(sql_string, id, appname, command)).Scan(&hit.Id, &hit.Appname, &hit.Command, &hit.Hits); err != nil {
			ctx.Commit()

			if ctx, err := a.database.Begin(); err != nil {
				return hit, err
			} else {

				sql_string = `insert into public.hit(id, appname, command) values($1, $2, $3)`
				if prepare, err := ctx.Prepare(sql_string); err != nil {
					return hit, err

				} else {

					if _, err := prepare.Exec(hit.Id, hit.Appname, hit.Command); err != nil {
						return hit, err
					}
				}
			}

		} else {

			hit.Hits += 1
			sql_string := `update public.hit set hits = %d where id = "%s" and appname = "%s" and command = "%s"`
			if _, err := ctx.Exec(fmt.Sprintf(sql_string, hit.Hits, hit.Id, hit.Appname, hit.Command)); err != nil {
				return hit, err
			}
			if err := ctx.Commit(); err != nil {
				return hit, err

			}
		}

	}

	return hit, nil
}

func AccessContain(access Access, access_list []Access) bool {
	var isvalid bool

	for _, acc := range access_list {
		isvalid = isvalid || acc == access
	}

	return isvalid
}
