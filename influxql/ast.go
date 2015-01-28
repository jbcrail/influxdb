package influxql

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DataType represents the primitive data types available in InfluxQL.
type DataType string

const (
	Unknown  = DataType("")
	Number   = DataType("number")
	Boolean  = DataType("boolean")
	String   = DataType("string")
	Time     = DataType("time")
	Duration = DataType("duration")
)

// InspectDataType returns the data type of a given value.
func InspectDataType(v interface{}) DataType {
	switch v.(type) {
	case float64:
		return Number
	case bool:
		return Boolean
	case string:
		return String
	case time.Time:
		return Time
	case time.Duration:
		return Duration
	default:
		return Unknown
	}
}

// Node represents a node in the InfluxDB abstract syntax tree.
type Node interface {
	node()
	String() string
}

func (_ *Query) node()     {}
func (_ Statements) node() {}

func (_ *AlterRetentionPolicyStatement) node()  {}
func (_ *CreateContinuousQueryStatement) node() {}
func (_ *CreateDatabaseStatement) node()        {}
func (_ *CreateRetentionPolicyStatement) node() {}
func (_ *CreateUserStatement) node()            {}
func (_ *DeleteStatement) node()                {}
func (_ *DropContinuousQueryStatement) node()   {}
func (_ *DropDatabaseStatement) node()          {}
func (_ *DropRetentionPolicyStatement) node()   {}
func (_ *DropSeriesStatement) node()            {}
func (_ *DropUserStatement) node()              {}
func (_ *GrantStatement) node()                 {}
func (_ *ShowContinuousQueriesStatement) node() {}
func (_ *ShowDatabasesStatement) node()         {}
func (_ *ShowFieldKeysStatement) node()         {}
func (_ *ShowFieldValuesStatement) node()       {}
func (_ *ShowRetentionPoliciesStatement) node() {}
func (_ *ShowMeasurementsStatement) node()      {}
func (_ *ShowSeriesStatement) node()            {}
func (_ *ShowTagKeysStatement) node()           {}
func (_ *ShowTagValuesStatement) node()         {}
func (_ *ShowUsersStatement) node()             {}
func (_ *RevokeStatement) node()                {}
func (_ *SelectStatement) node()                {}

func (_ *BinaryExpr) node()      {}
func (_ *BooleanLiteral) node()  {}
func (_ *Call) node()            {}
func (_ *Dimension) node()       {}
func (_ Dimensions) node()       {}
func (_ *DurationLiteral) node() {}
func (_ *Field) node()           {}
func (_ Fields) node()           {}
func (_ *Join) node()            {}
func (_ *Measurement) node()     {}
func (_ Measurements) node()     {}
func (_ *nilLiteral) node()      {}
func (_ *Merge) node()           {}
func (_ *NumberLiteral) node()   {}
func (_ *ParenExpr) node()       {}
func (_ *SortField) node()       {}
func (_ SortFields) node()       {}
func (_ *StringLiteral) node()   {}
func (_ *TagKeyIdent) node()     {}
func (_ *Target) node()          {}
func (_ *TimeLiteral) node()     {}
func (_ *VarRef) node()          {}
func (_ *Wildcard) node()        {}

// Query represents a collection of ordered statements.
type Query struct {
	Statements Statements
}

// String returns a string representation of the query.
func (q *Query) String() string { return q.Statements.String() }

// Statements represents a list of statements.
type Statements []Statement

// String returns a string representation of the statements.
func (a Statements) String() string {
	var str []string
	for _, stmt := range a {
		str = append(str, stmt.String())
	}
	return strings.Join(str, ";\n")
}

// Statement represents a single command in InfluxQL.
type Statement interface {
	Node
	stmt()
	RequiredPrivileges() ExecutionPrivileges
}

// ExecutionPrivilege is a privilege required for a user to execute
// a statement on a database or resource.
type ExecutionPrivilege struct {
	// Name of the database or resource.
	// If "", then the resource is the cluster.
	Name string

	// Privilege required.
	Privilege Privilege
}

// ExecutionPrivileges is a list of privileges required to execute a statement.
type ExecutionPrivileges []ExecutionPrivilege

func (_ *AlterRetentionPolicyStatement) stmt()  {}
func (_ *CreateContinuousQueryStatement) stmt() {}
func (_ *CreateDatabaseStatement) stmt()        {}
func (_ *CreateRetentionPolicyStatement) stmt() {}
func (_ *CreateUserStatement) stmt()            {}
func (_ *DeleteStatement) stmt()                {}
func (_ *DropContinuousQueryStatement) stmt()   {}
func (_ *DropDatabaseStatement) stmt()          {}
func (_ *DropRetentionPolicyStatement) stmt()   {}
func (_ *DropSeriesStatement) stmt()            {}
func (_ *DropUserStatement) stmt()              {}
func (_ *GrantStatement) stmt()                 {}
func (_ *ShowContinuousQueriesStatement) stmt() {}
func (_ *ShowDatabasesStatement) stmt()         {}
func (_ *ShowFieldKeysStatement) stmt()         {}
func (_ *ShowFieldValuesStatement) stmt()       {}
func (_ *ShowMeasurementsStatement) stmt()      {}
func (_ *ShowRetentionPoliciesStatement) stmt() {}
func (_ *ShowSeriesStatement) stmt()            {}
func (_ *ShowTagKeysStatement) stmt()           {}
func (_ *ShowTagValuesStatement) stmt()         {}
func (_ *ShowUsersStatement) stmt()             {}
func (_ *RevokeStatement) stmt()                {}
func (_ *SelectStatement) stmt()                {}

// Expr represents an expression that can be evaluated to a value.
type Expr interface {
	Node
	expr()
}

func (_ *BinaryExpr) expr()      {}
func (_ *BooleanLiteral) expr()  {}
func (_ *Call) expr()            {}
func (_ *DurationLiteral) expr() {}
func (_ *nilLiteral) expr()      {}
func (_ *NumberLiteral) expr()   {}
func (_ *ParenExpr) expr()       {}
func (_ *StringLiteral) expr()   {}
func (_ *TagKeyIdent) expr()     {}
func (_ *TimeLiteral) expr()     {}
func (_ *VarRef) expr()          {}
func (_ *Wildcard) expr()        {}

// Source represents a source of data for a statement.
type Source interface {
	Node
	source()
}

func (_ *Join) source()        {}
func (_ *Measurement) source() {}
func (_ *Merge) source()       {}

// SortField represents a field to sort results by.
type SortField struct {
	// Name of the field
	Name string

	// Sort order.
	Ascending bool
}

// String returns a string representation of a sort field
func (field *SortField) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString(field.Name)
	_, _ = buf.WriteString(" ")
	_, _ = buf.WriteString(strconv.FormatBool(field.Ascending))
	return buf.String()
}

// SortFields represents an ordered list of ORDER BY fields
type SortFields []*SortField

