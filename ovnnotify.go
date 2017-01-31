package libovndb

import (
	"github.com/socketplane/libovsdb"
)

type ovnNotifier struct {
	odbi *ovnDBImp
}

func (notify ovnNotifier) Update(context interface{}, tableUpdates libovsdb.TableUpdates) {
	notify.odbi.populateCache(tableUpdates)
}
func (notify ovnNotifier) Locked([]interface{}) {
}
func (notify ovnNotifier) Stolen([]interface{}) {
}
func (notify ovnNotifier) Echo([]interface{}) {
}
func (notify ovnNotifier) Disconnected(client *libovsdb.OvsdbClient) {
}
