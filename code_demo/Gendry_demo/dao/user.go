package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/mitchellh/mapstructure"

	"asong.cloud/Golang_Dream/code_demo/Gendry_demo/model"
)

const (
	tplTable = "users"
)

type UserDB struct {
	cli *sql.DB
}

func NewUserDB(cli *sql.DB) *UserDB {
	return &UserDB{
		cli: cli,
	}
}

func (db *UserDB) getFiledList() []string {
	list := []string{
		"id", "username", "nickname", "password",
		"salt", "avatar", "uptime",
	}
	return list
}

func (db *UserDB) buildWhere(condition *model.User) (map[string]interface{}, error) {
	res := make(map[string]interface{}, 0)
	err := mapstructure.Decode(condition, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 方式一：把条件放进来，利用mapstructure解码构造，只能构造等于条件，其他符合不支持
func (db *UserDB) GetMethodOne(ctx context.Context, condition *model.User) (*model.User, error) {
	cond, err := db.buildWhere(condition)
	if err != nil {
		return nil, err
	}
	builder.OmitEmpty(cond, db.getFiledList())
	sqlStr, values, err := builder.BuildSelect(tplTable, cond, db.getFiledList())
	if err != nil {
		return nil, err
	}
	fmt.Println(sqlStr, values)
	rows,err := db.cli.QueryContext(ctx, sqlStr, values...)
	defer func() {
		if rows != nil{
			_ = rows.Close()
		}
	}()
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,errors.New("not found")
		}
		return nil,err
	}
	user := model.NewEmptyUser()
	err = scanner.Scan(rows,&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 方式二： 自己构造好条件传进来
func (db *UserDB) GetMethodTwo(ctx context.Context,cond map[string]interface{}) (*model.User,error) {
	sqlStr,values, err := builder.BuildSelect(tplTable,cond,db.getFiledList())
	if err != nil{
		return nil, err
	}
	rows,err := db.cli.QueryContext(ctx,sqlStr,values...)
	defer func() {
		if rows != nil{
			_ = rows.Close()
		}
	}()
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,errors.New("not found")
		}
		return nil,err
	}
	user := model.NewEmptyUser()
	err = scanner.Scan(rows,&user)
	if err != nil{
		return nil,err
	}
	return user,nil
}

func (db *UserDB) Query(ctx context.Context,cond map[string]interface{}) ([]*model.User,error) {
	sqlStr,values,err := builder.BuildSelect(tplTable,cond,db.getFiledList())
	if err != nil{
		return nil, err
	}
	rows,err := db.cli.QueryContext(ctx,sqlStr,values...)
	defer func() {
		if rows != nil{
			_ = rows.Close()
		}
	}()
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,errors.New("not found")
		}
		return nil,err
	}
	user := make([]*model.User,0)
	err = scanner.Scan(rows,&user)
	if err != nil{
		return nil,err
	}
	return user,nil
}

func (db *UserDB) Add(ctx context.Context,cond map[string]interface{}) (int64,error) {
	sqlStr,values,err := builder.BuildInsert(tplTable,[]map[string]interface{}{cond})
	if err != nil{
		return 0,err
	}
	// TODO:DEBUG
	fmt.Println(sqlStr,values)
	res,err := db.cli.ExecContext(ctx,sqlStr,values...)
	if err != nil{
		return 0,err
	}
	return res.LastInsertId()
}

func (db *UserDB) Update(ctx context.Context,where map[string]interface{},data map[string]interface{}) error {
	sqlStr,values,err := builder.BuildUpdate(tplTable,where,data)
	if err != nil{
		return err
	}
	// TODO:DEBUG
	fmt.Println(sqlStr,values)
	res,err := db.cli.ExecContext(ctx,sqlStr,values...)
	if err != nil{
		return err
	}
	affectedRows,err := res.RowsAffected()
	if err != nil{
		return err
	}
	if affectedRows == 0{
		return errors.New("no record update")
	}
	return nil
}

func (db *UserDB)Delete(ctx context.Context,where map[string]interface{}) error {
	sqlStr,values,err := builder.BuildDelete(tplTable,where)
	if err != nil{
		return err
	}
	// TODO:DEBUG
	fmt.Println(sqlStr,values)
	res,err := db.cli.ExecContext(ctx,sqlStr,values...)
	if err != nil{
		return err
	}
	affectedRows,err := res.RowsAffected()
	if err != nil{
		return err
	}
	if affectedRows == 0{
		return errors.New("no record delete")
	}
	return nil
}

func (db *UserDB) CustomizeGet(ctx context.Context,sql string,data map[string]interface{}) (*model.User,error) {
	sqlStr,values,err := builder.NamedQuery(sql,data)
	if err != nil{
		return nil, err
	}
	// TODO:DEBUG
	fmt.Println(sql,values)
	rows,err := db.cli.QueryContext(ctx,sqlStr,values...)
	if err != nil{
		return nil,err
	}
	defer func() {
		if rows != nil{
			_ = rows.Close()
		}
	}()
	user := model.NewEmptyUser()
	err = scanner.Scan(rows,&user)
	if err != nil{
		return nil,err
	}
	return user,nil
}

func (db *UserDB) AggregateCount(ctx context.Context,where map[string]interface{},filed string) (int64,error) {
	res,err := builder.AggregateQuery(ctx,db.cli,tplTable,where,builder.AggregateCount(filed))
	if err != nil{
		return 0, err
	}
	numberOfRecords := res.Int64()
	return numberOfRecords,nil
}
