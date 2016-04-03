package api

// import "os"
// import def "github.com/peaberberian/GoBanks/database/definitions"
// import "github.com/peaberberian/GoBanks/file/qif"

// // just for tests
// func AddQifFile(filePath string, dateFormat string,
// 	db def.GoBanksDataBase, accountId int) (err error) {
// 	var f *os.File
// 	f, err = os.Open(filePath)
// 	if err != nil {
// 		return
// 	}
// 	ts, err := qif.ParseFile(f, accountId, dateFormat)
// 	if err != nil {
// 		return
// 	}
// 	for _, t := range ts {
// 		_, err = db.AddTransaction(t)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return nil
// }
