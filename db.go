package main

// func saveAllProxy(mProxy *mapProxy) {
// 	debugmsg("start saveAllProxy")
// 	var u, i int64
// 	// prepareInsert, _ := insertProxy()
// 	// prepareUpdate, _ := updateProxy()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		errmsg("saveAllProxy Begin", err)
// 		return
// 	}
// 	for _, p := range mProxy.values {
// 		if p.Update {
// 			u++
// 			_, err := tx.Exec(`
// 				UPDATE proxies SET
// 					host       = $2,
// 					port       = $3,
// 					work       = $4,
// 					anon       = $5,
// 					checks     = $6,
// 					create_at  = $7,
// 					update_at  = $8,
// 					response   = $9
// 				WHERE
// 					hostname = $1
// 				`,
// 				&p.Hostname,
// 				&p.Host,
// 				&p.Port,
// 				&p.IsWork,
// 				&p.IsAnon,
// 				&p.Checks,
// 				&p.CreateAt,
// 				&p.UpdateAt,
// 				&p.Response,
// 			)
// 			chkErr("saveAllProxy Update", err)
// 		}
// 		if p.Insert {
// 			i++
// 			_, err := tx.Exec(`
// 				INSERT INTO proxies (
// 					hostname,
// 					host,
// 					port,
// 					work,
// 					anon,
// 					checks,
// 					create_at,
// 					update_at,
// 					response
// 				) VALUES (
// 					$1,
// 					$2,
// 					$3,
// 					$4,
// 					$5,
// 					$6,
// 					$7,
// 					$8,
// 					$9
// 				)`,
// 				&p.Hostname,
// 				&p.Host,
// 				&p.Port,
// 				&p.IsWork,
// 				&p.IsAnon,
// 				&p.Checks,
// 				&p.CreateAt,
// 				&p.UpdateAt,
// 				&p.Response,
// 			)
// 			if err != nil {
// 				errmsg("saveAllLinks Insert", err)
// 			}
// 		}
// 	}
// 	chkErr("saveAllProxy commit", tx.Commit())
// 	debugmsg("update proxy", u)
// 	debugmsg("insert proxy", i)
// 	debugmsg("end getAllProxy")
// }

// func updateAllProxy(, mProxy *mapProxy) {
// 	debugmsg("start updateAllProxy")
// 	stmt, err := db.Prepare(`
// 		UPDATE proxies SET
// 			host       = $2,
// 			port       = $3,
// 			work       = $4,
// 			anon       = $5,
// 			checks     = $6,
// 			create_at  = $7,
// 			update_at  = $8,
// 			response   = $9
// 		WHERE
// 			hostname = $1
// 	`)
// 	if err != nil {
// 		errmsg("updateAllProxy Prepare", err)
// 		return
// 	}
// 	defer stmt.Close()
// 	for _, p := range mProxy.values {
// 		_, err := stmt.Exec(
// 			&p.Hostname,
// 			&p.Host,
// 			&p.Port,
// 			&p.IsWork,
// 			&p.IsAnon,
// 			&p.Checks,
// 			&p.CreateAt,
// 			&p.UpdateAt,
// 			&p.Response,
// 		)
// 		if err != nil {
// 			errmsg("updateAllProxy Exec", err)
// 		}
// 	}
// 	debugmsg("end updateAllProxy, update proxy", len(mProxy.values))
// }
