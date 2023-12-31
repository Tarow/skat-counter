//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var Round = newRoundTable("", "round", "")

type roundTable struct {
	sqlite.Table

	// Columns
	ID        sqlite.ColumnInteger
	GameID    sqlite.ColumnInteger
	CreatedAt sqlite.ColumnTimestamp
	Dealer    sqlite.ColumnInteger
	Declarer  sqlite.ColumnInteger
	Won       sqlite.ColumnBool
	Value     sqlite.ColumnInteger

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type RoundTable struct {
	roundTable

	EXCLUDED roundTable
}

// AS creates new RoundTable with assigned alias
func (a RoundTable) AS(alias string) *RoundTable {
	return newRoundTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new RoundTable with assigned schema name
func (a RoundTable) FromSchema(schemaName string) *RoundTable {
	return newRoundTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new RoundTable with assigned table prefix
func (a RoundTable) WithPrefix(prefix string) *RoundTable {
	return newRoundTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new RoundTable with assigned table suffix
func (a RoundTable) WithSuffix(suffix string) *RoundTable {
	return newRoundTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newRoundTable(schemaName, tableName, alias string) *RoundTable {
	return &RoundTable{
		roundTable: newRoundTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newRoundTableImpl("", "excluded", ""),
	}
}

func newRoundTableImpl(schemaName, tableName, alias string) roundTable {
	var (
		IDColumn        = sqlite.IntegerColumn("id")
		GameIDColumn    = sqlite.IntegerColumn("game_id")
		CreatedAtColumn = sqlite.TimestampColumn("created_at")
		DealerColumn    = sqlite.IntegerColumn("dealer")
		DeclarerColumn  = sqlite.IntegerColumn("declarer")
		WonColumn       = sqlite.BoolColumn("won")
		ValueColumn     = sqlite.IntegerColumn("value")
		allColumns      = sqlite.ColumnList{IDColumn, GameIDColumn, CreatedAtColumn, DealerColumn, DeclarerColumn, WonColumn, ValueColumn}
		mutableColumns  = sqlite.ColumnList{GameIDColumn, CreatedAtColumn, DealerColumn, DeclarerColumn, WonColumn, ValueColumn}
	)

	return roundTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		GameID:    GameIDColumn,
		CreatedAt: CreatedAtColumn,
		Dealer:    DealerColumn,
		Declarer:  DeclarerColumn,
		Won:       WonColumn,
		Value:     ValueColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
