package lake

import (
    rawSQL "database/sql"
    "errors"
    "github.com/jmoiron/sqlx"
    "strings"
    "time"
)

type DB interface{
    InsertMap(table string, params Parameter) (rowsAffected, lastInsertID int64, err error)
}

type exportDB struct {
    option *Option
    db *sqlx.DB
}

func (e *exportDB)InsertMap(table string, params Parameter) (rowsAffected, lastInsertID int64, err error)  {
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    builder := strings.Builder{}
    builder.WriteString("insert into ")
    builder.WriteString(table)
    builder.WriteByte('(')
    for k, v := range keys {
        builder.WriteString(v)
        if k != len(keys)-1 {
            builder.WriteByte(',')
        }
    }
    builder.WriteString(") values (")
    for k, v := range keys {
        builder.WriteByte(':')
        builder.WriteString(v)
        if k != len(keys)-1 {
            builder.WriteByte(',')
        }
    }
    builder.WriteByte(')')
    SQL := builder.String()
    if e.option.Mode == ModeDebug {
        logger.Infof("debug mode statement:%s", SQL)
        logger.Infof("params:%v", params)
    }
    var result rawSQL.Result
    result, err = e.db.NamedExec(SQL, map[string]interface{}(params))
    if err != nil {
        return
    }
    rowsAffected, err = result.RowsAffected()
    if err != nil {
        return
    }
    lastInsertID, err = result.LastInsertId()
    return
}

type Session interface{
    Select(key string, argv Parameter) (results []Result, err error)
    SelectOne(key string, argv Parameter) (result Result, err error)
    Insert(key string, argv Parameter) (act Affection, err error)
    Update(key string, argv Parameter) (rowsAffected int64, err error)
    Delete(key string, argv Parameter) (rowsAffected int64, err error)
    DB() DB
    Close() error
}

type mySQLSession struct {
    driveName string
    builder builder
    db *sqlx.DB
    exportDB DB
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
            exportDB: &exportDB{
                option: option,
                db:     db,
            },
        }
        return s, nil
    }
    return nil, errors.New("don't support drive " + driveName)
}

func (s *mySQLSession)DB() DB {
    return s.exportDB
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
    if err == rawSQL.ErrNoRows {
        err = nil
    }
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
    if err == rawSQL.ErrNoRows {
        err = nil
    }
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


