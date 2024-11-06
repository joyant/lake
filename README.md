# Lake: Dynamic SQL in Go, Inspired by MyBatis

Lake is a powerful Object-Relational Mapping (ORM) library for Go, designed to bring the flexibility and dynamism of MyBatis into the Go ecosystem. With Lake, you can define your SQL queries in XML, allowing for highly dynamic and flexible queries that can adapt to your application's needs on the fly.

## Features

- **Dynamic SQL**: Lake allows you to write dynamic SQL queries in XML, giving you the ability to adapt your queries to your application's needs in real time.

- **Go Integration**: Designed with Go in mind, Lake integrates seamlessly with your Go codebase, allowing you to leverage the power of dynamic SQL within your Go applications.

- **MyBatis Inspired**: Inspired by the flexibility and power of MyBatis, Lake brings similar capabilities to the Go ecosystem, providing a familiar interface for those coming from a Java/MyBatis background.

- **Flexible Query Definition**: With Lake, you're not limited to static SQL queries. You can easily define complex queries with conditional statements, loops, and more.

## Quick Start

Install Lake using `go get`:

```bash
go get github.com/joyant/lake
```

Import Lake in your Go code:

```go
import "github.com/joyant/lake"
```

Define your SQL queries in an XML file:

```xml
<select id="selectInID">
    select id, name, age from user
    where id in
    <foreach item="item" index="index" collection="ids" open="(" close=")" separator=",">
        #{item}
    </foreach>
</select>
```

Execute your query in Go:

```go
p := lake.Parameter{
    "ids": []int{1, 2, 3}
}
result, err := model.Session.SelectOne("selectInID", p)
```

## License

Lake is released under the [MIT License](https://github.com/joyant/lake/blob/main/LICENSE).
