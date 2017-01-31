package libovndb

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"math/rand"
	"strconv"

	"github.com/golang/glog"
	"github.com/socketplane/libovsdb"
	"time"
)

func init() {
	rand.Seed(MAX_TRANSACTION)
}

type OVNRow map[string]interface{}

func newNBCtlImp(client *ovnDBClient) *ovnDBImp {
	nbimp := ovnDBImp{client: client}
	nbimp.cache = make(map[string]map[string]libovsdb.Row)
	initial, err := nbimp.client.dbclient.MonitorAll(NBDB, "")
	if err != nil {
		glog.Fatalf("OVN DB monitor failed: %v", err)
		os.Exit(1)
	}
	nbimp.populateCache(*initial)
	notifier := ovnNotifier{&nbimp}
	nbimp.client.dbclient.Register(notifier)
	return &nbimp
}

func (odbi *ovnDBImp) lswListImp() *OvnCommand {
	condition := libovsdb.NewCondition("name", "!=", "")
	listOp := libovsdb.Operation{
		Op:    list,
		Table: LSWITCH,
		Where: []interface{}{condition},
	}

	operations := []libovsdb.Operation{listOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) lswAddImp(lsw string) *OvnCommand {
	namedUUID := "gopher_lsw_add" + strconv.Itoa(rand.Int())
	//row to insert
	lswitch := make(OVNRow)
	lswitch["name"] = lsw

	if odbi.getRowUUID(LSWITCH, lswitch) != "" {
		glog.V(OVNLOGLEVEL).Info("The logic switch existed, and will get nil command")
		return nil
	}

	// simple insert operation
	insertOp := libovsdb.Operation{
		Op:       insert,
		Table:    LSWITCH,
		Row:      lswitch,
		UUIDName: namedUUID,
	}
	operations := []libovsdb.Operation{insertOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) lswDelImp(lsw string) *OvnCommand {
	condition := libovsdb.NewCondition("name", "==", lsw)
	// simple insert operation
	delOp := libovsdb.Operation{
		Op:    del,
		Table: LSWITCH,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{delOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) getRowUUID(table string, row OVNRow) string {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for uuid, drows := range odbi.cache[table] {
		found := false
		for field, value := range row {
			if v, ok := drows.Fields[field]; ok {
				if v == value {
					found = true
				} else {
					found = false
					break
				}
			}
		}
		if found {
			return uuid
		}
	}
	return ""
}

func (odbi *ovnDBImp) getACLUUIDByRow(lsw, table string, row OVNRow) string {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for _, drows := range odbi.cache[LSWITCH] {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			acls := drows.Fields["acls"]
			if acls != nil {
				switch acls.(type) {
				case libovsdb.OvsSet:
					if as, ok := acls.(libovsdb.OvsSet); ok {
						for _, a := range as.GoSet {
							if va, ok := a.(libovsdb.UUID); ok {
								for field, value := range(row) {
									switch field {
									case "action":
										if odbi.cache[ACLS][va.GoUUID].Fields["action"].(string) != value {
											goto unmatched
										}
									case "direction":
										if odbi.cache[ACLS][va.GoUUID].Fields["direction"].(string) != value {
											goto unmatched
										}
									case "match":
										if odbi.cache[ACLS][va.GoUUID].Fields["match"].(string) != value {
											goto unmatched
										}
									case "priority":
										//Todo: Priority workaround: should know why it is float64, tracked at NTWK-2549
										if odbi.cache[ACLS][va.GoUUID].Fields["priority"].(float64) != value {
											goto unmatched
										}
									}

								}
								return va.GoUUID
							}
							unmatched:
						}
						return ""
					}
				case libovsdb.UUID:
					if va, ok := acls.(libovsdb.UUID); ok {
						for field, value := range(row) {
							switch field {
							case "action":
								if odbi.cache[ACLS][va.GoUUID].Fields["action"].(string) != value {
									goto out
								}
							case "direction":
								if odbi.cache[ACLS][va.GoUUID].Fields["direction"].(string) != value {
									goto out
								}
							case "match":
								if odbi.cache[ACLS][va.GoUUID].Fields["match"].(string) != value {
									goto out
								}
							case "priority":
								//Todo: Priority workaround: should know why it is float64, tracked at NTWK-2549
								if odbi.cache[ACLS][va.GoUUID].Fields["priority"].(float64) != value {
									goto out
								}
							}
						}
						return va.GoUUID
						out:
					}
				}
			}
		}
	}
	return ""
}

func (odbi *ovnDBImp) getRowUUIDContainsUUID(table, field, uuid string) string {
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for id, drows := range odbi.cache[table] {
		v := fmt.Sprintf("%s", drows.Fields[field])
		if strings.Contains(v, uuid) {
			return id
		}
	}
	return ""
}

func (odbi *ovnDBImp) lspAddImp(lsw, lsp string) *OvnCommand {
	namedUUID := "gopher_lsp_add" + strconv.Itoa(rand.Int())
	lsprow := make(OVNRow)
	lsprow["name"] = lsp

	if odbi.getRowUUID(LPORT, lsprow) != "" {
		glog.V(OVNLOGLEVEL).Info("The logic port existed, and will get nil command")
		return nil
	}

	// simple insert operation
	insertOp := libovsdb.Operation{
		Op:       insert,
		Table:    LPORT,
		Row:      lsprow,
		UUIDName: namedUUID,
	}

	mutateUUID := []libovsdb.UUID{{namedUUID}}
	mutateSet, _ := libovsdb.NewOvsSet(mutateUUID)
	mutation := libovsdb.NewMutation("ports", insert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsw)

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        mutate,
		Table:     LSWITCH,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) lspDelImp(lsp string) *OvnCommand {
	lsprow := make(OVNRow)
	lsprow["name"] = lsp
	lspUUID := odbi.getRowUUID(LPORT, lsprow)
	mutateUUID := []libovsdb.UUID{{lspUUID}}
	condition := libovsdb.NewCondition("name", "==", lsp)
	delOp := libovsdb.Operation{
		Op:    del,
		Table: LPORT,
		Where: []interface{}{condition},
	}
	mutateSet, _ := libovsdb.NewOvsSet(mutateUUID)
	mutation := libovsdb.NewMutation("ports", del, mutateSet)
	mucondition := libovsdb.NewCondition("_uuid", "==", libovsdb.UUID{odbi.getRowUUIDContainsUUID(LSWITCH, "ports", lspUUID)})
	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        mutate,
		Table:     LSWITCH,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{mucondition},
	}
	operations := []libovsdb.Operation{delOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) lspSetAddressImp(lsp string, addr ...string) *OvnCommand {
	var addresses []string
	for _, a := range addr {
		addresses = append(addresses, a)
	}

	mutateSet, _ := libovsdb.NewOvsSet(addresses)
	mutation := libovsdb.NewMutation("addresses", insert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsp)

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        mutate,
		Table:     LPORT,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) aclAddImp(lsw, direct, match, action string, priority int, external_ids *map[string]string, logflag bool) *OvnCommand {
	if external_ids != nil {
		glog.Fatalf("External id is not supprted in acl setting now")
	}

	namedUUID := "gopher_acl_add" + strconv.Itoa(rand.Int())
	aclrow := make(OVNRow)
	aclrow["action"] = action
	aclrow["direction"] = direct
	aclrow["match"] = match
	//Todo: Priority workaround: should know why it is float64, tracked at NTWK-2549
	aclrow["priority"] = float64(priority)
	if odbi.getACLUUIDByRow(lsw, ACLS, aclrow) != "" {
		glog.V(OVNLOGLEVEL).Info("The acl existed, and will get nil command")
		return nil
	}
	aclrow["log"] = logflag
	// simple insert operation
	insertOp := libovsdb.Operation{
		Op:       insert,
		Table:    ACLS,
		Row:      aclrow,
		UUIDName: namedUUID,
	}

	mutateUUID := []libovsdb.UUID{{namedUUID}}
	mutateSet, _ := libovsdb.NewOvsSet(mutateUUID)
	mutation := libovsdb.NewMutation("acls", insert, mutateSet)
	condition := libovsdb.NewCondition("name", "==", lsw)

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        mutate,
		Table:     LSWITCH,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{insertOp, mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}


func (odbi *ovnDBImp) aclDelImp(lsw, direct, match string, priority int) *OvnCommand {
	aclrow := make(OVNRow)

	wherecondition := []interface{}{}
	if direct != "" {
		directcondition := libovsdb.NewCondition("direction", "==", direct)
		wherecondition = append(wherecondition, directcondition)
		aclrow["direction"] = direct
	}
	if match != "" {
		matchcondition := libovsdb.NewCondition("match", "==", match)
		wherecondition = append(wherecondition, matchcondition)
		aclrow["match"] = match
	}
	if priority != 0 {
		pricondition := libovsdb.NewCondition("priority", "==", priority)
		wherecondition = append(wherecondition, pricondition)
		//Todo: Priority workaround: should know why it is float64, tracked at NTWK-2549
		aclrow["priority"] = float64(priority)
	}

	aclUUID := odbi.getACLUUIDByRow(lsw, ACLS, aclrow)
	if aclUUID == "" {
		glog.V(OVNLOGLEVEL).Info("The deleting acl not found in cache, and will get nil command")
		return nil
	}

	delOp := libovsdb.Operation{
		Op:    del,
		Table: ACLS,
		Where: wherecondition,
	}
/*
	mutateUUID := []libovsdb.UUID{{aclUUID}}
	mutateSet, _ := libovsdb.NewOvsSet(mutateUUID)*/
	mutation := libovsdb.NewMutation("acls", del, libovsdb.UUID{aclUUID})
	condition := libovsdb.NewCondition("name", "==", lsw)

	// Simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        mutate,
		Table:     LSWITCH,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp, delOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) ASAdd(name string, addrs []string) *OvnCommand {
	namedUUID := "gopher_as_add" + strconv.Itoa(rand.Int())
	adrow := make(OVNRow)
	adrow["name"] = name
	addresses, _ := libovsdb.NewOvsSet(addrs)
	adrow["addresses"] = addresses

	// simple insert operation
	Op := libovsdb.Operation{
		Op:       insert,
		Table:    Address_Set,
		Row:      adrow,
		UUIDName: namedUUID,
	}

	adnrow := make(OVNRow)
	adnrow["name"] = name

	if odbi.getRowUUID(Address_Set, adnrow) != "" {
		mutations := []interface{}{}
		mutation := libovsdb.NewMutation("addresses", "delete", " ")
		mutations = append(mutations, mutation)
		for _, ad := range addrs {
			mutation := libovsdb.NewMutation("addresses", "insert", ad)
			mutations = append(mutations, mutation)
		}
		condition := libovsdb.NewCondition("name", "==", name)

		// simple mutate operation
		Op = libovsdb.Operation{
			Op:        mutate,
			Table:     Address_Set,
			Mutations: mutations,
			Where:     []interface{}{condition},
		}
	}
	operations := []libovsdb.Operation{Op}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) ASDel(name string) *OvnCommand {
	condition := libovsdb.NewCondition("name", "==", name)
	delOp := libovsdb.Operation{
		Op:    del,
		Table: Address_Set,
		Where: []interface{}{condition},
	}
	operations := []libovsdb.Operation{delOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) LSSetOpt(lsp string, options map[string]string) *OvnCommand {
	mutatemap, _ := libovsdb.NewOvsMap(options)
	mutation := libovsdb.NewMutation("options", insert, mutatemap)
	condition := libovsdb.NewCondition("name", "==", lsp)

	// simple mutate operation
	mutateOp := libovsdb.Operation{
		Op:        mutate,
		Table:     LPORT,
		Mutations: []interface{}{mutation},
		Where:     []interface{}{condition},
	}
	operations := []libovsdb.Operation{mutateOp}
	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}
}

func (odbi *ovnDBImp) transact(ops ...libovsdb.Operation) ([]libovsdb.OperationResult, error) {
	// Only support one trans at same time now.
	odbi.tranmutex.Lock()
	defer odbi.tranmutex.Unlock()
	reply, err := odbi.client.dbclient.Transact(NBDB, ops...)

	if err != nil {
		return reply, err
	}

	if len(reply) < len(ops) {
		glog.V(OVNLOGLEVEL).Info("Number of Replies should be atleast equal to number of operations")
		for i, o := range reply {
			if o.Error != "" && i < len(ops) {
				glog.V(OVNLOGLEVEL).Info("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", ops[i])
				return nil, errors.New(fmt.Sprint("Transaction Failed due to an error :", o.Error, " details:", o.Details, " in ", ops[i]))
			}
		}
		return reply, errors.New(fmt.Sprint("Number of Replies should be atleast equal to number of operations"))
	}
	glog.V(OVNLOGLEVEL).Info("transaction reply : ", reply)
	return reply, nil
}

func (odbi *ovnDBImp) Execute(cmds ...*OvnCommand) error {
	if cmds == nil {
		glog.V(OVNLOGLEVEL).Infof("Inputting command is nil, will skip transaction.")
		return nil
	}
	var ops []libovsdb.Operation
	for _, cmd := range cmds {
		if cmd != nil {
			ops = append(ops, cmd.Operations...)
		}
	}
	reply, err := odbi.transact(ops...)
	glog.V(OVNLOGLEVEL).Infof("OVN replys: %v", reply)
	if err != nil {
		return err
	} /*else {
	    j := 0
	    index := 0
	    for i, _ := range reply {
	        if index >= len(cmds[j].Operations) {
	            index = 0
	            j++
	        }
	        cmds[j].Results[index] = reply[i].Rows
	        index++
	    }
	}*/
	//wait for cache updating. it's not efficient.
	//Todo: track at https://jirap.corp.ebay.com/browse/NTWK-2550
	time.Sleep(1 * time.Second)
	return nil
}

func (odbi *ovnDBImp) populateCache(updates libovsdb.TableUpdates) {
	glog.V(OVNLOGLEVEL).Info("New nofity arrived")
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for table, tableUpdate := range updates.Updates {
		if _, ok := odbi.cache[table]; !ok {
			odbi.cache[table] = make(map[string]libovsdb.Row)
		}
		for uuid, row := range tableUpdate.Rows {
			empty := libovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				odbi.cache[table][uuid] = row.New
			} else {
				delete(odbi.cache[table], uuid)
			}
		}
	}
}

func (odbi *ovnDBImp) ConvertGoSetToStringArray(oset libovsdb.OvsSet) []string {
	var ret = []string{}
	for _, s := range(oset.GoSet) {
		value, ok :=  s.(string)
		if ok {
			ret = append(ret, value)
		}
	}
	return ret
}

// Get all lport by lswitch
func (odbi *ovnDBImp) GetLogicPortsBySwitch(lsw string) []*LogcalPort {
	var lplist = []*LogcalPort{}
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for _, drows := range odbi.cache[LSWITCH] {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			ports := drows.Fields["ports"]
			if ports != nil {
				switch ports.(type) {
				case libovsdb.OvsSet:
					if ps, ok := ports.(libovsdb.OvsSet); ok {
						for _, p := range ps.GoSet {
							if vp, ok := p.(libovsdb.UUID); ok {
								tp := &LogcalPort{
									UUID: vp.GoUUID,
									Name: odbi.cache[LPORT][vp.GoUUID].Fields["name"].(string),
								}
								addr := odbi.cache[LPORT][vp.GoUUID].Fields["addresses"]
								switch addr.(type) {
								case string:
									tp.Addresses = []string{addr.(string)}
								case libovsdb.OvsSet:
									tp.Addresses = odbi.ConvertGoSetToStringArray(addr.(libovsdb.OvsSet))
								default:
									glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport address.")
								}
								lplist = append(lplist, tp)
							}
						}
					} else {
						glog.V(OVNLOGLEVEL).Info("Type libovsdb.OvsSet casting failed.")
					}
				case libovsdb.UUID:
					if vp, ok := ports.(libovsdb.UUID); ok {
						tp := &LogcalPort{
							UUID: vp.GoUUID,
							Name:      odbi.cache[LPORT][vp.GoUUID].Fields["name"].(string),
						}
						addr := odbi.cache[LPORT][vp.GoUUID].Fields["addresses"]
						switch addr.(type) {
						case string:
							tp.Addresses = []string{addr.(string)}
						case libovsdb.OvsSet:
							tp.Addresses = odbi.ConvertGoSetToStringArray(addr.(libovsdb.OvsSet))
						default:
							glog.V(OVNLOGLEVEL).Info("Unsupport type found in lport address.")
						}
						lplist = append(lplist, tp)
					} else {
						glog.V(OVNLOGLEVEL).Info("Type libovsdb.UUID casting failed.")
					}
				default:
					glog.V(OVNLOGLEVEL).Info("Unsupport type found in ovsdb rows.")
				}
			}
			break
		}
	}
	return lplist
}

// Get all acl by lswitch
func (odbi *ovnDBImp) GetACLsBySwitch(lsw string) []*ACL {
	//TODO: should be improvement here, when have lots of acls.
	acllist := make([]*ACL, 0, 0)
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for _, drows := range odbi.cache[LSWITCH] {
		if rlsw, ok := drows.Fields["name"].(string); ok && rlsw == lsw {
			acls := drows.Fields["acls"]
			if acls != nil {
				switch acls.(type) {
				case libovsdb.OvsSet:
					if as, ok := acls.(libovsdb.OvsSet); ok {
						for _, a := range as.GoSet {
							if va, ok := a.(libovsdb.UUID); ok {
								ta := &ACL{
									UUID: va.GoUUID,
									Action:    odbi.cache[ACLS][va.GoUUID].Fields["action"].(string),
									Direction: odbi.cache[ACLS][va.GoUUID].Fields["direction"].(string),
									Match:     odbi.cache[ACLS][va.GoUUID].Fields["match"].(string),
									//Todo: Priority workaround: should know why it is float64, tracked at NTWK-2549
									Priority:  int(odbi.cache[ACLS][va.GoUUID].Fields["priority"].(float64)),
								}
								acllist = append(acllist, ta)
							}
						}
					}
				case libovsdb.UUID:
					if va, ok := acls.(libovsdb.UUID); ok {
						ta := &ACL{
							UUID: va.GoUUID,
							Action:    odbi.cache[ACLS][va.GoUUID].Fields["action"].(string),
							Direction: odbi.cache[ACLS][va.GoUUID].Fields["direction"].(string),
							Match:     odbi.cache[ACLS][va.GoUUID].Fields["match"].(string),
							//Todo: Priority workaround: should know why it is float64, tracked at NTWK-2549
							Priority:  int(odbi.cache[ACLS][va.GoUUID].Fields["priority"].(float64)),
						}
						acllist = append(acllist, ta)
					}
				}
			}
			break
		}
	}
	return acllist
}

// Get all addressset
func (odbi *ovnDBImp) GetAddressSets() []*AddressSet {
	adlist := make([]*AddressSet, 0, 0)
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()
	for uuid, drows := range odbi.cache[Address_Set] {
		ta := &AddressSet{
			UUID: uuid,
			Name: drows.Fields["name"].(string),
		}
		addresses := []string{}
		//TODO: is it possible return interface type directly instead of GoSet
		//TODO: Tracked at NTWK-2524
		if goset, ok := drows.Fields["addresses"].(libovsdb.OvsSet); ok {
			for _, i := range goset.GoSet {
				addresses = append(addresses, i.(string))
			}
		}
		ta.Addresses = addresses

		adlist = append(adlist, ta)
	}
	return adlist
}