// String returns a string representation of sort fields
func (a SortFields) String() string {
	fields := make([]string, 0, len(a))
	for _, field := range a {
		fields = append(fields, field.String())
	}
	return strings.Join(fields, ", ")
}

// CreateDatabaseStatement represents a command for creating a new database.
type CreateDatabaseStatement struct {
	// Name of the database to be created.
	Name string
}

// String returns a string representation of the create database statement.
func (s *CreateDatabaseStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("CREATE DATABASE ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a CreateDatabaseStatement.
func (s *CreateDatabaseStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// DropDatabaseStatement represents a command to drop a database.
type DropDatabaseStatement struct {
	// Name of the database to be dropped.
	Name string
}

// String returns a string representation of the drop database statement.
func (s *DropDatabaseStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DROP DATABASE ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a DropDatabaseStatement.
func (s *DropDatabaseStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// DropRetentionPolicyStatement represents a command to drop a retention policy from a database.
type DropRetentionPolicyStatement struct {
	// Name of the policy to drop.
	Name string

	// Name of the database to drop the policy from.
	Database string
}

// String returns a string representation of the drop retention policy statement.
func (s *DropRetentionPolicyStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DROP RETENTION POLICY ")
	_, _ = buf.WriteString(s.Name)
	_, _ = buf.WriteString(" ON ")
	_, _ = buf.WriteString(s.Database)
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a DropRetentionPolicyStatement.
func (s *DropRetentionPolicyStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: s.Database, Privilege: WritePrivilege}}
}

// CreateUserStatement represents a command for creating a new user.
type CreateUserStatement struct {
	// Name of the user to be created.
	Name string

	// User's password
	Password string

	// User's privilege level.
	Privilege *Privilege
}

// String returns a string representation of the create user statement.
func (s *CreateUserStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("CREATE USER ")
	_, _ = buf.WriteString(s.Name)
	_, _ = buf.WriteString(" WITH PASSWORD ")
	_, _ = buf.WriteString(s.Password)

	if s.Privilege != nil {
		_, _ = buf.WriteString(" WITH ")
		_, _ = buf.WriteString(s.Privilege.String())
	}

	return buf.String()
}

// RequiredPrivilege returns the privilege(s) required to execute a CreateUserStatement.
func (s *CreateUserStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// DropUserStatement represents a command for dropping a user.
type DropUserStatement struct {
	// Name of the user to drop.
	Name string
}

// String returns a string representation of the drop user statement.
func (s *DropUserStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DROP USER ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

// RequiredPrivilege returns the privilege(s) required to execute a DropUserStatement.
func (s *DropUserStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// Privilege is a type of action a user can be granted the right to use.
type Privilege int

const (
	NoPrivileges Privilege = iota
	ReadPrivilege
	WritePrivilege
	AllPrivileges
)

// NewPrivilege returns an initialized *Privilege.
func NewPrivilege(p Privilege) *Privilege { return &p }

// String returns a string representation of a Privilege.
func (p Privilege) String() string {
	switch p {
	case NoPrivileges:
		return "NO PRIVILEGES"
	case ReadPrivilege:
		return "READ"
	case WritePrivilege:
		return "WRITE"
	case AllPrivileges:
		return "ALL PRIVILEGES"
	}
	return ""
}

// GrantStatement represents a command for granting a privilege.
type GrantStatement struct {
	// The privilege to be granted.
	Privilege Privilege

	// Thing to grant privilege on (e.g., a DB).
	On string

	// Who to grant the privilege to.
	User string
}

// String returns a string representation of the grant statement.
func (s *GrantStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("GRANT ")
	_, _ = buf.WriteString(s.Privilege.String())
	if s.On != "" {
		_, _ = buf.WriteString(" ON ")
		_, _ = buf.WriteString(s.On)
	}
	_, _ = buf.WriteString(" TO ")
	_, _ = buf.WriteString(s.User)
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a GrantStatement.
func (s *GrantStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// RevokeStatement represents a command to revoke a privilege from a user.
type RevokeStatement struct {
	// Privilege to be revoked.
	Privilege Privilege

	// Thing to revoke privilege to (e.g., a DB)
	On string

	// Who to revoke privilege from.
	User string
}

// String returns a string representation of the revoke statement.
func (s *RevokeStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("REVOKE ")
	_, _ = buf.WriteString(s.Privilege.String())
	if s.On != "" {
		_, _ = buf.WriteString(" ON ")
		_, _ = buf.WriteString(s.On)
	}
	_, _ = buf.WriteString(" FROM ")
	_, _ = buf.WriteString(s.User)
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a RevokeStatement.
func (s *RevokeStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// CreateRetentionPolicyStatement represents a command to create a retention policy.
type CreateRetentionPolicyStatement struct {
	// Name of policy to create.
	Name string

	// Name of database this policy belongs to.
	Database string

	// Duration data written to this policy will be retained.
	Duration time.Duration

	// Replication factor for data written to this policy.
	Replication int

	// Should this policy be set as default for the database?
	Default bool
}

// String returns a string representation of the create retention policy.
func (s *CreateRetentionPolicyStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("CREATE RETENTION POLICY ")
	_, _ = buf.WriteString(s.Name)
	_, _ = buf.WriteString(" ON ")
	_, _ = buf.WriteString(s.Database)
	_, _ = buf.WriteString(" DURATION ")
	_, _ = buf.WriteString(FormatDuration(s.Duration))
	_, _ = buf.WriteString(" REPLICATION ")
	_, _ = buf.WriteString(strconv.Itoa(s.Replication))
	if s.Default {
		_, _ = buf.WriteString(" DEFAULT")
	}
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a CreateRetentionPolicyStatement.
func (s *CreateRetentionPolicyStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// AlterRetentionPolicyStatement represents a command to alter an existing retention policy.
type AlterRetentionPolicyStatement struct {
	// Name of policy to alter.
	Name string

	// Name of the database this policy belongs to.
	Database string

	// Duration data written to this policy will be retained.
	Duration *time.Duration

	// Replication factor for data written to this policy.
	Replication *int

	// Should this policy be set as defalut for the database?
	Default bool
}

// String returns a string representation of the alter retention policy statement.
func (s *AlterRetentionPolicyStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("ALTER RETENTION POLICY ")
	_, _ = buf.WriteString(s.Name)
	_, _ = buf.WriteString(" ON ")
	_, _ = buf.WriteString(s.Database)

	if s.Duration != nil {
		_, _ = buf.WriteString(" DURATION ")
		_, _ = buf.WriteString(FormatDuration(*s.Duration))
	}

	if s.Replication != nil {
		_, _ = buf.WriteString(" REPLICATION ")
		_, _ = buf.WriteString(strconv.Itoa(*s.Replication))
	}

	if s.Default {
		_, _ = buf.WriteString(" DEFAULT")
	}

	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute an AlterRetentionPolicyStatement.
func (s *AlterRetentionPolicyStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// SelectStatement represents a command for extracting data from the database.
type SelectStatement struct {
	// Expressions returned from the selection.
	Fields Fields

	// Target (destination) for the result of the select.
	Target *Target

	// Expressions used for grouping the selection.
	Dimensions Dimensions

	// Data source that fields are extracted from.
	Source Source

	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// Clone returns a deep copy of the statement.
func (stmt *SelectStatement) Clone() *SelectStatement {
	other := &SelectStatement{
		Fields:     make(Fields, len(stmt.Fields)),
		Dimensions: make(Dimensions, len(stmt.Dimensions)),
		Source:     cloneSource(stmt.Source),
		SortFields: make(SortFields, len(stmt.SortFields)),
		Condition:  CloneExpr(stmt.Condition),
		Limit:      stmt.Limit,
	}
	if stmt.Target != nil {
		other.Target = &Target{Measurement: stmt.Target.Measurement, Database: stmt.Target.Database}
	}
	for i, f := range stmt.Fields {
		other.Fields[i] = &Field{Expr: CloneExpr(f.Expr), Alias: f.Alias}
	}
	for i, d := range stmt.Dimensions {
		other.Dimensions[i] = &Dimension{Expr: CloneExpr(d.Expr)}
	}
	// TODO: Copy sources.
	for i, f := range stmt.SortFields {
		other.SortFields[i] = &SortField{Name: f.Name, Ascending: f.Ascending}
	}
	return other
}

func cloneSource(s Source) Source {
	if s == nil {
		return nil
	}

	switch s := s.(type) {
	case *Measurement:
		return &Measurement{Name: s.Name}
	case *Join:
		other := &Join{Measurements: make(Measurements, len(s.Measurements))}
		for i, m := range s.Measurements {
			other.Measurements[i] = &Measurement{Name: m.Name}
		}
		return other
	case *Merge:
		other := &Merge{Measurements: make(Measurements, len(s.Measurements))}
		for i, m := range s.Measurements {
			other.Measurements[i] = &Measurement{Name: m.Name}
		}
		return other
	default:
		panic("unreachable")
	}
}

// String returns a string representation of the select statement.
func (s *SelectStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SELECT ")
	_, _ = buf.WriteString(s.Fields.String())

	if s.Target != nil {
		_, _ = buf.WriteString(" ")
		_, _ = buf.WriteString(s.Target.String())
	}
	_, _ = buf.WriteString(" FROM ")
	_, _ = buf.WriteString(s.Source.String())
	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.Dimensions) > 0 {
		_, _ = buf.WriteString(" GROUP BY ")
		_, _ = buf.WriteString(s.Dimensions.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = fmt.Fprintf(&buf, " LIMIT %d", s.Limit)
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute the SelectStatement.
func (s *SelectStatement) RequiredPrivileges() ExecutionPrivileges {
	ep := ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}

	if s.Target != nil {
		p := ExecutionPrivilege{Name: s.Target.Database, Privilege: WritePrivilege}
		ep = append(ep, p)
	}
	return ep
}

// Aggregated returns true if the statement uses aggregate functions.
func (s *SelectStatement) Aggregated() bool {
	var v bool
	WalkFunc(s.Fields, func(n Node) {
		if _, ok := n.(*Call); ok {
			v = true
		}
	})
	return v
}

/*

BinaryExpr

SELECT mean(xxx.value) + avg(yyy.value) FROM xxx JOIN yyy WHERE xxx.host = 123

from xxx where host = 123
select avg(value) from yyy where host = 123

SELECT xxx.value FROM xxx WHERE xxx.host = 123
SELECT yyy.value FROM yyy

---

SELECT MEAN(xxx.value) + MEAN(cpu.load.value)
FROM xxx JOIN yyy
GROUP BY host
WHERE (xxx.region == "uswest" OR yyy.region == "uswest") AND xxx.otherfield == "XXX"

select * from (
	select mean + mean from xxx join yyy
	group by time(5m), host
) (xxx.region == "uswest" OR yyy.region == "uswest") AND xxx.otherfield == "XXX"

(seriesIDS for xxx.region = 'uswest' union seriesIDs for yyy.regnion = 'uswest') | seriesIDS xxx.otherfield = 'XXX'

WHERE xxx.region == "uswest" AND xxx.otherfield == "XXX"
WHERE yyy.region == "uswest"


*/

// Substatement returns a single-series statement for a given variable reference.
func (s *SelectStatement) Substatement(ref *VarRef) (*SelectStatement, error) {
	// Copy dimensions and properties to new statement.
	other := &SelectStatement{
		Fields:     Fields{{Expr: ref}},
		Dimensions: s.Dimensions,
		Limit:      s.Limit,
		SortFields: s.SortFields,
	}

	// If there is only one series source then return it with the whole condition.
	if _, ok := s.Source.(*Measurement); ok {
		other.Source = s.Source
		other.Condition = s.Condition
		return other, nil
	}

	// Find the matching source.
	name := MatchSource(s.Source, ref.Val)
	if name == "" {
		return nil, fmt.Errorf("field source not found: %s", ref.Val)
	}
	other.Source = &Measurement{Name: name}

	// Filter out conditions.
	if s.Condition != nil {
		other.Condition = filterExprBySource(name, s.Condition)
	}

	return other, nil
}

// filters an expression to exclude expressions unrelated to a source.
func filterExprBySource(name string, expr Expr) Expr {
	switch expr := expr.(type) {
	case *VarRef:
		if !strings.HasPrefix(expr.Val, name) {
			return nil
		}

	case *BinaryExpr:
		lhs := filterExprBySource(name, expr.LHS)
		rhs := filterExprBySource(name, expr.RHS)

		// If an expr is logical then return either LHS/RHS or both.
		// If an expr is arithmetic or comparative then require both sides.
		if expr.Op == AND || expr.Op == OR {
			if lhs == nil && rhs == nil {
				return nil
			} else if lhs != nil && rhs == nil {
				return lhs
			} else if lhs == nil && rhs != nil {
				return rhs
			}
		} else {
			if lhs == nil || rhs == nil {
				return nil
			}
		}
		return &BinaryExpr{Op: expr.Op, LHS: lhs, RHS: rhs}

	case *ParenExpr:
		exp := filterExprBySource(name, expr.Expr)
		if exp == nil {
			return nil
		}
		return &ParenExpr{Expr: exp}
	}
	return expr
}

// MatchSource returns the source name that matches a field name.
// Returns a blank string if no sources match.
func MatchSource(src Source, name string) string {
	switch src := src.(type) {
	case *Measurement:
		if strings.HasPrefix(name, src.Name) {
			return src.Name
		}
	case *Join:
		for _, m := range src.Measurements {
			if strings.HasPrefix(name, m.Name) {
				return m.Name
			}
		}
	case *Merge:
		for _, m := range src.Measurements {
			if strings.HasPrefix(name, m.Name) {
				return m.Name
			}
		}
	}
	return ""
}

// Target represents a target (destination) policy, measurment, and DB.
type Target struct {
	// Measurement to write into.
	Measurement string

	// Database to write into.
	Database string
}

// String returns a string representation of the Target.
func (t *Target) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("INTO ")
	_, _ = buf.WriteString(t.Measurement)

	if t.Database != "" {
		_, _ = buf.WriteString(" ON ")
		_, _ = buf.WriteString(t.Database)
	}

	return buf.String()
}

// DeleteStatement represents a command for removing data from the database.
type DeleteStatement struct {
	// Data source that values are removed from.
	Source Source

	// An expression evaluated on data point.
	Condition Expr
}

// String returns a string representation of the delete statement.
func (s *DeleteStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DELETE ")
	_, _ = buf.WriteString(s.Source.String())
	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	return s.String()
}

// RequiredPrivilege returns the privilege required to execute a DeleteStatement.
func (s *DeleteStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: WritePrivilege}}
}

// ShowSeriesStatement represents a command for listing series in the database.
type ShowSeriesStatement struct {
	// Measurement(s) the series are listed for.
	Source Source

	// An expression evaluated on a series name or tag.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// String returns a string representation of the list series statement.
func (s *ShowSeriesStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW SERIES")

	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = buf.WriteString(" LIMIT ")
		_, _ = buf.WriteString(strconv.Itoa(s.Limit))
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivilege returns the privilege required to execute a ShowSeriesStatement.
func (s *ShowSeriesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// DropSeriesStatement represents a command for removing a series from the database.
type DropSeriesStatement struct {
	Name string
}

// String returns a string representation of the drop series statement.
func (s *DropSeriesStatement) String() string { return fmt.Sprintf("DROP SERIES %s", s.Name) }

// RequiredPrivilege returns the privilige reqired to execute a DropSeriesStatement.
func (s DropSeriesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: WritePrivilege}}
}

// ShowContinuousQueriesStatement represents a command for listing continuous queries.
type ShowContinuousQueriesStatement struct{}

// String returns a string representation of the list continuous queries statement.
func (s *ShowContinuousQueriesStatement) String() string { return "SHOW CONTINUOUS QUERIES" }

// RequiredPrivilege returns the privilege required to execute a ShowContinuousQueriesStatement.
func (s *ShowContinuousQueriesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// ShowDatabasesStatement represents a command for listing all databases in the cluster.
type ShowDatabasesStatement struct{}

// String returns a string representation of the list databases command.
func (s *ShowDatabasesStatement) String() string { return "SHOW DATABASES" }

// RequiredPrivilege returns the privilege required to execute a ShowDatabasesStatement
func (s *ShowDatabasesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// CreateContinuousQueriesStatement represents a command for creating a continuous query.
type CreateContinuousQueryStatement struct {
	// Name of the continuous query to be created.
	Name string

	// Name of the database to create the continuous query on.
	Database string

	// Source of data (SELECT statement).
	Source *SelectStatement
}

// String returns a string representation of the statement.
func (s *CreateContinuousQueryStatement) String() string {
	return fmt.Sprintf("CREATE CONTINUOUS QUERY %s ON %s BEGIN %s END", s.Name, s.Database, s.Source.String())
}

// RequiredPrivilege returns the privilege required to execute a CreateContinuousQueryStatement.
func (s *CreateContinuousQueryStatement) RequiredPrivileges() ExecutionPrivileges {
	ep := ExecutionPrivileges{{Name: s.Database, Privilege: ReadPrivilege}}

	// Selecting into a database that's different from the source?
	if s.Source.Target.Database != "" {
		// Change source database privilege requirement to read.
		ep[0].Privilege = ReadPrivilege

		// Add destination database privilege requirement and set it to write.
		p := ExecutionPrivilege{
			Name:      s.Source.Target.Database,
			Privilege: WritePrivilege,
		}
		ep = append(ep, p)
	}

	return ep
}

// DropContinuousQueriesStatement represents a command for removing a continuous query.
type DropContinuousQueryStatement struct {
	Name string
}

// String returns a string representation of the statement.
func (s *DropContinuousQueryStatement) String() string {
	return fmt.Sprintf("DROP CONTINUOUS QUERY %s", s.Name)
}

// RequiredPrivileges returns the privilege(s) required to execute a DropContinuousQueryStatement
func (s *DropContinuousQueryStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: WritePrivilege}}
}

// ShowMeasurementsStatement represents a command for listing measurements.
type ShowMeasurementsStatement struct {
	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// String returns a string representation of the statement.
func (s *ShowMeasurementsStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW MEASUREMENTS")

	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = buf.WriteString(" LIMIT ")
		_, _ = buf.WriteString(strconv.Itoa(s.Limit))
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowMeasurementsStatement
func (s *ShowMeasurementsStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// ShowRetentionPoliciesStatement represents a command for listing retention policies.
type ShowRetentionPoliciesStatement struct {
	// Name of the database to list policies for.
	Database string
}

// String returns a string representation of a ShowRetentionPoliciesStatement.
func (s *ShowRetentionPoliciesStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW RETENTION POLICIES ")
	_, _ = buf.WriteString(s.Database)
	return buf.String()
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowRetentionPoliciesStatement
func (s *ShowRetentionPoliciesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// ShowTagKeysStatement represents a command for listing tag keys.
type ShowTagKeysStatement struct {
	// Data source that fields are extracted from.
	Source Source

	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// String returns a string representation of the statement.
func (s *ShowTagKeysStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW TAG KEYS")

	if s.Source != nil {
		_, _ = buf.WriteString(" FROM ")
		_, _ = buf.WriteString(s.Source.String())
	}
	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = buf.WriteString(" LIMIT ")
		_, _ = buf.WriteString(strconv.Itoa(s.Limit))
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowTagKeysStatement
func (s *ShowTagKeysStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// ShowTagValuesStatement represents a command for listing tag values.
type ShowTagValuesStatement struct {
	// Data source that fields are extracted from.
	Source Source

	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// String returns a string representation of the statement.
func (s *ShowTagValuesStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW TAG VALUES")

	if s.Source != nil {
		_, _ = buf.WriteString(" FROM ")
		_, _ = buf.WriteString(s.Source.String())
	}
	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = buf.WriteString(" LIMIT ")
		_, _ = buf.WriteString(strconv.Itoa(s.Limit))
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowTagValuesStatement
func (s *ShowTagValuesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// ShowUsersStatement represents a command for listing users.
type ShowUsersStatement struct{}

// String retuns a string representation of the ShowUsersStatement.
func (s *ShowUsersStatement) String() string {
	return "SHOW USERS"
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowUsersStatement
func (s *ShowUsersStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: AllPrivileges}}
}

// ShowFieldKeyStatement represents a command for listing field keys.
type ShowFieldKeysStatement struct {
	// Data source that fields are extracted from.
	Source Source

	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// String returns a string representation of the statement.
func (s *ShowFieldKeysStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW FIELD KEYS")

	if s.Source != nil {
		_, _ = buf.WriteString(" FROM ")
		_, _ = buf.WriteString(s.Source.String())
	}
	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = buf.WriteString(" LIMIT ")
		_, _ = buf.WriteString(strconv.Itoa(s.Limit))
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowFieldKeysStatement
func (s *ShowFieldKeysStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// ShowFieldValuesStatement represents a command for listing field values.
type ShowFieldValuesStatement struct {
	// Data source that fields are extracted from.
	Source Source

	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// String returns a string representation of the statement.
func (s *ShowFieldValuesStatement) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("SHOW FIELD VALUES")

	if s.Source != nil {
		_, _ = buf.WriteString(" FROM ")
		_, _ = buf.WriteString(s.Source.String())
	}
	if s.Condition != nil {
		_, _ = buf.WriteString(" WHERE ")
		_, _ = buf.WriteString(s.Condition.String())
	}
	if len(s.SortFields) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		_, _ = buf.WriteString(s.SortFields.String())
	}
	if s.Limit > 0 {
		_, _ = buf.WriteString(" LIMIT ")
		_, _ = buf.WriteString(strconv.Itoa(s.Limit))
	}
	if s.Offset > 0 {
		_, _ = buf.WriteString(" OFFSET ")
		_, _ = buf.WriteString(strconv.Itoa(s.Offset))
	}
	return buf.String()
}

// RequiredPrivileges returns the privilege(s) required to execute a ShowFieldValuesStatement
func (s *ShowFieldValuesStatement) RequiredPrivileges() ExecutionPrivileges {
	return ExecutionPrivileges{{Name: "", Privilege: ReadPrivilege}}
}

// Fields represents a list of fields.
type Fields []*Field

// String returns a string representation of the fields.
func (a Fields) String() string {
	var str []string
	for _, f := range a {
		str = append(str, f.String())
	}
	return strings.Join(str, ", ")
}

// Field represents an expression retrieved from a select statement.
type Field struct {
	Expr  Expr
	Alias string
}

// Name returns the name of the field. Returns alias, if set.
// Otherwise uses the function name or variable name.
func (f *Field) Name() string {
	// Return alias, if set.
	if f.Alias != "" {
		return f.Alias
	}

	// Return the function name or variable name, if available.
	switch expr := f.Expr.(type) {
	case *Call:
		return expr.Name
	case *VarRef:
		return expr.Val
	}

	// Otherwise return a blank name.
	return ""
}

// String returns a string representation of the field.
func (f *Field) String() string {
	if f.Alias == "" {
		return f.Expr.String()
	}
	return fmt.Sprintf("%s AS %s", f.Expr.String(), f.Alias)
}

// Dimensions represents a list of dimensions.
type Dimensions []*Dimension

// String returns a string representation of the dimensions.
func (a Dimensions) String() string {
	var str []string
	for _, d := range a {
		str = append(str, d.String())
	}
	return strings.Join(str, ", ")
}

// Normalize returns the interval and tag dimensions separately.
// Returns 0 if no time interval is specified.
// Returns an error if multiple time dimensions exist or if non-VarRef dimensions are specified.
func (a Dimensions) Normalize() (time.Duration, []string, error) {
	var dur time.Duration
	var tags []string

	for _, dim := range a {
		switch expr := dim.Expr.(type) {
		case *Call:
			// Ensure the call is time() and it only has one duration argument.
			// If we already have a duration
			if strings.ToLower(expr.Name) != "time" {
				return 0, nil, errors.New("only time() calls allowed in dimensions")
			} else if len(expr.Args) != 1 {
				return 0, nil, errors.New("time dimension expected one argument")
			} else if lit, ok := expr.Args[0].(*DurationLiteral); !ok {
				return 0, nil, errors.New("time dimension must have one duration argument")
			} else if dur != 0 {
				return 0, nil, errors.New("multiple time dimensions not allowed")
			} else {
				dur = lit.Val
			}

		case *VarRef:
			tags = append(tags, expr.Val)

		default:
			return 0, nil, errors.New("only time and tag dimensions allowed")
		}
	}

	return dur, tags, nil
}

// Dimension represents an expression that a select statement is grouped by.
type Dimension struct {
	Expr Expr
}

// String returns a string representation of the dimension.
func (d *Dimension) String() string { return d.Expr.String() }

// Measurements represents a list of measurements.
type Measurements []*Measurement

// String returns a string representation of the measurements.
func (a Measurements) String() string {
	var str []string
	for _, m := range a {
		str = append(str, m.String())
	}
	return strings.Join(str, ", ")
}

// Measurement represents a single measurement used as a datasource.
type Measurement struct {
	Name string
}

// String returns a string representation of the measurement.
func (m *Measurement) String() string { return m.Name }

// Join represents two datasources joined together.
type Join struct {
	Measurements Measurements
}

// String returns a string representation of the join.
func (j *Join) String() string {
	return fmt.Sprintf("join(%s)", j.Measurements.String())
}

// Merge represents a datasource created by merging two datasources.
type Merge struct {
	Measurements Measurements
}

// String returns a string representation of the merge.
func (m *Merge) String() string {
	return fmt.Sprintf("merge(%s)", m.Measurements.String())
}

// VarRef represents a reference to a variable.
type VarRef struct {
	Val string
}

// String returns a string representation of the variable reference.
func (r *VarRef) String() string { return r.Val }

// Call represents a function call.
type Call struct {
	Name string
	Args []Expr
}

// String returns a string representation of the call.
func (c *Call) String() string {
	// Join arguments.
	var str []string
	for _, arg := range c.Args {
		str = append(str, arg.String())
	}

	// Write function name and args.
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(str, ", "))
}

// NumberLiteral represents a numeric literal.
type NumberLiteral struct {
	Val float64
}

// String returns a string representation of the literal.
func (l *NumberLiteral) String() string { return strconv.FormatFloat(l.Val, 'f', 3, 64) }

// BooleanLiteral represents a boolean literal.
type BooleanLiteral struct {
	Val bool
}

// String returns a string representation of the literal.
func (l *BooleanLiteral) String() string {
	if l.Val {
		return "true"
	}
	return "false"
}

// isTrueLiteral returns true if the expression is a literal "true" value.
func isTrueLiteral(expr Expr) bool {
	if expr, ok := expr.(*BooleanLiteral); ok {
		return expr.Val == true
	}
	return false
}

// isFalseLiteral returns true if the expression is a literal "false" value.
func isFalseLiteral(expr Expr) bool {
	if expr, ok := expr.(*BooleanLiteral); ok {
		return expr.Val == false
	}
	return false
}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Val string
}

// String returns a string representation of the literal.
func (l *StringLiteral) String() string { return QuoteString(l.Val) }

// TagKeyIdent represents a special TAG KEY identifier.
type TagKeyIdent struct{}

// String returns a string representation of the TagKeyIdent.
func (t *TagKeyIdent) String() string { return "TAG KEY" }

// TimeLiteral represents a point-in-time literal.
type TimeLiteral struct {
	Val time.Time
}

// String returns a string representation of the literal.
func (l *TimeLiteral) String() string {
	return `"` + l.Val.UTC().Format(DateTimeFormat) + `"`
}

// DurationLiteral represents a duration literal.
type DurationLiteral struct {
	Val time.Duration
}

// String returns a string representation of the literal.
func (l *DurationLiteral) String() string { return FormatDuration(l.Val) }

// nilLiteral represents a nil literal.
// This is not available to the query language itself. It's only used internally.
type nilLiteral struct{}

// String returns a string representation of the literal.
func (l *nilLiteral) String() string { return `nil` }

// BinaryExpr represents an operation between two expressions.
type BinaryExpr struct {
	Op  Token
	LHS Expr
	RHS Expr
}

// String returns a string representation of the binary expression.
func (e *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.LHS.String(), e.Op.String(), e.RHS.String())
}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Expr Expr
}

// String returns a string representation of the parenthesized expression.
func (e *ParenExpr) String() string { return fmt.Sprintf("(%s)", e.Expr.String()) }

// Wildcard represents a wild card expression.
type Wildcard struct{}

// String returns a string representation of the wildcard.
func (e *Wildcard) String() string { return "*" }

// CloneExpr returns a deep copy of the expression.
func CloneExpr(expr Expr) Expr {
	if expr == nil {
		return nil
	}
	switch expr := expr.(type) {
	case *BinaryExpr:
		return &BinaryExpr{Op: expr.Op, LHS: CloneExpr(expr.LHS), RHS: CloneExpr(expr.RHS)}
	case *BooleanLiteral:
		return &BooleanLiteral{Val: expr.Val}
	case *Call:
		args := make([]Expr, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = CloneExpr(arg)
		}
		return &Call{Name: expr.Name, Args: args}
	case *DurationLiteral:
		return &DurationLiteral{Val: expr.Val}
	case *NumberLiteral:
		return &NumberLiteral{Val: expr.Val}
	case *ParenExpr:
		return &ParenExpr{Expr: CloneExpr(expr.Expr)}
	case *StringLiteral:
		return &StringLiteral{Val: expr.Val}
	case *TimeLiteral:
		return &TimeLiteral{Val: expr.Val}
	case *VarRef:
		return &VarRef{Val: expr.Val}
	case *Wildcard:
		return &Wildcard{}
	}
	panic("unreachable")
}

// TimeRange returns the minimum and maximum times specified by an expression.
// Returns zero times if there is no bound.
func TimeRange(expr Expr) (min, max time.Time) {
	WalkFunc(expr, func(n Node) {
		if n, ok := n.(*BinaryExpr); ok {
			// Extract literal expression & operator on LHS.
			// Check for "time" on the left-hand side first.
			// Otherwise check for for the right-hand side and flip the operator.
			value, op := timeExprValue(n.LHS, n.RHS), n.Op
			if value.IsZero() {
				if value = timeExprValue(n.RHS, n.LHS); value.IsZero() {
					return
				} else if op == LT {
					op = GT
				} else if op == LTE {
					op = GTE
				} else if op == GT {
					op = LT
				} else if op == GTE {
					op = LTE
				}
			}

			// Update the min/max depending on the operator.
			// The GT & LT update the value by +/- 1µs not make them "not equal".
			switch op {
			case GT:
				if min.IsZero() || value.After(min) {
					min = value.Add(time.Microsecond)
				}
			case GTE:
				if min.IsZero() || value.After(min) {
					min = value
				}
			case LT:
				if max.IsZero() || value.Before(max) {
					max = value.Add(-time.Microsecond)
				}
			case LTE:
				if max.IsZero() || value.Before(max) {
					max = value
				}
			case EQ:
				if min.IsZero() || value.After(min) {
					min = value
				}
				if max.IsZero() || value.Before(max) {
					max = value
				}
			}
		}
	})
	return
}

// timeExprValue returns the time literal value of a "time == <TimeLiteral>" expression.
// Returns zero time if the expression is not a time expression.
func timeExprValue(ref Expr, lit Expr) time.Time {
	if ref, ok := ref.(*VarRef); ok && strings.ToLower(ref.Val) == "time" {
		switch lit := lit.(type) {
		case *TimeLiteral:
			return lit.Val
		case *DurationLiteral:
			return time.Unix(0, int64(lit.Val)).UTC()
		}
	}
	return time.Time{}
}

// Visitor can be called by Walk to traverse an AST hierarchy.
// The Visit() function is called once per node.
type Visitor interface {
	Visit(Node) Visitor
}

// Walk traverses a node hierarchy in depth-first order.
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Query:
		Walk(v, n.Statements)

	case Statements:
		for _, s := range n {
			Walk(v, s)
		}

	case *SelectStatement:
		Walk(v, n.Fields)
		Walk(v, n.Dimensions)
		Walk(v, n.Source)
		Walk(v, n.Condition)

	case *ShowSeriesStatement:
		Walk(v, n.Source)
		Walk(v, n.Condition)

	case Fields:
		for _, c := range n {
			Walk(v, c)
		}

	case *Field:
		Walk(v, n.Expr)

	case Dimensions:
		for _, c := range n {
			Walk(v, c)
		}

	case *Dimension:
		Walk(v, n.Expr)

	case *BinaryExpr:
		Walk(v, n.LHS)
		Walk(v, n.RHS)

	case *ParenExpr:
		Walk(v, n.Expr)

	case *Call:
		for _, expr := range n.Args {
			Walk(v, expr)
		}
	}
}

// WalkFunc traverses a node hierarchy in depth-first order.
func WalkFunc(node Node, fn func(Node)) {
	Walk(walkFuncVisitor(fn), node)
}

type walkFuncVisitor func(Node)

func (fn walkFuncVisitor) Visit(n Node) Visitor { fn(n); return fn }

// Rewriter can be called by Rewrite to replace nodes in the AST hierarchy.
// The Rewrite() function is called once per node.
type Rewriter interface {
	Rewrite(Node) Node
}

// Rewrite recursively invokes the rewriter to replace each node.
// Nodes are traversed depth-first and rewritten from leaf to root.
func Rewrite(r Rewriter, node Node) Node {
	switch n := node.(type) {
	case *Query:
		n.Statements = Rewrite(r, n.Statements).(Statements)

	case Statements:
		for i, s := range n {
			n[i] = Rewrite(r, s).(Statement)
		}

	case *SelectStatement:
		n.Fields = Rewrite(r, n.Fields).(Fields)
		n.Dimensions = Rewrite(r, n.Dimensions).(Dimensions)
		n.Source = Rewrite(r, n.Source).(Source)
		n.Condition = Rewrite(r, n.Condition).(Expr)

	case Fields:
		for i, f := range n {
			n[i] = Rewrite(r, f).(*Field)
		}

	case *Field:
		n.Expr = Rewrite(r, n.Expr).(Expr)

	case Dimensions:
		for i, d := range n {
			n[i] = Rewrite(r, d).(*Dimension)
		}

	case *Dimension:
		n.Expr = Rewrite(r, n.Expr).(Expr)

	case *BinaryExpr:
		n.LHS = Rewrite(r, n.LHS).(Expr)
		n.RHS = Rewrite(r, n.RHS).(Expr)

	case *ParenExpr:
		n.Expr = Rewrite(r, n.Expr).(Expr)

	case *Call:
		for i, expr := range n.Args {
			n.Args[i] = Rewrite(r, expr).(Expr)
		}
	}

	return r.Rewrite(node)
}

// RewriteFunc rewrites a node hierarchy.
func RewriteFunc(node Node, fn func(Node) Node) Node {
	return Rewrite(rewriterFunc(fn), node)
}

type rewriterFunc func(Node) Node

func (fn rewriterFunc) Rewrite(n Node) Node { return fn(n) }

// Eval evaluates expr against a map.
func Eval(expr Expr, m map[string]interface{}) interface{} {
	if expr == nil {
		return nil
	}

	switch expr := expr.(type) {
	case *BinaryExpr:
		return evalBinaryExpr(expr, m)
	case *BooleanLiteral:
		return expr.Val
	case *NumberLiteral:
		return expr.Val
	case *ParenExpr:
		return Eval(expr.Expr, m)
	case *StringLiteral:
		return expr.Val
	case *VarRef:
		return m[expr.Val]
	default:
		return nil
	}
}

func evalBinaryExpr(expr *BinaryExpr, m map[string]interface{}) interface{} {
	lhs := Eval(expr.LHS, m)
	rhs := Eval(expr.RHS, m)

	// Evaluate if both sides are simple types.
	switch lhs := lhs.(type) {
	case bool:
		rhs, _ := rhs.(bool)
		switch expr.Op {
		case AND:
			return lhs && rhs
		case OR:
			return lhs || rhs
		}
	case float64:
		rhs, _ := rhs.(float64)
		switch expr.Op {
		case EQ:
			return lhs == rhs
		case NEQ:
			return lhs != rhs
		case LT:
			return lhs < rhs
		case LTE:
			return lhs <= rhs
		case GT:
			return lhs > rhs
		case GTE:
			return lhs >= rhs
		case ADD:
			return lhs + rhs
		case SUB:
			return lhs - rhs
		case MUL:
			return lhs * rhs
		case DIV:
			if rhs == 0 {
				return float64(0)
			}
			return lhs / rhs
		}
	case string:
		rhs, _ := rhs.(string)
		switch expr.Op {
		case EQ:
			return lhs == rhs
		case NEQ:
			return lhs != rhs
		}
	}
	return nil
}

// Reduce evaluates expr using the available values in valuer.
// References that don't exist in valuer are ignored.
func Reduce(expr Expr, valuer Valuer) Expr {
	expr = reduce(expr, valuer)

	// Unwrap parens at top level.
	if expr, ok := expr.(*ParenExpr); ok {
		return expr.Expr
	}
	return expr
}

func reduce(expr Expr, valuer Valuer) Expr {
	if expr == nil {
		return nil
	}

	switch expr := expr.(type) {
	case *BinaryExpr:
		return reduceBinaryExpr(expr, valuer)
	case *Call:
		return reduceCall(expr, valuer)
	case *ParenExpr:
		return reduceParenExpr(expr, valuer)
	case *VarRef:
		return reduceVarRef(expr, valuer)
	default:
		return CloneExpr(expr)
	}
}

func reduceBinaryExpr(expr *BinaryExpr, valuer Valuer) Expr {
	// Reduce both sides first.
	op := expr.Op
	lhs := reduce(expr.LHS, valuer)
	rhs := reduce(expr.RHS, valuer)

	// Do not evaluate if one side is nil.
	if lhs == nil || rhs == nil {
		return &BinaryExpr{LHS: lhs, RHS: rhs, Op: expr.Op}
	}

	// If we have a logical operator (AND, OR) and one side is a boolean literal
	// then we need to have special handling.
	if op == AND {
		if isFalseLiteral(lhs) || isFalseLiteral(rhs) {
			return &BooleanLiteral{Val: false}
		} else if isTrueLiteral(lhs) {
			return rhs
		} else if isTrueLiteral(rhs) {
			return lhs
		}
	} else if op == OR {
		if isTrueLiteral(lhs) || isTrueLiteral(rhs) {
			return &BooleanLiteral{Val: true}
		} else if isFalseLiteral(lhs) {
			return rhs
		} else if isFalseLiteral(rhs) {
			return lhs
		}
	}

	// Evaluate if both sides are simple types.
	switch lhs := lhs.(type) {
	case *BooleanLiteral:
		return reduceBinaryExprBooleanLHS(op, lhs, rhs)
	case *DurationLiteral:
		return reduceBinaryExprDurationLHS(op, lhs, rhs)
	case *nilLiteral:
		return reduceBinaryExprNilLHS(op, lhs, rhs)
	case *NumberLiteral:
		return reduceBinaryExprNumberLHS(op, lhs, rhs)
	case *StringLiteral:
		return reduceBinaryExprStringLHS(op, lhs, rhs)
	case *TimeLiteral:
		return reduceBinaryExprTimeLHS(op, lhs, rhs)
	default:
		return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
	}
}

func reduceBinaryExprBooleanLHS(op Token, lhs *BooleanLiteral, rhs Expr) Expr {
	switch rhs := rhs.(type) {
	case *BooleanLiteral:
		switch op {
		case EQ:
			return &BooleanLiteral{Val: lhs.Val == rhs.Val}
		case NEQ:
			return &BooleanLiteral{Val: lhs.Val != rhs.Val}
		case AND:
			return &BooleanLiteral{Val: lhs.Val && rhs.Val}
		case OR:
			return &BooleanLiteral{Val: lhs.Val || rhs.Val}
		}
	case *nilLiteral:
		return &BooleanLiteral{Val: false}
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
}

func reduceBinaryExprDurationLHS(op Token, lhs *DurationLiteral, rhs Expr) Expr {
	switch rhs := rhs.(type) {
	case *DurationLiteral:
		switch op {
		case ADD:
			return &DurationLiteral{Val: lhs.Val + rhs.Val}
		case SUB:
			return &DurationLiteral{Val: lhs.Val - rhs.Val}
		case EQ:
			return &BooleanLiteral{Val: lhs.Val == rhs.Val}
		case NEQ:
			return &BooleanLiteral{Val: lhs.Val != rhs.Val}
		case GT:
			return &BooleanLiteral{Val: lhs.Val > rhs.Val}
		case GTE:
			return &BooleanLiteral{Val: lhs.Val >= rhs.Val}
		case LT:
			return &BooleanLiteral{Val: lhs.Val < rhs.Val}
		case LTE:
			return &BooleanLiteral{Val: lhs.Val <= rhs.Val}
		}
	case *NumberLiteral:
		switch op {
		case MUL:
			return &DurationLiteral{Val: lhs.Val * time.Duration(rhs.Val)}
		case DIV:
			if rhs.Val == 0 {
				return &DurationLiteral{Val: 0}
			}
			return &DurationLiteral{Val: lhs.Val / time.Duration(rhs.Val)}
		}
	case *TimeLiteral:
		switch op {
		case ADD:
			return &TimeLiteral{Val: rhs.Val.Add(lhs.Val)}
		}
	case *nilLiteral:
		return &BooleanLiteral{Val: false}
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
}

func reduceBinaryExprNilLHS(op Token, lhs *nilLiteral, rhs Expr) Expr {
	switch op {
	case EQ, NEQ:
		return &BooleanLiteral{Val: false}
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
}

func reduceBinaryExprNumberLHS(op Token, lhs *NumberLiteral, rhs Expr) Expr {
	switch rhs := rhs.(type) {
	case *NumberLiteral:
		switch op {
		case ADD:
			return &NumberLiteral{Val: lhs.Val + rhs.Val}
		case SUB:
			return &NumberLiteral{Val: lhs.Val - rhs.Val}
		case MUL:
			return &NumberLiteral{Val: lhs.Val * rhs.Val}
		case DIV:
			if rhs.Val == 0 {
				return &NumberLiteral{Val: 0}
			}
			return &NumberLiteral{Val: lhs.Val / rhs.Val}
		case EQ:
			return &BooleanLiteral{Val: lhs.Val == rhs.Val}
		case NEQ:
			return &BooleanLiteral{Val: lhs.Val != rhs.Val}
		case GT:
			return &BooleanLiteral{Val: lhs.Val > rhs.Val}
		case GTE:
			return &BooleanLiteral{Val: lhs.Val >= rhs.Val}
		case LT:
			return &BooleanLiteral{Val: lhs.Val < rhs.Val}
		case LTE:
			return &BooleanLiteral{Val: lhs.Val <= rhs.Val}
		}
	case *nilLiteral:
		return &BooleanLiteral{Val: false}
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
}

func reduceBinaryExprStringLHS(op Token, lhs *StringLiteral, rhs Expr) Expr {
	switch rhs := rhs.(type) {
	case *StringLiteral:
		switch op {
		case EQ:
			return &BooleanLiteral{Val: lhs.Val == rhs.Val}
		case NEQ:
			return &BooleanLiteral{Val: lhs.Val != rhs.Val}
		case ADD:
			return &StringLiteral{Val: lhs.Val + rhs.Val}
		}
	case *nilLiteral:
		switch op {
		case EQ, NEQ:
			return &BooleanLiteral{Val: false}
		}
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
}

func reduceBinaryExprTimeLHS(op Token, lhs *TimeLiteral, rhs Expr) Expr {
	switch rhs := rhs.(type) {
	case *DurationLiteral:
		switch op {
		case ADD:
			return &TimeLiteral{Val: lhs.Val.Add(rhs.Val)}
		case SUB:
			return &TimeLiteral{Val: lhs.Val.Add(-rhs.Val)}
		}
	case *TimeLiteral:
		switch op {
		case SUB:
			return &DurationLiteral{Val: lhs.Val.Sub(rhs.Val)}
		case EQ:
			return &BooleanLiteral{Val: lhs.Val.Equal(rhs.Val)}
		case NEQ:
			return &BooleanLiteral{Val: !lhs.Val.Equal(rhs.Val)}
		case GT:
			return &BooleanLiteral{Val: lhs.Val.After(rhs.Val)}
		case GTE:
			return &BooleanLiteral{Val: lhs.Val.After(rhs.Val) || lhs.Val.Equal(rhs.Val)}
		case LT:
			return &BooleanLiteral{Val: lhs.Val.Before(rhs.Val)}
		case LTE:
			return &BooleanLiteral{Val: lhs.Val.Before(rhs.Val) || lhs.Val.Equal(rhs.Val)}
		}
	case *nilLiteral:
		return &BooleanLiteral{Val: false}
	}
	return &BinaryExpr{Op: op, LHS: lhs, RHS: rhs}
}

func reduceCall(expr *Call, valuer Valuer) Expr {
	// Evaluate "now()" if valuer is set.
	if strings.ToLower(expr.Name) == "now" && len(expr.Args) == 0 && valuer != nil {
		if v, ok := valuer.Value("now()"); ok {
			v, _ := v.(time.Time)
			return &TimeLiteral{Val: v}
		}
	}

	// Otherwise reduce arguments.
	args := make([]Expr, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = reduce(arg, valuer)
	}
	return &Call{Name: expr.Name, Args: args}
}

func reduceParenExpr(expr *ParenExpr, valuer Valuer) Expr {
	subexpr := reduce(expr.Expr, valuer)
	if subexpr, ok := subexpr.(*BinaryExpr); ok {
		return &ParenExpr{Expr: subexpr}
	}
	return subexpr
}

func reduceVarRef(expr *VarRef, valuer Valuer) Expr {
	// Ignore if there is no valuer.
	if valuer == nil {
		return &VarRef{Val: expr.Val}
	}

	// Retrieve the value of the ref.
	// Ignore if the value doesn't exist.
	v, ok := valuer.Value(expr.Val)
	if !ok {
		return &VarRef{Val: expr.Val}
	}

	// Return the value as a literal.
	switch v := v.(type) {
	case bool:
		return &BooleanLiteral{Val: v}
	case time.Duration:
		return &DurationLiteral{Val: v}
	case float64:
		return &NumberLiteral{Val: v}
	case string:
		return &StringLiteral{Val: v}
	case time.Time:
		return &TimeLiteral{Val: v}
	default:
		return &nilLiteral{}
	}
}

// Valuer is the interface that wraps the Value() method.
//
// Value returns the value and existence flag for a given key.
type Valuer interface {
	Value(key string) (interface{}, bool)
}

// nowValuer returns only the value for "now()".
type nowValuer struct {
	Now time.Time
}

func (v *nowValuer) Value(key string) (interface{}, bool) {
	if key == "now()" {
		return v.Now, true
	}
	return nil, false
}