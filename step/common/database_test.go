package common_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/step/common"
)

func setUp(t *testing.T) (*logrus.Entry, datasource.Datasourcer, sqlmock.Sqlmock) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ds.MockedDb = db
	return log, ds, mock
}

func TestToSkipOk(t *testing.T) {
	type e struct {
		query string
		value int
	}

	type a struct {
		query    string
		inverted bool
		value    int
	}

	tests := []struct {
		name     string
		expected []e
		args     []a
		wantOk   bool
		wantErr  bool
	}{
		{
			name: "skipSimple",
			expected: []e{
				{
					query: "SELECT COUNT\\(id\\) from dtable WHERE title like '%'",
					value: 10,
				},
				{
					query: "SELECT COUNT\\(id\\) from stable WHERE title like '%'",
					value: 1,
				},
			},
			args: []a{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "noSkipSimple",
			expected: []e{
				{
					query: "SELECT COUNT\\(id\\) from dtable WHERE title like '%'",
					value: 0,
				},
			},
			args: []a{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "skipValue",
			expected: []e{
				{
					query: "SELECT COUNT\\(id\\) from dtable WHERE title like '%'",
					value: 0,
				},
				{
					query: "SELECT COUNT\\(id\\) from stable WHERE title like '%'",
					value: 1,
				},
			},
			args: []a{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: false,
					value:    10,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "skipInverted",
			expected: []e{
				{
					query: "SELECT COUNT\\(id\\) from dtable WHERE title like '%'",
					value: 2,
				},
				{
					query: "SELECT COUNT\\(id\\) from stable WHERE title like '%'",
					value: 1,
				},
			},
			args: []a{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: true,
					value:    2,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "noSkipInverted",
			expected: []e{
				{
					query: "SELECT COUNT\\(id\\) from dtable WHERE title like '%'",
					value: 0,
				},
				{
					query: "SELECT COUNT\\(id\\) from stable WHERE title like '%'",
					value: 10,
				},
			},
			args: []a{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: true,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: true,
					value:    0,
				},
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "noSkipSecond",
			expected: []e{
				{
					query: "SELECT COUNT\\(id\\) from dtable WHERE title like '%'",
					value: 10,
				},
				{
					query: "SELECT COUNT\\(id\\) from stable WHERE title like '%'",
					value: 1,
				},
			},
			args: []a{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: true,
					value:    0,
				},
			},
			wantOk:  false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, ds, mock := setUp(t)

			for _, exp := range tt.expected {
				rows := sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(exp.value)
				mock.ExpectQuery(exp.query).WillReturnRows(rows)
			}

			skq := make([]common.SkipQuery, 0, len(tt.args))
			for _, arg := range tt.args {
				skq = append(skq, common.NewSkipQuery(arg.query, arg.inverted, arg.value))
			}

			ok, err := common.ToSkipDatabase(
				context.Background(),
				log,
				ds,
				false,
				false,
				skq,
			)

			switch {
			case tt.wantErr && err == nil:
				t.Error("ToSkipDatabase should failed")
			case !tt.wantErr && err != nil:
				t.Error("ToSkipDatabase should not failed")
			case tt.wantOk && !ok:
				t.Error("ToSkipDatabase should return true")
			case !tt.wantOk && ok:
				t.Error("ToSkipDatabase should return false")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestToSkipError(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ds.MockedDb = db
	mock.ExpectQuery("SELECT COUNT\\(id\\) from stable WHERE title like '%'").WillReturnError(fmt.Errorf("fake error"))

	_, err = common.ToSkipDatabase(
		context.Background(),
		log,
		ds,
		false,
		false,
		[]common.SkipQuery{
			common.NewSkipQuery("SELECT COUNT(id) from stable WHERE title like '%'", false, 0),
			common.NewSkipQuery("SELECT COUNT(id) from dtable WHERE title like '%'", false, 0),
		})
	if err == nil {
		t.Errorf("ToSkip should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestToSkipOpenError(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	ds.ErrorOpenDb = fmt.Errorf("Fake OpenDatabase error")

	_, err := common.ToSkipDatabase(
		context.Background(),
		log,
		ds,
		false,
		false,
		[]common.SkipQuery{
			common.NewSkipQuery("SELECT COUNT(id) from stable WHERE title like '%'", false, 0),
			common.NewSkipQuery("SELECT COUNT(id) from dtable WHERE title like '%'", false, 0),
		})
	if err == nil {
		t.Errorf("ToSkip should return error")
	}
}

func TestParseRenderQueries(t *testing.T) {
	type e struct {
		query    string
		inverted bool
		value    int
	}

	tests := []struct {
		name          string
		expected      []e
		queries       []string
		wantParseErr  bool
		wantRenderErr bool
	}{
		{
			name: "renderSimple",
			queries: []string{
				"SELECT COUNT(id) from dtable WHERE title like '%'",
				"!SELECT COUNT(id) from stable WHERE title like '%'",
				"=10:SELECT COUNT(id) from dtable WHERE title like '%'",
				"!=10:SELECT COUNT(id) from stable WHERE title like '%'",
			},
			expected: []e{
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: true,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from dtable WHERE title like '%'",
					inverted: false,
					value:    10,
				},
				{
					query:    "SELECT COUNT(id) from stable WHERE title like '%'",
					inverted: true,
					value:    10,
				},
			},
			wantParseErr:  false,
			wantRenderErr: false,
		},
		{
			name: "renderTemplate",
			queries: []string{
				"SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
				"=10:SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!=10:SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
			},
			expected: []e{
				{
					query:    "SELECT COUNT(id) from public.dtable WHERE title like '%'",
					inverted: false,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from public.stable WHERE title like '%'",
					inverted: true,
					value:    0,
				},
				{
					query:    "SELECT COUNT(id) from public.dtable WHERE title like '%'",
					inverted: false,
					value:    10,
				},
				{
					query:    "SELECT COUNT(id) from public.stable WHERE title like '%'",
					inverted: true,
					value:    10,
				},
			},
			wantParseErr:  false,
			wantRenderErr: false,
		},
		{
			name:          "ParseNoQuery",
			queries:       []string{},
			expected:      []e{},
			wantParseErr:  true,
			wantRenderErr: false,
		},
		{
			name: "ParseEmptyQuery",
			queries: []string{
				"SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"",
				"=10SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!=10:SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
			},
			expected:      []e{},
			wantParseErr:  true,
			wantRenderErr: false,
		},
		{
			name: "ParseValueSeparatorError",
			queries: []string{
				"SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
				"=10SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!=10:SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
			},
			expected:      []e{},
			wantParseErr:  true,
			wantRenderErr: false,
		},
		{
			name: "ParseValueError",
			queries: []string{
				"SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
				"=1a:SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!=10:SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
			},
			expected:      []e{},
			wantParseErr:  true,
			wantRenderErr: false,
		},
		{
			name: "ParseTemplateError",
			queries: []string{
				"SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!SELECT COUNT(id) from {{.Schema.stable WHERE title like '%'",
				"=10:SELECT COUNT(id) from {{.Schema}}.dtable WHERE title like '%'",
				"!=10:SELECT COUNT(id) from {{.Schema}}.stable WHERE title like '%'",
			},
			expected:      []e{},
			wantParseErr:  true,
			wantRenderErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logrus.New()
			log := logger.WithField("appname", "kamino")
			tmpValues := datasource.TmplValues{Schema: "public"}

			tqueries, err := common.ParseQueries(log, tt.queries)
			if tt.wantParseErr && err == nil {
				t.Fatal("ParseQueries should failed")
			} else if !tt.wantParseErr && err != nil {
				t.Fatal("ParseQueries should not failed")
			}

			if err == nil {
				if len(tqueries) != len(tt.queries) {
					t.Fatal("ParseQueries should return the same number of queries")
				}

				queries, err := common.RenderQueries(log, tqueries, tmpValues)
				if tt.wantRenderErr && err == nil {
					t.Fatal("RenderQueries should failed")
				} else if !tt.wantRenderErr && err != nil {
					t.Fatal("RenderQueries should not failed")
				}

				if err == nil {
					if len(queries) != len(tt.queries) {
						t.Fatal("RenderQueries should return the same number of queries")
					}

					for i, e := range tt.expected {
						if !queries[i].CompareSkipQuery(e.query, e.inverted, e.value) {
							t.Errorf("%s was not correctly parsed or rendered: %v", tt.queries[i], queries[i])
						}
					}
				}
			}
		})
	}
}
