package crud

import "strings"

func CrudTemplate() string {

	return `
// Code generated by 'freedom new-po'
package po
{{.Import}}
{{.Content}}

// TakeChanges .
func (obj *{{.Name}})TakeChanges() map[string]interface{} {
	if obj.changes == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range obj.changes {
		result[k] = v
	}
	obj.changes = nil
	return result
}

// updateChanges .
func (obj *{{.Name}}) setChanges(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

{{range .Fields}}
// Set{{.Value}} .
func (obj *{{.StructName}}) Set{{.Value}} ({{.Arg}} {{.Type}}) {
	obj.{{.Value}} = {{.Arg}} 
	obj.setChanges("{{.Column}}", {{.Arg}})
}
{{ end }}

{{range .NumberFields}}
// Add{{.Value}} .
func (obj *{{.StructName}}) Add{{.Value}} ({{.Arg}} {{.Type}}) {
	obj.{{.Value}} += {{.Arg}} 
	obj.setChanges("{{.Column}}", gorm.Expr("{{.Column}} + ?", {{.Arg}}))
}
{{ end }}
`
}

func FunTemplatePackage() string {
	source := `
	// Code generated by 'freedom new-po'
	package repository
	import (
		"github.com/8treenet/freedom"
		"github.com/jinzhu/gorm"
		"time"
		"{{.PackagePath}}"
		"fmt"
		"strings"
	)
	
	// GORMRepository .
	type GORMRepository interface {
		db() *gorm.DB
		GetWorker() freedom.Worker
	}

	// NewORMDescBuilder .
	func NewORMDescBuilder(column string, columns ...string) *Reorder {
		return newReorder("desc", column, columns...)
	}

	// NewORMAscBuilder .
	func NewORMAscBuilder(column string, columns ...string) *Reorder {
		return newReorder("asc", column, columns...)
	}

	// NewORMBuilder .
	func NewORMBuilder() *Builder {
		return &Builder{}
	}

	// NewDescOrder .
	func newReorder(sort, field string, args ...string) *Reorder {
		fields := []string{field}
		fields = append(fields, args...)
		orders := []string{}
		for index := 0; index < len(fields); index++ {
			orders = append(orders, sort)
		}
		return &Reorder{
			fields: fields,
			orders: orders,
		}
	}
	// Reorder .
	type Reorder struct {
		fields []string
		orders []string
	}
	
	// NewPageBuilder .
	func (o *Reorder) NewPageBuilder(page, pageSize int) *Builder {
		pager := new(Builder)
		pager.reorder = o
		pager.page = page
		pager.pageSize = pageSize
		return pager
	}
	
	// NewBuilder .
	func (o *Reorder) NewBuilder() *Builder {
		pager := new(Builder)
		pager.reorder = o
		return pager
	}
	
	// Order .
	func (o *Reorder) Order() interface{} {
		args := []string{}
		for index := 0; index < len(o.fields); index++ {
			args = append(args, fmt.Sprintf("$$wave%s$$wave %s", o.fields[index], o.orders[index]))
		}
	
		return strings.Join(args, ",")
	}
	
	// Builder
	type Builder struct {
		reorder       *Reorder
		pageSize      int
		page          int
		totalPage     int
		selectColumn  []string
	}
	
	// TotalPage .
	func (p *Builder) TotalPage() int {
		return p.totalPage
	}
	
	func (b *Builder) Order() interface{} {
		if b.reorder != nil {
			return b.Order()
		}
		return ""
	}
	
	// Execute .
	func (p *Builder) Execute(db *gorm.DB, object interface{}) (e error) {
		pageFind := false
		if p.reorder != nil {
			db = db.Order(p.reorder.Order())
		} else {
			db = db.Set("gorm:order_by_primary_key", "DESC")
		}
		if p.page != 0 && p.pageSize != 0 {
			pageFind = true
			db = db.Offset((p.page - 1) * p.pageSize).Limit(p.pageSize)
		}
	
		if len(p.selectColumn) > 0 {
			db = db.Select(p.selectColumn)
		}
	
		resultDB := db.Find(object)
		if resultDB.Error != nil {
			return resultDB.Error
		}
	
		if !pageFind {
			return
		}
	
		var count int
		e = resultDB.Offset(0).Limit(1).Count(&count).Error
		if e == nil && count != 0 {
			//计算分页
			if count%p.pageSize == 0 {
				p.totalPage = count / p.pageSize
			} else {
				p.totalPage = count/p.pageSize + 1
			}
		}
		return
	}
	
	// SetPage .
	func (b *Builder) SetPage(page, pageSize int) *Builder {
		b.page = page
		b.pageSize = pageSize
		return b
	}
	
	// SelectColumn .
	func (b *Builder) SelectColumn(column ...string) *Builder {
		b.selectColumn = append(b.selectColumn, column...)
		return b
	}
	
	func ormErrorLog(repo GORMRepository, model, method string, e error, expression ...interface{}) {
		if e == nil || e == gorm.ErrRecordNotFound {
			return
		}
		repo.GetWorker().Logger().Errorf("Orm error, model: %s, method: %s, expression :%v, reason for error:%v", model, method, expression, e)
	}
`
	return strings.ReplaceAll(source, "$$wave", "`")
}
func FunTemplate() string {
	return `
	// find{{.Name}} .
	func find{{.Name}}(repo GORMRepository, result interface{}, builders ...*Builder) (e error) {
		now := time.Now()
		defer func() {
			freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}", e, now)
			ormErrorLog(repo, "{{.Name}}", "find{{.Name}}", e, result)
		}()
		db := repo.db()
		if len(builders) == 0 {
			e = db.Where(result).Last(result).Error
			return
		}
		e = builders[0].Execute(db.Limit(1), result)
		return
	}
	
	// find{{.Name}}ListByPrimarys .
	func find{{.Name}}ListByPrimarys(repo GORMRepository, results interface{}, primarys ...interface{}) (e error) {
		now := time.Now()
		e = repo.db().Find(results, primarys).Error
		freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}ListByPrimarys", e, now)
		ormErrorLog(repo, "{{.Name}}", "find{{.Name}}sByPrimarys", e, primarys)
		return
	}
	
	// find{{.Name}}ByWhere .
	func find{{.Name}}ByWhere(repo GORMRepository, query string, args []interface{}, result interface{}, builders ...*Builder) (e error) {
		now := time.Now()
		defer func() {
			freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}ByWhere", e, now)
			ormErrorLog(repo, "{{.Name}}", "find{{.Name}}ByWhere", e, query, args)
		}()
		db := repo.db()
		if query != "" {
			db = db.Where(query, args...)
		}
		if len(builders) == 0 {
			e = db.Last(result).Error
			return
		}
	
		e = builders[0].Execute(db.Limit(1), result)
		return
	}
	
	// find{{.Name}}ByMap .
	func find{{.Name}}ByMap(repo GORMRepository, query map[string]interface{}, result interface{}, builders ...*Builder) (e error) {
		now := time.Now()
		defer func() {
			freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}ByMap", e, now)
			ormErrorLog(repo, "{{.Name}}", "find{{.Name}}ByMap", e, query)
		}()

		db := repo.db().Where(query)
		if len(builders) == 0 {
			e = db.Last(result).Error
			return
		}
	
		e = builders[0].Execute(db.Limit(1), result)
		return
	}
	
	// find{{.Name}}List .
	func find{{.Name}}List(repo GORMRepository, query po.{{.Name}}, results interface{}, builders ...*Builder) (e error) {
		now := time.Now()
		defer func() {
			freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}List", e, now)
			ormErrorLog(repo, "{{.Name}}", "find{{.Name}}s", e, query)
		}()
		db := repo.db().Where(query)
	
		if len(builders) == 0 {
			e = db.Find(results).Error
			return
		}
		e = builders[0].Execute(db, results)
		return
	}
	
	// find{{.Name}}ListByWhere .
	func find{{.Name}}ListByWhere(repo GORMRepository, query string, args []interface{}, results interface{}, builders ...*Builder) (e error) {
		now := time.Now()
		defer func() {
			freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}ListByWhere", e, now)
			ormErrorLog(repo, "{{.Name}}", "find{{.Name}}sByWhere", e, query, args)
		}()
		db := repo.db()
		if query != "" {
			db = db.Where(query, args...)
		}
	
		if len(builders) == 0 {
			e = db.Find(results).Error
			return
		}
		e = builders[0].Execute(db, results)
		return
	}
	
	// find{{.Name}}ListByMap .
	func find{{.Name}}ListByMap(repo GORMRepository, query map[string]interface{}, results interface{}, builders ...*Builder) (e error) {
		now := time.Now()
		defer func() {
			freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "find{{.Name}}ListByMap", e, now)
			ormErrorLog(repo, "{{.Name}}", "find{{.Name}}sByMap", e, query)
		}()

		db := repo.db().Where(query)
	
		if len(builders) == 0 {
			e = db.Find(results).Error
			return
		}
		e = builders[0].Execute(db, results)
		return
	}
	
	// create{{.Name}} .
	func create{{.Name}}(repo GORMRepository, object *po.{{.Name}}) (rowsAffected int64, e error) {
		now := time.Now()
		db := repo.db().Create(object)
		rowsAffected = db.RowsAffected
		e = db.Error
		freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "create{{.Name}}", e, now)
		ormErrorLog(repo, "{{.Name}}", "create{{.Name}}", e, *object)
		return
	}

	// save{{.Name}} .
	func save{{.Name}}(repo GORMRepository, object *po.{{.Name}}) (affected int64, e error) {
		now := time.Now()
		db := repo.db().Model(object).Updates(object.TakeChanges())
		e = db.Error
		affected = db.RowsAffected
		freedom.Prometheus().OrmWithLabelValues("{{.Name}}", "save{{.Name}}", e, now)
		ormErrorLog(repo, "{{.Name}}", "save{{.Name}}", e, *object)
		return
	}
`
}
