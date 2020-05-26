package main

import (
    "errors"
    "github.com/jmoiron/sqlx"
    "time"
)

type Session interface{
    Select(key string, argv Parameter) (results []Result, err error)
    SelectOne(key string, argv Parameter) (results Result, err error)
    Insert(key string, argv Parameter) (act Affection, err error)
    Update(key string, argv Parameter) (rowsAffected int64, err error)
    Delete(key string, argv Parameter) (rowsAffected int64, err error)
    Close() error
}

type Result map[string]interface{}

type Affection struct {
    Rows int64
    LastInsertID int64
}

type mySQLSession struct {
    driveName string
    builder builder
    db *sqlx.DB
    option *Option
}

func NewSession(driveName, dataSourceName string, option *Option) (Session, error) {
    if option == nil {
        option = &Option{}
    }
    if option.Mode == "" {
        option.Mode = ModeDebug
    }
    if option.ConnMaxLifetime <= 0 {
        option.ConnMaxLifetime = time.Second * 3600
    }
    if option.MaxIdleConns <= 0 {
        option.MaxIdleConns = 5
    }
    if option.MaxOpenConns <= 0 {
        option.MaxOpenConns = 5
    }
    if driveName == MySQLDrive {
        db, err := sqlx.Connect(driveName, dataSourceName)
        if err != nil {
            return nil, err
        }
        db.SetConnMaxLifetime(option.ConnMaxLifetime)
        db.SetMaxIdleConns(option.MaxIdleConns)
        db.SetMaxOpenConns(option.MaxOpenConns)
        s := &mySQLSession{
            driveName: driveName,
            builder:   &mySQLBuilder{},
            db:        db,
            option:    option,
        }
        return s, nil
    }
    return nil, errors.New("don't support drive " + driveName)
}

func (s *mySQLSession)Select(key string, argv Parameter) (results []Result, err error) {
    var sta stack
    sta, err = ms.find(key)
    if err != nil {
        return
    }
    sta.params = argv
    var stmt []byte
    stmt, err = sta.Parse()
    if err != nil {
        return
    }
    sql, params, err := s.builder.build(stmt, argv)
    if err != nil {
        return results, err
    }
    if s.option.Mode == ModeDebug {
        logger.Infof("debug mode, stmt:%s\nparams:%v\n", s.builder.lastSQL(sql, params), params)
    }
    rows, err := s.db.Queryx(sql, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        m := make(Result)
        err = rows.MapScan(m)
        if err != nil {
            return nil, err
        }
        for k, v := range m {
            switch v.(type) {
            case []uint8:
                m[k] = string(v.([]uint8))
            }
        }
        results = append(results, m)
    }
    return
}

func (s *mySQLSession)SelectOne(key string, argv Parameter) (result Result, err error) {
    var sta stack
    sta, err = ms.find(key)
    if err != nil {
        return
    }
    sta.params = argv
    var stmt []byte
    stmt, err = sta.Parse()
    if err != nil {
        return
    }
    sql, params, err := s.builder.build(stmt, argv)
    if err != nil {
        return
    }
    if s.option.Mode == ModeDebug {
        logger.Infof("debug mode, stmt:%s\nparams:%v\n", s.builder.lastSQL(sql, params), params)
    }
    rows := s.db.QueryRowx(sql, params...)
    result = make(Result)
    err = rows.MapScan(result)
    if err != nil {
        return nil, err
    }
    for k, v := range result {
        switch v.(type) {
        case []uint8:
            result[k] = string(v.([]uint8))
        }
    }
    return
}

func (s *mySQLSession)Insert(key string, argv Parameter) (act Affection, err error) {
    var sta stack
    sta, err = ms.find(key)
    if err != nil {
        return
    }
    sta.params = argv
    var stmt []byte
    stmt, err = sta.Parse()
    if err != nil {
        return
    }
    sql, params, err := s.builder.build(stmt, argv)
    if err != nil {
        return
    }
    if s.option.Mode == ModeDebug {
        logger.Infof("debug mode, stmt:%s\nparams:%v\n", s.builder.lastSQL(sql, params), params)
    }
    result, err := s.db.Exec(sql, params...)
    if err != nil {
        return
    }
    act.Rows, err = result.RowsAffected()
    if err != nil {
        return
    }
    act.LastInsertID, err = result.LastInsertId()
    return
}

func (s *mySQLSession)Update(key string, argv Parameter) (rowsAffected int64, err error) {
    var sta stack
    sta, err = ms.find(key)
    if err != nil {
        return
    }
    sta.params = argv
    var stmt []byte
    stmt, err = sta.Parse()
    if err != nil {
        return
    }
    sql, params, err := s.builder.build(stmt, argv)
    if err != nil {
        return
    }
    if s.option.Mode == ModeDebug {
        logger.Infof("debug mode, stmt:%s\nparams:%v\n", s.builder.lastSQL(sql, params), params)
    }
    result, err := s.db.Exec(sql, params...)
    if err != nil {
        return
    }
    rowsAffected, err = result.RowsAffected()
    return
}

func (s *mySQLSession)Delete(key string, argv Parameter) (rowsAffected int64, err error) {
    var sta stack
    sta, err = ms.find(key)
    if err != nil {
        return
    }
    sta.params = argv
    var stmt []byte
    stmt, err = sta.Parse()
    if err != nil {
        return
    }
    sql, params, err := s.builder.build(stmt, argv)
    if err != nil {
        return
    }
    if s.option.Mode == ModeDebug {
        logger.Infof("debug mode, stmt:%s\nparams:%v\n", s.builder.lastSQL(sql, params), params)
    }
    result, err := s.db.Exec(sql, params...)
    if err != nil {
        return
    }
    rowsAffected, err = result.RowsAffected()
    return
}

func (s *mySQLSession)Close() error {
    return s.db.Close()
}


